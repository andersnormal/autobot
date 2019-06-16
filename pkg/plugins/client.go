package plugins

import (
	"os/exec"

	plugin "github.com/hashicorp/go-plugin"
)

// ClientConfig ...
func ClientConfig(m PluginMeta) *plugin.ClientConfig {
	return &plugin.ClientConfig{
		Cmd:              exec.Command(m.Path),
		HandshakeConfig:  Handshake,
		VersionedPlugins: VersionedPlugins,
		Managed:          true,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		AutoMTLS:         true,
	}
}

// Client ...
func Client(m PluginMeta) *plugin.Client {
	return plugin.NewClient(ClientConfig(m))
}
