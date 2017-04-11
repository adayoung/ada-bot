package discord

import (
	"time"
)

type message struct {
	ChannelID string
	Message   string
}

type messageQueue struct {
	Active   bool
	Messages []message
}

var mq messageQueue
var mqc chan bool

func dispatchMessages() {
	var m message
	for _ = range(mqc) {
		for {
			if len(mq.Messages) == 0 {
				break
			}
			m, mq.Messages = mq.Messages[0], mq.Messages[1:]
			_, _ = dg.ChannelMessageSend(m.ChannelID, m.Message)
			time.Sleep(time.Millisecond * 500)
		}
	}
}
