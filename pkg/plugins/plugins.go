package plugins

import (
	"strings"

	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
)

const (
	// DefaultProtocolVersion ...
	DefaultProtocolVersion = 1
)

const (
	// MetricsPluginName ...
	MetricsPluginName = "metrics"

	// UnknownPluginName ...
	UnknownPluginName = "unknown"
)

const (
	// MetricsPluginPrefix ...
	MetricsPluginPrefix = "plugin-metrics"
)

// PluginsMap ...
type PluginsMap map[string]string

// PluginName ...
func (mm PluginsMap) PluginName(name string) (string, bool) {
	for k, plugin := range mm {
		if strings.HasPrefix(name, k) {
			return plugin, true
		}
	}

	return "", false
}

// Plugins ...
var Plugins = PluginsMap{
	MetricsPluginPrefix: MetricsPluginName,
}

// MetricsPlugin ...
type MetricsPlugin interface {
}

// MetricsPluginFunc ...
type MetricsPluginFunc func() MetricsPlugin

// GRPCMetricsPluginFunc ...
type GRPCMetricsPluginFunc func() pb.MetricsServer

// ServeOpts ...
type ServeOpts struct {
	MetricsPluginFunc MetricsPluginFunc

	// Wrapped gRPC functions ...
	GRPCMetricsPluginFunc GRPCMetricsPluginFunc
}

// Handshake ...
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  DefaultProtocolVersion,
	MagicCookieKey:   "DO_AWESOME",
	MagicCookieValue: "foo",
}

// Serve ...
func Serve(opts *ServeOpts) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig:  Handshake,
		VersionedPlugins: pluginSet(opts),
		GRPCServer:       plugin.DefaultGRPCServer,
	})
}

func pluginSet(opts *ServeOpts) map[int]plugin.PluginSet {
	plugins := map[int]plugin.PluginSet{
		DefaultProtocolVersion: plugin.PluginSet{},
	}

	if opts.GRPCMetricsPluginFunc != nil {
		plugins[DefaultProtocolVersion][MetricsPluginName] = &GRPCMetricsPlugin{
			GRPCMetrics: opts.GRPCMetricsPluginFunc,
		}
	}

	return plugins
}
