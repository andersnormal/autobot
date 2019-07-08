package main

import (
	"log"
  "os"
  "context"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"
)

func main() {
  name := os.Args[0]
  
  // have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, ctx, err := plugins.WithContext(ctx, pb.NewPlugin(name))
	if err != nil {
		log.Fatalf("could not create plugin: %v", err)
	}

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
		if in.GetMessage() != nil {
			return &pb.Event{
				Event: &pb.Event_Reply{
					Reply: &pb.Message{
						Text:    "hello world",
						Channel: in.GetMessage().GetChannel(),
					},
				},
			}, nil
		}

		return nil, nil
	}
}
