package bot_reactions

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/storage"
)

type Logger struct {
	Trigger string
}

func (l *Logger) Help() string {
	return ""
}

func (l *Logger) HelpDetail(m *discordgo.Message) string {
	return l.Help()
}

func (l *Logger) Reaction(m *discordgo.Message, a *discordgo.Member, u bool) string {
	saveMessage(m, a, u)
	return "" // We don't talk, we just listen -sagenod-
}

func init() {
	storage.OnReady(initDB)

	logger := &Logger{
		Trigger: "*",
	}
	addReaction(logger.Trigger, logger)
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
		// initDBComplete = true
	} else {
		log.Printf("error: %v", err) // We won't store messages, that's what!
	}

	sqlTable = `
		CREATE TABLE IF NOT EXISTS "discord_messages_updates" (
			"id" serial NOT NULL PRIMARY KEY,
			"message_id" varchar(24) NOT NULL,
			"content" varchar(2000) NOT NULL,
			"timestamp" timestamp NOT NULL
		);
	`
	if _, err := storage.DB.Exec(sqlTable); err == nil {
		initDBComplete = true
	} else {
		log.Printf("error: %v", err) // We won't store messages, that's what!
	}

	// TODO: Add foreign key constraint to enforce:
	// discord_messages_updates.message_id -> discord_messages.message_id
}

func saveMessage(m *discordgo.Message, member *discordgo.Member, u bool) {
	if !initDBComplete {
		return // We're not ready to save events
	}

	if m.Author.Bot {
		return // Do not log messages posted by other bots
	}

	var _member string
	var _guildID string
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

	if timestamp, err := m.Timestamp.Parse(); err == nil {
		if !u { // New Message
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
		} else { // Updated Message
			message := `INSERT INTO discord_messages_updates (
				message_id, content, timestamp) VALUES (?, ?, ?)`
			message = storage.DB.Rebind(message)
			if _, err := storage.DB.Exec(message, m.ID, m.Content, timestamp); err != nil {
				log.Printf("error: %v", err) // We won't store messages, that's what!
			}
		}
	} else {
		log.Printf("error: %v", err) // Error at m.Timestamp.Parse() call
	}
}
