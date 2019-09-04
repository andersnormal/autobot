package cmd

import (
	c "github.com/andersnormal/autobot/pkg/config"

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

	// set log format
	cmd.Flags().StringVar(&cfg.LogFormat, "log-format", c.DefaultLogFormat, "log format (default is 'text')")

	// set log level
	cmd.Flags().StringVar(&cfg.LogLevel, "log-level", c.DefaultLogLevel, "log level (default is 'warn'")

	// addr to listen on
	cmd.Flags().StringVar(&cfg.Addr, "addr", c.DefaultAddr, "address to listen")

	// plugins ...
	cmd.Flags().StringSliceVar(&cfg.Plugins, "plugins", c.DefaultPlugins, "plugins directory (default is 'plugins')")

	// configs to bind
	viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
	viper.BindPFlag("log_format", cmd.Flags().Lookup("log-format"))
	viper.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))
	viper.BindPFlag("plugins", cmd.Flags().Lookup("plugins"))
}
