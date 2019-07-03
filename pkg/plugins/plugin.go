package plugins

import (
	"log"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/stan.go"
)

const (
	AutobotName          = "AUTOBOT_NAME"
	AutobotClusterID     = "AUTOBOT_CLUSTER_ID"
	AutobotClusterURL    = "AUTOBOT_CLUSTER_URL"
	AutobotTopicMessages = "AUTOBOT_TOPIC_MESSAGES"
	AutobotTopicReplies  = "AUTOBOT_TOPIC_REPLIES"
)

// Plugin ...
type Plugin interface {
	// SubscribeMessages ...
	SubscribeMessages() <-chan *pb.Event
	// PublishMessages ...
	PublishMessages() chan<- *pb.Event
	// SubscribeReplies ...
	SubscribeReplies() <-chan *pb.Event
	// PublishReplies ...
	PublishReplies() chan<- *pb.Event
	// ReplyMessageWithFunc ...
	ReplyMessageWithFunc(func(*pb.Event) (*pb.Event, error), ...ReplyMessageWithFuncOpt) error
	// Run ...
	Run(func() error)
	// Wait ...
	Wait() error
}

// ReplyMessageWithFuncOpt ...
type ReplyMessageWithFuncOpt func(*ReplyMessageWithFuncOpts)

// Opts ...
type ReplyMessageWithFuncOpts struct {
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Env Env
}

// SubscribeFunc ...
type SubscribeFunc = func(*pb.Event) (*pb.Event, error)

type plugin struct {
	name string

	opts *Opts
	conn stan.Conn
	env  Env

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Plugins ...
func New(name string, opts ...Opt) (Plugin, error) {
	options := new(Opts)

	p := new(plugin)
	// setting a default env for a plugin
	p.opts = options
	p.name = name
	p.env = NewDefaultEnv()

	// configure plugin
	configure(p, opts...)

	// configure client
	if err := configureClient(p); err != nil {
		return nil, err
	}

	return p, nil
}

// SubscribeMessages ...
func (p *plugin) SubscribeMessages() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(subMessagesFunc(p.conn, p.env, sub, p.exit))

	return sub
}

// SubscribeReplies ...
func (p *plugin) SubscribeReplies() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.Run(subRepliesFunc(p.conn, p.env, sub, p.exit))

	return sub
}

// PublishMessages ...
func (p *plugin) PublishMessages() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(pubMessagesFunc(p.conn, p.env, pub, p.exit))

	return pub
}

// PublishReplies ...
func (p *plugin) PublishReplies() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.Run(pubRepliesFunc(p.conn, p.env, pub, p.exit))

	return pub
}

// Wait ...
func (p *plugin) Wait() error {
	<-p.exit

	return p.err
}

// ReplyMessageWithFunc ...
func (p *plugin) ReplyMessageWithFunc(fn SubscribeFunc, opts ...ReplyMessageWithFuncOpt) error {
	p.Run(func() error {
		// create publish channel ...
		pubReply := p.PublishReplies()
		subMsg := p.SubscribeMessages()

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

func pubMessagesFunc(conn stan.Conn, env Env, pub <-chan *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				msg, err := proto.Marshal(e)
				if err != nil {
					return err
				}

				if err := conn.Publish(env.AutobotTopicMessages(), msg); err != nil {
					return err
				}
			case <-exit:
				return nil
			}
		}
	}
}

func pubRepliesFunc(conn stan.Conn, env Env, pub <-chan *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		for {
			select {
			case e, ok := <-pub:
				if !ok {
					return nil
				}

				msg, err := proto.Marshal(e)
				if err != nil {
					return err
				}

				if err := conn.Publish(env.AutobotTopicReplies(), msg); err != nil {
					return err
				}

				log.Printf("published: %v", msg)

			case <-exit:
				return nil
			}
		}
	}
}

func subMessagesFunc(conn stan.Conn, env Env, sub chan<- *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		s, err := conn.Subscribe(env.AutobotTopicMessages(), func(m *stan.Msg) {
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

		<-exit

		// close channel
		close(sub)

		return nil
	}
}

func subRepliesFunc(conn stan.Conn, env Env, sub chan<- *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		s, err := conn.Subscribe(env.AutobotTopicReplies(), func(m *stan.Msg) {
			event := new(pb.Event)
			if err := proto.Unmarshal(m.Data, event); err != nil {
				// no nothing now
				return
			}

			log.Printf("got new event: %v", event)

			sub <- event
		})
		if err != nil {
			return err
		}

		defer s.Close()

		<-exit

		// close channel
		close(sub)

		return nil
	}
}

func configureClient(p *plugin) error {
	sc, err := stan.Connect(p.env.AutobotClusterID(), p.name)
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
