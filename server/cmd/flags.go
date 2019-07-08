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

	// nats prefix ...
	cmd.Flags().StringVar(&cfg.NatsPrefix, "nats-prefix", c.DefaultNatsPrefix, "NATS channel prefix")

	// nats url ...
	cmd.Flags().StringVar(&cfg.NatsClusterURL, "nats-cluster-url", c.DefaultNatsClusterURL, "NATS cluster url")

	// nats url ...
	cmd.Flags().StringVar(&cfg.NatsClusterID, "nats-cluster-id", c.DefaultNatsClusterID, "NATS cluster id")
}
