package bot_reactions

import (
	"fmt"
	"log"
	"regexp"
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

var user_regexp *regexp.Regexp = regexp.MustCompile("<@!?([0-9]+)>")

func (r *RandomQ) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	if m.Content == fmt.Sprintf("%s%s", settings.Settings.Discord.BotPrefix, r.Trigger) {
		return randomQuote(m.ChannelID, nil)
	} else {
		user_id := user_regexp.FindStringSubmatch(m.Content[len(settings.Settings.Discord.BotPrefix)+len(r.Trigger)+1:])
		if user_id != nil {
			return randomQuote(m.ChannelID, &user_id[1])
		}
	}
	return ""
}

func init() {
	randomq := &RandomQ{
		Trigger: "random",
	}
	addReaction(randomq.Trigger, randomq)
}

func randomQuote(channelID string, user *string) string {
	query := "SELECT user_id, content, timestamp from discord_messages WHERE channel_id=?"
	if user != nil {
		query = fmt.Sprintf("%s AND user_id=?", query)
	}
	query = fmt.Sprintf("%s AND content NOT LIKE '%s%%'", query, settings.Settings.Discord.BotPrefix)
	query = fmt.Sprintf("%s ORDER BY random() LIMIT 1", query)
	query = storage.DB.Rebind(query)

	var user_id string
	var content string
	var timestamp time.Time
	var result bool = true

	if user != nil {
		if err := storage.DB.QueryRow(query, channelID, user).Scan(&user_id, &content, &timestamp); err != nil {
			result = false
			log.Printf("error: %v", err) // Error with ... something
		}
	} else {
		if err := storage.DB.QueryRow(query, channelID).Scan(&user_id, &content, &timestamp); err != nil {
			result = false
			log.Printf("error: %v", err) // Error with ... something
		}
	}

	if result == true {
		return fmt.Sprintf("%s\n\t -- <@%s> on %s", content, user_id, timestamp.Format("Monday, Jan _2, 2006"))
	}
	return fmt.Sprintf("No quotes retrieved for <@%s>", *user)
}
