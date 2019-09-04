package config

import (
	"os"
	"syscall"
	"time"
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
	// Plugins ...
	Plugins []string
	// FileChmod ...
	FileChmod os.FileMode
	// Env ...
	Env []string
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
	// HTTPPort ...
	HTTPPort int
	// Port ...
	Port int
}
