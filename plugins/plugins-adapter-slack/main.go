package main

import (
	"context"

	"github.com/andersnormal/autobot/pkg/plugins"
	pb "github.com/andersnormal/autobot/proto"
)

type slackPlugin struct {
	ctx context.Context
}

// Register ...
func (s *slackPlugin) Register(ctx context.Context, req *pb.Register_Request) (*pb.Register_Response, error) {
	return &pb.Register_Response{}, nil
}

// Register ...
func (s *slackPlugin) Unregister(ctx context.Context, req *pb.Unregister_Request) (*pb.Unregister_Response, error) {
	return &pb.Unregister_Response{}, nil
}

func newPlugin(ctx context.Context) *slackPlugin {
	p := new(slackPlugin)

	return p
}

func main() {
	// create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a new plugin
	plugin := newPlugin(ctx)

	// start serving the api
	plugins.Serve(&plugins.ServeOpts{
		GRPCAdapterPluginFunc: func() pb.AdapterServer {
			return plugin
		},
	})
}
