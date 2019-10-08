package message_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/andersnormal/autobot/pkg/plugins/message"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	tests := []struct {
		desc string
		uuid func() string
		in   interface{}
	}{
		{
			desc: "marshal string",
			uuid: genUUID,
			in:   "foo",
		},
		{
			desc: "marshal integer",
			uuid: genUUID,
			in:   123456,
		},
		{
			desc: "marshal bool",
			uuid: genUUID,
			in:   true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			assert := assert.New(t)

			b, err := json.Marshal(tt.in)
			assert.NoError(err)

			m := JSONMarshaler{NewUUID: tt.uuid}

			msg, err := m.Marshal(tt.in)
			assert.NoError(err)

			assert.Equal(Payload(b), msg.Payload)
			assert.Equal(msg.UUID, genUUID())
		})
	}
}

func TestUnmarshal(t *testing.T) {
	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	tests := []struct {
		desc string
		uuid func() string
		out  interface{}
		in   *Message
	}{
		{
			desc: "unmarshal string",
			uuid: genUUID,
			in:   &Message{UUID: genUUID(), Payload: []byte(`"foo"`)},
			out:  "foo",
		},
		{
			desc: "unmarshal float64",
			uuid: genUUID,
			in:   &Message{UUID: genUUID(), Payload: []byte(`12345`)},
			out:  float64(12345),
		},
		{
			desc: "unmarshal bool",
			uuid: genUUID,
			in:   &Message{UUID: genUUID(), Payload: []byte(`true`)},
			out:  true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			assert := assert.New(t)

			m := JSONMarshaler{NewUUID: tt.uuid}

			var v interface{}
			err := m.Unmarshal(tt.in, &v)
			assert.NoError(err)

			assert.Equal(tt.out, v)
		})
	}
}
