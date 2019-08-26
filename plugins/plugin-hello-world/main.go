package main

import (
	"context"
	"log"

	"github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"
)

func main() {
	// have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create env ...
	env := runtime.DefaultEnv()

	// plugin ....
	plugin, ctx := plugins.WithContext(ctx, env)

	// use the schedule function from the plugin
	if err := plugin.ReplyWithFunc(msgFunc()); err != nil {
		log.Fatalf("could not create plugin: %v", err)
	}

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}

func msgFunc() plugins.SubscribeFunc {
	return func(in *pb.Event) (*pb.Event, error) {
		if in.GetPrivate() != nil {
			return pb.NewReply(in.GetPrivate().Reply("hello world")), nil
		}

		return nil, nil
	}
}
