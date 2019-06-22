package nats

import (
	"context"
	"net"
	"time"

	natsd "github.com/nats-io/gnatsd/server"
	stand "github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/nats-streaming-server/stores"

	log "github.com/sirupsen/logrus"
)

const (
	defaultNatsHTTPPort = 8223
	defaultNatsPort     = 4223
)

// New returns a new server
func New(opts ...Opt) Nats {
	options := new(Opts)

	n := new(nats)
	n.opts = options

	n.logger = log.WithFields(log.Fields{})

	configure(n, opts...)

	return n
}

// WithID ...
func WithID(id string) func(o *Opts) {
	return func(o *Opts) {
		o.ID = id
	}
}

// WithTimeout ...
func WithTimeout(t time.Duration) func(o *Opts) {
	return func(o *Opts) {
		o.Timeout = t
	}
}

// WithDataDir ...
func WithDataDir(dir string) func(o *Opts) {
	return func(o *Opts) {
		o.Dir = dir
	}
}

// WithDebug ...
func WithDebug() func(o *Opts) {
	return func(o *Opts) {
		o.Debug = true
	}
}

// WithVerbose ...
func WithVerbose() func(o *Opts) {
	return func(o *Opts) {
		o.Verbose = true
	}
}

// Stop is stopping the queue
func (n *nats) Stop() error {
	n.log().Info("shutting down nats...")

	if n.ss != nil {
		n.ss.Shutdown()
	}

	if n.ns != nil {
		n.ns.Shutdown()
	}

	return nil
}

// Start is starting the queue
func (n *nats) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		// creating NATS ...
		nopts := new(natsd.Options)
		nopts.HTTPPort = 8223
		nopts.Port = defaultNatsPort
		nopts.NoSigs = true

		n.ns = n.startNatsd(nopts) // wait for the Nats server to come available
		if !n.ns.ReadyForConnections(n.opts.Timeout * time.Second) {
			return NewError("could not start Nats server in %s seconds", n.opts.Timeout)
		}

		// verbose
		n.log().Infof("started NATS server")

		// Get NATS Streaming Server default options
		opts := stand.GetDefaultOptions()
		opts.StoreType = stores.TypeFile
		opts.FilestoreDir = n.opts.Dir
		opts.ID = n.opts.ID

		// set custom logger
		logger := NewLogger()
		logger.SetLogger(n.log())
		opts.CustomLogger = logger

		// Do not handle signals
		opts.HandleSignals = false
		opts.EnableLogging = true
		opts.Debug = n.opts.Debug
		opts.Trace = n.opts.Verbose

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

		// verbose
		n.log().Infof("started cluster %s", n.ss.ClusterID())

		// wait for the server to be ready
		time.Sleep(n.opts.Timeout)

		// call to be ready
		ready()

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

func (n *nats) startNatsd(nopts *natsd.Options) *natsd.Server {
	// Create the NATS Server
	ns := natsd.New(nopts)

	// Start it as a go routine
	go ns.Start()

	return ns
}

func (n *nats) log() *log.Entry {
	return n.logger
}

func configure(n *nats, opts ...Opt) error {
	for _, o := range opts {
		o(n.opts)
	}

	return nil
}
