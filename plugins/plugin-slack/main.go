package main

import (
	"context"
	"os"
	"path"
	"strconv"

	"github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/log"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	ll "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// this is provided via the env variable from the server
	slackToken = "SLACK_TOKEN"
)

var slackRuntime = &runtime.Runtime{
	RunE: runE,
}

func init() {
	runtime.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("autobot")
	viper.AutomaticEnv()

	// set some default flags
	pflag.String("name", path.Base(os.Args[0]), "plugin name")
	pflag.String("log_format", runtime.DefaultLogFormat, "log format")
	pflag.String("log_level", runtime.DefaultLogLevel, "log level")
	pflag.BoolP("verbose", "v", true, "verbose")
	pflag.BoolP("debug", "d", true, "debug")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// unmarshal to config
	if err := viper.Unmarshal(runtime.Env); err != nil {
		ll.Fatalf(errors.Wrap(err, "cannot unmarshal runtime").Error())
	}
}

func runE(env *runtime.Environment) error {
	// have a root context ...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin, ctx := plugins.WithContext(ctx, env)

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
		case <-ctx.Done():
			plugin.Log().Fatal(ctx.Err())
		}
	}

	return nil
}

func main() {
	if err := slackRuntime.Execute(); err != nil {
		panic(err)
	}
}

// FromUserIDWithContext ...
func FromUserIDWithContext(ctx context.Context, api *slack.Client, userID string) (*pb.Message_User, error) {
	u := new(pb.Message_User)

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
func FromChannelIDWithContext(ctx context.Context, api *slack.Client, channelID string) (*pb.Message_Channel, error) {
	c := new(pb.Message_Channel)

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
func FromChannelID(channelID string) *pb.Message_Channel {
	return &pb.Message_Channel{Id: channelID}
}

// FromMsgWithContext ...
func FromMsgWithContext(ctx context.Context, api *slack.Client, msg *slack.MessageEvent) (*pb.Message, error) {
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

	return m, nil
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

// FromOutboxEvent ...
func FromOutboxEvent(rtm *slack.RTM, e *pb.Message) *slack.OutgoingMessage {
	return rtm.NewOutgoingMessage(e.GetText(), e.GetChannel().GetId())
}
