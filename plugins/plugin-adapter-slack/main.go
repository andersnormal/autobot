package main

import (
	"log"
	"os"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nlopes/slack"
)

const (
	slackToken = "SLACK_TOKEN"
)


func main() {
	name := os.Args[0]

	// plugin ....
	plugin, err := plugins.New(pb.NewPlugin(name))
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
	pubMsg := plugin.PublishInbox()
	subReply := plugin.SubscribeOutbox()

	// Run in plugin loop ...
	plugin.Run(func() error {
		// subscribe ...
		for {
			select {
			case e, ok := <-rtm.IncomingEvents:
				// channel is closed ... should be an error?
				if !ok {
					return nil
				}

				switch ev := e.Data.(type) {
				case *slack.HelloEvent:
					// Ignore hello

				case *slack.ConnectedEvent:
					// ignore connect event

				case *slack.MessageEvent:
					// should this go to the main loop ?
					go func() {
						pubMsg <- FromMsg(ev)
					}()

				case *slack.PresenceChangeEvent:
					// ignore for now

				case *slack.LatencyReport:
					// ignore for now

				case *slack.RTMError:
					// ignore for now
					// log.Printf("Error: %s\n", ev.Error())

				case *slack.InvalidAuthEvent:
					return plugins.ErrPluginAuthentication

				default:

					// Ignore other events..
					// fmt.Printf("Unexpected: %v\n", msg.Data)
				}
			case e, ok := <-subReply:
				if !ok {
					// should there be a different error?
					return nil
				}

				if e.GetReply() != nil {
					go func() {
						rtm.SendMessage(FromMessageEvent(rtm, e.GetReply()))
					}()
				}
			}
		}
	})

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}
