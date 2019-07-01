package plugins

import (
	"log"
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
	AutobotTopicMessages = "AUTOBOT_TOPIC_MESSAGES"
	AutobotTopicReplies  = "AUTOBOT_TOPIC_REPLIES"
)

// Plugin ...
type Plugin interface {
	SubscribeMessages() <-chan *pb.Event
	PublishMessages() chan<- *pb.Event
	SubscribeReplies() <-chan *pb.Event
	PublishReplies() chan<- *pb.Event

	// Wait ...
	Wait() error
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
}

type plugin struct {
	name string

	opts *Opts
	conn stan.Conn

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

	p.run(subMessagesFunc(p.conn, sub, p.exit))

	return sub
}

// SubscribeReplies ...
func (p *plugin) SubscribeReplies() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.run(subRepliesFunc(p.conn, sub, p.exit))

	return sub
}

// PublishMessages ...
func (p *plugin) PublishMessages() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.run(pubMessagesFunc(p.conn, pub, p.exit))

	return pub
}

// PublishReplies ...
func (p *plugin) PublishReplies() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.run(pubRepliesFunc(p.conn, pub, p.exit))

	return pub
}

// Wait ...
func (p *plugin) Wait() error {
	<-p.exit

	return p.err
}

func (p *plugin) run(f func() error) {
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

func pubMessagesFunc(conn stan.Conn, pub <-chan *pb.Event, exit <-chan struct{}) func() error {
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

				if err := conn.Publish(os.Getenv(AutobotTopicMessages), msg); err != nil {
					return err
				}
			case <-exit:
				return nil
			}
		}
	}
}

func pubRepliesFunc(conn stan.Conn, pub <-chan *pb.Event, exit <-chan struct{}) func() error {
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

				if err := conn.Publish(os.Getenv(AutobotTopicReplies), msg); err != nil {
					return err
				}

				log.Printf("published: %v", msg)

			case <-exit:
				return nil
			}
		}
	}
}

func subMessagesFunc(conn stan.Conn, sub chan<- *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		s, err := conn.Subscribe(os.Getenv(AutobotTopicMessages), func(m *stan.Msg) {
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

func subRepliesFunc(conn stan.Conn, sub chan<- *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		s, err := conn.Subscribe(os.Getenv(AutobotTopicReplies), func(m *stan.Msg) {
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
	sc, err := stan.Connect(os.Getenv(AutobotClusterID), p.name)
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
