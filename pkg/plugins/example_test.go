package plugins_test

import (
	"context"
	"fmt"
	"log"

	. "github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"
)

func ExamplePlugin_Events() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	for {
		select {
		case e := <-plugin.Events():
			fmt.Printf("received event: %v", e)
		case <-ctx.Done():
			return
		}
	}
}

func ExamplePlugin_SubscribeInbox() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	for {
		select {
		case e := <-plugin.SubscribeInbox():
			fmt.Printf("received a new message: %v", e)
		case <-ctx.Done():
			return
		}
	}
}

func ExamplePlugin_SubscribeOutbox() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	for {
		select {
		case e := <-plugin.SubscribeOutbox():
			fmt.Printf("received message to be send out: %v", e)
		case <-ctx.Done():
			return
		}
	}
}

func ExamplePlugin_PublishInbox() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	inbox := plugin.PublishInbox()

	inbox <- &pb.Bot{Bot: &pb.Bot_Message{
		Message: &pb.Message{
			Text: "foo != bar",
		},
	}}
}

func ExamplePlugin_PublishOutbox() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	inbox := plugin.PublishOutbox()

	inbox <- &pb.Bot{Bot: &pb.Bot_Message{
		Message: &pb.Message{
			Text: "foo != bar",
		},
	}}
}

func ExamplePlugin() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	// here you can interact with the plugin data

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}

func ExamplePlugin_ReplyWithFunc() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	err := plugin.ReplyWithFunc(func(event *pb.Bot) (*pb.Bot, error) {
		log.Printf("received message: %v", event)

		return event, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}

func ExamplePlugin_Wait() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}

func ExampleWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}
