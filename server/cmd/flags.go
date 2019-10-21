package cmd

import (
	c "github.com/andersnormal/autobot/pkg/config"

	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command, cfg *c.Config) {
	// set the config file
	cmd.Flags().StringVar(&cfg.File, "config", "", "config file (default is $HOME/.autobot.yaml")
	// enable verbose output
	cmd.Flags().BoolVar(&cfg.Verbose, "verbose", c.DefaultVerbose, "enable verbose output")
	// enable debug options
	cmd.Flags().BoolVar(&cfg.Debug, "debug", c.DefaultDebug, "enable debug")
	// set log format
	cmd.Flags().StringVar(&cfg.LogFormat, "log-format", c.DefaultLogFormat, "log format (default is 'text')")
	// set log level
	cmd.Flags().StringVar(&cfg.LogLevel, "log-level", c.DefaultLogLevel, "log level (default is 'warn'")
	// clustering ...
	cmd.Flags().BoolVar(&cfg.Nats.Clustering, "clustering", cfg.Nats.Clustering, "enable clustering")
	// clustering bootstrap ...
	cmd.Flags().BoolVar(&cfg.Nats.Bootstrap, "bootstrap", cfg.Nats.Bootstrap, "bootstrap cluster")
	// clustering node id ...
	cmd.Flags().StringVar(&cfg.Nats.ClusterNodeID, "node-id", cfg.Nats.ClusterNodeID, "node id")
	// clustering peers ...
	cmd.Flags().StringSliceVar(&cfg.Nats.ClusterPeers, "peers", cfg.Nats.ClusterPeers, "peers")
	// clustering peers ...
	cmd.Flags().StringVar(&cfg.Nats.ClusterURL, "url", cfg.Nats.ClusterURL, "cluster url")
}
