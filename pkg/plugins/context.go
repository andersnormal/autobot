package plugins

import (
	"context"
	"encoding/json"
	"sync"

	pb "github.com/andersnormal/autobot/proto"
)

var _ Context = (*cbContext)(nil)

// Context ...
type Context interface {
	// Message ...
	Message() *pb.Message
	// Send ...
	Send(*pb.Message) error
	// Context ...
	Context() context.Context
}

type cbContext struct {
	ctx context.Context

	plugin *Plugin
	msg    *pb.Message

	sync.Mutex
}

// Message ...
func (ctx *cbContext) Message() *pb.Message {
	return ctx.msg
}

// Send ...
func (ctx *cbContext) Send(msg *pb.Message) error {
	sc, err := ctx.plugin.getConn()
	if err != nil {
		return err
	}

	m, err := ctx.plugin.marshaler.Marshal(msg)
	if err != nil {
		return err
	}

	// add some metadata
	m.Metadata.Src(ctx.plugin.runtime.Name)

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := sc.Publish(ctx.plugin.runtime.Inbox, b); err != nil {
		return err
	}

	return nil
}

// Context ...
func (ctx *cbContext) Context() context.Context {
	return ctx.ctx
}
