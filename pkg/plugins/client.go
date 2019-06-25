package plugins

import (
	"os/exec"

	pb "github.com/andersnormal/autobot/proto"

	plugin "github.com/hashicorp/go-plugin"
)

// ClientConfig ...
func ClientConfig(p *pb.Plugin) *plugin.ClientConfig {
	return &plugin.ClientConfig{
		Cmd:              exec.Command(p.GetMeta().GetPath()),
		HandshakeConfig:  Handshake,
		VersionedPlugins: VersionedPlugins,
		Managed:          true,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		AutoMTLS:         true,
	}
}

// Client ...
func Client(p *pb.Plugin) *plugin.Client {
	return plugin.NewClient(ClientConfig(p))
}
