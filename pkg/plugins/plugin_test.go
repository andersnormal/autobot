package plugins

import (
	"context"
	"testing"
	"time"

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
			botMessage := new(pb.Message)

			err := plugin.marshaler.Unmarshal(m.Data, botMessage)
			assert.NoError(err)

			sub <- botMessage
		}, stan.DurableName(env.Name), stan.DeliverAllAvailable())
		assert.NoError(err)

		defer s.Close()

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

func TestSubscribeInbox(t *testing.T) {
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

		select {
		case e, ok := <-plugin.SubscribeInbox():
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
			botMessage := new(pb.Message)

			err := plugin.marshaler.Unmarshal(m.Data, botMessage)
			assert.NoError(err)

			sub <- botMessage
		}, stan.DurableName(env.Name), stan.DeliverAllAvailable())
		assert.NoError(err)

		defer s.Close()

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

func TestSubscribeOutbox(t *testing.T) {
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

		select {
		case e, ok := <-plugin.SubscribeOutbox():
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "foo")
		case <-time.After(5 * time.Second):
			assert.Fail("did not received message")
		}
	})
}

func TestReplyFunc(t *testing.T) {
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

		plugin.ReplyWithFunc(func(c Context) error {
			msg := c.Message()
			replyMessage := msg.Reply("echo: " + msg.GetText())
			return plugin.PublishOutbox(replyMessage)
		})

		err := plugin.PublishInbox(&pb.Message{Text: "foo"})
		assert.NoError(err)

		select {
		case e, ok := <-plugin.SubscribeOutbox():
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "echo: foo")
		case <-time.After(5 * time.Second):
			assert.Fail("did not received message")
		}
	})
}

func TestReplyFunc_FIFO(t *testing.T) {
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

		plugin.ReplyWithFunc(func(c Context) error {
			if c.Message().GetText() == "delayed message" {
				time.Sleep(300 * time.Millisecond)
				return plugin.PublishOutbox(c.Message())
			}

			return plugin.PublishOutbox(c.Message())
		})

		err := plugin.PublishInbox(&pb.Message{Text: "delayed message"})
		assert.NoError(err)

		err = plugin.PublishInbox(&pb.Message{Text: "will still arrive last"})
		assert.NoError(err)

		outbox := plugin.SubscribeOutbox()

		select {
		case e, ok := <-outbox:
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "delayed message")
		case <-time.After(5 * time.Second):
			assert.Fail("did not receive the first message")
		}

		select {
		case e, ok := <-outbox:
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "will still arrive last")
		case <-time.After(5 * time.Second):
			assert.Fail("did not received the last message")
		}
	})
}

func TestAsyncReplyFunc(t *testing.T) {
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

		plugin.AsyncReplyWithFunc(func(c Context) error {
			msg := c.Message()
			replyMessage := msg.Reply("echo: " + msg.GetText())
			return plugin.PublishOutbox(replyMessage)
		})

		err := plugin.PublishInbox(&pb.Message{Text: "foo"})
		assert.NoError(err)

		select {
		case e, ok := <-plugin.SubscribeOutbox():
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "echo: foo")
		case <-time.After(5 * time.Second):
			assert.Fail("did not received message")
		}
	})
}

func TestAsyncReplyWithFunc_HandlesMessagesConcurrently(t *testing.T) {
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

		plugin.AsyncReplyWithFunc(func(c Context) error {
			if c.Message().GetText() == "delayed message" {
				time.Sleep(300 * time.Millisecond)
				return plugin.PublishOutbox(c.Message())
			}

			return plugin.PublishOutbox(c.Message())
		})

		err := plugin.PublishInbox(&pb.Message{Text: "delayed message"})
		assert.NoError(err)

		err = plugin.PublishInbox(&pb.Message{Text: "sent last, finishes first"})
		assert.NoError(err)

		outbox := plugin.SubscribeOutbox()

		select {
		case e, ok := <-outbox:
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "sent last, finishes first")
		case <-time.After(5 * time.Second):
			assert.Fail("did not receive the first message")
		}

		select {
		case e, ok := <-outbox:
			assert.True(ok)
			assert.IsType(e, &pb.Message{})
			assert.Equal(e.(*pb.Message).GetText(), "delayed message")
		case <-time.After(5 * time.Second):
			assert.Fail("did not receive the second message")
		}
	})
}
