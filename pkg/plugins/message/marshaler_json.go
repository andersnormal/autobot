package message

import (
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go"
)

// Marshaler ...
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// JSONMarshaler ...
type JSONMarshaler struct {
	NewUUID func() string
}

// Marshal ...
func (m JSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	// this is converting to a cloud event
	e := cloudevents.NewEvent()

	e.SetID(m.NewUUID())
	e.SetType("us.andersnormal.autobot.message")
	e.SetSource("github.com/andersnormal/autobot/pkg/plugins/message")
	e.SetDataContentType(cloudevents.ApplicationJSON)

	if err := e.SetData(v); err != nil {
		return nil, err
	}

	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Unmarshal ...
func (JSONMarshaler) Unmarshal(b []byte, v interface{}) error {
	e := cloudevents.NewEvent()

	err := json.Unmarshal(b, &e)
	if err != nil {
		return err
	}

	err = e.DataAs(v)
	if err != nil {
		return err
	}

	return nil
}
