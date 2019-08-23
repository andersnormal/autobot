package filters

import (
	pb "github.com/andersnormal/autobot/proto"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Filter is operating on the events flowing in and out
// of the inbox /outbox via the subscriptions or publications
// to the message queue.
type Filter interface {
	Filter(*pb.Event) (*pb.Event, error)
}

type filter struct {
	funcs []FilterFunc
}

// FilterFunc is the wrapper for the plugin filtering option.
type FilterFunc func(e *pb.Event) (*pb.Event, error)

var (
	// DefaultOutboxFilterOpts is a slice of default filters
	// for messages to the outbox.
	DefaultOutboxFilterOpts = []FilterFunc{}
	// DefaultInboxFilterOpts is a slice of default filters
	// for messages to the inbox.
	DefaultInboxFilterOpts = []FilterFunc{WithFilterUUID()}
)

// New is creating a new filter.
// This is mostly used internally.
func New(funcs ...FilterFunc) Filter {
	f := new(filter)
	f.funcs = append(f.funcs, funcs...)

	return f
}

// Filter is handeling an inflowing message according to the middleware.
func (f *filter) Filter(e *pb.Event) (*pb.Event, error) {
	// doing this in reverse to keep order of the middleware.
	for i := len(f.funcs) - 1; i >= 0; i-- {
		var err error
		// we track the error from the relevant func.
		// the last that errors and the error gets passed along
		e, err = f.funcs[i](e)
		if err != nil {
			return nil, errors.Wrap(errors.Cause(err), "filter")
		}
	}

	return e, nil
}

// WithFilterPlugin is filtering an event for the plugin
// as configured by its meta information.
func WithFilterPlugin(name string) FilterFunc {
	return func(e *pb.Event) (*pb.Event, error) {
		name := name

		if e.GetPlugin() != nil && e.GetPlugin().GetName() != name {
			return nil, nil
		}

		return e, nil
	}
}

// WithFilterUUID is adding a uuid to events.
func WithFilterUUID() FilterFunc {
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
