package plugins

import (
	"context"
	"testing"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/stretchr/testify/assert"
)

func TestContext_Context(t *testing.T) {
	assert := assert.New(t)

	cbCtx := &cbContext{
		ctx: context.Background(),
	}

	assert.NotNil(cbCtx.Context())
}

func TestContext_Message(t *testing.T) {
	assert := assert.New(t)

	cbCtx := &cbContext{
		ctx: context.Background(),
		msg: &pb.Message{},
	}

	assert.NotNil(cbCtx.Message())
}
