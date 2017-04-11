package discord

import (
	"time"
)

type message struct {
	ChannelID string
	Message   string
}

var messageQueue chan message

func dispatchMessages() {
	var m message
	for m = range messageQueue {
		_, _ = dg.ChannelMessageSend(m.ChannelID, m.Message)
		time.Sleep(time.Millisecond * 500)
	}
}
