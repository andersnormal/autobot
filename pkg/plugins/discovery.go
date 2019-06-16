package plugins

import (
	"crypto/sha256"
	"io"
	"os"
	"path"
)

// Version ...
type Version string

// PluginMeta ...
type PluginMeta struct {
	// Name ...
	Name string

	// Version ...
	Version Version

	// Path ...
	Path string

	// PluginName ...
	PluginName string
}

// NewPluginMeta ...
func NewPluginMeta(p string) (PluginMeta, error) {
	meta := PluginMeta{}

	n := path.Base(p)
	name, found := Plugins.PluginName(n)
	if !found {
		return meta, ErrNoPlugin
	}

	meta.PluginName = name
	meta.Path = p
	meta.Name = n

	return meta, nil
}

// SHA256 ...
func (m PluginMeta) SHA256() ([]byte, error) {
	f, err := os.Open(m.Path)
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
