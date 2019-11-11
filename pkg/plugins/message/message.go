package message

import (
	"sync"

	pb "github.com/andersnormal/autobot/proto"
)

// Message is the message to be send in the inbox or outbox of autobot.
type Message struct {
	// Payload is the messages payload.
	Payload *pb.Message

	// ack is closed when the message has been acknowledged.
	ack chan struct{}

	// nack is closed when the message is not transmitted.
	nack chan struct{}

	ackSentType ackType
	sync.Mutex
}

// New is creating a new message
func New(payload *pb.Message) *Message {
	return &Message{
		Payload: payload,
		ack:     make(chan struct{}),
		nack:    make(chan struct{}),
	}
}

type ackType int

const (
	unknown ackType = iota
	ack
	nack
)

// Ack sends message acknowledgement.
//
// Ack is not blocking.
func (m *Message) Ack() bool {
	m.Lock()
	defer m.Unlock()

	if m.ackSentType == nack {
		return false
	}

	if m.ackSentType != unknown {
		return true
	}

	m.ackSentType = ack

	close(m.ack)

	return true
}

// Acked returns a channel that is closed when the message is ack'ed.
func (m *Message) Acked() <-chan struct{} {
	return m.ack
}

// Nack sends a negative acknowledgement
func (m *Message) Nack() bool {
	m.Lock()
	defer m.Unlock()

	if m.ackSentType == ack {
		return false
	}

	if m.ackSentType != unknown {
		return true
	}

	m.ackSentType = nack

	close(m.nack)

	return true
}

// Nacked returns a channel that is closed thwn the message is nack'ed.
func (m *Message) Nacked() <-chan struct{} {
	return m.nack
}
