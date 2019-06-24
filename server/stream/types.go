package stream

import (
	"sync"

	"github.com/andersnormal/autobot/pkg/nats"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/andersnormal/pkg/server"
	"github.com/nats-io/stan.go"
)

type Stream interface {
	server.Listener
}

type stream struct {
	nats nats.Nats
	conn stan.Conn

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup

	pub chan pb.Message
	sub chan pb.Message

	opts *Opts
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
}
