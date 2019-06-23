package plugins

import (
	plugin "github.com/hashicorp/go-plugin"
)

// VersionedPlugins ...
var VersionedPlugins = map[int]plugin.PluginSet{
	1: {
		AdapterPluginName: &GRPCAdapterPlugin{},
	},
}
