package main

import (
	"context"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
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
	plugin, ctx := plugins.WithContext(ctx, plugins.Name(name), plugins.Debug())

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

// FromUserIDWithContext ...
func FromUserIDWithContext(ctx context.Context, api *slack.Client, userID string) (*pb.User, error) {
	u := new(pb.User)

	// get user
	user, err := api.GetUserInfoContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	u.Id = user.ID
	u.Name = user.Name

	return u, nil
}

// FromChannelWithContext ...
func FromChannelIDWithContext(ctx context.Context, api *slack.Client, channelID string) (*pb.Channel, error) {
	c := new(pb.Channel)

	// get channel
	channel, err := api.GetChannelInfoContext(ctx, channelID)
	if err != nil {
		return nil, err
	}

	c.Id = channel.ID
	c.Name = channel.Name

	return c, nil
}

// FromChannelID ...
func FromChannelID(channelID string) *pb.Channel {
	return &pb.Channel{Id: channelID}
}

// FromMsgWithContext ...
func FromMsgWithContext(ctx context.Context, api *slack.Client, msg *slack.MessageEvent) (*pb.Event, error) {
	m := new(pb.Message)

	// basic ...
	m.Type = msg.Type

	// resolve user
	user, err := FromUserIDWithContext(ctx, api, msg.User)
	if err != nil {
		return nil, err
	}

	m.From = user

	// get list of direct message channels
	im, err := api.GetIMChannelsContext(ctx)
	if err != nil {
		return nil, err
	}

	var isPrivate bool
	for _, dm := range im {
		isPrivate = dm.ID == msg.Channel

		if isPrivate {
			m.Channel = FromChannelID(dm.ID)

			break
		}
	}

	if !isPrivate {
		// resolve channel ...
		channel, err := FromChannelIDWithContext(ctx, api, msg.Channel)
		if err != nil {
			return nil, err
		}

		m.Channel = channel
	}

	m.TextFormat = pb.Message_PLAIN_TEXT
	m.Text = msg.Text

	// timestamp ...
	ts, err := FromTimestamp(msg.Timestamp)
	if err != nil {
		return nil, err
	}
	m.Timestamp = ts

	// standard event ...
	e := &pb.Event{Event: &pb.Event_Message{
		Message: m,
	}}

	// if it is private
	if isPrivate {
		e.Event = &pb.Event_Private{
			Private: m,
		}
	}

	return e, nil
}

// FromTimestamp ...
func FromTimestamp(ts string) (*timestamp.Timestamp, error) {
	f, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return nil, err
	}

	secs := int64(f)
	nsecs := int32((f - float64(secs)) * 1e9)

	return &timestamp.Timestamp{
		Seconds: secs,
		Nanos:   nsecs,
	}, nil
}

// FromReplyEvent ...
func FromReplyEvent(rtm *slack.RTM, e *pb.Message) *slack.OutgoingMessage {
	return rtm.NewOutgoingMessage(e.GetText(), e.GetChannel().GetId())
}
