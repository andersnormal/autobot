package main

import (
	"context"
	"log"
  "os"
  "path"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nlopes/slack"
)

const (
	// this is provided via the env variable from the server
	slackToken = "SLACK_TOKEN"
)

func main() {
	// extract the plugin name
	name := path.Base(os.Args[0])

	// have root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, ctx := plugins.WithContext(ctx, plugins.Name(name))

	// create publish channel ...
	pubMsg := plugin.PublishInbox()
	subReply := plugin.SubscribeOutbox(
		plugins.WithFilterPlugin(),
	)

	events := plugin.Events()

	// Run in plugin loop ...
	plugin.Run(func() error {
		// subscribe ...

	OUTER:
		for {
			restart := make(chan struct{})

			// collect options ...
			opts := []slack.Option{
				slack.OptionDebug(plugin.Debug()),
			}

			// enable verbosity ...
			if plugin.Verbose() {
				opts = append(opts, slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))
			}

			// create client ...
			api := slack.New(
				os.Getenv(slackToken),
				opts...,
			)

			// create connection ...
			rtm := api.NewRTM()
			plugin.Run(func() error {
				rtm.ManageConnection()

				<-restart

				return nil
			})

			for {
				select {
				case event := <-events:
					if event == plugins.RefreshConfigEvent {
						// reload the plugin
						continue OUTER
					}
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
							log.Printf("could not parse message from: %v", err)

							continue
						}

						// run publish in go routine
						plugin.Run(publishEvent(pubMsg, msg))

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
						// reply in go routine
						plugin.Run(sendReply(rtm, e.GetReply()))
					}
				case <-ctx.Done():
					return nil
				}
			}
		}
	})

	// wait for the routines to either finish,
	// while the process is killed via the context
	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}

func publishEvent(pubMsg chan<- *pb.Event, msg *pb.Event) func() error {
	return func() error {
		pubMsg <- msg

		return nil
	}
}

func sendReply(rtm *slack.RTM, msg *pb.Message) func() error {
	return func() error {
		rtm.SendMessage(FromReplyEvent(rtm, msg))

		return nil
	}
}
