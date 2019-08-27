package cmd

import (
	"context"
	"time"

	"github.com/andersnormal/autobot/pkg/discovery"
	"github.com/andersnormal/autobot/pkg/nats"
	"github.com/andersnormal/autobot/pkg/run"
	"github.com/andersnormal/autobot/server/api"

	"github.com/andersnormal/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type root struct {
	logger *log.Entry
	nats   nats.Nats
}

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

	// NATS ...
	if !cfg.Nats.Disabled {
		root.nats = nats.New(
			nats.WithDebug(),
			nats.WithVerbose(),
			nats.WithDataDir(cfg.NatsFilestoreDir()),
			nats.WithID("autobot"),
			nats.WithTimeout(2500*time.Millisecond),
		)

		// create Nats
		s.Listen(root.nats, true)
	}

	// get plugins ...
	plugins, err := cfg.Plugins()
	if err != nil {
		root.logger.Fatalf("error getting plugins: %v", err)
	}

	// create plugin discovery
	reg := discovery.New()
	s.Listen(reg, true)

	// create apis ...
	a := api.New(cfg)
	s.Listen(a, true)

	// env ...
	env := cfg.Env()

	// run plugins ...
	r := run.New(plugins, env, root.logger)
	s.Listen(r, true)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
