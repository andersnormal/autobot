package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// Server is the interface to the server
type Server interface {
	// Run is running a new routine
	Listen(listener Listener, ready bool)
	// Name
	Name() string
	// Env
	Env() string
	// Waits for the server to fail,
	// or gracefully shutdown if called
	Wait() error
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	// Name
	Name string
	// Env
	Env string
	// ReloadSignal
	ReloadSignal syscall.Signal
	// TermSignal
	TermSignal syscall.Signal
	// KillSignal
	KillSignal syscall.Signal
}

type listeners map[Listener]bool

// server holds the instance info of the server
type server struct {
	errGroup *errgroup.Group
	errCtx   context.Context
	cancel   context.CancelFunc

	listeners map[Listener]bool

	ready chan bool
	sys   chan os.Signal

	opts *Opts
}

// NewServer ....
func NewServer(ctx context.Context, opts ...Opt) Server {
	options := &Opts{}

	s := new(server)
	s.opts = options

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.errGroup, s.errCtx = errgroup.WithContext(ctx)
	s.listeners = make(listeners)
	s.ready = make(chan bool)

	configure(s, opts...)
	configureSignals(s)

	return s
}

// Listen ...
func (s *server) Listen(listener Listener, ready bool) {
	if _, found := s.listeners[listener]; found {
		return
	}

	s.listeners[listener] = false
	g := s.errGroup

	g.Go(listener.Start(s.errCtx, func() { s.ready <- true }))

	if ready {
		<-s.ready
	}
}

// Env ...
func (s *server) Env() string {
	return s.opts.Env
}

// Name ...
func (s *server) Name() string {
	return s.opts.Name
}

// Wait ...
func (s *server) Wait() error {
	// create ticker for interrupt signals
	ticker := time.NewTicker(1 * time.Second)

	ctx := s.errCtx

	for {
		select {
		case <-ticker.C:
		case <-s.sys:
			s.cancel()
		case <-ctx.Done():
			for listener := range s.listeners {
				listener.Stop()
			}

			if err := ctx.Err(); err != nil {
				return err
			}

			return nil
		}
	}
}

// WithName ...
func WithName(name string) func(o *Opts) {
	return func(o *Opts) {
		o.Name = name
	}
}

// WithEnv ...
func WithEnv(env string) func(o *Opts) {
	return func(o *Opts) {
		o.Env = env
	}
}

// Listener is the interface to a listener,
// so starting and shutdown of a listener,
// or any routine.
type Listener interface {
	Start(ctx context.Context, ready func()) func() error
	Stop() error
}

func configureSignals(s *server) {
	s.sys = make(chan os.Signal, 1)

	signal.Notify(s.sys, s.opts.ReloadSignal, s.opts.KillSignal, s.opts.TermSignal)
}

func configure(s *server, opts ...Opt) error {
	for _, o := range opts {
		o(s.opts)
	}

	s.opts.TermSignal = syscall.SIGTERM
	s.opts.ReloadSignal = syscall.SIGHUP
	s.opts.KillSignal = syscall.SIGINT

	return nil
}
