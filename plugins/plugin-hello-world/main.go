package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"
)

func main() {
	name := path.Base(os.Args[0])

	// have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, ctx := plugins.WithContext(ctx, plugins.Name(name), plugins.Debug())

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
			return &pb.Event{
				Plugin: in.GetPlugin(),
				Event: &pb.Event_Reply{
					Reply: &pb.Message{
						Text:    "hello world",
						Channel: in.GetPrivate().GetChannel(),
					},
				},
			}, nil
		}

		return nil, nil
	}
}
