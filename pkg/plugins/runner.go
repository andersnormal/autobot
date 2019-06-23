package plugins

import (
	"context"
	"sync"

	"github.com/andersnormal/pkg/server"
)

var _ server.Listener = (*runner)(nil)

type runner struct {
	ctx     context.Context
	cancel  context.CancelFunc
	plugins []PluginMeta

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

func NewRunner(plugins []PluginMeta) server.Listener {
	r := new(runner)

	// pass in the plugins
	r.plugins = plugins
	r.exit = make(chan struct{})

	return r
}

// Start ...
func (r *runner) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		r.ctx, r.cancel = context.WithCancel(ctx)

		// get to run the plugins
		for _, p := range r.plugins {
			r.run(runFunc(r.ctx, p))
		}

		// call for ready
		ready()

		// run this in a loop, and wait for exit
		// or ctx being canceled
		<-r.exit

		return r.err
	}
}

// Stop ...
func (r *runner) Stop() error {
	// todo: add timeout here
	r.cancel()

	r.wg.Wait()

	return nil
}

func (r *runner) run(f func() error) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()

		if err := f(); err != nil {
			r.errOnce.Do(func() {
				r.err = err
				r.exit <- struct{}{}
			})
		}
	}()
}

func runFunc(ctx context.Context, p PluginMeta) func() error {
	return func() error {
		client := Client(p)
		defer client.Kill()

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		// Request the plugin
		raw, err := rpcClient.Dispense(p.PluginName)
		if err != nil {
			return err
		}

		if _, ok := raw.(*GRPCAdapter); ok {
			//  have to start here and hold
		}

		return nil
	}
}
