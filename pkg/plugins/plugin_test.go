package plugins

import (
	"context"
	"testing"
	"time"

	"github.com/andersnormal/autobot/pkg/plugins/message"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	at "github.com/andersnormal/autobot/pkg/testing"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
)

func TestPublishInbox(t *testing.T) {
	env := &runtime.Environment{
		ClusterID:  runtime.DefaultClusterID,
		ClusterURL: runtime.DefaultClusterURL,
		Inbox:      runtime.DefaultClusterInbox,
		Outbox:     runtime.DefaultClusterOutbox,
		LogFormat:  runtime.DefaultLogFormat,
		LogLevel:   runtime.DefaultLogLevel,
	}

	env.Name = "skrimish"

	at.WithAutobot(t, env, func(t *testing.T) {
		assert := assert.New(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// create test plugin ....
		plugin, _ := WithContext(ctx, env)

		err := plugin.PublishInbox(&pb.Message{Text: "foo"})
		assert.NoError(err)

		sc, err := plugin.getConn()
		assert.NoError(err)

		sub := make(chan Event)

		s, err := sc.QueueSubscribe(env.Inbox, env.Name, func(m *stan.Msg) {
			// this is recreating the messsage from the inbox
			msg, err := message.FromByte(m.Data)
			assert.NoError(err)

			botMessage := new(pb.Message)
			err = plugin.marshaler.Unmarshal(msg, botMessage)
			assert.NoError(err)

			sub <- botMessage
		}, stan.DurableName(env.Name), stan.StartWithLastReceived())
		assert.NoError(err)

		defer func() { s.Close() }()

		select {
		case e, ok := <-sub:
			assert.True(ok)

			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "foo")
		case <-time.After(5 * time.Second):
			assert.Fail("did not received message")
		}
	})
}

func TestPublishOutbox(t *testing.T) {
	env := &runtime.Environment{
		ClusterID:  runtime.DefaultClusterID,
		ClusterURL: runtime.DefaultClusterURL,
		Inbox:      runtime.DefaultClusterInbox,
		Outbox:     runtime.DefaultClusterOutbox,
		LogFormat:  runtime.DefaultLogFormat,
		LogLevel:   runtime.DefaultLogLevel,
	}

	env.Name = "skrimish"

	at.WithAutobot(t, env, func(t *testing.T) {
		assert := assert.New(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// create test plugin ....
		plugin, _ := WithContext(ctx, env)

		err := plugin.PublishOutbox(&pb.Message{Text: "foo"})
		assert.NoError(err)

		sc, err := plugin.getConn()
		assert.NoError(err)

		sub := make(chan Event)

		s, err := sc.QueueSubscribe(env.Outbox, env.Name, func(m *stan.Msg) {
			// this is recreating the messsage from the inbox
			msg, err := message.FromByte(m.Data)
			assert.NoError(err)

			botMessage := new(pb.Message)
			err = plugin.marshaler.Unmarshal(msg, botMessage)
			assert.NoError(err)

			sub <- botMessage
		}, stan.DurableName(env.Name), stan.StartWithLastReceived())
		assert.NoError(err)

		defer func() { s.Close() }()

		select {
		case e, ok := <-sub:
			assert.True(ok)

			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "foo")
		case <-time.After(5 * time.Second):
			assert.Fail("did not received message")
		}
	})
}
