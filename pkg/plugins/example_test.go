package plugins_test

import (
	"context"
	"fmt"

	. "github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
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

func ExampleWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := runtime.DefaultEnv()
	plugin, ctx := WithContext(ctx, env)

	if err := plugin.Wait(); err != nil {
		panic(err)
	}
}
