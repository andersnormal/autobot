package dialer

import (
	"github.com/andersnormal/autobot/cli/config"

	"google.golang.org/grpc"
)

func NewDialer(cfg *config.Config) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if cfg.TLS {
		// tba
	}

	if !cfg.TLS {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(cfg.ServerAddr, opts...)
}
