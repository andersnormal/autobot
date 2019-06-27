package main

import (
	"context"
	"fmt"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"
)

type slackPlugin struct {
}

func main() {
	// create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// plugin ....
	plugin := plugins.New(plugins.DefaultEnv())

	// subscribee ...
	plugin.Subscribe(ctx, func(e *pb.Event) {
		fmt.Printf("received: %v", e)
	})
}
