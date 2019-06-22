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
	// CommandPluginName ...
	CommandPluginName = "cmd"
	// AdapterPluginName ...
	AdapterPluginName = "adapter"
)

const (
	// CommandPluginName ...
	CommandPluginPrefix = DefaultPrefix + CommandPluginName
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
	CommandPluginPrefix: CommandPluginName,
}

// AdapterPlugin ...
type AdapterPlugin interface{}

// CommandPlugin
type CommandPlugin interface{}

// AdapterPluginFunc ...
type AdapterPluginFunc func() AdapterPlugin

// CommandPluginFunc ...
type CommandPluginFunc func() CommandPlugin

// GRPCAdapterPluginFunc ...
type GRPCAdapterPluginFunc func(broker *plugin.GRPCBroker) pb.AdapterServer

// GRPCCommandPluginFunc ...
type GRPCCommandPluginFunc func() pb.CommandServer

// ServeOpts ...
type ServeOpts struct {
	AdapterPluginFunc AdapterPluginFunc
	CommandPluginFunc CommandPluginFunc

	// Wrapped gRPC functions ...
	GRPCAdapterPluginFunc GRPCAdapterPluginFunc
	GRPCCommandPluginFunc GRPCCommandPluginFunc
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

	if opts.GRPCCommandPluginFunc != nil {
		plugins[DefaultProtocolVersion][AdapterPluginName] = &GRPCCommandPlugin{
			GRPCCommand: opts.GRPCCommandPluginFunc,
		}
	}

	return plugins
}
