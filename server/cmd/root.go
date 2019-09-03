package cmd

import (
	"fmt"
	"os"

	"github.com/andersnormal/autobot/server/config"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultEnvPrefix = "AUTOBOT"
)

var (
	cfg *config.Config
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "autobot",
	Short: "",
	Long:  `Not yet`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: runE,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// init config
	cfg = config.New()

	// initialize cobra
	cobra.OnInitialize(initConfig)

	// adding flags
	addFlags(RootCmd, cfg)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	if cfg.File != "" {
		viper.SetConfigFile(cfg.File)

		// do not forget to read in the config
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf(errors.Wrap(err, "cannot read config").Error())
		}
	}

	// unmarshal to config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf(errors.Wrap(err, "cannot unmarshal config").Error())
	}

	// set the default format, which is basically text
	log.SetFormatter(&log.TextFormatter{})

	// reset log format
	if cfg.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if cfg.Verbose {
		cfg.LogLevel = "info"
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	// set the configured log level
	if level, err := log.ParseLevel(cfg.LogLevel); err == nil {
		log.SetLevel(level)
	}
}
