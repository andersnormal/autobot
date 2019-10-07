package proto_test

import (
	"testing"

	. "github.com/andersnormal/autobot/proto"

	"github.com/stretchr/testify/assert"
)

func TestReply(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  "foo",
			out: "bar",
		},
	}

	for _, tt := range tests {
		assert := assert.New(t)

		msg := &Message{
			Text: tt.in,
		}

		reply := msg.Reply(tt.out)

		assert.NotNil(reply)
		assert.Equal(reply.Text, tt.out)
	}
}
