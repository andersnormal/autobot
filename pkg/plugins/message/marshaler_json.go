package message

import (
	"encoding/json"
)

// JSONMarshaler ...
type JSONMarshaler struct {
	NewUUID func() string
}

// Marshal ...
func (m JSONMarshaler) Marshal(v interface{}) (*Message, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg := New(m.NewUUID(), b)

	return msg, nil
}

// Unmarshal ...
func (JSONMarshaler) Unmarshal(msg *Message, v interface{}) (err error) {
	return json.Unmarshal(msg.Payload, v)
}
