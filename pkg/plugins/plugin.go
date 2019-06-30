package plugins

import (
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
)

// SubscribeFunc ...
type SubscribeFunc = func(name string, env *Env, sub <-chan *pb.Event, exit <-chan struct{}) func() error

// PublishFunc ...
type PublishFunc = func(name string, env *Env, sub <-chan *pb.Event, exit <-chan struct{}) func() error

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

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Env *Env
}

// WithEnv ...
func WithEnv(e *Env) func(o *Opts) {
	return func(o *Opts) {
		o.Env = e
	}
}

type plugin struct {
	env  *Env
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
	p.env = DefaultEnv()
	p.opts = options
	p.name = name

	// configure plugin
	configure(p, opts...)

	return p
}

// Subscribe ...
func (p *plugin) Subscribe() <-chan *pb.Event {
	sub := make(chan *pb.Event)

	p.run(DefaultSubscribeFunc(p.name, p.env, sub, p.exit))

	return sub
}

// Publish ...
func (p *plugin) Publish() chan<- *pb.Event {
	pub := make(chan *pb.Event)

	p.run(DefaultPublishFunc(p.name, p.env, pub, p.exit))

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

func pubFunc(name string, env *Env, pub <-chan *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(name, env.Get(AutobotClusterURL))
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

				if err := sc.Publish(env.Get(AutobotTopicEvents), []byte("")); err != nil {
					return err
				}
			case <-exit:
				return nil
			}
		}
	}
}

func subFunc(name string, env *Env, sub chan<- *pb.Event, exit <-chan struct{}) func() error {
	return func() error {
		sc, err := stan.Connect(env.Get(AutobotClusterID), name)
		if err != nil {
			return err
		}
		defer sc.Close()

		s, err := sc.Subscribe(env.Get(AutobotTopicEvents), func(m *stan.Msg) {
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

func configure(p *plugin, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	if p.opts.Env != nil {
		p.env = p.opts.Env
	}

	return nil
}
