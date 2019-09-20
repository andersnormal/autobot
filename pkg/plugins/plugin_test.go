package plugins

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const waitTimeout = 5 * time.Second

func TestOutbox(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverCfg := withTestConfig()

	withTestAutobot(ctx, serverCfg, func() {
		// create test plugin ....
		plugin := newTestPlugin(ctx, "outbox-test", serverCfg)

		// create channels...
		write := plugin.PublishOutbox()
		read := plugin.SubscribeOutbox()

		received := make(chan string, 1)

		var g errgroup.Group
		g.Go(func() error {

			var e Event
			select {
			case e = <-read:
			case <-time.After(waitTimeout):
				return errors.New("timed out")
			}

			switch ev := e.(type) {
			case *pb.Message:
				received <- ev.GetText()
			default:
			}

			return nil
		})

		write <- &pb.Message{
			Text: "message to outbox",
		}

		g.Wait()

		select {
		case msg := <-received:
			assert.Equal(t, "message to outbox", msg)
		case <-time.After(waitTimeout):
			assert.FailNow(t, "timed out waiting for message to arrive at the outbox")
		}
	})
}

func TestInbox(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverCfg := withTestConfig()

	withTestAutobot(ctx, serverCfg, func() {
		// create test plugin ....
		plugin := newTestPlugin(ctx, "inbox-test", serverCfg)

		// create channels...
		write := plugin.PublishInbox()
		read := plugin.SubscribeInbox()

		received := make(chan string, 1)

		var g errgroup.Group
		g.Go(func() error {

			var e Event
			select {
			case e = <-read:
			case <-time.After(waitTimeout):
				return errors.New("timed out")
			}

			switch ev := e.(type) {
			case *pb.Message:
				received <- ev.GetText()
			default:
			}

			return nil
		})

		// this is very odd...
		time.Sleep(1 * time.Second)
		write <- &pb.Message{
			Text: "message to inbox",
		}

		g.Wait()

		select {
		case msg := <-received:
			assert.Equal(t, "message to inbox", msg)
		case <-time.After(waitTimeout):
			assert.FailNow(t, "timed out waiting for message to arrive at the inbox")
		}
	})
}

func TestReplyFunc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverCfg := withTestConfig()

	withTestAutobot(ctx, serverCfg, func() {
		// create test plugin ....
		plugin := newTestPlugin(ctx, "reply-func-plugin", serverCfg)

		// create channels...
		inbox := plugin.PublishInbox()
		outbox := plugin.SubscribeOutbox()

		replies := make(chan string, 1)

		var g errgroup.Group
		g.Go(func() error {

			var e Event
			select {
			case e = <-outbox:
			case <-time.After(waitTimeout):
				return errors.New("timed out")
			}

			switch ev := e.(type) {
			case *pb.Message:
				replies <- ev.GetText()
			default:
			}

			return nil
		})

		plugin.ReplyWithFunc(func(msg *pb.Message) (*pb.Message, error) {
			return msg.Reply("echo: " + msg.GetText()), nil
		})

		inbox <- &pb.Message{
			Text: "hello world",
		}

		g.Wait()

		select {
		case msg := <-replies:
			assert.Equal(t, "echo: hello world", msg)
		case <-time.After(waitTimeout):
			assert.FailNow(t, "timed out waiting for reply")
		}
	})
}
