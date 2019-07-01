package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nlopes/slack"
)

const (
	slackToken = "SLACK_TOKEN"
)

type slackPlugin struct {
}

func main() {
	// create root context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// plugin ....
	plugin, err := plugins.New("slack-adapter")
	if err != nil {
		log.Fatalf("could not create plugin: %v", err)
	}

	// create client ...
	api := slack.New(
		os.Getenv(slackToken),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	// create connection ...
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// create publish channel ...
	pubMsg := plugin.PublishMessages()
	pubReply := plugin.PublishReplies()
	subMsg := plugin.SubscribeMessages()
	subReply := plugin.SubscribeReplies()

	// process messages ...
	go func() {
		// subscribe ...
		for {
			select {
			case e, ok := <-rtm.IncomingEvents:
				log.Printf("Event received: ")

				if !ok {
					return
				}

				switch ev := e.Data.(type) {
				case *slack.HelloEvent:
					// Ignore hello

				case *slack.ConnectedEvent:
					fmt.Println("Infos:", ev.Info)
					fmt.Println("Connection counter:", ev.ConnectionCount)
					// Replace C2147483705 with your Channel ID
					rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

				case *slack.MessageEvent:
					fmt.Printf("Message: %v\n", ev)
					go func() {
						pubMsg <- FromMsg(ev)
					}()

				case *slack.PresenceChangeEvent:
					fmt.Printf("Presence Change: %v\n", ev)

				case *slack.LatencyReport:
					fmt.Printf("Current latency: %v\n", ev.Value)

				case *slack.RTMError:
					fmt.Printf("Error: %s\n", ev.Error())

				case *slack.InvalidAuthEvent:
					fmt.Printf("Invalid credentials")
					return

				default:

					// Ignore other events..
					// fmt.Printf("Unexpected: %v\n", msg.Data)
				}
			case e, ok := <-subReply:
				if !ok {
					return
				}

				log.Printf("got reply: %v", e)

				if e.GetReply() != nil {
					go func() {
						rtm.SendMessage(FromMessageEvent(rtm, e.GetReply()))
					}()
				}

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
