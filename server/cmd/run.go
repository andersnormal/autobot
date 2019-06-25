package cmd

import (
	"context"
	"time"

	"github.com/andersnormal/autobot/pkg/nats"

	"github.com/andersnormal/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runE(cmd *cobra.Command, args []string) error {
	// create a new root
	root := new(root)

	// init logger
	root.logger = log.WithFields(log.Fields{
		"debug":   cfg.Debug,
		"verbose": cfg.Verbose,
	})

	// create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create server
	s := server.NewServer(ctx)

	// create Nats
	nats := nats.New(
		nats.WithDebug(),
		nats.WithVerbose(),
		nats.WithDataDir(cfg.NatsFilestoreDir()),
		nats.WithID("autobot"),
		nats.WithTimeout(2500*time.Millisecond),
	)
	s.Listen(nats, true)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
