package plugin

import (
	"fmt"
	"io"
	"os/exec"
)

// Plugin ...
type Plugin interface {
	Run() error
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

// Env ...
type Env map[string]string

// Strings ...
func (ev Env) Strings() []string {
	var env []string
	for k, v := range ev {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

type plugin struct {
	cmd  *exec.Cmd
	opts *Opts
}

type Opt func(*Opts)

type Opts struct {
	Dir    string
	Env    Env
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
