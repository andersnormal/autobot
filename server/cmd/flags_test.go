package cmd

import (
	"testing"

	"github.com/andersnormal/autobot/pkg/config"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddFlags_Defaults(t *testing.T) {
	assert := assert.New(t)

	cfg := config.New()
	TestCmd := &cobra.Command{}

	assert.NotPanics(func() {
		addFlags(TestCmd, cfg)
	})

	c, err := TestCmd.Flags().GetString("config")
	assert.NoError(err)
	assert.Equal(c, "")

	v, err := TestCmd.Flags().GetBool("verbose")
	assert.NoError(err)
	assert.Equal(v, config.DefaultVerbose)

	d, err := TestCmd.Flags().GetBool("debug")
	assert.NoError(err)
	assert.Equal(d, config.DefaultDebug)

	lf, err := TestCmd.Flags().GetString("log-format")
	assert.NoError(err)
	assert.Equal(lf, config.DefaultLogFormat)

	ll, err := TestCmd.Flags().GetString("log-level")
	assert.NoError(err)
	assert.Equal(ll, config.DefaultLogLevel)

	clustering, err := TestCmd.Flags().GetBool("clustering")
	assert.NoError(err)
	assert.Equal(clustering, cfg.Nats.Clustering)

	bootstrap, err := TestCmd.Flags().GetBool("bootstrap")
	assert.NoError(err)
	assert.Equal(bootstrap, cfg.Nats.Bootstrap)

	nodeID, err := TestCmd.Flags().GetString("node-id")
	assert.NoError(err)
	assert.Equal(nodeID, cfg.Nats.ClusterNodeID)

	peers, err := TestCmd.Flags().GetStringSlice("peers")
	assert.NoError(err)
	assert.ElementsMatch(cfg.Nats.ClusterPeers, peers)

	url, err := TestCmd.Flags().GetString("url")
	assert.NoError(err)
	assert.Equal(cfg.Nats.ClusterURL, url)
}
