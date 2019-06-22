package nats

import (
	"net"
	"time"

	s "github.com/andersnormal/pkg/server"
	natsd "github.com/nats-io/gnatsd/server"
	stand "github.com/nats-io/nats-streaming-server/server"
	log "github.com/sirupsen/logrus"
)

type Nats interface {
	// ClusterID ...
	ClusterID() string
	// Addr ...
	Addr() net.Addr
	// MonitorAddr ...
	MonitorAddr() *net.TCPAddr

	s.Listener
}

type nats struct {
	ns  *natsd.Server
	ss  *stand.StanServer

	opts *Opts

	// logger instance
	logger *log.Entry
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	ID      string
  Timeout time.Duration
  Verbose bool
  Debug bool
  Dir string
}
