package plugins_test

import (
	"context"
	"fmt"
	"log"

	. "github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	pb "github.com/andersnormal/autobot/proto"
)

func ExamplePlugin_SubscribeInbox() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.Default()
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

	env := runtime.Default()
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

	env := runtime.Default()
	plugin, ctx := WithContext(ctx, env)

	msg := &pb.Message{
		Text: "foo != bar",
	}

	if err := plugin.PublishInbox(msg); err != nil {
		panic(err)
	}
}

func ExamplePlugin_PublishOutbox() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.Default()
	plugin, ctx := WithContext(ctx, env)

	msg := &pb.Message{
		Text: "foo != bar",
	}

	if err := plugin.PublishOutbox(msg); err != nil {
		panic(err)
	}
}

func ExamplePlugin() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.Default()
	plugin, ctx := WithContext(ctx, env)

	// here you can interact with the plugin data

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}

func ExamplePlugin_ReplyWithFunc() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.Default()
	plugin, ctx := WithContext(ctx, env)

	err := plugin.ReplyWithFunc(func(ctx Context) error {
		log.Printf("received message: %v", ctx.Message())

		return nil
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

	env := runtime.Default()
	plugin, ctx := WithContext(ctx, env)

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}

func ExampleWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.Default()
	plugin, ctx := WithContext(ctx, env)

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}
