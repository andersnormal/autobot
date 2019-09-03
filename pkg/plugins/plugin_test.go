package plugins

import (
	"context"
	"testing"
	"time"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const (
	waitTimeout = 3 * time.Second
)

var env = runtime.DefaultEnv().WithName("autobot").WithClusterURL("nats://localhost:4222")

func TestInboxPub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	plugin, ctx := WithContext(ctx, env)

	inboxMessageCh := make(chan string, 1)

	plugin.ReplyWithFunc(func(msg *pb.Message) (*pb.Message, error) {
		inboxMessageCh <- msg.GetText()
		return nil, nil
	})

	plugin.PublishInbox() <- &pb.Message{
		Text: "hello world",
	}

	select {
	case msg := <-inboxMessageCh:
		assert.Equal(t, "hello world", msg)
	case <-time.After(waitTimeout):
		assert.FailNow(t, "timed out waiting for message...")
	}
}

func TestOutboxPub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	plugin, ctx := WithContext(ctx, env)

	outboxMessageCh := make(chan string, 1)

	var g errgroup.Group
	g.Go(func() error {
		e := <-plugin.SubscribeOutbox()

		switch ev := e.(type) {
		case *pb.Message:
			outboxMessageCh <- ev.GetText()
		}

		return nil
	})

	plugin.PublishOutbox() <- &pb.Message{
		Text: "hello world",
	}

	g.Wait()

	select {
	case msg := <-outboxMessageCh:
		assert.Equal(t, "hello world", msg)
	case <-time.After(waitTimeout):
		assert.FailNow(t, "timed out waiting for message...")
	}
}

func TestReplyFunc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	plugin, ctx := WithContext(ctx, env)

	outboxMessageCh := make(chan string, 1)

	var g errgroup.Group
	g.Go(func() error {
		e := <-plugin.SubscribeOutbox()

		switch ev := e.(type) {
		case *pb.Message:
			outboxMessageCh <- ev.GetText()
		}

		return nil
	})

	plugin.ReplyWithFunc(func(msg *pb.Message) (*pb.Message, error) {
		return msg.Reply("reply message"), nil
	})

	plugin.PublishInbox() <- &pb.Message{
		Text: "hello world",
	}

	g.Wait()

	select {
	case msg := <-outboxMessageCh:
		assert.Equal(t, "reply message", msg)
	case <-time.After(waitTimeout):
		assert.FailNow(t, "timed out waiting for message...")
	}
}
