package main

import (
	pb "github.com/andersnormal/autobot/proto"

	"github.com/nlopes/slack"
)

// FromMsg ...
func FromMsg(msg *slack.MessageEvent) *pb.Message {
	m := new(pb.Message)

	e := &pb.Message_Message{
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
