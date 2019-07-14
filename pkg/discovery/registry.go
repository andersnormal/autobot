package discovery

import (
	"context"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	s "github.com/andersnormal/pkg/server"
	"github.com/golang/protobuf/proto"
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
	opts    *Opts
	conn    *nats.Conn
	addr    string
	plugins []*pb.Plugin

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
					// this is not resillient, should perhaps reflect on error in reply
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
	action := new(pb.Action)
	if err := proto.Unmarshal(msg.Data, action); err != nil {
		return err
	}

	res, err := proto.Marshal(pb.NewEmpty())
	if err != nil {
		return err
	}

	if action.GetRegister() == nil {
		// try to marshal into []byte ...
		errr, err := proto.Marshal(pb.NewErrRegister("no valid action"))
		if err != nil {
			return err
		}

		res = errr
	}

	r.plugins = append(r.plugins, action.GetRegister().GetPlugin())

	conn, err := r.getConn()
	if err != nil {
		return err
	}

	err = conn.PublishMsg(&nats.Msg{
		Subject: msg.Reply,
		Data:    res,
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
