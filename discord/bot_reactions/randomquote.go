package bot_reactions

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/storage"
)

type RandomQ struct {
	Trigger string
}

func (r *RandomQ) Help() string {
	return "Random Quote!"
}

func (r *RandomQ) HelpDetail(m *discordgo.Message) string {
	return r.Help()
}

var userRegexp *regexp.Regexp = regexp.MustCompile("<@!?([0-9]+)>")

func (r *RandomQ) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) string {
	if a.GuildID == "" {
		return fmt.Sprintf("Oops, %srandom is not available on direct messages :ghost:", settings.Settings.Discord.BotPrefix)
	}

	if m.Content == fmt.Sprintf("%s%s", settings.Settings.Discord.BotPrefix, r.Trigger) {
		return randomQuote(a.GuildID, nil, nil)
	} else {
		request := strings.TrimSpace(m.Content[len(settings.Settings.Discord.BotPrefix)+len(r.Trigger):])
		userID := userRegexp.FindStringSubmatch(request)
		if userID != nil {
			return randomQuote(a.GuildID, &userID[1], nil)
		} else {
			return randomQuote(a.GuildID, nil, &request)
		}
	}
}

func init() {
	randomq := &RandomQ{
		Trigger: "random",
	}
	addReaction(randomq.Trigger, "CREATE", randomq)
}

func randomQuote(guildID string, user *string, member *string) string {
	query := "SELECT member, content, channel_id, timestamp from discord_messages WHERE guild_id=?"
	if user != nil {
		query = fmt.Sprintf("%s AND user_id=?", query)
	} else if member != nil {
		query = fmt.Sprintf("%s AND member ILIKE ?", query)
	}

	query = fmt.Sprintf("%s AND bot_command = 'false'", query)
	query = fmt.Sprintf("%s AND character_length(content) > 6", query)
	query = fmt.Sprintf("%s ORDER BY random() LIMIT 1", query)
	query = storage.DB.Rebind(query)

	var userID string
	var content string
	var channelID string
	var timestamp time.Time
	var result bool = true

	if user != nil {
		if err := storage.DB.QueryRow(query, guildID, user).Scan(&userID, &content, &channelID, &timestamp); err != nil {
			result = false
			log.Printf("error: %v", err) // Error with ... something
		}
	} else if member != nil {
		if err := storage.DB.QueryRow(query, guildID, member).Scan(&userID, &content, &channelID, &timestamp); err != nil {
			result = false
			log.Printf("error: %v", err) // Error with ... something
		}
	} else {
		if err := storage.DB.QueryRow(query, guildID).Scan(&userID, &content, &channelID, &timestamp); err != nil {
			result = false
			log.Printf("error: %v", err) // Error with ... something
		}
	}

	if result == true && len(content) > 0 {
		return fmt.Sprintf("%s\n\t -- %s on <#%s> at %s", content, userID, channelID, timestamp.Format("Monday, Jan _2, 2006"))
	}

	if user != nil {
		return fmt.Sprintf("No quotes retrieved for <@%s>", *user)
	} else if member != nil {
		return fmt.Sprintf("No quotes retrieved for %s", *member)
	}

	return ""
}
