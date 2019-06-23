package stream

import (
	"sync"

	"github.com/andersnormal/voskhod/server/nats"
)

type Stream interface {
}

type stream struct {
	nats nats.Nats

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup

	opts *Opts
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
}
