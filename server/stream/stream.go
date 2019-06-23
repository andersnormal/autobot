package stream

import (
	"github.com/andersnormal/autobot/pkg/nats"

	"github.com/nats-io/stan.go"
)

// New ...a
func New(nats nats.Nats, opts ...Opt) Stream {
	options := new(Opts)

	s := new(stream)
	s.opts = options

	configure(s, opts...)

	return s
}

func (s *stream) Subscribe() (chan *stan.Msg, error) {
	// connect to cluster
	sc, err := stan.Connect(s.nats.ClusterID(), "server")
	if err != nil {
		return nil, err
	}

	defer sc.Close()

	msg := make(chan *stan.Msg)
	// Simple Async Subscriber
	sub, err := sc.Subscribe("events", func(m *stan.Msg) {

	})

	// Unsubscribe
	defer sub.Unsubscribe()

	// noop
	return
}

func (s *stream) Publish(m []byte) error {
	// connect to cluster
	sc, err := stan.Connect(s.nats.ClusterID(), "server")
	if err != nil {
		return err
	}

	if err := sc.Publish("messages", m); err != nil {
		return err
	}

	return nil
}

func configure(s *stream, opts ...Opt) error {
	for _, o := range opts {
		o(s.opts)
	}

	return nil
}
