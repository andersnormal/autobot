package main

import (
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nlopes/slack"
)

// FromMsg ...
func FromMsg(msg *slack.MessageEvent) *pb.Event {
	m := new(pb.Event)

	e := &pb.Event_Message{
		Message: &pb.MessageEvent{
			Text:     msg.Text,
			Channel:  msg.Channel,
			User:     msg.User,
			Username: msg.Username,
			Topic:    msg.Topic,
		},
	}

	m.Event = e

	return m
}

// FromMessageEvent ...
func FromMessageEvent(rtm *slack.RTM, e *pb.MessageEvent) *slack.OutgoingMessage {
	return rtm.NewOutgoingMessage(e.GetText(), e.GetChannel())
}
