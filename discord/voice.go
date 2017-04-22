package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var vc *discordgo.VoiceConnection

func JoinVoice(gID, cID string) {
	var err error
	if vc, err = dg.ChannelVoiceJoin(gID, cID, true, true); err != nil {
		log.Printf("warning: %v", err) // Not a fatal error
	}
}

func LeaveVoice() {
	if vc != nil {
		if err := vc.Disconnect(); err != nil {
			log.Printf("warning: %v", err) // Not a fatal error
		}
	}
}
