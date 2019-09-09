package nats

import (
	"context"
	"net"
	"time"

	"github.com/andersnormal/autobot/pkg/config"

	s "github.com/andersnormal/pkg/server"
	natsd "github.com/nats-io/gnatsd/server"
	stand "github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/nats-streaming-server/stores"
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
	ns *natsd.Server
	ss *stand.StanServer

	opts *Opts
	cfg  *config.Config

	// logger instance
	logger *log.Entry
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Timeout time.Duration
}

// New returns a new server
func New(cfg *config.Config, opts ...Opt) Nats {
	options := new(Opts)

	n := new(nats)
	n.opts = options
	n.cfg = cfg

	n.logger = log.WithFields(log.Fields{})

	configure(n, opts...)

	return n
}

// Timeout ...
func Timeout(t time.Duration) func(o *Opts) {
	return func(o *Opts) {
		o.Timeout = t
	}
}

// Start is starting the queue
func (n *nats) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		defer func() { ready() }()

		// create logger ...
		logger := NewLogger()
		logger.SetLogger(n.logger)

		// creating NATS ...
		nopts := new(natsd.Options)
		nopts.HTTPPort = n.cfg.Nats.HTTPPort
		nopts.Port = n.cfg.Nats.Port
		nopts.NoSigs = true

		n.ns = n.startNatsd(nopts, logger) // wait for the Nats server to come available
		if !n.ns.ReadyForConnections(n.opts.Timeout) {
			return NewError("could not start Nats server in %s seconds", n.opts.Timeout)
		}

		// Get NATS Streaming Server default options
		opts := stand.GetDefaultOptions()
		opts.StoreType = stores.TypeFile
		opts.FilestoreDir = n.cfg.NatsFilestoreDir()
		opts.ID = n.cfg.Nats.ClusterID

		// set custom logger
		opts.CustomLogger = logger

		// Do not handle signals
		opts.HandleSignals = false
		opts.EnableLogging = true
		opts.Debug = n.cfg.Debug
		opts.Trace = n.cfg.Debug

		// clustering mode
		opts.Clustering.Clustered = n.cfg.Nats.Clustering
		opts.Clustering.Bootstrap = n.cfg.Nats.Bootstrap
		opts.Clustering.NodeID = n.cfg.Nats.ClusterNodeID
		opts.Clustering.Peers = n.cfg.Nats.ClusterPeers
		opts.Clustering.LogCacheSize = stand.DefaultLogCacheSize
		opts.Clustering.LogSnapshots = stand.DefaultLogSnapshots
		opts.Clustering.RaftLogPath = n.cfg.RaftLogPath()
		opts.NATSServerURL = n.cfg.Nats.ClusterURL

		// Now we want to setup the monitoring port for NATS Streaming.
		// We still need NATS Options to do so, so create NATS Options
		// using the NewNATSOptions() from the streaming server package.
		snopts := stand.NewNATSOptions()
		snopts.HTTPPort = 8222
		snopts.NoSigs = true

		// Now run the server with the streaming and streaming/nats options.
		ss, err := stand.RunServerWithOpts(opts, snopts)
		if err != nil {
			return err
		}
		n.ss = ss

		ready()

		<-ctx.Done()

		n.ns.Shutdown()
		n.ss.Shutdown()

		// noop
		return nil
	}
}

// ClusterID ...
func (n *nats) ClusterID() string {
	return n.ss.ClusterID()
}

// MonitorAddr ...
func (n *nats) MonitorAddr() *net.TCPAddr {
	return n.ns.MonitorAddr()
}

// Addr ...
func (n *nats) Addr() net.Addr {
	return n.ns.Addr()
}

func (n *nats) startNatsd(nopts *natsd.Options, l natsd.Logger) *natsd.Server {
	// Create the NATS Server
	ns := natsd.New(nopts)
	ns.SetLogger(l, n.cfg.Debug, n.cfg.Debug)

	// Start it as a go routine
	go ns.Start()

	return ns
}

func configure(n *nats, opts ...Opt) error {
	for _, o := range opts {
		o(n.opts)
	}

	return nil
}
