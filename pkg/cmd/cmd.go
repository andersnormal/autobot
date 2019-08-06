package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/andersnormal/autobot/pkg/plugins"
)

// Env ...
type Env map[string]string

func (ev Env) Strings() []string {
	var env []string
	for k, v := range ev {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

// DefaultEnv ...
func DefaultEnv() Env {
	env := Env{
		plugins.AutobotClusterID:  "",
		plugins.AutobotClusterURL: "",
	}

	return env
}

// Cmd ...
type Cmd interface {
	Run(context.Context) func() error

	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type cmd struct {
	cmd  *exec.Cmd
	env  Env
	opts *Opts
}

type Opt func(*Opts)

type Opts struct {
	Dir string
	Env Env

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New ...
func New(c *exec.Cmd, env Env, opts ...Opt) Cmd {
	options := new(Opts)

	p := new(cmd)
	p.opts = options
	p.cmd = c
	p.env = env

	configure(p, opts...)

	return p
}

// Stdin ...
func (p *cmd) Stdin() io.Reader {
	return p.cmd.Stdin
}

// Stdout ...
func (p *cmd) Stdout() io.Writer {
	return p.cmd.Stdout
}

// Stderr ...
func (p *cmd) Stderr() io.Writer {
	return p.cmd.Stderr
}

// Run ... context via exec.CommandContext
func (p *cmd) Run(ctx context.Context) func() error {
	return func() error {
		// set env ...
		p.cmd.Env = p.env.Strings()

		// run the command, and wait
		// todo: restart
		if err := p.cmd.Run(); err != nil {
			return err
		}

		return nil
	}
}

func configure(p *cmd, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	if p.opts.Stdin == nil {
		p.opts.Stdin = os.Stdin
	}

	if p.opts.Stdout == nil {
		p.opts.Stdout = os.Stdout
	}

	if p.opts.Stderr == nil {
		p.opts.Stderr = os.Stderr
	}

	return nil
}
