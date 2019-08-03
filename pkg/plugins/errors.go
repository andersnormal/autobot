package plugins

import (
	"errors"
)

var (
	// ErrMissingName is signaling that the plugin needs a unique name
	ErrMissingName = errors.New("plugin: no plugin name")
	// ErrPluginAuthentication ...
	ErrPluginAuthentication = errors.New("plugin: failed to authenticate")
	// ErrPluginRegister ...
	ErrPluginRegister = errors.New("plugin: did not register")
)
