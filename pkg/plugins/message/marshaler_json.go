package message

import (
	cloudevents "github.com/cloudevents/sdk-go"
)

// JSONMarshaler ...
type JSONMarshaler struct {
	NewUUID func() string
}

// Marshal ...
func (m JSONMarshaler) Marshal(v interface{}) (cloudevents.Event, error) {
	// this is converting to a cloud event
	e := cloudevents.NewEvent()

	e.SetID(m.NewUUID())
	e.SetType("us.andersnormal.autobot.message")
	e.SetSource("github.com/andersnormal/autobot/pkg/plugins/message")
	e.SetDataContentType(cloudevents.ApplicationJSON)

	if err := e.SetData(v); err != nil {
		return e, err
	}

	return e, nil
}

// Unmarshal ...
func (JSONMarshaler) Unmarshal(e cloudevents.Event, v interface{}) (err error) {
	return e.DataAs(v)
}
