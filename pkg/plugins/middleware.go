package plugins

import (
	pb "github.com/andersnormal/autobot/proto"
)

// Middleware is operating on the events flowing in and out
// of the inbox /outbox via the subscriptions or publications
// to the message queue.
type Middleware interface {
	Handle(*pb.Event) *pb.Event
}

type middleware struct {
	funcs []func(*pb.Event) *pb.Event
}

// NewMiddleware is creating a new middleware.
// This is mostly used internally.
func NewMiddleware(p *Plugin, opts ...MiddlewareOpt) Middleware {
	m := new(middleware)

	for _, o := range opts {
		m.funcs = append(m.funcs, o(p))
	}

	return m
}

// Handle is handeling an inflowing message according to the middleware.
func (m *middleware) Handle(e *pb.Event) *pb.Event {
	// doing this in reverse to keep order of the middleware
	for i := len(m.funcs) - 1; i >= 0; i-- {
		e = m.funcs[i](e)
	}

	return e
}

// MiddlewareOpt is an option that can be passed to the middleware functions.
type MiddlewareOpt func(p *Plugin) func(*pb.Event) *pb.Event
