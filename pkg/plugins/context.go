package plugins

import (
	"context"
	"sync"

	pb "github.com/andersnormal/autobot/proto"
)

var _ Context = (*cbContext)(nil)

// Context is providing information and operations
// upon the reply execution context.
type Context interface {
	// Message should return the message of the current context.
	Message() *pb.Message
	// Send should send a new message with in the current context.
	Send(*pb.Message) error
	// AsyncSend should asynchronously send a new message in the current context.
	AsyncSend(msg *pb.Message)
	// Context should return the current execution context.
	Context() context.Context
}

type cbContext struct {
	ctx context.Context

	plugin *Plugin
	msg    *pb.Message

	sync.Mutex
}

// Message is returning the emboddied message of the context.
func (ctx *cbContext) Message() *pb.Message {
	return ctx.msg
}

// Send is sending a new message within the message context.
func (ctx *cbContext) Send(msg *pb.Message) error {
	sc, err := ctx.plugin.getConn()
	if err != nil {
		return err
	}

	b, err := ctx.plugin.marshaler.Marshal(msg)
	if err != nil {
		return err
	}

	if err := sc.Publish(ctx.plugin.runtime.Outbox, b); err != nil {
		return err
	}

	return nil
}

// AsyncSend is asynchronously sending a new message in the context
// of this message.
func (ctx *cbContext) AsyncSend(msg *pb.Message) {
	ctx.plugin.run(func() error {
		return ctx.Send(msg)
	})
}

// Context is returning a context that can be used to bound an
// operation to the execution of this message.
func (ctx *cbContext) Context() context.Context {
	return ctx.ctx
}
