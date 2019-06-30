package cmd

import (
	"context"
	"time"

	c "github.com/andersnormal/autobot/pkg/cmd"
	"github.com/andersnormal/autobot/pkg/nats"
	p "github.com/andersnormal/autobot/pkg/plugins"
	"github.com/andersnormal/autobot/pkg/run"

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

	// get plugins ...
	plugins, err := cfg.Plugins()
	if err != nil {
		root.logger.Fatalf("error getting plugins: %v", err)
	}

	// create env ...
	env := c.Env{
		p.AutobotClusterID:   nats.ClusterID(),
		p.AutobotClusterURL:  nats.Addr().String(),
		p.AutobotTopicEvents: "events",
  }

	// run plugins ...
	r := run.New(plugins, env)
	s.Listen(r, true)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
