package plugins

import (
	pb "github.com/andersnormal/autobot/proto"
)

// Filter is operating on the events flowing in and out
// of the inbox /outbox via the subscriptions or publications
// to the message queue.
type Filter interface {
	Filter(*pb.Event) *pb.Event
}

type filter struct {
	funcs []func(*pb.Event) *pb.Event
}

// FilterOpt is an option that can be passed to the filter functions.
type FilterOpt func(*Plugin) func(*pb.Event) *pb.Event

// NewFilter is creating a new filter.
// This is mostly used internally.
func NewFilter(p *Plugin, opts ...FilterOpt) Filter {
	f := new(filter)

	for _, o := range opts {
		f.funcs = append(f.funcs, o(p))
	}

	return f
}

// Filter is handeling an inflowing message according to the middleware.
func (f *filter) Filter(e *pb.Event) *pb.Event {
	// doing this in reverse to keep order of the middleware
	for i := len(f.funcs) - 1; i >= 0; i-- {
		e = f.funcs[i](e)
	}

	return e
}

// WithFilterPlugin is filtering an event for the plugin
// as configured by its meta information.
func WithFilterPlugin() FilterOpt {
	return func(p *Plugin) func(e *pb.Event) *pb.Event {
		return func(e *pb.Event) *pb.Event {
			if e.GetPlugin() != nil && e.GetPlugin().GetName() != p.meta.GetName() {
				return nil
			}

			return e
		}
	}
}
