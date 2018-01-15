package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var vc *discordgo.VoiceConnection

// JoinVoice makes the bot join a voice channel with the supplied ID
func JoinVoice(gID, cID string) {
	var err error
	if vc, err = dg.ChannelVoiceJoin(gID, cID, true, true); err != nil {
		log.Printf("warning: discordgo.Session.ChannelVoiceJoin: %v", err) // Not a fatal error
	}
}

// LeaveVoice makes the bot leave the voice channel
func LeaveVoice() {
	if vc != nil {
		if err := vc.Disconnect(); err != nil {
			log.Printf("warning: discordgo.VoiceConnection.Disconnect: %v", err) // Not a fatal error
		}
	}
}
