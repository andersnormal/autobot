package plugins

import (
	pb "github.com/andersnormal/autobot/proto"

	"github.com/google/uuid"
)

// Filter is operating on the events flowing in and out
// of the inbox /outbox via the subscriptions or publications
// to the message queue.
type Filter interface {
	Filter(*pb.Event) (*pb.Event, error)
}

type filter struct {
	funcs []func(*pb.Event) (*pb.Event, error)
}

// FilterOpt is an option that can be passed to the filter functions.
type FilterOpt func(*Plugin) FilterFunc

// FilterFunc is the wrapper for the plugin filtering option.
type FilterFunc func(e *pb.Event) (*pb.Event, error)

var (
	// DefaultOutboxFilterOpts is a slice of default filters
	// for messages to the outbox.
	DefaultOutboxFilterOpts = []FilterOpt{}
	// DefaultInboxFilterOpts is a slice of default filters
	// for messages to the inbox.
	DefaultInboxFilterOpts = []FilterOpt{WithFilterUUID()}
)

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
func (f *filter) Filter(e *pb.Event) (*pb.Event, error) {
	// doing this in reverse to keep order of the middleware
	for i := len(f.funcs) - 1; i >= 0; i-- {
		var err error
		e, err = f.funcs[i](e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

// WithFilterPlugin is filtering an event for the plugin
// as configured by its meta information.
func WithFilterPlugin() FilterOpt {
	return func(p *Plugin) FilterFunc {
		return func(e *pb.Event) (*pb.Event, error) {
			if e.GetPlugin() != nil && e.GetPlugin().GetName() != p.meta.GetName() {
				return nil, nil
			}

			return e, nil
		}
	}
}

// WithFilterUUID is adding a uuid to events.
func WithFilterUUID() FilterOpt {
	return func(p *Plugin) FilterFunc {
		return func(e *pb.Event) (*pb.Event, error) {
			if e != nil {
				uuid, err := uuid.NewRandom()
				if err != nil {
					return nil, err
				}
				e.Uuid = uuid.String()
			}

			return e, nil
		}
	}
}
