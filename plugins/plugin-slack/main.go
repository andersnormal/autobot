package main

import (
	"context"
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

	// have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, ctx, err := plugins.WithContext(ctx, pb.NewPlugin(name))
	log.Fatalf("could not create plugin: %v", err)

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
					msg, err := FromMsgWithContext(ctx, api, ev)
					if err != nil {
						log.Printf("could not parse message from: %v", ev.User)

						continue
					}

					go func() {
						pubMsg <- msg
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
			case <-ctx.Done():
				return nil
			}
		}
	})

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}
