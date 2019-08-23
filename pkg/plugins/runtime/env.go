package runtime

import (
	"flag"
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
)

// Env describes a run time environment for a plugin.
// This contains information about the used NATS cluster,
// the cluster id and the topic for plugin discovery.
type Env struct {
	Name             string
	ClusterID        string
	ClusterURL       string
	ClusterDiscovery string
	Debug            bool
	Inbox            string
	Outbox           string
}

// DefaultEnv returns the default environment for a plugin.
// This reads in the environment variables and the command line
// parameters to configure the plugin runtime environment.
// The command line flags override any environment variable.
func DefaultEnv() Env {
	env := Env{}

	envflag.StringVar(&env.Name, path.Base(os.Args[0]), "autobot", "cluster id")
	envflag.StringVar(&env.ClusterID, AutobotClusterID, "autobot", "cluster id")
	envflag.StringVar(&env.ClusterURL, AutobotClusterURL, "nats://localhost:4222", "cluster url")
	envflag.StringVar(&env.ClusterDiscovery, AutobotClusterDiscovery, "autobot.discovery", "cluster url")

	flag.StringVar(&env.Name, "name", "", "name")
	flag.StringVar(&env.ClusterID, "cluster-id", "", "cluster id")
	flag.StringVar(&env.ClusterURL, "cluster-name", "", "cluster name")
	flag.StringVar(&env.ClusterDiscovery, "cluster-discovery", "", "discovery topic")

	envflag.Parse()
	flag.Parse()

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
