package plugins

import (
	"os"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/stan.go"
)

const (
	AutobotName          = "AUTOBOT_NAME"
	AutobotClusterID     = "AUTOBOT_CLUSTER_ID"
	AutobotClusterURL    = "AUTOBOT_CLUSTER_URL"
	AutobotChannelInbox  = "AUTOBOT_CHANNEL_INBOX"
	AutobotChannelOutbox = "AUTOBOT_CHANNEL_OUTBOX"
)

// Plugin ...
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

	conn stan.Conn
	meta *pb.Plugin

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Plugins ...
func New(meta *pb.Plugin, opts ...Opt) (Plugin, error) {
	options := new(Opts)

	p := new(plugin)
	// setting a default env for a plugin
	p.opts = options
	p.meta = meta

	// configure plugin
	configure(p, opts...)

	// configure client
	if err := configureClient(p); err != nil {
		return nil, err
	}

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
	<-p.exit

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
				p.exit <- struct{}{}
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

				if err := p.conn.Publish(os.Getenv(AutobotChannelInbox), msg); err != nil {
					return err
				}
			case <-p.exit:
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

				if err := p.conn.Publish(os.Getenv(AutobotChannelOutbox), msg); err != nil {
					return err
				}

			case <-p.SubscribeOutbox():
				return nil
			}
		}
	}
}

func (p *plugin) subMessagesFunc(sub chan<- *pb.Event) func() error {
	return func() error {
		s, err := p.conn.QueueSubscribe(os.Getenv(AutobotChannelInbox), p.meta.GetName(), func(m *stan.Msg) {
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

		<-p.exit

		// close channel
		close(sub)

		return nil
	}
}

func (p *plugin) subRepliesFunc(sub chan<- *pb.Event) func() error {
	return func() error {
		s, err := p.conn.QueueSubscribe(os.Getenv(AutobotChannelOutbox), p.meta.GetName(), func(m *stan.Msg) {
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

		<-p.exit

		// close channel
		close(sub)

		return nil
	}
}

func configureClient(p *plugin) error {
	sc, err := stan.Connect(os.Getenv(AutobotClusterID), p.meta.GetName())
	if err != nil {
		return err
	}

	p.conn = sc

	return nil
}

func configure(p *plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	return nil
}
