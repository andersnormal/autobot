package plugins

import (
	"context"
	"sync"
	"time"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// Event symbolizes events that can occur in the plugin.
// Though they are generally triggered from the sever.
type Event int32

const (
	ReadyEvent         Event = 0
	RefreshConfigEvent Event = 1
)

// Plugin describes a plugin in general.
// It should not be instantiated directly.
type Plugin struct {
	opts *Opts

	sc stan.Conn
	nc *nats.Conn

	ready          chan struct{}
	events         chan Event
	eventsChannels []chan Event
	resp           string

	env runtime.Env

	ctx     context.Context
	cancel  func()
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup

	pp *pb.Plugin

	sync.RWMutex
}

// Opt is an Option
type Opt func(*Opts)

// Opts are the available options
type Opts struct {
}

// SubscribeFunc ...
type SubscribeFunc = func(*pb.Event) (*pb.Event, error)

// WithContext is creating a new plugin and a context to run operations in routines.
// When the context is canceled, all concurrent operations are canceled.
func WithContext(ctx context.Context, env runtime.Env, opts ...Opt) (*Plugin, context.Context) {
	p := newPlugin(opts...)

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.ctx = ctx
	p.env = env

	return p, ctx
}

// Events returns a channel which is used to publish important
// events to the plugin. For example that the server pushes a new config.
// or that the plugin needs to restart.
func (p *Plugin) Events() <-chan Event {
	out := make(chan Event)

	p.eventsChannels = append(p.eventsChannels, out)

	return out
}

func (p *Plugin) multiplexEvents() func() error {
	return func() error {
		for {
			select {
			case e := <-p.events:
				for _, c := range p.eventsChannels {
					// this pushes the events to a routine
					// to avoid the loop to block in execution
					p.Run(func() error {
						c <- e

						return nil
					})
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

// newPlugin ...
func newPlugin(opts ...Opt) *Plugin {
	options := new(Opts)

	p := new(Plugin)
	// setting a default env for a plugin
	p.opts = options
	p.ready = make(chan struct{})
	p.events = make(chan Event)

	// configure plugin
	configure(p, opts...)

	// starts the multiplexder
	p.Run(p.multiplexEvents())

	return p
}

// SubscribeInbox is subscribing to the inbox of messages.
// This if for plugins that want to consume the message that other
// plugins publish to Autobot.
func (p *Plugin) SubscribeInbox(opts ...FilterOpt) <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(p.subInboxFunc(sub, opts...))

	return sub
}

// SubscribeOutbox is subscribing to the outbox of messages.
// These are the message that ought to be published to an external service (e.g. Slack, MS Teams).
func (p *Plugin) SubscribeOutbox(opts ...FilterOpt) <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(p.subOutboxFunc(sub, opts...))

	return sub
}

// PublishInbox is publishing message to the inbox in the controller.
// The returned channel pushes all of the send message to the inbox in the controller.
func (p *Plugin) PublishInbox(opts ...FilterOpt) chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(p.pubInboxFunc(pub, append(DefaultInboxFilterOpts, opts...)...))

	return pub
}

// PublishOutbox is publishing message to the outbox in the controller.
// The returned channel pushes all the send messages to the outbox in the controller.
func (p *Plugin) PublishOutbox(opts ...FilterOpt) chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(p.pubOutboxFunc(pub, append(DefaultOutboxFilterOpts, opts...)...))

	return pub
}

// Wait is waiting for the underlying WaitGroup. All run go routines are
// hold here to before one exists with an error. Then the provided context is canceled.
func (p *Plugin) Wait() error {
	p.wg.Wait()

	if p.cancel != nil {
		p.cancel()
	}

	if p.sc != nil {
		p.sc.Close()
	}

	if p.nc != nil {
		p.nc.Close()
	}

	return p.err
}

// ReplyWithFunc is a wrapper function to provide a function which may
// send replies to received message for this plugin.
func (p *Plugin) ReplyWithFunc(fn SubscribeFunc, opts ...FilterOpt) error {
	p.Run(func() error {
		// create publish channel ...
		pubReply := p.PublishOutbox()
		subMsg := p.SubscribeInbox()

		for {
			select {
			case e, ok := <-subMsg:
				if !ok {
					return nil
				}

				r, err := fn(e)
				if err != nil {
					return err
				}

				pubReply <- r
			case <-p.ctx.Done():
				return nil
			}
		}
	})

	return nil
}

// Run is a wrapper similar to errgroup to schedule functions
// as go routines in a waitGroup.
func (p *Plugin) Run(f func() error) {
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		if err := f(); err != nil {
			p.errOnce.Do(func() {
				p.err = err
				if p.cancel != nil {
					p.cancel()
				}
			})
		}
	}()
}

func (p *Plugin) getConn() (stan.Conn, error) {
	p.Lock()
	defer p.Unlock()

	if p.sc != nil {
		return p.sc, nil
	}

	c, err := p.newConn()
	if err != nil {
		return nil, err
	}

	p.sc = c

	// start watcher ...
	p.Run(p.watchAction())

	if err := p.register(); err != nil {
		return nil, err
	}

	select {
	case <-time.After(2 * time.Second):
		return nil, ErrPluginRegister
	case <-p.ready:
		p.events <- ReadyEvent
	}

	return p.sc, nil
}

func (p *Plugin) newConn() (stan.Conn, error) {
	nc, err := nats.Connect(
		p.env.ClusterURL,
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	)
	if err != nil {
		return nil, err
	}

	p.nc = nc

	sc, err := stan.Connect(p.env.ClusterID, p.pp.GetName(), stan.SetConnectionLostHandler(p.lostHandler()))
	if err != nil {
		return nil, err
	}

	// creates a new re-used response inbox for the plugin
	p.resp = p.nc.NewRespInbox()

	return sc, nil
}

func (p *Plugin) watchAction() func() error {
	return func() error {
		exit := make(chan error)
		resp := make(chan *pb.Action)

		sub, err := p.nc.Subscribe(p.resp, func(msg *nats.Msg) {
			r := new(pb.Action)
			if err := proto.Unmarshal(msg.Data, r); err != nil {
				// no nothing nows
				exit <- err
			}

			resp <- r
		})
		if err != nil {
			return err
		}

		defer sub.Unsubscribe()

		for {
			select {
			case err := <-exit:
				return err
			case <-p.ctx.Done():
				return nil
			case action := <-resp:
				if err := p.handleAction(action); err != nil {
					return err
				}
			}
		}
	}
}

func (p *Plugin) handleAction(action *pb.Action) error {
	// identify action ...
	switch a := action.GetAction().(type) {
	case *pb.Action_Config:
		return p.handleConfig(a.Config)
	case *pb.Action_Restart:
		return p.handleRestart(a.Restart)
	default:
	}

	return nil
}

func (p *Plugin) handleRestart(res *pb.Restart) error {
	p.cancel()

	return nil
}

func (p *Plugin) handleConfig(cfg *pb.Config) error {
	// reconfiguring plugin ...
	p.env = p.env.WithConfig(cfg)

	p.ready <- struct{}{}
	p.events <- RefreshConfigEvent

	return nil
}

func (p *Plugin) sendAction(a *pb.Action) error {
	// try to marshal into []byte ...
	msg, err := proto.Marshal(a)
	if err != nil {
		return err
	}

	err = p.nc.PublishMsg(&nats.Msg{
		Subject: p.env.ClusterDiscovery,
		Reply:   p.resp,
		Data:    msg,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) register() error {
	// action ...
	action := &pb.Action{
		Action: &pb.Action_Register{
			Register: &pb.Register{
				Plugin: p.pp,
			},
		},
	}

	// send the action
	if err := p.sendAction(action); err != nil {
		return err
	}

	// set this to be registered
	return nil
}

func (p *Plugin) pubInboxFunc(pub <-chan *pb.Event, opts ...FilterOpt) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := NewFilter(p, opts...)

		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				if e.Plugin == nil {
					e.Plugin = p.pp
				}

				// filtering an event
				e, err := f.Filter(e)
				if err != nil || e == nil {
					return err
				}

				// try to marshal into []byte ...
				msg, err := proto.Marshal(e)
				if err != nil {
					return err
				}

				if err := sc.Publish(p.env.Inbox, msg); err != nil {
					return err
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

func (p *Plugin) pubOutboxFunc(pub <-chan *pb.Event, opts ...FilterOpt) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := NewFilter(p, opts...)

		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				if e == nil {
					continue
				}

				if e.Plugin == nil {
					e.Plugin = p.pp
				}

				msg, err := proto.Marshal(e)
				if err != nil {
					return err
				}

				// filtering an event
				e, err = f.Filter(e)
				if err != nil {
					return err
				}

				// if this has not been filtered, or else
				if e == nil {
					continue
				}

				if err := sc.Publish(p.env.Outbox, msg); err != nil {
					return err
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

func (p *Plugin) subInboxFunc(sub chan<- *pb.Event, opts ...FilterOpt) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := NewFilter(p, opts...)

		s, err := sc.QueueSubscribe(p.env.Inbox, p.pp.GetName(), func(m *stan.Msg) {
			event := new(pb.Event)
			if err := proto.Unmarshal(m.Data, event); err != nil {
				// no nothing now
				return
			}

			// filtering an event
			event, err := f.Filter(event)
			if err != nil {
				return
			}

			// if this has not been filtered, or else
			if event != nil {
				sub <- event
			}
		}, stan.DurableName(p.pp.GetName()))
		if err != nil {
			return err
		}

		defer s.Close()

		<-p.ctx.Done()

		// close channel
		close(sub)

		return nil
	}
}

func (p *Plugin) subOutboxFunc(sub chan<- *pb.Event, opts ...FilterOpt) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := NewFilter(p, opts...)

		s, err := sc.QueueSubscribe(p.env.Outbox, p.pp.GetName(), func(m *stan.Msg) {
			event := new(pb.Event)
			if err := proto.Unmarshal(m.Data, event); err != nil {
				// no nothing now
				return
			}

			// filtering an event
			event, err := f.Filter(event)
			if err != nil {
				return
			}

			// if this has not been filtered, or else
			if event != nil {
				sub <- event
			}
		}, stan.DurableName(p.pp.GetName()))
		if err != nil {
			return err
		}

		defer s.Close()

		<-p.ctx.Done()

		// close channel
		close(sub)

		return nil
	}
}

func (p *Plugin) lostHandler() func(_ stan.Conn, reason error) {
	return func(_ stan.Conn, reason error) {
		p.cancel()

		p.err = reason
	}
}

func configure(p *Plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	return nil
}
