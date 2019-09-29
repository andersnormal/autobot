package main

import (
	"context"
	"os"
	"path"

	"github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var helloRuntime = &runtime.Runtime{
	RunE: runE,
}

func init() {
	runtime.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("autobot")
	viper.AutomaticEnv()

	// set some default flags
	pflag.String("name", path.Base(os.Args[0]), "plugin name")
	pflag.String("log_format", runtime.DefaultLogFormat, "log format")
	pflag.String("log_level", runtime.DefaultLogLevel, "log level")
	pflag.BoolP("verbose", "v", true, "verbose")
	pflag.BoolP("debug", "d", true, "debug")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// unmarshal to config
	if err := viper.Unmarshal(runtime.Env); err != nil {
		log.Fatalf(errors.Wrap(err, "cannot unmarshal runtime").Error())
	}
}

func runE(env *runtime.Environment) error {
	// have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, ctx := plugins.WithContext(ctx, env)

	// log ..
	plugin.Log().Infof("starting hello world plugin ...")

	// Processing incoming messages ...
	msgFunc := func(ctx plugins.Context) error {
		ctx.Send(ctx.Message().Reply("hello world"))

		return nil
	}

	// use the schedule function from the plugin
	if err := plugin.ReplyWithFunc(msgFunc); err != nil {
		return err
	}

	if err := plugin.Wait(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := helloRuntime.Execute(); err != nil {
		log.Fatal(err)
	}
}
