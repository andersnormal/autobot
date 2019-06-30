package cmd

import (
	"io"
	"os/exec"
)

// Cmd ...
type Cmd interface {
	Run() error
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type cmd struct {
  cmd  *exec.Cmd
  env Env
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
