package bot_reactions

import (
	"log"

	"github.com/bwmarrin/discordgo"

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

func (l *Logger) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	saveMessage(*m)
	return ""
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
	sql_table := `
    CREATE TABLE IF NOT EXISTS "discord_messages" (
      "id" serial NOT NULL PRIMARY KEY,
      "message_id" varchar(24) NOT NULL UNIQUE,
      "channel_id" varchar(24) NOT NULL,
      "content" varchar(2000) NOT NULL,
      "timestamp" timestamp NOT NULL,
      "user_id" varchar(24) NOT NULL
    );
  `
	if _, err := storage.DB.Exec(sql_table); err == nil {
		initDBComplete = true
	} else {
		log.Printf("error: %v", err) // We won't store messages, that's what!
	}
}

func saveMessage(m discordgo.Message) {
	if !initDBComplete {
		return // We're not ready to save events
	}
	if timestamp, err := m.Timestamp.Parse(); err == nil {
		message := `INSERT INTO discord_messages (
      message_id, channel_id, content, timestamp, user_id
      ) VALUES (?, ?, ?, ?, ?)
      `
		message = storage.DB.Rebind(message)
		if _, err := storage.DB.Exec(message, m.ID, m.ChannelID,
			m.Content, timestamp, m.Author.ID); err != nil {
			log.Printf("error: %v", err) // We won't store messages, that's what!
		}
	} else {
		log.Printf("error: %v", err) // Error at m.Timestamp.Parse() call
	}
}
