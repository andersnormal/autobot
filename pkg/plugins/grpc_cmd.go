package plugins

import (
	"context"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCCommandPlugin ...
type GRPCCommandPlugin struct {
	plugin.Plugin
	GRPCCommand func() pb.CommandServer
}

// GRPCClient ...
func (p *GRPCCommandPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCAdapter{
		client: pb.NewAdapterClient(c),
		ctx:    ctx,
	}, nil
}

// GRPCServer ...
func (p *GRPCCommandPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterCommandServer(s, p.GRPCCommand())

	return nil
}

// GRPCCommand ...
type GRPCCommand struct {
	// PluginClient ...
	PluginClient *plugin.Client

	// TestServer ...
	TestServer *grpc.Server

	client pb.CommandClient
	ctx    context.Context

	mu sync.Mutex
}

func (g *GRPCCommand) Process() error {
	req := new(pb.Process_Request)

	_, err := g.client.Process(g.ctx, req)

	if err != nil {
		return err
	}

	return nil
}
