package config

import (
	"os"
	"path"
	"syscall"
	"time"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
)

// Config contains a configuration for Autobot
type Config struct {
	// File is a config file provided
	File string
	// Verbose toggles the verbosity
	Verbose bool
	// LogLevel is the level with with to log for this config
	LogLevel string `mapstructure:"log_level"`
	// LogFormat is the format that is used for logging
	LogFormat string `mapstructure:"log_format"`
	// ReloadSignal
	ReloadSignal syscall.Signal
	// TermSignal
	TermSignal syscall.Signal
	// KillSignal
	KillSignal syscall.Signal
	// Timeout of the runtime
	Timeout time.Duration
	// StatusAddr is the addr of the debug listener
	StatusAddr string
	// Addr is the address to listen on
	Addr string
	// Debug ...
	Debug bool
	// DataDir ...
	DataDir string
	// FileChmod ...
	FileChmod os.FileMode
	// Nats ...
	Nats *Nats
}

// Nats ...
type Nats struct {
	// Disabked ...
	Disabled bool
	// Clustering ...
	Clustering bool
	// Bootstrap ...
	Bootstrap bool
	// ClusterPeers ...
	ClusterPeers []string
	// ClusterNodeID ...
	ClusterNodeID string
	// ClusterID ...
	ClusterID string
	// ClusterURL ...
	ClusterURL string
	// Inbox ...
	Inbox string
	// Outbox ...
	Outbox string
	// Discovery ...
	Discovery string
	// DataDir is the directory for Nats
	DataDir string
	// LogDir is the directory for the Raft logs
	LogDir string
	// HTTPPort ...
	HTTPPort int
	// SPort ...
	SPort int
	// SPort ...
	NPort int
}

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
	// DefaultNatsClusterURL ...
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
	DefaultSNatsPort = 4223
	// DefaultNatsPort ...
	DefaultNNatsPort = 4222
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
			NPort:     DefaultNNatsPort,
			SPort:     DefaultSNatsPort,
			LogDir:    DefaultNatsRaftLogDir,
		},
		FileChmod: DefaultFileChmod,
	}
}

// NatsFilestoreDir returns the nats data dir
func (c *Config) NatsFilestoreDir() string {
	return path.Join(c.DataDir, c.Nats.DataDir)
}

// RaftLogPath returns the raft log dir
func (c *Config) RaftLogPath() string {
	return path.Join(c.DataDir, c.Nats.LogDir)
}
