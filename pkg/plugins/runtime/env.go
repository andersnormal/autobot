package runtime

import (
	"os"
	"path"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/ianschenck/envflag"
)

const (
	// AutobotName is the name of the plugin.
	// This is also the name of the environment variable.
	AutobotName = "AUTOBOT_NAME"
	// AutobotClusterID is the id of the NATS cluster.
	// This is also the name of the environment variable.
	AutobotClusterID = "AUTOBOT_CLUSTER_ID"
	// AutobotClusterURL is the URL of the NATS cluster.
	// This is also the name of the environment variable.
	AutobotClusterURL = "AUTOBOT_CLUSTER_URL"
	// AutobotClusterDiscovery is the name of the topic for the plugin discovery.
	// This is also the name of the environment variable.
	AutobotClusterDiscovery = "AUTOBOT_CHANNEL_DISCOVERY"
	// AutobotDebug signals that the plugin needs to provide debug output.
	AutobotDebug = "AUTOBOT_DEBUG"
	// AutobotVerbose signals that the plugin needs to provide more verbosity.
	AutobotVerbose = "AUTOBOT_VERBOSE"
	// AutobotLogFormat is the log format to be used for the log output.
	AutobotLogFormat = "AUTOBOT_LOG_FORMAT"
	// AutobotLogLevel is the log level to use for the log output.
	AutobotLogLevel = "AUTOBOT_LOG_LEVEL"
)

// Env describes a run time environment for a plugin.
// This contains information about the used NATS cluster,
// the cluster id and the topic for plugin discovery.
type Env struct {
	Name             string
	ClusterID        string
	ClusterURL       string
	ClusterDiscovery string
	LogFormat        string
	LogLevel         string
	Debug            bool
	Verbose          bool
	Inbox            string
	Outbox           string
}

// DefaultEnv returns the default environment for a plugin.
// This reads in the environment variables and the command line
// parameters to configure the plugin runtime environment.
// The command line flags override any environment variable.
func DefaultEnv() Env {
	env := Env{}

	envflag.StringVar(&env.Name, AutobotName, path.Base(os.Args[0]), "bot name")
	envflag.StringVar(&env.ClusterID, AutobotClusterID, "autobot", "cluster id")
	envflag.StringVar(&env.ClusterURL, AutobotClusterURL, "nats://localhost:4222", "cluster url")
	envflag.StringVar(&env.ClusterDiscovery, AutobotClusterDiscovery, "autobot.discovery", "cluster discovery topic")
	envflag.BoolVar(&env.Verbose, AutobotVerbose, true, "verbosity")
	envflag.BoolVar(&env.Debug, AutobotDebug, false, "debug output")
	envflag.StringVar(&env.LogFormat, AutobotLogFormat, "text", "log format")
	envflag.StringVar(&env.LogLevel, AutobotLogLevel, "info", "log level")

	envflag.Parse()

	return env
}

// WithConfig returns a new runtime environment with
// a proto.Config mapped to the environment properties.
func (e Env) WithConfig(cfg *pb.Config) Env {
	e.Inbox = cfg.GetInbox()
	e.Outbox = cfg.GetOutbox()
	e.Debug = cfg.GetDebug()

	return e
}

// WithInbox returns a new environment with this inbox topic name.
func (e Env) WithInbox(inbox string) Env {
	e.Inbox = inbox

	return e
}

// WithOutbox returns a new environment with this outbox topic name..
func (e Env) WithOutbox(outbox string) Env {
	e.Outbox = outbox

	return e
}
