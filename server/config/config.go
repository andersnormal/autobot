package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/andersnormal/autobot/pkg/cmd"
	"github.com/andersnormal/autobot/pkg/utils"
	pb "github.com/andersnormal/autobot/proto"
)

const (
	defaultInbox     = "inbox"
	defaultOutbox    = "outbox"
	defaultDiscovery = "discovery"
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

	// DefaultAddr is the default addrs to listen on
	DefaultAddr = ":443"

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

	// DefaultNatsPrefix ...
	DefaultNatsPrefix = "autobot"

	// DefaultFileChmod ...
	DefaultFileChmod = 0600

	// DefaultGRPCAddr is the default grpc address.
	DefaultGRPCAddr = ":8888"
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
		Addr:         DefaultAddr,
		Nats: &Nats{
			ClusterID: DefaultNatsClusterID,
			DataDir:   DefaultNatsDataDir,
		},
		Plugins:   DefaultPlugins,
		FileChmod: DefaultFileChmod,
		GRPCAddr:  DefaultGRPCAddr,
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

// Inbox ...
func (c *Config) Inbox() string {
	return strings.Join([]string{c.Nats.Prefix, defaultInbox}, ".")
}

// Outbox ...
func (c *Config) Outbox() string {
	return strings.Join([]string{c.Nats.Prefix, defaultOutbox}, ".")
}

// Discovery ...
func (c *Config) Discovery() string {
	return strings.Join([]string{c.Nats.Prefix, defaultDiscovery}, ".")
}

// Env ...
func (c *Config) PluginEnv() cmd.Env {
	env := cmd.DefaultEnv()

	env["AUTOBOT_CLUSTER_URL"] = c.Nats.ClusterURL
	env["AUTOBOT_CLUSTER_ID"] = c.Nats.ClusterID
	env["AUTOBOT_CLUSTER_DISCOVERY"] = c.Discovery()

	for _, e := range c.Env {
		s := strings.Split(e, "=")
		env[s[0]] = s[1]
	}

	return env
}

// Plugins ...
func (c *Config) LoadPlugins() ([]*pb.Plugin, error) {
	var pp []*pb.Plugin

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
				pp = append(pp, pb.NewPlugin(p))
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return pp, nil
}
