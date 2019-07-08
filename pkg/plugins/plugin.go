package plugins

import (
	"context"
	"os"
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
)

// Plugin is the interface to a new plugin.
//
//  create a root context ...
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
type Plugin interface {
	// SubscribeInbox ...
	SubscribeInbox() <-chan *pb.Event
	// PublishInbox ...
	PublishInbox() chan<- *pb.Event
	// SubscribeOutbox ...
	SubscribeOutbox() <-chan *pb.Event
	// PublishOutbox ...
	PublishOutbox() chan<- *pb.Event
	// ReplyWithFunc ...
	ReplyWithFunc(func(*pb.Event) (*pb.Event, error), ...ReplyWithFuncOpt) error
	// Run ...
	Run(func() error)
	// Wait ...
	Wait() error
}

// ReplyWithFuncOpt ...
type ReplyWithFuncOpt func(*ReplyWithFuncOpts)

// ReplyWithFuncOpts ...
type ReplyWithFuncOpts struct {
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct{}

// SubscribeFunc ...
type SubscribeFunc = func(*pb.Event) (*pb.Event, error)

type plugin struct {
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

// WithContext ...
func WithContext(ctx context.Context, meta *pb.Plugin, opts ...Opt) (Plugin, context.Context, error) {
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
func newPlugin(meta *pb.Plugin, opts ...Opt) (*plugin, error) {
	options := new(Opts)

	p := new(plugin)
	// setting a default env for a plugin
	p.opts = options
	p.meta = meta

	// configure plugin
	configure(p, opts...)

	// connect client ...
	if err := configureClient(p); err != nil {
		return nil, err
	}

	// watch non context func ...
	go func() {

	}()

	return p, nil
}

// SubscribeInbox ...
func (p *plugin) SubscribeInbox() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(p.subMessagesFunc(sub))

	return sub
}

// SubscribeOutbox ...
func (p *plugin) SubscribeOutbox() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(p.subRepliesFunc(sub))

	return sub
}

// PublishInbox ...
func (p *plugin) PublishInbox() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(p.pubMessagesFunc(pub))

	return pub
}

// PublishOutbox ...
func (p *plugin) PublishOutbox() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(p.pubRepliesFunc(pub))

	return pub
}

// Wait ...
func (p *plugin) Wait() error {
	p.wg.Wait()

	if p.cancel != nil {
		p.cancel()
	}

	p.sc.Close()
	p.nc.Close()

	return p.err
}

// ReplyWithFunc ...
func (p *plugin) ReplyWithFunc(fn SubscribeFunc, opts ...ReplyWithFuncOpt) error {
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
func (p *plugin) Run(f func() error) {
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

func (p *plugin) pubMessagesFunc(pub <-chan *pb.Event) func() error {
	return func() error {
		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				e.Plugin = p.meta

				// try to marshal into []byte
				msg, err := proto.Marshal(e)
				if err != nil {
					return err
				}

				if err := p.sc.Publish(os.Getenv(AutobotChannelInbox), msg); err != nil {
					return err
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

func (p *plugin) pubRepliesFunc(pub <-chan *pb.Event) func() error {
	return func() error {
		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				e.Plugin = p.meta

				msg, err := proto.Marshal(e)
				if err != nil {
					return err
				}

				if err := p.sc.Publish(os.Getenv(AutobotChannelOutbox), msg); err != nil {
					return err
				}
			case <-p.ctx.Done():
				return nil
			}
		}
	}
}

func (p *plugin) subMessagesFunc(sub chan<- *pb.Event) func() error {
	return func() error {
		s, err := p.sc.QueueSubscribe(os.Getenv(AutobotChannelInbox), p.meta.GetName(), func(m *stan.Msg) {
			event := new(pb.Event)
			if err := proto.Unmarshal(m.Data, event); err != nil {
				// no nothing now
				return
			}

			sub <- event
		})
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

func (p *plugin) subRepliesFunc(sub chan<- *pb.Event) func() error {
	return func() error {
		s, err := p.sc.QueueSubscribe(os.Getenv(AutobotChannelOutbox), p.meta.GetName(), func(m *stan.Msg) {
			event := new(pb.Event)
			if err := proto.Unmarshal(m.Data, event); err != nil {
				// no nothing now
				return
			}

			sub <- event
		})
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

func configureClient(p *plugin) error {
	nc, err := nats.Connect(
		os.Getenv(AutobotClusterURL),
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	)
	if err != nil {
		return err
	}

	p.nc = nc

	sc, err := stan.Connect(os.Getenv(AutobotClusterID), p.meta.GetName())
	if err != nil {
		return err
	}

	p.sc = sc

	return nil
}

func configure(p *plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	return nil
}
