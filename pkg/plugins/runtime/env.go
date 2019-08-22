package runtime

import (
	"flag"
	"os"
	"path"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/ianschenck/envflag"
)

const (
	// AutobotName ...
	AutobotName = "AUTOBOT_NAME"
	// AutobotClusterID ...
	AutobotClusterID = "AUTOBOT_CLUSTER_ID"
	// AutobotClusterURL ...
	AutobotClusterURL = "AUTOBOT_CLUSTER_URL"
	// AutobotClusterDiscovery ...
	AutobotClusterDiscovery = "AUTOBOT_CHANNEL_DISCOVERY"
)

// Env ...
type Env struct {
	Name             string
	ClusterID        string
	ClusterURL       string
	ClusterDiscovery string
	Debug            bool
	Inbox            string
	Outbox           string
}

type opts struct {
}

type defaultEnv struct {
	opts opts
}

// DefaultEnv ...
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

// WithConfig ...
func (e Env) WithConfig(cfg *pb.Config) Env {
	e.Inbox = cfg.GetInbox()
	e.Outbox = cfg.GetOutbox()
	e.Debug = cfg.GetDebug()

	return e
}

// WithInbox ...
func (e Env) WithInbox(inbox string) Env {
	e.Inbox = inbox

	return e
}

// WithOutbox ...
func (e Env) WithOutbox(outbox string) Env {
	e.Outbox = outbox

	return e
}
