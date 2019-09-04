package config

import (
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/andersnormal/autobot/pkg/discovery"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	"github.com/andersnormal/autobot/pkg/utils"
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
			Inbox:     runtime.DefaultClusterOutbox,
			Outbox:    runtime.DefaultClusterOutbox,
			Discovery: runtime.DefaultClusterDiscovery,
		},
		Plugins:   DefaultPlugins,
		FileChmod: DefaultFileChmod,
	}
}

// NatsFilestoreDir returns the
func (c *Config) NatsFilestoreDir() string {
	return path.Join(c.DataDir, c.Nats.DataDir)
}

// Cwd ...
func (c *Config) Cwd() (string, error) {
	return os.Getwd()
}

// PWD ...
func (c *Config) Pwd() (string, error) {
	return filepath.Abs(os.Args[0])
}

// Dir ...
func (c *Config) Dir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// Plugins ...
func (c *Config) LoadPlugins() ([]*discovery.Plugin, error) {
	var pp []*discovery.Plugin

	// current dir of the bin
	dir, err := c.Dir()
	if err != nil {
		return nil, err
	}

	// currend pwd ...
	pwd, err := c.Pwd()
	if err != nil {
		return nil, err
	}

	// create dir if not exists
	if err := utils.CreateDirIfNotExist(dir, c.FileChmod); err != nil {
		return nil, err
	}

	// get all plugins from all directories
	for _, p := range c.Plugins {
		// walk the plugins dir and fetch the a
		err = filepath.Walk(path.Join(dir, p), func(p string, info os.FileInfo, err error) error {
			// do not start current process
			if pwd == p || info == nil {
				return nil
			}

			// only add files
			if !info.IsDir() && path.Ext(info.Name()) == "" {
				p := &discovery.Plugin{
					Name: path.Base(p),
					Path: p,
				}

				pp = append(pp, p)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return pp, nil
}
