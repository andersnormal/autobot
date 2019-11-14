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

// ThreadedReply is producing a reply from a parent message to thread
func (m *Message) ThreadedReply(text string) *Message {
	msg := proto.Clone(m).(*Message)
	msg.Text = text
	msg.ReplyInThread = true

	return msg
}
