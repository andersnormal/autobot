package runtime_test

import (
	"context"
	"fmt"

	. "github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
)

func ExampleDefaultEnv() {
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
