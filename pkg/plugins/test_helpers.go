package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/andersnormal/autobot/pkg/config"
	"github.com/andersnormal/autobot/pkg/nats"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"
	"github.com/andersnormal/pkg/server"
	"github.com/nats-io/stan.go"
	"golang.org/x/sync/errgroup"
)

// runtime can only be initialized only once
var defaultRuntime = runtime.Default()

func withTestConfig() *config.Config {
	cfg := config.New()

	cfg.Nats.Port = 4224
	cfg.Nats.HTTPPort = 8224
	cfg.Nats.ClusterURL = "nats://localhost:4224"

	return cfg
}

func withTestAutobot(ctx context.Context, cfg *config.Config, f func()) {
	// some config overrides ...
	cfg.Nats.Port = 4224
	cfg.Nats.ClusterURL = "nats://localhost:4224"
	cfg.Nats.HTTPPort = 8224
	cfg.StatusAddr = ":8444"

	// create server
	ctx, cancel := context.WithCancel(ctx)

	s, ctx := server.WithContext(ctx)

	// only will use temp dir for tests...
	cfg.DataDir, _ = ioutil.TempDir("/tmp", "")
	defer os.RemoveAll(cfg.DataDir)
	defer os.Remove(cfg.DataDir)

	n := nats.New(cfg, nats.Timeout(5*time.Second))

	s.Listen(n, true)

	var g errgroup.Group

	g.Go(func() error {
		s.Wait()
		return nil
	})

	time.Sleep(3 * time.Second)
	f()

	// wait for server to close ...
	cancel()
	g.Wait()
	time.Sleep(3 * time.Second)
}

func newTestPlugin(ctx context.Context, name string, serverCfg *config.Config) *Plugin {
	defaultRuntime.Name = name
	defaultRuntime.ClusterURL = serverCfg.Nats.ClusterURL

	plugin, _ := WithContext(ctx, defaultRuntime, WithSubscriptionOpts(
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
	))

	return plugin
}
