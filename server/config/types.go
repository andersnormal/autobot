package config

import (
	"os"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config contains a configuration for Autobot
type Config struct {
	// Verbose toggles the verbosity
	Verbose bool

	// LogLevel is the level with with to log for this config
	LogLevel log.Level

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

	// NatsDataDir is the directory for Nats
	NatsDataDir string

	// PluginsDir ...
	PluginsDir string

	// FileChmod ...
	FileChmod os.FileMode
}
