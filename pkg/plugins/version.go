package plugins

import (
	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
)

// VersionedPlugins ...
var VersionedPlugins = map[int]plugin.PluginSet{
	1: {
		pb.AdapterPluginName: &GRPCAdapterPlugin{},
	},
}
