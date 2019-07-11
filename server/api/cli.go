package api

import (
	"context"

	pb "github.com/andersnormal/autobot/proto"
	"github.com/andersnormal/autobot/server/config"
)

type CLI struct {
	// config to be used with the api
	cfg *config.Config
}

// ListPlugins ...
func (c *CLI) ListPlugins(ctx context.Context, req *pb.ListDomains_Request) (*pb.ListDomains_Response, error) {
	return &pb.ListDomains_Response{}, nil
}
