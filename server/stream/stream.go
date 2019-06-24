package stream

import (
	"context"
	"time"

	"github.com/andersnormal/autobot/pkg/nats"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
)

// New ...a
func New(nats nats.Nats, sub chan pb.Message, pub chan pb.Message, opts ...Opt) Stream {
	options := new(Opts)

	s := new(stream)
	s.nats = nats
	s.opts = options
	s.sub = sub
	s.pub = pub

	configure(s, opts...)

	return s
}

// Stop ...
func (s *stream) Stop() error {
	return nil
}

// Start ...
func (s *stream) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		// wait for the server to be ready
		time.Sleep(5 * time.Second)

		// connect to cluster
		sc, err := stan.Connect(s.nats.ClusterID(), "server")
		if err != nil {
			return err
		}

		s.conn = sc

		// Close connection
		defer func() { sc.Close() }()

		// start subscriptions
		s.run(ctx, s.publishMessages)
		s.run(ctx, s.subscribeMessages)

		// call to be ready
		ready()

		// wait for error to appear here
		if err := s.wait(); err != nil {
			return err
		}

		return nil
	}
}

func (s *stream) subscribeMessages(ctx context.Context) error {
	// Simple Async Subscriber
	sub, err := s.conn.Subscribe("messages", func(m *stan.Msg) {
		// format to message
		s.sub <- pb.Message{}
	})

	if err != nil {
		return err
	}

	// Unsubscribe
	defer sub.Unsubscribe()

	// wait for the ca
	<-ctx.Done()

	// noop
	return nil
}

func (s *stream) publishMessages(ctx context.Context) error {
	for {
		select {
		case _, ok := <-s.pub:
			if !ok {
				return nil
			}

			s.conn.Publish("messages", []byte("test"))
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *stream) wait() error {
	s.wg.Wait()

	return s.err
}

func (s *stream) run(ctx context.Context, fn func(context.Context) error) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		if err := fn(ctx); err != nil {
			s.errOnce.Do(func() {
				s.err = err

				close(s.exit)
			})
		}
	}()
}

func configure(s *stream, opts ...Opt) error {
	for _, o := range opts {
		o(s.opts)
	}

	return nil
}
