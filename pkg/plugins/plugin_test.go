package plugins

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"
)

func TestInboxPubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv().WithName("test").WithClusterURL("nats://localhost:4222")
	plugin, ctx := WithContext(ctx, env)

	plugin.PublishInbox() <- &pb.Message{
		Text: "hello world",
	}

	var receivedInboxMessage string

	err := plugin.ReplyWithFunc(func(msg *pb.Message) (*pb.Message, error) {
		receivedInboxMessage = msg.GetText()
		return nil, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	assert.Eventually(t, func() bool {
		if receivedInboxMessage == "hello world" {
			return true
		}

		return false
	}, 2*time.Second, 10*time.Millisecond)
}
