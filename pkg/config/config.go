package config

import (
	"path"
	"syscall"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
)

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = "warn"
	// DefaultLogFormat is the default format of the logger
	DefaultLogFormat = "text"
	// DefaultTermSignal is the signal to term the agent.
	DefaultTermSignal = syscall.SIGTERM
	// DefaultReloadSignal is the default signal for reload.
	DefaultReloadSignal = syscall.SIGHUP
	// DefaultKillSignal is the default signal for termination.
	DefaultKillSignal = syscall.SIGINT
	// DefaultVerbose is the default verbosity.
	DefaultVerbose = false
	// DefaultStatusAddr is the default addrs for debug listener
	DefaultStatusAddr = ":8443"
	// DefaultDebug is the default debug status.
	DefaultDebug = false
	// DefaultDataDir ...
	DefaultDataDir = "data"
	// DefaultNats ...
	DefaultNats = true
	// DefaultNatsURL ...
	DefaultNatsClusterURL = "nats://localhost:4222"
	// DefaultNatsDataDir is the default directory for nats data
	DefaultNatsDataDir = "nats"
	// DefaultNatsClusterID ...
	DefaultNatsClusterID = "autobot"
	// DefaultNatsRaftLogDir ...
	DefaultNatsRaftLogDir = "raft"
	// DefaultNatsHTTPPort ...
	DefaultNatsHTTPPort = 8223
	// DefaultNatsPort ...
	DefaultNatsPort = 4223
	// DefaultFileChmod ...
	DefaultFileChmod = 0600
)

// New returns a new Config
func New() *Config {
	return &Config{
		Verbose:      DefaultVerbose,
		LogLevel:     DefaultLogLevel,
		LogFormat:    DefaultLogFormat,
		ReloadSignal: DefaultReloadSignal,
		TermSignal:   DefaultTermSignal,
		KillSignal:   DefaultKillSignal,
		StatusAddr:   DefaultStatusAddr,
		Debug:        DefaultDebug,
		DataDir:      DefaultDataDir,
		Nats: &Nats{
			ClusterID: DefaultNatsClusterID,
			DataDir:   DefaultNatsDataDir,
			HTTPPort:  DefaultNatsHTTPPort,
			Inbox:     runtime.DefaultClusterInbox,
			Outbox:    runtime.DefaultClusterOutbox,
			Port:      DefaultNatsPort,
			LogDir:    DefaultNatsRaftLogDir,
		},
		FileChmod: DefaultFileChmod,
	}
}

// NatsFilestoreDir returns the
func (c *Config) NatsFilestoreDir() string {
	return path.Join(c.DataDir, c.Nats.DataDir)
}

// NatsFilestoreDir returns the
func (c *Config) RaftLogPath() string {
	return path.Join(c.DataDir, c.Nats.LogDir)
}
