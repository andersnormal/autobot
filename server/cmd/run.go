package cmd

import (
	"context"
	"time"

	"github.com/andersnormal/autobot/pkg/cmd"
	"github.com/andersnormal/autobot/pkg/nats"
	"github.com/andersnormal/autobot/pkg/run"

	"github.com/andersnormal/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type root struct {
	logger *log.Entry
	nats   nats.Nats
}

func runE(c *cobra.Command, args []string) error {
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
		root.nats = nats.New(cfg, nats.Timeout(5*time.Second))

		// create Nats
		s.Listen(root.nats, true)
	}

	// get plugins ...
	plugins, err := cfg.LoadPlugins()
	if err != nil {
		root.logger.Fatalf("error getting plugins: %v", err)
	}

	// env ...
	env := cmd.DefaultEnv()
	env.Set("AUTOBOT_CLUSTER_URL", cfg.Nats.ClusterURL)
	env.Set("AUTOBOT_CLUSTER_ID", cfg.Nats.ClusterID)
	env.Set("AUTOBOT_CLUSTER_DISCOVERY", cfg.Nats.Discovery)

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
