package plugins

import (
	"context"
	"sync"

	"github.com/andersnormal/autobot/pkg/nats"
	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCAdapterPlugin ...
type GRPCAdapterPlugin struct {
	plugin.Plugin
	GRPCAdapter func() pb.AdapterServer
}

// GRPCClient ...
func (p *GRPCAdapterPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCAdapter{
		client: pb.NewAdapterClient(c),
		ctx:    ctx,
		broker: broker,
	}, nil
}

// GRPCServer ...
func (p *GRPCAdapterPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterAdapterServer(s, p.GRPCAdapter())

	return nil
}

// GRPCAdapter ...
type GRPCAdapter struct {
	// PluginClient ...
	PluginClient *plugin.Client

	// TestServer ...
	TestServer *grpc.Server

	client pb.AdapterClient
	ctx    context.Context

	broker *plugin.GRPCBroker

	mu sync.Mutex
}

func (g *GRPCAdapter) Register(nats nats.Nats) error {
	req := new(pb.Register_Request)
	req.Topic = &pb.Topic{
		ClusterUrl: nats.Addr().String(),
		ClusterId:  nats.ClusterID(),
		Name:       "messages",
	}

	_, err := g.client.Register(g.ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (g *GRPCAdapter) Unregister() error {
	_, err := g.client.Unregister(g.ctx, new(pb.Unregister_Request))
	if err != nil {
		return err
	}

	return nil
}
