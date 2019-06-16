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

	// addr to listen on
	cmd.Flags().StringVar(&cfg.Addr, "addr", c.DefaultAddr, "address to listen")
}
