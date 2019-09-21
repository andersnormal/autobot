package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/andersnormal/autobot/pkg/plugins/filters"
	"github.com/andersnormal/autobot/pkg/plugins/message"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

// Event symbolizes events that can occur in the plugin.
// Though they maybe triggered from somewhere else.
type Event interface{}

const (
	ErrUnknown = iota
	ErrParse
)

// MessageError represents an error that may occurs
type MessageError struct {
	Code int
	Msg  string
}

// returns
func (e MessageError) Error() string {
	return e.Msg
}

// Plugin describes a plugin in general.
// It should not be instantiated directly.
type Plugin struct {
	opts *Opts

	sc stan.Conn
	nc *nats.Conn

	resp      string
	marshaler message.Marshaler

	runtime runtime.Runtime

	ctx     context.Context
	cancel  func()
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup

	logger *log.Entry

	sync.RWMutex
}

// Opt is an Option
type Opt func(*Opts)

// Opts are the available options
type Opts struct {
	SubscriptionOpts []stan.SubscriptionOption
}

// SubscribeFunc ...
type SubscribeFunc = func(*pb.Message) (*pb.Message, error)

// WithContext is creating a new plugin and a context to run operations in routines.
// When the context is canceled, all concurrent operations are canceled.
func WithContext(ctx context.Context, env runtime.Runtime, opts ...Opt) (*Plugin, context.Context) {
	p := newPlugin(env, opts...)

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.ctx = ctx

	return p, ctx
}

// Log is returning the logger to log to the log formatter.
func (p *Plugin) Log() *log.Entry {
	return p.logger
}

// newPlugin ...
func newPlugin(env runtime.Runtime, opts ...Opt) *Plugin {
	options := new(Opts)

	p := new(Plugin)
	// setting a default env for a plugin
	p.opts = options
	p.runtime = env

	// this is the basic marshaler
	p.marshaler = message.ProtoMarshaler{NewUUID: message.NewUUID}

	// configure plugin
	configure(p, opts...)

	// logging ...
	configureLogging(p)

	return p
}

// WithSubscriptionOpts ...
func WithSubscriptionOpts(subOpts ...stan.SubscriptionOption) func(*Opts) {
	return func(o *Opts) {
		o.SubscriptionOpts = subOpts
	}
}

// SubscribeInbox is subscribing to the inbox of messages.
// This if for plugins that want to consume the message that other
// plugins publish to Autobot.
func (p *Plugin) SubscribeInbox(funcs ...filters.FilterFunc) <-chan Event {
	sub := make(chan Event)

	p.run(p.subFunc(p.runtime.Inbox, sub, funcs...))

	return sub
}

// SubscribeOutbox is subscribing to the outbox of messages.
// These are the message that ought to be published to an external service (e.g. Slack, MS Teams).
func (p *Plugin) SubscribeOutbox(funcs ...filters.FilterFunc) <-chan Event {
	sub := make(chan Event)

	p.run(p.subFunc(p.runtime.Outbox, sub, funcs...))

	return sub
}

// PublishInbox is publishing message to the inbox in the controller.
// The returned channel pushes all of the send message to the inbox in the controller.
func (p *Plugin) PublishInbox(funcs ...filters.FilterFunc) chan<- *pb.Message { // male this an actual interface
	pub := make(chan *pb.Message)

	p.run(p.pubInboxFunc(pub, append(filters.DefaultInboxFilterOpts, funcs...)...))

	return pub
}

// PublishOutbox is publishing message to the outbox in the controller.
// The returned channel pushes all the send messages to the outbox in the controller.
func (p *Plugin) PublishOutbox(funcs ...filters.FilterFunc) chan<- *pb.Message {
	pub := make(chan *pb.Message)

	p.run(p.pubOutboxFunc(pub, append(filters.DefaultOutboxFilterOpts, funcs...)...))

	return pub
}

// Wait is waiting for the underlying WaitGroup. All run go routines are
// hold here to before one exists with an error. Then the provided context is canceled.
func (p *Plugin) Wait() error {
	p.wg.Wait()

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
func (p *Plugin) ReplyWithFunc(fn SubscribeFunc, funcs ...filters.FilterFunc) error {
	p.run(func() error {
		// create publish channel ...
		pubReply := p.PublishOutbox()
		subMsg := p.SubscribeInbox()

		for {
			select {
			case e, ok := <-subMsg:
				if !ok {
					return nil
				}

				switch ev := e.(type) {
				case *MessageError:
					return ev
				case *pb.Message:
					r, err := fn(ev)
					if err != nil {
						return err
					}

					pubReply <- r
				default:
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	})

	return nil
}

func (p *Plugin) run(f func() error) {
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
	p.run(p.watchcat())

	return p.sc, nil
}

func (p *Plugin) newConn() (stan.Conn, error) {
	nc, err := nats.Connect(
		p.runtime.ClusterURL,
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	)
	if err != nil {
		return nil, err
	}

	p.nc = nc

	sc, err := stan.Connect(
		p.runtime.ClusterID,
		p.runtime.Name,
		stan.NatsConn(nc),
		stan.SetConnectionLostHandler(p.lostHandler()),
	)
	if err != nil {
		return nil, err
	}

	// creates a new re-used response inbox for the plugin
	p.resp = p.nc.NewRespInbox()

	return sc, nil
}

func (p *Plugin) watchcat() func() error {
	return func() error {
		exit := make(chan error)
		resp := make(chan *pb.Bot)

		sub, err := p.nc.Subscribe(p.resp, func(msg *nats.Msg) {
			r := new(pb.Bot)
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
				if p.sc != nil {
					p.sc.Close()
				}

				if p.nc != nil {
					p.nc.Close()
				}

				return nil
			case <-resp:
				// nothing to do yet
			}
		}
	}
}

func (p *Plugin) publishEvent(a *pb.Bot) error {
	// try to marshal into []byte ...
	msg, err := proto.Marshal(a)
	if err != nil {
		return err
	}

	err = p.nc.PublishMsg(&nats.Msg{
		Subject: p.runtime.ClusterDiscovery,
		Reply:   p.resp,
		Data:    msg,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) pubInboxFunc(pub <-chan *pb.Message, funcs ...filters.FilterFunc) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := filters.New(funcs...)

		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				// filtering an event
				e, err := f.Filter(e)
				if err != nil || e == nil {
					return err
				}

				msg, err := p.marshaler.Marshal(e)
				if err != nil {
					return err
				}

				// add some metadata
				msg.Metadata.Src(p.runtime.Name)

				b, err := json.Marshal(msg)
				if err != nil {
					return err
				}

				if err := sc.Publish(p.runtime.Inbox, b); err != nil {
					return err
				}

				fmt.Println("finished publishing...")
			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

func (p *Plugin) pubOutboxFunc(pub <-chan *pb.Message, funcs ...filters.FilterFunc) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := filters.New(funcs...)

		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				if e == nil {
					continue
				}

				msg, err := p.marshaler.Marshal(e)
				if err != nil {
					return err
				}

				// add some metadata
				msg.Metadata.Src(p.runtime.Name)

				b, err := json.Marshal(msg)
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

				fmt.Println("publishing to", p.runtime.Outbox)

				fmt.Println("publishing and will wait for an ack...")

				if err := sc.Publish(p.runtime.Outbox, b); err != nil {
					return err
				}

				fmt.Println("received an ack!")

			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

func (p *Plugin) subInboxFunc(sub chan<- Event, funcs ...filters.FilterFunc) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := filters.New(funcs...)

		// we are using a queue subscription to only deliver the work to one of the plugins,
		// because they subscribe to a group by the plugin name.
		s, err := sc.QueueSubscribe(p.runtime.Inbox, p.runtime.Name, func(m *stan.Msg) {

			// this is recreating the messsage from the inbox
			msg, err := message.FromByte(m.Data)
			if err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			botMessage := new(pb.Message)
			if err := p.marshaler.Unmarshal(msg, botMessage); err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			// filtering an event
			event, err := f.Filter(botMessage)
			if err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			// if this has not been filtered, or else
			if event != nil {
				sub <- botMessage
			}

			m.Ack()
			fmt.Println("just acked message")
		},
			p.opts.SubscriptionOpts...,
		)
		if err != nil {
			return err
		}

		<-p.ctx.Done()

		// close channel
		close(sub)

		defer sc.Close()
		defer s.Close()

		return nil
	}
}

func (p *Plugin) subFunc(subject string, sub chan<- Event, funcs ...filters.FilterFunc) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := filters.New(funcs...)

		// we are using a queue subscription to only deliver the work to one of the plugins,
		// because they subscribe to a group by the plugin name.
		s, err := sc.QueueSubscribe(subject, p.runtime.Name, func(m *stan.Msg) {
			// this is recreating the messsage from the inbox
			msg, err := message.FromByte(m.Data)
			if err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			botMessage := new(pb.Message)
			if err := p.marshaler.Unmarshal(msg, botMessage); err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			// filtering an event
			event, err := f.Filter(botMessage)
			if err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			// if this has not been filtered, or else
			if event != nil {
				sub <- botMessage
			}

			m.Ack()
			fmt.Println("just acked message")
		},
			p.opts.SubscriptionOpts...,
		)
		if err != nil {
			return err
		}

		<-p.ctx.Done()

		// close channel
		close(sub)

		defer sc.Close()
		defer s.Close()

		return nil
	}
}

func (p *Plugin) subOutboxFunc(sub chan<- Event, funcs ...filters.FilterFunc) func() error {
	return func() error {
		sc, err := p.getConn()
		if err != nil {
			return err
		}

		f := filters.New(funcs...)

		// we are using a queue subscription to only deliver the work to one of the plugins,
		// because they subscribe to a group by the plugin name.
		s, err := sc.QueueSubscribe(p.runtime.Outbox, p.runtime.Name, func(m *stan.Msg) {
			// this is recreating the messsage from the inbox
			msg, err := message.FromByte(m.Data)
			if err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			botMessage := new(pb.Message)
			if err := p.marshaler.Unmarshal(msg, botMessage); err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			// filtering an event
			event, err := f.Filter(botMessage)
			if err != nil {
				sub <- &MessageError{Code: ErrParse, Msg: err.Error()}

				return
			}

			// if this has not been filtered, or else
			if event != nil {
				sub <- botMessage
			}

			m.Ack()
			fmt.Println("just acked message - outbox")
		},
			p.opts.SubscriptionOpts...,
		)
		if err != nil {
			return err
		}

		<-p.ctx.Done()

		// close channel
		close(sub)

		defer sc.Close()
		defer s.Close()

		return nil
	}
}

func (p *Plugin) lostHandler() func(_ stan.Conn, reason error) {
	return func(_ stan.Conn, reason error) {
		p.cancel()

		p.err = reason
	}
}

func configureLogging(p *Plugin) error {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)

	if p.runtime.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	p.logger = log.WithFields(
		log.Fields{
			"autobot_name": p.runtime.Name,
			"cluster_url":  p.runtime.ClusterURL,
			"cluster_id":   p.runtime.ClusterID,
		},
	)

	if level, err := log.ParseLevel(p.runtime.LogLevel); err == nil {
		log.SetLevel(level)
	}

	return nil
}

func configure(p *Plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	return nil
}
