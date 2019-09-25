package plugins

import (
	"context"
	"sync"

	pb "github.com/andersnormal/autobot/proto"
)

var _ Context = (*cbContext)(nil)

// Context ...
type Context interface {
	// Message ...
	Message() *pb.Message
	// Send ...
	Send(*pb.Message)
	// Context ...
	Context() context.Context
}

type cbContext struct {
	ctx context.Context

	send chan<- *pb.Message
	msg  *pb.Message

	sync.Mutex
}

// Message ...
func (ctx *cbContext) Message() *pb.Message {
	return ctx.msg
}

// Send ...
func (ctx *cbContext) Send(msg *pb.Message) {
	ctx.send <- msg
}

// SendAsync ...
func (ctx *cbContext) SendAsync(msg *pb.Message) {
	go func() {
		ctx.send <- msg
	}()
}

// Context ...
func (ctx *cbContext) Context() context.Context {
	return ctx.ctx
}
