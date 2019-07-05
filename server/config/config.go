package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/andersnormal/autobot/pkg/cmd"
	"github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/utils"
	pb "github.com/andersnormal/autobot/proto"

	log "github.com/sirupsen/logrus"
)

const (
	defaultInbox  = "inbox"
	defaultOutbox = "outbox"
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

	// DefaultBotName ...
	DefaultBotName = "autobot"

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

	// DefaultPluginsDir is the default directory to find plugins
	DefaultPluginsDir = "plugins"

	// DefaultFileChmod ...
	DefaultFileChmod = 0600
)

// New returns a new Config
func New() *Config {
	return &Config{
		Verbose:        DefaultVerbose,
		LogLevel:       DefaultLogLevel,
		ReloadSignal:   DefaultReloadSignal,
		TermSignal:     DefaultTermSignal,
		KillSignal:     DefaultKillSignal,
		StatusAddr:     DefaultStatusAddr,
		Debug:          DefaultDebug,
		DataDir:        DefaultDataDir,
		Addr:           DefaultAddr,
		Nats:           DefaultNats,
		NatsClusterID:  DefaultNatsClusterID,
		NatsClusterURL: DefaultNatsClusterURL,
		NatsDataDir:    DefaultNatsDataDir,
		PluginsDir:     DefaultPluginsDir,
		FileChmod:      DefaultFileChmod,
		BotName:        DefaultBotName,
	}
}

// NatsFilestoreDir returns the
func (c *Config) NatsFilestoreDir() string {
	return path.Join(c.DataDir, c.NatsDataDir)
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
	return strings.Join([]string{c.NatsPrefix, defaultInbox}, ".")
}

// Outbox ...
func (c *Config) Outbox() string {
	return strings.Join([]string{c.NatsPrefix, defaultOutbox}, ".")
}

// Env ...
func (c *Config) Env() cmd.Env {
	env := cmd.DefaultEnv()

	env[plugins.AutobotClusterURL] = c.NatsClusterURL
	env[plugins.AutobotClusterID] = c.NatsClusterID
	env[plugins.AutobotChannelInbox] = c.Inbox()
	env[plugins.AutobotChannelOutbox] = c.Outbox()
	env[plugins.AutobotName] = c.BotName

	for _, e := range c.PluginEnv {
		s := strings.Split(e, "=")
		env[s[0]] = s[1]
	}

	return env
}

// Plugins ...
func (c *Config) Plugins() ([]*pb.Plugin, error) {
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

	// walk the plugins dir and fetch the a
	err = filepath.Walk(path.Join(dir, c.PluginsDir), func(p string, info os.FileInfo, err error) error {
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

	return pp, nil
}
