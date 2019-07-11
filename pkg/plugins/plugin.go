package plugins

import (
	"context"
	"os"
	"strconv"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// These are the environment variables that are provided
// by the server to started plugins.
const (
	// AutobotName is the name of the autobot.
	AutobotName = "AUTOBOT_NAME"
	// AutobotClusterID is the id of the started NATS Streaming Server.
	AutobotClusterID = "AUTOBOT_CLUSTER_ID"
	// AutobotClusterURL is the URL of the started NATS Streaming Server.
	AutobotClusterURL = "AUTOBOT_CLUSTER_URL"
	// AutobotChannelInbox is the name of the inbox that the plugin should subscribe to.
	AutobotChannelInbox = "AUTOBOT_CHANNEL_INBOX"
	// AutobotChannelOutbox is the name of the outbox that the plugin should publish to.
	AutobotChannelOutbox = "AUTOBOT_CHANNEL_OUTBOX"
	// AutobotVerbose should enable the verbose behavior in the plugin.
	AutobotVerbose = "AUTOBOT_VERBOSE"
	// AutobotDebug should enable the debug behavior in the plugin.
	AutobotDebug = "AUTOBOT_DEBUG"
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

	sc   stan.Conn
	nc   *nats.Conn
	meta *pb.Plugin

	ctx     context.Context
	cancel  func()
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct{}

// SubscribeFunc ...
type SubscribeFunc = func(*pb.Event) (*pb.Event, error)

// WithContext is creating a new plugin and a context to run operations in routines.
// When the context is canceled, all concurrent operations are canceled.
func WithContext(ctx context.Context, meta *pb.Plugin, opts ...Opt) (*Plugin, context.Context, error) {
	p, err := newPlugin(meta, opts...)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.ctx = ctx

	return p, ctx, nil
}

// newPlugin ...
func newPlugin(meta *pb.Plugin, opts ...Opt) (*Plugin, error) {
	options := new(Opts)

	p := new(Plugin)
	// setting a default env for a plugin
	p.opts = options
	p.meta = meta

	// configure plugin
	configure(p, opts...)

	// connect client ...
	if err := configureClient(p); err != nil {
		return nil, err
	}

	return p, nil
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

	p.sc.Close()
	p.nc.Close()

	return p.err
}

// ReplyWithFunc ...
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

// Run ...
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
func (p *Plugin) Debug() (bool, error) {
	b, err := strconv.ParseBool(os.Getenv(AutobotDebug))
	if err != nil {
		return false, err
	}

	return b, nil
}

// Verbose ...
func (p *Plugin) Verbose() (bool, error) {
	b, err := strconv.ParseBool(os.Getenv(AutobotVerbose))
	if err != nil {
		return false, err
	}

	return b, nil
}

// Inbox ...
func (p *Plugin) Inbox() string {
	return os.Getenv(AutobotChannelInbox)
}

// Outbox ...
func (p *Plugin) Outbox() string {
	return os.Getenv(AutobotChannelOutbox)
}

// ClusterID ...
func (p *Plugin) ClusterID() string {
	return os.Getenv(AutobotClusterID)
}

// ClusterURL ...
func (p *Plugin) ClusterURL() string {
	return os.Getenv(AutobotClusterURL)
}

func (p *Plugin) pubInboxFunc(pub <-chan *pb.Event, opts ...FilterOpt) func() error {
	return func() error {
		f := NewFilter(p, opts...)

		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				if e.Plugin == nil {
					e.Plugin = p.meta
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

				if err := p.sc.Publish(p.Inbox(), msg); err != nil {
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
					e.Plugin = p.meta
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

				if err := p.sc.Publish(p.Outbox(), msg); err != nil {
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
		f := NewFilter(p, opts...)

		s, err := p.sc.QueueSubscribe(p.Inbox(), p.meta.GetName(), func(m *stan.Msg) {
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
		}, stan.DurableName(p.meta.GetName()))
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
		f := NewFilter(p, opts...)

		s, err := p.sc.QueueSubscribe(p.Outbox(), p.meta.GetName(), func(m *stan.Msg) {
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
		}, stan.DurableName(p.meta.GetName()))
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

func configureClient(p *Plugin) error {
	nc, err := nats.Connect(
		p.ClusterURL(),
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	)
	if err != nil {
		return err
	}

	p.nc = nc

	sc, err := stan.Connect(p.ClusterID(), p.meta.GetName(), stan.SetConnectionLostHandler(p.lostHandler()))
	if err != nil {
		return err
	}

	p.sc = sc

	return nil
}

func configure(p *Plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	return nil
}
