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

func ExamplePlugin_ReplyWithFunc() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	err := plugin.ReplyWithFunc(func(event *pb.Event) (*pb.Event, error) {
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

func ExampleWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}
