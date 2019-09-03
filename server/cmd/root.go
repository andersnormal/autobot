package cmd

import (
	"fmt"
	"os"

	"github.com/andersnormal/autobot/server/config"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	viper.AutomaticEnv() // read in environment variables that match

	if cfg.File != "" {
		viper.SetConfigFile(cfg.File)
	} else {
		// find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(errors.Wrap(err, "cannot find home dir").Error())
		}

		// enforcing the type here if nothing is found
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".autobot.yaml")
	}

	// do not forget to read in the config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf(errors.Wrap(err, "cannot read config").Error())
	}

	// unmarshal to config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf(errors.Wrap(err, "cannot unmarshal config").Error())
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Only log the warning severity or above.
	log.SetLevel(cfg.LogLevel)

	// if we should output verbose
	if cfg.Verbose || cfg.Debug {
		log.SetLevel(log.InfoLevel)
	}
}
