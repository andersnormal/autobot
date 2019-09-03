package discovery

import (
	"context"
	"sync"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"

	s "github.com/andersnormal/pkg/server"
	"github.com/nats-io/nats.go"
)

// Registry ...
type Registry interface {
	s.Listener
}

type registry struct {
	opts    *Opts
	conn    *nats.Conn
	plugins []*pb.Plugin

	sync.RWMutex
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	ClusterDiscvoery string
	ClusterURL       string
}

// New ...
func New(opts ...Opt) Registry {
	options := new(Opts)

	r := new(registry)
	r.opts = options

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

		sub, err := conn.Subscribe(r.opts.ClusterDiscvoery, func(msg *nats.Msg) {
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
			case <-rr:
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
	c, err := nats.Connect(r.opts.ClusterURL)
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

	if r.opts.ClusterDiscvoery == "" {
		r.opts.ClusterDiscvoery = runtime.DefaultClusterDiscovery
	}

	if r.opts.ClusterURL == "" {
		r.opts.ClusterURL = runtime.DefaultClusterURL
	}

	return nil
}
