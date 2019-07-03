package proto

import (
	"crypto/sha256"
	"io"
	"os"
	"path"
)

// NewPlugin ...
func NewPlugin(p string) *Plugin {
	return &Plugin{
		Name: path.Base(p),
		Path: p,
	}
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
