package cmd

import (
	c "github.com/andersnormal/autobot/server/config"

	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command, cfg *c.Config) {
	// enable verbose output
	cmd.Flags().BoolVar(&cfg.Verbose, "verbose", c.DefaultVerbose, "enable verbose output")

	// enable debug options
	cmd.Flags().BoolVar(&cfg.Debug, "debug", c.DefaultDebug, "enable debug")

	// bot name
	cmd.Flags().StringVar(&cfg.BotName, "bot-name", c.DefaultBotName, "bot name")

	// addr to listen on
	cmd.Flags().StringVar(&cfg.Addr, "addr", c.DefaultAddr, "address to listen")

	// plugins ...
	cmd.Flags().StringSliceVar(&cfg.PluginsDirs, "plugins", c.DefaultPluginsDirs, "plugins directory")

	// env ...
	cmd.Flags().StringSliceVar(&cfg.PluginEnv, "env", []string{}, "env")

	// nats disabled ...
	cmd.Flags().BoolVar(&cfg.Nats.Disabled, "nats-disabled", false, "disable embedded NATS")

	// nats clustering ...
	cmd.Flags().BoolVar(&cfg.Nats.Clustering, "nats-clustering", false, "NATS clustering mode")

	// nats bootstrp ...
	cmd.Flags().BoolVar(&cfg.Nats.Bootstrap, "nats-bootstrap", false, "NATS bootstrap leader")

	// nats cluster peers ...
	cmd.Flags().StringSliceVar(&cfg.Nats.ClusterPeers, "nats-cluster-peers", []string{}, "NATS cluster peers")

	// nats cluster node id ...
	cmd.Flags().StringVar(&cfg.Nats.ClusterNodeID, "nats-cluster-node-id", "optimus_prime", "NATS cluster node id")

	// nats prefix ...
	cmd.Flags().StringVar(&cfg.Nats.Prefix, "nats-prefix", c.DefaultNatsPrefix, "NATS channel prefix")

	// nats url ...
	cmd.Flags().StringVar(&cfg.Nats.ClusterURL, "nats-cluster-url", c.DefaultNatsClusterURL, "NATS cluster url")

	// nats url ...
	cmd.Flags().StringVar(&cfg.Nats.ClusterID, "nats-cluster-id", c.DefaultNatsClusterID, "NATS cluster id")

	// grpc addr ...
	cmd.Flags().StringVar(&cfg.GRPCAddr, "grpc-addr", c.DefaultGRPCAddr, "grpc listen address")
}
