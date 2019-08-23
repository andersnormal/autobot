package proto

import (
	proto "github.com/golang/protobuf/proto"
)

// Reply is automatically creating a reply message event
// from the proto.Message or proto.Message_Private.
func (m *Message) Reply(text string) *Event {
	reply := proto.Clone(m).(*Message)
	reply.Text = text

	return &Event{
		Event: &Event_Reply{
			Reply: reply,
		},
	}
}
