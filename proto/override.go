package proto

import (
	"github.com/golang/protobuf/proto"
)

// Reply is producing a reply from a message
func (m *Message) Reply(text string) *Message {
	msg := proto.Clone(m).(*Message)
	msg.Text = text

	return msg
}
