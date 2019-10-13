package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/andersnormal/autobot/pkg/config"
	"github.com/andersnormal/autobot/pkg/nats"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"

	"github.com/andersnormal/pkg/server"
	"golang.org/x/sync/errgroup"
)

func withTestAutobot(ctx context.Context, t *testing.T, env *runtime.Environment, f func(*testing.T, context.CancelFunc)) {
	cfg := config.New()
	cfg.Verbose = true
	cfg.Debug = true
	// cfg.Nats.HTTPPort = 0

	// create server
	ctx, cancel := context.WithCancel(ctx)

	s, ctx := server.WithContext(ctx)

	// only will use temp dir for tests...
	cfg.DataDir, _ = ioutil.TempDir("/tmp", "")
	defer os.RemoveAll(cfg.DataDir)
	defer cancel()

	n := nats.New(cfg, nats.Timeout(5*time.Second))

	s.Listen(n, true)

	var g errgroup.Group

	g.Go(func() error {
		if err := s.Wait(); err != nil {
			return err
		}

		return nil
	})

	<-n.Ready()

	f(t, cancel)

	// wait for server to close ...
	g.Wait()
}
