package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
)

// Env ...
type Env map[string]string

func (ev Env) Strings() []string {
	var env []string
	for k, v := range ev {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

// Set ...
func (ev Env) Set(name string, value string) Env {
	ev[name] = value

	return ev
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
	// DefaultNatsURL ...
	DefaultNatsClusterURL = "nats://localhost:4222"
	// DefaultNatsDataDir is the default directory for nats data
	DefaultNatsDataDir = "nats"
	// DefaultNatsClusterID ...
	DefaultNatsClusterID = "autobot"
	// DefaultNatsHTTPPort ...
	DefaultNatsHTTPPort = 8223
	// DefaultNatsPort ...
	DefaultNatsPort = 4223
	// DefaultFileChmod ...
	DefaultFileChmod = 0600
)

var (
	// DefaultPlugins is the default directory to find plugins
	DefaultPlugins = []string{"plugins"}
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
			HTTPPort:  DefaultNatsHTTPPort,
			Port:      DefaultNatsPort,
			ClusterID: DefaultNatsClusterID,
			DataDir:   DefaultNatsDataDir,
			Inbox:     runtime.DefaultClusterInbox,
			Outbox:    runtime.DefaultClusterOutbox,
			Discovery: runtime.DefaultClusterDiscovery,
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
	// // create dir if not exists
	// if err := utils.CreateDirIfNotExist("raft", c.FileChmod); err != nil {
	// 	return "raft"
	// }

	return "/raft"
}

// Cwd ...
func (c *Config) Cwd() (string, error) {
	return os.Getwd()
}

// Pwd ...
func (c *Config) Pwd() (string, error) {
	return filepath.Abs(os.Args[0])
}

// Dir ...
func (c *Config) Dir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
