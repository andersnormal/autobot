package runtime

import (
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// DefaultClusterID ...
	DefaultClusterID = "autobot"
	// DefaultClusterURL ...
	DefaultClusterURL = "nats://localhost:4222"
	// DefaultClusterDiscovery ...
	DefaultClusterDiscovery = "autobot.discovery"
	// DefaultClusterInbox ...
	DefaultClusterInbox = "autobot.inbox"
	// DefaultClusterOutbox ...
	DefaultClusterOutbox = "autobot.outbox"
)

var initializers []func()
var runtime Runtime

// OnInitialize sets the passed functions to be run when each command's
// Execute method is called.
func OnInitialize(y ...func()) {
	initializers = append(initializers, y...)
}

func runInitializers() {
	for _, fn := range initializers {
		fn()
	}
}

func init() {
	runtime = Runtime{}

	viper.SetEnvPrefix("autobot")
	viper.AutomaticEnv()

	// set some default flags
	pflag.String("name", path.Base(os.Args[0]), "plugin name")
	pflag.String("cluster_id", DefaultClusterID, "cluster id")
	pflag.String("cluster_url", DefaultClusterURL, "cluster url")
	pflag.String("cluster_discovery", DefaultClusterDiscovery, "cluster discovery")
	pflag.String("inbox", DefaultClusterInbox, "cluster inbox")
	pflag.String("outbox", DefaultClusterOutbox, "cluster outbot")
	pflag.String("log_format", "text", "log format")
	pflag.String("log_level", "text", "log level")
	pflag.BoolP("verbose", "v", true, "verbose")
	pflag.BoolP("debug", "d", true, "verbose")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// unmarshal to config
	if err := viper.Unmarshal(&runtime); err != nil {
		log.Fatalf(errors.Wrap(err, "cannot unmarshal runtime").Error())
	}
}

// Runtime describes a runtime environment for a plugin.
// This contains information about the used NATS cluster,
// the cluster id and the topic for plugin discovery.
type Runtime struct {
	ClusterDiscovery string `mapstructure:"cluster_discovery"`
	ClusterID        string `mapstructure:"cluster_id"`
	ClusterURL       string `mapstructure:"cluster_url"`
	Debug            bool
	Inbox            string
	LogFormat        string `mapstructure:"log_format"`
	LogLevel         string `mapstructure:"log_level"`
	Name             string
	Outbox           string
	Verbose          bool
}

// Default returns the default environment for a plugin.
// This reads in the environment variables and the command line
// parameters to configure the plugin runtime environment.
// The command line flags override any environment variable.
func Default() Runtime {
	runInitializers()

	return runtime
}
