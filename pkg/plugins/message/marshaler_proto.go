package message

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// Marshaller ...
type Marshaler interface {
	Marshal(interface{}) (*Message, error)
	Unmarshal(*Message, interface{}) error
}

// ProtoMarshaler ...
type ProtoMarshaler struct {
	NewUUID func() string
}

// Marshal ...
func (m ProtoMarshaler) Marshal(v interface{}) (*Message, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("autobot: interface must be a pointer")
	}

	protoMsg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("autobot: interface must be conform to proto.Message")
	}

	b, err := proto.Marshal(protoMsg)
	if err != nil {
		return nil, err
	}

	msg := New(m.NewUUID(), b)

	return msg, nil
}

// Unmarshal ...
func (ProtoMarshaler) Unmarshal(msg *Message, v interface{}) error {
	return proto.Unmarshal(msg.Payload, v.(proto.Message))
}
