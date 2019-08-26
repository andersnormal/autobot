package proto

import (
	"crypto/sha256"
	"io"
	"os"
	"path"
)

// NewError ...
func NewError(c Error_Code, msg string) *Error {
	return &Error{
		Code:    c,
		Message: msg,
	}
}

// NewUnknownError ...
func NewUnknownError(msg string) *Error {
	return NewError(Error_UNKNOWN, msg)
}

// NewEmpty ...
func NewEmpty() *Empty {
	return &Empty{}
}

// NewErrRegister ...
func NewErrRegister(msg string) *Error {
	return NewError(Error_REGISTER, msg)
}

// NewPlugin ...
func NewPlugin(p string) *Plugin {
	return &Plugin{
		Name: path.Base(p),
		Path: p,
	}
}

// NewReply ...
func NewReply(r *Message) *Event {
	return &Event{
		Event: &Event_Reply{
			Reply: r,
		},
	}
}

// NewConfig ...
func NewConfig(cfg *Config) *Event {
	return &Event{
		Event: &Event_Config{
			Config: cfg,
		},
	}
}

// NewRegister ...
func NewRegister(pp *Plugin) *Event {
	return &Event{
		Event: &Event_Register{
			Register: &Register{
				Plugin: pp,
			},
		},
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
