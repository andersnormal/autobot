package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage_New(t *testing.T) {
	assert := assert.New(t)

	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	msg := New(genUUID(), Payload([]byte("foo")))

	assert.NotNil(msg)
	assert.Equal(genUUID(), msg.UUID)
	assert.Equal(Payload([]byte("foo")), msg.Payload)
}

func TestMessage_Ack(t *testing.T) {
	assert := assert.New(t)

	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	msg := New(genUUID(), Payload([]byte("foo")))

	assert.NotNil(msg)

	ack := msg.Ack()

	assert.True(ack)

	_, ok := (<-msg.ack)

	assert.False(ok)
}

func TestMessage_Acked(t *testing.T) {
	assert := assert.New(t)

	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	msg := New(genUUID(), Payload([]byte("foo")))

	assert.NotNil(msg)

	testChan := make(chan struct{})

	go func() {
		select {
		case _, ok := <-msg.Acked():
			assert.False(ok)

			close(testChan)
		case <-time.After(1 * time.Second):
			assert.True(false)
		}
	}()

	msg.Ack()

	<-testChan
}

func TestMessage_Nack(t *testing.T) {
	assert := assert.New(t)

	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	msg := New(genUUID(), Payload([]byte("foo")))

	assert.NotNil(msg)

	nack := msg.Nack()

	assert.True(nack)

	_, ok := (<-msg.nack)

	assert.False(ok)
}

func TestMessage_Nacked(t *testing.T) {
	assert := assert.New(t)

	genUUID := func() string {
		return "00000000-0000-0000-0000-000000000000"
	}

	msg := New(genUUID(), Payload([]byte("foo")))

	assert.NotNil(msg)

	testChan := make(chan struct{})

	go func() {
		select {
		case _, ok := <-msg.Nacked():
			assert.False(ok)

			close(testChan)
		case <-time.After(1 * time.Second):
			assert.True(false)
		}
	}()

	msg.Nack()

	<-testChan
}
