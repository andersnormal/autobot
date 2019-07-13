package discovery

import (
	"context"
	"sync"

	s "github.com/andersnormal/pkg/server"
	"github.com/nats-io/nats.go"
)

const (
	defaultRegistryTopic = "autobot.discovery"
)

// Registry ...
type Registry interface {
	s.Listener
}

type registry struct {
	opts *Opts
	conn *nats.Conn
	addr string

	sync.RWMutex
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	RegistryTopic string
}

// New ...
func New(addr string, opts ...Opt) Registry {
	options := new(Opts)

	r := new(registry)
	r.opts = options
	r.addr = addr

	configure(r, opts...)

	return r
}

// Start ...
func (r *registry) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		conn, err := r.getConn()
		if err != nil {
			return err
		}

		rr := make(chan *nats.Msg)

		sub, err := conn.Subscribe(r.opts.RegistryTopic, func(msg *nats.Msg) {
			rr <- msg
		})
		if err != nil {
			return err
		}

		defer sub.Unsubscribe()
		defer conn.Close()

		ready()

		for {
			select {
			case msg := <-rr:
				if err := r.register(msg); err != nil {
					return err
				}
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// Stop ...
func (r *registry) Stop() error {
	return nil
}

func (r *registry) register(msg *nats.Msg) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}

	err = conn.PublishMsg(&nats.Msg{
		Subject: msg.Reply,
		Data:    []byte("test"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *registry) getConn() (*nats.Conn, error) {
	r.Lock()
	defer r.Unlock()

	if r.conn != nil {
		return r.conn, nil
	}

	c, err := r.newConn()
	if err != nil {
		return nil, err
	}

	r.conn = c

	return r.conn, nil
}

func (r *registry) newConn() (*nats.Conn, error) {
	// todo: support TLS
	c, err := nats.Connect(r.addr)
	if err != nil {
		return nil, err
	}

	r.conn = c

	return r.conn, nil
}

func configure(r *registry, opts ...Opt) error {
	for _, o := range opts {
		o(r.opts)
	}

	if r.opts.RegistryTopic == "" {
		r.opts.RegistryTopic = defaultRegistryTopic
	}

	return nil
}
