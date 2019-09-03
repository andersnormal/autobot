package cmd

import (
	c "github.com/andersnormal/autobot/server/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addFlags(cmd *cobra.Command, cfg *c.Config) {
	// set the config file
	cmd.Flags().StringVar(&cfg.File, "config", "", "config file (default is $HOME/.autobot.yaml")

	// enable verbose output
	cmd.Flags().BoolVar(&cfg.Verbose, "verbose", c.DefaultVerbose, "enable verbose output")

	// enable debug options
	cmd.Flags().BoolVar(&cfg.Debug, "debug", c.DefaultDebug, "enable debug")

	// addr to listen on
	cmd.Flags().StringVar(&cfg.Addr, "addr", c.DefaultAddr, "address to listen")

	// plugins ...
	cmd.Flags().StringSliceVar(&cfg.PluginsDirs, "plugins", c.DefaultPluginsDirs, "plugins directory")

	// set the link between flags
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
	viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))
}
