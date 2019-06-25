package proto

import (
	"crypto/sha256"
	"io"
	"os"
	"path"
	"strings"
)

const (
	// DefaultPrefix ...
	DefaultPrefix = "plugin-"

	// AdapterPluginName ...
	AdapterPluginName = "adapter"

	// UnknownPluginName ...
	UnknownPluginName = "unknown"
)

const (
	// AdapterPluginPrefix ...
	AdapterPluginPrefix = DefaultPrefix + AdapterPluginName
)

// PluginsMap ...
type PluginsMap map[string]string

// PluginName ...
func (mm PluginsMap) PluginIdentifier(name string) string {
	for k, plugin := range mm {
		if strings.HasPrefix(name, k) {
			return plugin
		}
	}

	return UnknownPluginName
}

// Plugins ...
var Plugins = PluginsMap{
	AdapterPluginPrefix: AdapterPluginName,
}

// NewPlugin ...
func NewPlugin(p string) *Plugin {
	plugin := new(Plugin)

	// create meta
	meta := new(Plugin_Meta)
	meta.Name = path.Base(p)
	meta.Path = p

	// identify the plugin
	meta.Identifier = Plugins.PluginIdentifier(meta.GetName())

	// set meta information
	plugin.Meta = meta

	return plugin
}

// SHA256 ...
func (p *Plugin) SHA256() ([]byte, error) {
	f, err := os.Open(p.Meta.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
