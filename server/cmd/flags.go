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
	cmd.Flags().StringVar(&cfg.PluginsDir, "plugins", c.DefaultPluginsDir, "plugins directory")

	// env ...
	cmd.Flags().StringSliceVar(&cfg.PluginEnv, "env", []string{}, "env")

	// replies topic ...
	cmd.Flags().StringVar(&cfg.NatsRepliesTopic, "replies-topic", c.DefaultNatsRepliesTopic, "replies topic name")

	// message topic ...
	cmd.Flags().StringVar(&cfg.NatsMessagesTopic, "message-topic", c.DefaultNatsMessagesTopic, "message topic name")
}
