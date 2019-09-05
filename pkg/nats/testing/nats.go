package testing

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/andersnormal/autobot/pkg/config"
	"github.com/andersnormal/autobot/pkg/nats"
)

// NewTestServer ...
func NewTestServer(cfg *config.Config, opts ...nats.Opt) nats.Nats {
	// only will use temp dir for tests...
	cfg.DataDir, _ = ioutil.TempDir("/tmp", "")
	s := nats.New(cfg, nats.Timeout(5*time.Second))

	return s
}

// WithAutobot ...
func WithTestServer(ctx context.Context, cfg *config.Config, f func()) {
	s := NewTestServer(cfg)

	go s.Start(ctx, func() {})

	f()
}
