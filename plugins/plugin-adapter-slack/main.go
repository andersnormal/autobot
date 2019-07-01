package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andersnormal/autobot/pkg/plugins"

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
	plugin := plugins.New("slack-adapter")

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
	pub := plugin.PublishMessages()

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
						pub <- FromMsg(ev)
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
			case e, ok := <-plugin.SubscribeMessages():
				if !ok {
					return
				}

				log.Printf("received event: %v", e)
			}
		}
	}()

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}
