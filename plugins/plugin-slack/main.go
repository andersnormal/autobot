package main

import (
	"fmt"
	"os"
	"path"

	"github.com/andersnormal/autobot/plugins/plugin-slack/cmd"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"

	"github.com/pkg/errors"
	ll "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	runtime.OnInitialize(initConfig)
}

func initConfig() {
	// set some default flags
	pflag.String("name", path.Base(os.Args[0]), "plugin name")
	pflag.String("log_format", runtime.DefaultLogFormat, "log format")
	pflag.String("log_level", runtime.DefaultLogLevel, "log level")
	pflag.BoolP("verbose", "v", false, "verbose")
	pflag.BoolP("debug", "d", false, "debug")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// unmarshal to config
	if err := viper.Unmarshal(runtime.Env()); err != nil {
		ll.Fatalf(errors.Wrap(err, "cannot unmarshal runtime").Error())
	}

	fmt.Println(runtime.Env())
}

func main() {
	if err := cmd.Slack.Execute(); err != nil {
		panic(err)
	}
}
