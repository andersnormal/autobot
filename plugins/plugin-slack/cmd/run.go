package cmd

import (
	"context"
	"os"

	"github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/log"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nlopes/slack"
)

const (
	// this is provided via the env variable from the server
	slackToken = "SLACK_TOKEN"
)

// Slack ...
var Slack = &runtime.Runtime{
	RunE: runE,
}

func runE(env *runtime.Environment) error {
	// have a root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, rctx := plugins.WithContext(ctx, env)

	// create publish channel ...
	subReply := plugin.SubscribeOutbox()

	// log ..
	plugin.Log().Infof("starting slack plugin ...")

	// subscribe ...
	// collect options ...
	opts := []slack.Option{
		slack.OptionDebug(env.Debug),
	}

	// enable verbosity ...
	if env.Verbose {
		opts = append(opts, slack.OptionLog(&log.InternalLogger{Logger: plugin.Log()}))
	}

	// create client ...
	api := slack.New(
		os.Getenv(slackToken),
		opts...,
	)

	// create connection ...
	rtm := api.NewRTM()
	go rtm.ManageConnection()

OUTER:
	for {
		select {
		case e, ok := <-rtm.IncomingEvents:
			// channel is closed ... should be an error?
			if !ok {
				break OUTER
			}

			switch ev := e.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				// ignore connect event

			case *slack.MessageEvent:
				msg, err := FromMsgWithContext(ctx, api, ev)
				if err != nil {
					continue
				}

				if err := plugin.PublishInbox(msg); err != nil {
					plugin.Log().Errorf("could not publish message: %v", err)
				}

			case *slack.PresenceChangeEvent:
				// ignore for now

			case *slack.LatencyReport:
				// ignore for now

			case *slack.RTMError:
				// ignore for now
				// log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				plugin.Log().Fatal(plugins.ErrPluginAuthentication)

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		case e, ok := <-subReply:
			if !ok {
				// should there be a different error?
				break OUTER
			}

			switch ev := e.(type) {
			case *pb.Message:
				rtm.SendMessage(FromOutboxEvent(rtm, ev))
			}
		case <-rctx.Done():
			plugin.Log().Fatal(rctx.Err())
		}
	}

	return nil
}
