package plugins

import (
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
)

// Plugin ...
type Plugin interface {
	Subscribe() <-chan *pb.Event
	Wait() error
}

type plugin struct {
	env *Env

	sub chan *pb.Event

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Plugins ...
func New(e *Env) Plugin {
	p := new(plugin)

	p.env = e
	p.sub = make(chan *pb.Event)

	return p
}

func (p *plugin) Subscribe() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.run(p.subFunc(sub))

	return sub
}

func (p *plugin) subFunc(sub chan *pb.Event) func() error {
	return func() error {
		sc, err := stan.Connect(p.env.Get(AutobotClusterID), p.env.Get(AutobotClusterURL))
		if err != nil {
			return err
		}
		defer sc.Close()

		s, err := sc.Subscribe(p.env.Get(AutobotTopic), func(m *stan.Msg) {
			sub <- &pb.Event{}
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
