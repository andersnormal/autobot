package main

import (
	"context"
	"log"

	"github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
)

func main() {
	// have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create env ...
	env := runtime.Default()

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
		log.Fatalf("could not create plugin: %v", err)
	}

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}
