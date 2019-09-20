package plugins

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const waitTimeout = 5 * time.Second

func TestInbox(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverCfg := withTestConfig()

	withTestAutobot(ctx, serverCfg, func() {
		// create test plugin ....
		plugin := newTestPlugin(ctx, "inbox-test", serverCfg)

		received := make(chan string, 1)

		in := plugin.SubscribeInbox()

		var g errgroup.Group
		g.Go(func() error {

			var e Event
			var ok bool

			fmt.Println("we are definitely waiting")
			select {
			case e, ok = <-in:
				fmt.Println("received some shit from inbox")
				if !ok {
					fmt.Println("not ok!")
				}
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

		plugin.PublishInbox() <- &pb.Message{
			Text: "message to inbox",
		}

		g.Wait()
		cancel()
		plugin.Wait()

		select {
		case msg := <-received:
			assert.Equal(t, "message to inbox", msg)
		case <-time.After(waitTimeout):
			assert.FailNow(t, "timed out waiting for message to arrive to the inbox")
		}
	})
}

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
