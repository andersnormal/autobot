package discovery

import (
	"crypto/sha256"
	"io"
	"os"
)

// Plugin ...
type Plugin struct {
	// UUID ...
	UUID string
	// Name ...
	Name string
	// Identifier...
	Identifier string
	// Version ...
	Version string
	// Path ...
	Path string
}

// SHA256 ...
func (p *Plugin) SHA256() ([]byte, error) {
	f, err := os.Open(p.Path)
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
