package plugins

import (
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
)

// SubscribeFunc ...
type SubscribeFunc = func(env *Env, sub <-chan *pb.Event, exit <-chan struct{}) func() error

// PublishFunc ...
type PublishFunc = func(env *Env, sub <-chan *pb.Event, exit <-chan struct{}) func() error

var (
	// DefaultSubscribeFunc ...
	DefaultSubscribeFunc = subFunc
	// DefaultPublishFunc ...
	DefaultPublishFunc = pubFunc
)

// Plugin ...
type Plugin interface {
	Subscribe() <-chan *pb.Event
	Publish() chan<- *pb.Event
	Wait() error
}

type plugin struct {
	env *Env

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Plugins ...
func New(e *Env) Plugin {
	p := new(plugin)

	p.env = e

	return p
}

// Subscribe ...
func (p *plugin) Subscribe() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.run(DefaultSubscribeFunc(p.env, sub, p.exit))

	return sub
}

// Publish ...
func (p *plugin) Publish() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.run(DefaultPublishFunc(p.env, pub, p.exit))

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

func pubFunc(env *Env, pub <-chan *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(env.Get(AutobotClusterID), env.Get(AutobotClusterURL))
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

				if err := sc.Publish(env.Get(AutobotTopic), []byte("")); err != nil {
					return err
				}
			case <-exit:
				return nil
			}
		}
	}
}

func subFunc(env *Env, sub chan<- *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(env.Get(AutobotClusterID), env.Get(AutobotClusterURL))
		if err != nil {
			return err
		}
		defer sc.Close()

		s, err := sc.Subscribe(env.Get(AutobotTopic), func(m *stan.Msg) {
			sub <- &pb.Event{}
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
