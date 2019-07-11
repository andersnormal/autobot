package api

import (
	"context"
	"net"

	pb "github.com/andersnormal/autobot/proto"
	"github.com/andersnormal/autobot/server/config"

	s "github.com/andersnormal/pkg/server"
	"google.golang.org/grpc"
)

type Server interface {
	s.Listener
}

type server struct {
	s   *grpc.Server
	cfg *config.Config

	opts *Opts
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
}

// New returns a new server
func New(cfg *config.Config, opts ...Opt) Server {
	options := new(Opts)

	s := new(server)
	s.opts = options
	s.cfg = cfg

	configure(s, opts...)

	return s
}

func (s *server) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		s.s = grpc.NewServer()
		pb.RegisterAutobotServer(s.s, &API{s.cfg})

		lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
		if err != nil {
			return err
		}

		ready()

		if err = s.s.Serve(lis); err != nil {
			return err
		}

		return nil
	}
}

func (s *server) Stop() error {
	s.s.GracefulStop()

	return nil
}

func configure(s *server, opts ...Opt) error {
	for _, o := range opts {
		o(s.opts)
	}

	return nil
}
