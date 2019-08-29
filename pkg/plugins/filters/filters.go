package filters

import (
	pb "github.com/andersnormal/autobot/proto"

	"github.com/pkg/errors"
)

// Filter is operating on the events flowing in and out
// of the inbox /outbox via the subscriptions or publications
// to the message queue.
type Filter interface {
	Filter(*pb.Message) (*pb.Message, error)
}

type filter struct {
	funcs []FilterFunc
}

// FilterFunc is the wrapper for the plugin filtering option.
type FilterFunc func(e *pb.Message) (*pb.Message, error)

var (
	// DefaultOutboxFilterOpts is a slice of default filters
	// for messages to the outbox.
	DefaultOutboxFilterOpts = []FilterFunc{}
	// DefaultInboxFilterOpts is a slice of default filters
	// for messages to the inbox.
	DefaultInboxFilterOpts = []FilterFunc{}
)

// New is creating a new filter.
// This is mostly used internally.
func New(funcs ...FilterFunc) Filter {
	f := new(filter)
	f.funcs = append(f.funcs, funcs...)

	return f
}

// Filter is handeling an inflowing message according to the middleware.
func (f *filter) Filter(e *pb.Message) (*pb.Message, error) {
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
