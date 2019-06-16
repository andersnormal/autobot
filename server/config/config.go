package config

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = log.WarnLevel

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

	// DefaultAddr is the default addrs to listen on
	DefaultAddr = ":443"

	// DefaultDebug is the default debug status.
	DefaultDebug = false

	// DefaultDataDir ...
	DefaultDataDir = "data"
)

// New returns a new Config
func New() *Config {
	return &Config{
		Verbose:      DefaultVerbose,
		LogLevel:     DefaultLogLevel,
		ReloadSignal: DefaultReloadSignal,
		TermSignal:   DefaultTermSignal,
		KillSignal:   DefaultKillSignal,
		StatusAddr:   DefaultStatusAddr,
		Debug:        DefaultDebug,
		DataDir:      DefaultDataDir,
		Addr:         DefaultAddr,
	}
}
