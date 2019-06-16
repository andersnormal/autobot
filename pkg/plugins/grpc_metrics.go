package plugins

import (
	"context"
	"sync"

	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCMetricsPlugin ...
type GRPCMetricsPlugin struct {
	plugin.Plugin
	GRPCMetrics func() pb.MetricsServer
}

// GRPCClient ...
func (p *GRPCMetricsPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCMetrics{
		client: pb.NewMetricsClient(c),
		ctx:    ctx,
	}, nil
}

// GRPCServer ...
func (p *GRPCMetricsPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterMetricsServer(s, p.GRPCMetrics())

	return nil
}

// GRPCMetrics ...
type GRPCMetrics struct {
	// PluginClient ...
	PluginClient *plugin.Client

	// TestServer ...
	TestServer *grpc.Server

	client pb.MetricsClient
	ctx    context.Context

	mu sync.Mutex
}

func (g *GRPCMetrics) Payload() (Payload, error) {
	resp, err := g.client.GetPayload(g.ctx, new(pb.GetPayload_Request))

	if err != nil {
		return Payload{}, err
	}

	return Payload{resp.GetName()}, nil
}

func (g *GRPCMetrics) Stop() error {
	_, err := g.client.Stop(g.ctx, new(pb.Stop_Request))
	if err != nil {
		return err
	}

	return nil
}
