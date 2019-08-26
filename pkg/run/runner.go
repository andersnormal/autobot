package run

import (
	"context"
	"io"
	"os/exec"
	"sync"

	"github.com/andersnormal/autobot/pkg/cmd"
	pb "github.com/andersnormal/autobot/proto"
	"github.com/andersnormal/pkg/server"

	"github.com/cenkalti/backoff"
	log "github.com/sirupsen/logrus"
)

// Runner ...
type Runner interface {
	server.Listener
}

type runner struct {
	opts *Opts

	plugins []*pb.Plugin
	env     cmd.Env

	logger *log.Entry

	cancel  func()
	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
}

// New ...
func New(plugins []*pb.Plugin, env cmd.Env, logger *log.Entry, opts ...Opt) Runner {
	options := new(Opts)

	r := new(runner)
	r.plugins = plugins
	r.env = env
	r.opts = options
	r.logger = logger

	configure(r, opts...)

	return r
}

// Start ...
func (r *runner) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		ctx, cancel := context.WithCancel(ctx)
		r.cancel = cancel
		defer cancel()

		for _, p := range r.plugins {
			r.exec(ctx, p)
		}

		ready()

		if err := r.wait(); err != nil {
			return err
		}

		return nil
	}
}

// Stop is stopping the queue
func (r *runner) Stop() error {
	return nil
}

func (r *runner) wait() error {
	r.wg.Wait()

	return r.err
}

func (r *runner) exec(ctx context.Context, p *pb.Plugin) {
	r.run(func() error {
		err := backoff.Retry(func() error {
			c := exec.CommandContext(ctx, p.GetPath())

			stdout, err := c.StdoutPipe()
			if err != nil {
				return err
			}

			stderr, err := c.StderrPipe()
			if err != nil {
				return err
			}

			w := r.logger.Writer()
			defer w.Close()

			go io.Copy(w, stdout)
			go io.Copy(w, stderr)

			cc := cmd.New(c, r.env)

			if err := cc.Run(ctx); err != nil {
				return err
			}

			return nil
		}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *runner) run(f func() error) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()

		if err := f(); err != nil {
			r.errOnce.Do(func() {
				r.err = err
				r.cancel()
			})
		}
	}()
}

func configure(r *runner, opts ...Opt) error {
	for _, o := range opts {
		o(r.opts)
	}

	return nil
}
