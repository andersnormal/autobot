package plugins

import (
	"context"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCAdapterPlugin ...
type GRPCAdapterPlugin struct {
	plugin.Plugin
	GRPCAdapter func(broker *plugin.GRPCBroker) pb.AdapterServer
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
	pb.RegisterAdapterServer(s, p.GRPCAdapter(broker))

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
	server *grpc.Server

	mu sync.Mutex
}

func (g *GRPCAdapter) Serve(s pb.BotServer) (uint32, error) {
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		g.server = grpc.NewServer(opts...)
		pb.RegisterBotServer(g.server, s)

		return g.server
	}

	brokerID := g.broker.NextId()
	go g.broker.AcceptAndServe(brokerID, serverFunc)

	return brokerID, nil
}

func (g *GRPCAdapter) Register(id uint32) error {
	req := new(pb.Register_Request)
	req.Id = id

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
