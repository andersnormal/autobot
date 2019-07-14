package plugins

import (
	"errors"
)

var (
	// ErrPluginAuthentication ...
	ErrPluginAuthentication = errors.New("plugin: failed to authenticate")
	// ErrPluginRegister ...
	ErrPluginRegister = errors.New("plugin: did not register")
)
