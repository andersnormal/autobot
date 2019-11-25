package testing

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
)

func WithAutobot(t *testing.T, env *runtime.Environment, f func(*testing.T)) {
	cfg := config.New()
	cfg.Verbose = true
	cfg.Debug = true
	cfg.Nats.HTTPPort = 0

	// create server
	ctx, cancel := context.WithCancel(context.Background())
	s, _ := server.WithContext(ctx)

	// only will use temp dir for tests...
	cfg.DataDir, _ = ioutil.TempDir(os.TempDir(), "")
	defer os.RemoveAll(cfg.DataDir)

	n := nats.New(cfg, nats.Timeout(5*time.Second))
	s.Listen(n, true)

	go func() {
		<-n.Ready()

		f(t)

		cancel()
	}()

	// wait for server to close ...
	_ = s.Wait()

	// <-time.After(3 * time.Second)
}
