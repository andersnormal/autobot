package plugins

import (
	"context"
	"os"
	"sync"
	"time"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// These are the environment variables that are provided
// by the server to started plugins.
const (
	// AutobotClusterID ...
	AutobotClusterID = "AUTOBOT_CLUSTER_ID"
	// AutobotClusterURL ...
	AutobotClusterURL = "AUTOBOT_CLUSTER_URL"
	// AutobotChannelDiscovery is the name of the topic to register a plugin.
	AutobotChannelDiscovery = "AUTOBOT_CHANNEL_DISCOVERY"
)

// Event symbolizes events that can occur in the plugin.
// Though they are generally triggered from the sever.
type Event int32

const (
	ReadyEvent         Event = 0
	RefreshConfigEvent Event = 1
)

// Plugin describes a plugin
//
//  // create a root context ...
//  ctx, cancel := context.WithCancel(context.Background())
//  defer cancel()
//
//  // plugin ....
//  plugin, ctx, err := plugins.WithContext(ctx, pb.NewPlugin(name))
//  if err != nil {
//    log.Fatalf("could not create plugin: %v", err)
//  }
//
//  if err := plugin.Wait(); err != nil {
//    log.Fatalf("stopped plugin: %v", err)
//  }
type Plugin struct {
	opts *Opts

	sc stan.Conn
	nc *nats.Conn

	ready          chan struct{}
	events         chan Event
	eventsChannels []chan Event
	resp           string

	cfg *pb.Config

	ctx     context.Context
	cancel  func()
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup

	pp *pb.Plugin

	sync.RWMutex
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Name                    string
	ClusterID               string
	ClusterURL              string
	ClusterDiscoveryChannel string
}

// SubscribeFunc ...
type SubscribeFunc = func(*pb.Event) (*pb.Event, error)

// WithContext is creating a new plugin and a context to run operations in routines.
// When the context is canceled, all concurrent operations are canceled.
func WithContext(ctx context.Context, opts ...Opt) (*Plugin, context.Context) {
	p := newPlugin(opts...)

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.ctx = ctx

	return p, ctx
}

// Name is providing a name for the plugin,
// which is similar to a service name that is uniquely identifiying this plugin
func Name(name string) func(o *Opts) {
	return func(o *Opts) {
		o.Name = name
	}
}

// ClusterID is providing an id for the cluster to use.
// If the plugin is not currated by the controller or provided as environment variable
// this can be used to set the NATS streaming cluster id.
func ClusterID(clusterID string) func(o *Opts) {
	return func(o *Opts) {
		o.ClusterID = clusterID
	}
}

// ClusterURL is providing the URL for the cluster to use.
// If the plugin is not currated by the controller or provided as environment variable
// this can be used to set the NATS streaming cluster URL.
func ClusterURL(clusterURL string) func(o *Opts) {
	return func(o *Opts) {
		o.ClusterID = clusterURL
	}
}

// ClusterDiscoveryChannel is providing the name of the discovery channel in the cluster.
// If the plugin is not currated by the controller or provided as environment variable
// this can be used to set the NATS streaming cluster discovery channel.
func ClusterDiscoveryChannel(clusterDiscoveryChannel string) func(o *Opts) {
	return func(o *Opts) {
		o.ClusterDiscoveryChannel = clusterDiscoveryChannel
	}
}

// Events returns a channel which is used to publish important
// events in the plugin. For example that the server pushes a new config.
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

// SubscribeInbox ...
func (p *Plugin) SubscribeInbox(opts ...FilterOpt) <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(p.subInboxFunc(sub, opts...))

	return sub
}

// SubscribeOutbox ...
func (p *Plugin) SubscribeOutbox(opts ...FilterOpt) <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(p.subOutboxFunc(sub, opts...))

	return sub
}

// PublishInbox ...
func (p *Plugin) PublishInbox(opts ...FilterOpt) chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(p.pubInboxFunc(pub, append(DefaultInboxFilterOpts, opts...)...))

	return pub
}

// PublishOutbox ...
func (p *Plugin) PublishOutbox(opts ...FilterOpt) chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(p.pubOutboxFunc(pub, append(DefaultOutboxFilterOpts, opts...)...))

	return pub
}

// Wait ...
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
//
//  ctx, cancel := context.WithCancel(context.Background())
//  defer cancel()
//
//  // plugin ....
//  plugin, ctx := plugins.WithContext(ctx, pb.NewPlugin(name))
//
//  // use the schedule function from the plugin
//  if err := plugin.ReplyWithFunc(msgFunc()); err != nil {
//    log.Fatalf("could not create plugin: %v", err)
//  }
//
//  if err := plugin.Wait(); err != nil {
//    log.Fatalf("stopped plugin: %v", err)
//  }
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

// Debug ...
func (p *Plugin) Debug() bool {
	return p.cfg.GetDebug()
}

// Verbose ...
func (p *Plugin) Verbose() bool {
	return p.cfg.GetVerbose()
}

// Inbox ...
func (p *Plugin) Inbox() string {
	return p.cfg.GetInbox()
}

// Outbox ...
func (p *Plugin) Outbox() string {
	return p.cfg.GetOutbox()
}

// Discovery ...
func (p *Plugin) Discovery() string {
	return p.cfg.GetDiscovery()
}

// ClusterID ...
func (p *Plugin) ClusterID() string {
	return p.cfg.GetClusterId()
}

// ClusterURL ...
func (p *Plugin) ClusterURL() string {
	return p.cfg.GetClusterId()
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
		p.cfg.GetClusterUrl(),
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	)
	if err != nil {
		return nil, err
	}

	p.nc = nc

	sc, err := stan.Connect(p.cfg.GetClusterId(), p.pp.GetName(), stan.SetConnectionLostHandler(p.lostHandler()))
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
	p.cfg = cfg

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
		Subject: p.cfg.GetDiscovery(),
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

				if err := sc.Publish(p.Inbox(), msg); err != nil {
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

				if err := sc.Publish(p.Outbox(), msg); err != nil {
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

		s, err := sc.QueueSubscribe(p.Inbox(), p.pp.GetName(), func(m *stan.Msg) {
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

		s, err := sc.QueueSubscribe(p.Outbox(), p.pp.GetName(), func(m *stan.Msg) {
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

	// create plugin meta
	plugin := new(pb.Plugin)
	plugin.Name = p.opts.Name

	p.pp = plugin

	// create config
	cfg := new(pb.Config)
	cfg.ClusterId = p.opts.ClusterID
	cfg.ClusterUrl = p.opts.ClusterURL
	cfg.Discovery = p.opts.ClusterDiscoveryChannel

	if cfg.GetClusterId() == "" {
		cfg.ClusterId = os.Getenv(AutobotClusterID)
	}

	if cfg.GetClusterUrl() == "" {
		cfg.ClusterUrl = os.Getenv(AutobotClusterURL)
	}

	if cfg.GetDiscovery() == "" {
		cfg.Discovery = os.Getenv(AutobotChannelDiscovery)
	}

	p.cfg = cfg

	return nil
}
