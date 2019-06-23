package plugins

import (
	"strings"

	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
)

const (
	// DefaultProtocolVersion ...
	DefaultProtocolVersion = 1
	// DefaultPrefix ...
	DefaultPrefix = "plugin-"
)

const (
	// AdapterPluginName ...
	AdapterPluginName = "adapter"
)

const (
	// AdapterPluginPrefix ...
	AdapterPluginPrefix = DefaultPrefix + AdapterPluginName
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
	AdapterPluginPrefix: AdapterPluginName,
}

// AdapterPlugin ...
type AdapterPlugin interface{}

// AdapterPluginFunc ...
type AdapterPluginFunc func() AdapterPlugin

// GRPCAdapterPluginFunc ...
type GRPCAdapterPluginFunc func(broker *plugin.GRPCBroker) pb.AdapterServer

// ServeOpts ...
type ServeOpts struct {
	AdapterPluginFunc AdapterPluginFunc

	// Wrapped gRPC functions ...
	GRPCAdapterPluginFunc GRPCAdapterPluginFunc
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

	if opts.GRPCAdapterPluginFunc != nil {
		plugins[DefaultProtocolVersion][AdapterPluginName] = &GRPCAdapterPlugin{
			GRPCAdapter: opts.GRPCAdapterPluginFunc,
		}
	}

	return plugins
}
