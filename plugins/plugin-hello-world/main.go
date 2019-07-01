package main

import (
	"log"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"
)

func main() {
	// create root context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// plugin ....
	plugin, err := plugins.New("hello-world")
	if err != nil {
		log.Fatalf("could not create plugin: %v", err)
	}

	// create publish channel ...
	pubReply := plugin.PublishReplies()
	subMsg := plugin.SubscribeMessages()

	// process messages ...
	go func() {
		// subscribe ...
		for {
			select {
			case e, ok := <-subMsg:
				if !ok {
					return
				}

				if e.GetMessage() != nil {
					log.Printf("got event: %v", e.GetMessage())
					go func() {
						reply := &pb.Event{
							Event: &pb.Event_Reply{
								Reply: &pb.MessageEvent{
									Text:     "hello world",
									Channel:  e.GetMessage().GetChannel(),
									User:     e.GetMessage().GetUser(),
									Username: e.GetMessage().GetUsername(),
									Topic:    e.GetMessage().GetTopic(),
								},
							},
						}

						pubReply <- reply
					}()
				}
			}
		}
	}()

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}
