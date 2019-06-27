package plugin

import (
	"io"
	"os"
	"os/exec"
)

// New ...
func New(cmd *exec.Cmd, opts ...Opt) Plugin {
	options := new(Opts)
	options.Env = make(Env)

	p := new(plugin)
	p.opts = options

	configure(p, opts...)

	return p
}

// Stdin ...
func (p *plugin) Stdin() io.Reader {
	return p.cmd.Stdin
}

// Stdout ...
func (p *plugin) Stdout() io.Writer {
	return p.cmd.Stdout
}

// Stderr ...
func (p *plugin) Stderr() io.Writer {
	return p.cmd.Stderr
}

// Run ... context via exec.CommandContext
func (p *plugin) Run() error {
	// run the command, and wait
	// todo: restart
	if err := p.cmd.Run(); err != nil {
		return err
	}

	return nil
}

func configure(p *plugin, opts ...Opt) error {
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
