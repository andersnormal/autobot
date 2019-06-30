package plugins

import (
	"os"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
)

const (
	AutobotClusterID     = "AUTOBOT_CLUSTER_ID"
	AutobotClusterURL    = "AUTOBOT_CLUSTER_URL"
	AutobotTopicMessages = "AUTOBOT_TOPIC_MESSAGEES"
	AutobotTopicReplies  = "AUTOBOT_TOPIC_REPLIES"
)

// Plugin ...
type Plugin interface {
	SubscribeMessages() <-chan *pb.Message
	PublishMessages() chan<- *pb.Message
	SubscribeReplies() <-chan *pb.Reply
	PublishReplies() chan<- *pb.Reply

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

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Plugins ...
func New(name string, opts ...Opt) Plugin {
	options := new(Opts)

	p := new(plugin)
	// setting a default env for a plugin
	p.opts = options
	p.name = name

	// configure plugin
	configure(p, opts...)

	return p
}

// SubscribeMessages ...
func (p *plugin) SubscribeMessages() <-chan *pb.Message {
	sub := make(chan *pb.Message)

	p.run(subMessagesFunc(p.name, sub, p.exit))

	return sub
}

// SubscribeReplies ...
func (p *plugin) SubscribeReplies() <-chan *pb.Reply {
	sub := make(chan *pb.Reply)

	p.run(subRepliesFunc(p.name, sub, p.exit))

	return sub
}

// PublishMessages ...
func (p *plugin) PublishMessages() chan<- *pb.Message {
	pub := make(chan *pb.Message)

	p.run(pubMessagesFunc(p.name, pub, p.exit))

	return pub
}

// PublishReplies ...
func (p *plugin) PublishReplies() chan<- *pb.Reply {
	pub := make(chan *pb.Reply)

	p.run(pubRepliesFunc(p.name, pub, p.exit))

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

func pubMessagesFunc(name string, pub <-chan *pb.Message, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(os.Getenv(AutobotClusterID), name)
		if err != nil {
			return err
		}
		defer sc.Close()

		for {
			select {
			case _, ok := <-pub:
				if !ok {
					return nil
				}

				if err := sc.Publish(os.Getenv(AutobotTopicMessages), []byte("")); err != nil {
					return err
				}
			case <-exit:
				return nil
			}
		}
	}
}

func pubRepliesFunc(name string, pub <-chan *pb.Reply, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(os.Getenv(AutobotClusterID), name)
		if err != nil {
			return err
		}
		defer sc.Close()

		for {
			select {
			case _, ok := <-pub:
				if !ok {
					return nil
				}

				if err := sc.Publish(os.Getenv(AutobotTopicMessages), []byte("")); err != nil {
					return err
				}
			case <-exit:
				return nil
			}
		}
	}
}

func subMessagesFunc(name string, sub chan<- *pb.Message, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(os.Getenv(AutobotClusterID), name)
		if err != nil {
			return err
		}
		defer sc.Close()

		s, err := sc.Subscribe(os.Getenv(AutobotTopicMessages), func(m *stan.Msg) {
			sub <- &pb.Message{}
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

func subRepliesFunc(name string, sub chan<- *pb.Reply, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(os.Getenv(AutobotClusterID), name)
		if err != nil {
			return err
		}
		defer sc.Close()

		s, err := sc.Subscribe(os.Getenv(AutobotTopicMessages), func(m *stan.Msg) {
			sub <- &pb.Reply{}
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

func configure(p *plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	return nil
}
