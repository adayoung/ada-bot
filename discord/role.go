package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
)

func setRole(s *discordgo.Session, m *discordgo.Message, guildID string, roleName string) {
	// 1. Locate roleID on server
	if guild, err := s.Guild(guildID); err == nil {
		roleID := "" // init?
		if guild.Roles != nil {
			for _, gRole := range guild.Roles {
				if gRole.Name == roleName {
					roleID = gRole.ID
				}
			}
		}

		if len(roleID) > 0 {
			// 2. We found a role! Let's go through members and apply it
			if guild.Members != nil {
				gRoleCounter := 0
				for _, gMember := range guild.Members {
					if gMember.User != nil {
						if err := s.GuildMemberRoleAdd(guild.ID, gMember.User.ID, roleID); err == nil {
							gRoleCounter += 1
						} else {
							log.Printf("warning: %v", err) // Non-fatal error at s.GuildMemberRoleAdd() call
						}
					}
				}
				PostMessage(m.ChannelID, fmt.Sprintf("I've set the role of %d users to %s! :dancer:", gRoleCounter, roleName))
			}
			// 3. Make that role persistent for this guild
			settings.Settings.Discord.DefaultRoles[guild.ID] = roleID
			settings.Settings.Save()
		} else {
			PostMessage(m.ChannelID, "I couldn't locate a role with that name :shrug:")
		}
	} else {
		log.Printf("warning: %v", err) // Non-fatal error at s.Guild() call
	}
}
