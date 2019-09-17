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
)

// runtime can only be initialized only once
var defaultRuntime = runtime.Default()

func withTestAutobot(cfg *config.Config, f func()) {
	// create server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, ctx := server.WithContext(ctx)

	// only will use temp dir for tests...
	cfg.DataDir, _ = ioutil.TempDir("/tmp", "")
	defer os.RemoveAll(cfg.DataDir)
	defer os.Remove(cfg.DataDir)

	n := nats.New(cfg, nats.Timeout(5*time.Second))

	s.Listen(n, true)

	go s.Wait()

	time.Sleep(4 * time.Second)
	f()
}

func newTestPlugin(ctx context.Context, name string, serverCfg *config.Config) *Plugin {
	defaultRuntime.Name = name
	defaultRuntime.ClusterURL = serverCfg.Nats.ClusterURL

	plugin, _ := WithContext(ctx, defaultRuntime, WithSubscriptionOpts(
		stan.StartAtSequence(1),
	))

	return plugin
}
