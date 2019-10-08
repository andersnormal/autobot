package message_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/andersnormal/autobot/pkg/plugins/message"
	pb "github.com/andersnormal/autobot/proto"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestProtoMarshaler_Marshal(t *testing.T) {
	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	tests := []struct {
		desc string
		uuid func() string
		in   interface{}
		err  error
	}{
		{
			desc: "create new message",
			uuid: genUUID,
			in:   &pb.Message{Text: "test"},
		},
		{
			desc: "not a pointer",
			uuid: genUUID,
			in:   pb.Message{Text: "test"},
			err:  errors.New("autobot: interface must be a pointer"),
		},
		{
			desc: "not conform to proto.Message",
			uuid: genUUID,
			in:   &Message{},
			err:  errors.New("autobot: interface must be conform to proto.Message"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			assert := assert.New(t)

			var b []byte
			var err error
			if tt.err == nil {
				b, err = proto.Marshal(tt.in.(proto.Message))
				assert.NoError(err)
			}

			m := ProtoMarshaler{NewUUID: tt.uuid}

			msg, err := m.Marshal(tt.in)

			if tt.err != nil {
				assert.Error(err)
				assert.EqualError(err, tt.err.Error())

				return
			}

			assert.NoError(err)

			assert.Equal(Payload(b), msg.Payload)
			assert.Equal(msg.UUID, genUUID())
		})
	}
}

func TestProtoMarshaler_Unmarshal(t *testing.T) {
	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	msg := &pb.Message{Text: "foo"}
	b, _ := proto.Marshal(msg)

	tests := []struct {
		desc string
		uuid func() string
		out  *pb.Message
		in   *Message
	}{
		{
			desc: "unmarshal proto.Message",
			uuid: genUUID,
			in:   &Message{UUID: genUUID(), Payload: b},
			out:  msg,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			assert := assert.New(t)

			m := ProtoMarshaler{NewUUID: tt.uuid}

			var v pb.Message
			err := m.Unmarshal(tt.in, &v)
			assert.NoError(err)

			assert.Equal(tt.out.GetText(), v.GetText())
		})
	}
}
