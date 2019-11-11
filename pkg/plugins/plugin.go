package plugins

import (
	"context"
	"sync"

	"github.com/andersnormal/autobot/pkg/plugins/filters"
	"github.com/andersnormal/autobot/pkg/plugins/message"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

// Event symbolizes events that can occur in the plugin.
// Though they maybe triggered from somewhere else.
type Event interface{}

const (
	// ErrUnknown ...
	ErrUnknown = iota
	// ErrParse ...
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

	marshaler message.Marshaler
	runtime   *runtime.Environment

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
type Opts struct{}

// SubscribeFunc ...
type SubscribeFunc = func(Context) error

// WithContext is creating a new plugin and a context to run operations in routines.
// When the context is canceled, all concurrent operations are canceled.
func WithContext(ctx context.Context, env *runtime.Environment, opts ...Opt) (*Plugin, context.Context) {
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
func newPlugin(env *runtime.Environment, opts ...Opt) *Plugin {
	options := new(Opts)

	p := new(Plugin)
	// setting a default env for a plugin
	p.opts = options
	p.runtime = env

	// this is the basic marshaler
	p.marshaler = message.JSONMarshaler{NewUUID: message.NewUUID}

	// configure plugin
	configure(p, opts...)

	// logging ...
	configureLogging(p)

	return p
}

// SubscribeInbox is subscribing to the inbox of messages.
// This if for plugins that want to consume the message that other
// plugins publish to Autobot.
func (p *Plugin) SubscribeInbox(funcs ...filters.FilterFunc) <-chan Event {
	sub := make(chan Event)

	p.run(p.subInboxFunc(sub, funcs...))

	return sub
}

// SubscribeOutbox is subscribing to the outbox of messages.
// These are the message that ought to be published to an external service (e.g. Slack, MS Teams).
func (p *Plugin) SubscribeOutbox(funcs ...filters.FilterFunc) <-chan Event {
	sub := make(chan Event)

	p.run(p.subOutboxFunc(sub, funcs...))

	return sub
}

// PublishInbox is publishing message to the inbox in the controller.
func (p *Plugin) PublishInbox(msg *pb.Message) error { // male this an actual interface
	return p.publish(p.runtime.Inbox, msg)
}

// PublishOutbox is publishing message to the outbox in the controller.
func (p *Plugin) PublishOutbox(msg *pb.Message) error {
	return p.publish(p.runtime.Outbox, msg)
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

// AsyncReplyWithFunc is a wrapper function to provide a message handler function
// which will asynchronously reply to the messages received by a plugin.
// The difference between this and the synchronous ReplyWithFunc wrapper is that
// while both receive messages in sequence, ReplyWithFunc will block on reading
// the subsequent message from the inbox while the current message is being handled.
func (p *Plugin) AsyncReplyWithFunc(fn SubscribeFunc, funcs ...filters.FilterFunc) error {
	p.run(func() error {

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
					// creating a new context for the message
					// this could be moved to a sync.Pool
					ctx := &cbContext{
						plugin: p,
						msg:    ev,
					}

					// launch goroutine to handle async reply handling
					p.asyncHandleMessage(ctx, fn)
				default:
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	})

	return nil
}

// asyncHandleMessage is a helper used to
func (p *Plugin) asyncHandleMessage(c Context, fn SubscribeFunc) {
	p.run(func() error {
		err := fn(c)
		// TODO: in 1.1 we will enqueue to a dead letter topic and proceed
		// perhaps, unless a specific error is returned.
		return err
	})
}

// ReplyWithFunc is a wrapper function to provide a function which may
// send replies to received message for this plugin.
func (p *Plugin) ReplyWithFunc(fn SubscribeFunc, funcs ...filters.FilterFunc) error {
	p.run(func() error {
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
					// creating a new context for the message
					// this could be moved to a sync.Pool
					ctx := &cbContext{
						plugin: p,
						msg:    ev,
					}

					err := fn(ctx)
					if err != nil {
						return err
					}
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

	return sc, nil
}

func (p *Plugin) watchcat() func() error {
	return func() error {
		<-p.ctx.Done()

		if p.sc != nil {
			// defer to make sure subscribers have a chance
			// to unsubcribe cleanly.
			defer p.sc.Close()
		}

		if p.nc != nil {
			defer p.nc.Close()
		}

		return nil
	}
}

func (p *Plugin) publish(topic string, msg *pb.Message) error {
	sc, err := p.getConn()
	if err != nil {
		return err
	}

	b, err := p.marshaler.Marshal(msg)
	if err != nil {
		return err
	}

	if err := sc.Publish(topic, b); err != nil {
		return err
	}
	return nil
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
			botMessage := new(pb.Message)
			if err := p.marshaler.Unmarshal(m.Data, botMessage); err != nil {
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
				sub <- event
			}
		}, stan.DurableName(p.runtime.Name), stan.StartWithLastReceived())
		if err != nil {
			return err
		}

		<-p.ctx.Done()
		s.Unsubscribe()
		s.Close()
		// close channel
		close(sub)
		sub = nil

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
			botMessage := new(pb.Message)
			if err := p.marshaler.Unmarshal(m.Data, botMessage); err != nil {
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
				sub <- event
			}
		}, stan.DurableName(p.runtime.Name), stan.StartWithLastReceived())
		if err != nil {
			return err
		}

		<-p.ctx.Done()
		s.Unsubscribe()
		s.Close()

		// close channel
		close(sub)
		sub = nil

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
