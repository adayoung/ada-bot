package botReactions

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/storage"
)

type logger struct {
	Trigger string
}

func (l *logger) Help() string {
	return ""
}

func (l *logger) HelpDetail() string {
	return l.Help()
}

func (l *logger) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	saveMessage(m, a, mType)
	return Reaction{Text: ""} // We don't talk, we just listen -sagenod-
}

func init() {
	storage.OnReady(initDB)

	_logger := &logger{
		Trigger: "*",
	}
	addReaction(_logger.Trigger, "CREATE", _logger)
	addReaction(_logger.Trigger, "UPDATE", _logger)
	addReaction(_logger.Trigger, "DELETE", _logger)
}

var initDBComplete bool = false

func initDB() {
	sqlTable := `
		CREATE TABLE IF NOT EXISTS "discord_messages" (
			"id" serial NOT NULL PRIMARY KEY,
			"message_id" varchar(24) NOT NULL UNIQUE,
			"channel_id" varchar(24) NOT NULL,
			"guild_id" varchar(24) NOT NULL,
			"content" varchar(2000) NOT NULL,
			"timestamp" timestamp NOT NULL,
			"user_id" varchar(24) NOT NULL,
			"member" varchar(32) NOT NULL,
			"bot_command" boolean DEFAULT false
		);
	`
	if _, err := storage.DB.Exec(sqlTable); err == nil {
		initDBComplete = true
	} else {
		log.Printf("error: %v", err) // We won't store messages, that's what!
	}
}

func saveMessage(m *discordgo.Message, member *discordgo.Member, mType string) {
	if !initDBComplete {
		return // We're not ready to save events
	}

	if m.Author != nil {
		if m.Author.Bot {
			return // Do not log messages posted by other bots
		}
	}

	var _member string
	var _guildID string

	if m.Author != nil {
		if member.GuildID != "" {
			_guildID = member.GuildID
			if member.Nick != "" {
				_member = member.Nick
			} else {
				_member = member.User.Username
			}
		} else {
			_member = m.Author.Username
		}
	}

	if mType == "CREATE" { // New Message
		if timestamp, err := m.Timestamp.Parse(); err == nil {
			message := `INSERT INTO discord_messages (
				message_id, channel_id, guild_id, content,
				timestamp, user_id, member, bot_command
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`
			message = storage.DB.Rebind(message)
			if _, err := storage.DB.Exec(message, m.ID, m.ChannelID, _guildID,
				m.Content, timestamp, m.Author.ID, _member,
				strings.HasPrefix(m.Content, settings.Settings.Discord.BotPrefix)); err != nil {
				log.Printf("error: %v", err) // We won't store messages, that's what!
			}
		} else {
			log.Printf("error: %v", err) // Error at m.Timestamp.Parse() call
		}
	}

	if mType == "UPDATE" { // Updated Message
		if timestamp, err := m.EditedTimestamp.Parse(); err == nil {
			message := "UPDATE discord_messages SET content=?, timestamp=? WHERE message_id=?"
			message = storage.DB.Rebind(message)
			if _, err := storage.DB.Exec(message, m.Content, timestamp, m.ID); err != nil {
				log.Printf("error: %v", err) // Oops, something wrong with updating message
			}
		} else {
			log.Printf("error: %v", err) // Error at m.Timestamp.Parse() call
		}
	}

	if mType == "DELETE" { // Deleted Message
		message := "DELETE FROM discord_messages WHERE message_id=?"
		message = storage.DB.Rebind(message)
		if _, err := storage.DB.Exec(message, m.ID); err != nil {
			log.Printf("error: %v", err) // Oops, something wrong with deleting message
		}
		return
	}
}
