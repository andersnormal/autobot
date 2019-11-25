package config_test

import (
	"os"
	"path"
	"testing"

	. "github.com/andersnormal/autobot/pkg/config"
	"github.com/andersnormal/autobot/pkg/plugins/runtime"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := New()

	assert.NotNil(cfg)

	assert.Equal(cfg.Verbose, DefaultVerbose)
	assert.Equal(cfg.LogLevel, DefaultLogLevel)
	assert.Equal(cfg.LogFormat, DefaultLogFormat)
	assert.Equal(cfg.ReloadSignal, DefaultReloadSignal)
	assert.Equal(cfg.TermSignal, DefaultTermSignal)
	assert.Equal(cfg.KillSignal, DefaultKillSignal)
	assert.Equal(cfg.StatusAddr, DefaultStatusAddr)
	assert.Equal(cfg.Debug, DefaultDebug)
	assert.Equal(cfg.DataDir, DefaultDataDir)
	assert.Equal(cfg.FileChmod, os.FileMode(DefaultFileChmod))

	assert.Equal(cfg.Nats.ClusterID, DefaultNatsClusterID)
	assert.Equal(cfg.Nats.DataDir, DefaultNatsDataDir)
	assert.Equal(cfg.Nats.HTTPPort, DefaultNatsHTTPPort)
	assert.Equal(cfg.Nats.Inbox, runtime.DefaultClusterInbox)
	assert.Equal(cfg.Nats.Inbox, runtime.DefaultClusterInbox)
	assert.Equal(cfg.Nats.SPort, DefaultSNatsPort)
	assert.Equal(cfg.Nats.NPort, DefaultNNatsPort)
	assert.Equal(cfg.Nats.LogDir, DefaultNatsRaftLogDir)
}

func TestConfig_NatsFilestoreDir(t *testing.T) {
	assert := assert.New(t)

	cfg := New()

	assert.NotNil(cfg)
	assert.Equal(path.Join(DefaultDataDir, DefaultNatsDataDir), cfg.NatsFilestoreDir())
}

func TestConfig_RaftLogPath(t *testing.T) {
	assert := assert.New(t)

	cfg := New()

	assert.NotNil(cfg)
	assert.Equal(path.Join(cfg.DataDir, DefaultNatsRaftLogDir), cfg.RaftLogPath())
}
