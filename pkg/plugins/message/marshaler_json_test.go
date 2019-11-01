package message_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/andersnormal/autobot/pkg/plugins/message"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestJSONMarshaler_Marshal(t *testing.T) {
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

			assert.Equal(b, msg.Data)
			assert.Equal(msg.ID(), genUUID())
		})
	}
}

func TestJSONMarshaler_Unmarshal(t *testing.T) {
	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	tests := []struct {
		desc string
		uuid func() string
		out  interface{}
		in   func() cloudevents.Event
	}{
		{
			desc: "unmarshal string",
			uuid: genUUID,
			in: func() cloudevents.Event {
				e := cloudevents.NewEvent()

				e.SetID(genUUID())
				e.SetData("foo")
				e.SetDataContentType(cloudevents.ApplicationJSON)

				return e
			},
			out: "foo",
		},
		{
			desc: "unmarshal float64",
			uuid: genUUID,
			in: func() cloudevents.Event {
				e := cloudevents.NewEvent()

				e.SetID(genUUID())
				e.SetData(12345)
				e.SetDataContentType(cloudevents.ApplicationJSON)

				return e
			},
			out: float64(12345),
		},
		{
			desc: "unmarshal bool",
			uuid: genUUID,
			in: func() cloudevents.Event {
				e := cloudevents.NewEvent()

				e.SetID(genUUID())
				e.SetData(true)
				e.SetDataContentType(cloudevents.ApplicationJSON)

				return e
			},
			out: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			assert := assert.New(t)

			m := JSONMarshaler{NewUUID: tt.uuid}

			var v interface{}
			err := m.Unmarshal(tt.in(), &v)
			assert.NoError(err)

			assert.Equal(tt.out, v)
		})
	}
}
