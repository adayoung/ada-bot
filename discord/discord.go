package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"time"

	"github.com/adayoung/ada-bot/discord/bot_reactions"
)

var BotID string
var dg *discordgo.Session

func init() {
	rand.Seed(time.Now().Unix())
}

func InitDiscordSession(token string, q_length int, wait_ms string) error {
	// Create a new Discord session using the provided login information.
	var err error
	if dg, err = discordgo.New(fmt.Sprintf("Bot %s", token)); err == nil {
		if u, err := dg.User("@me"); err == nil {
			BotID = u.ID

			dg.AddHandler(ready)
			// Add handlers for messages received
			dg.AddHandler(messageCreate)
			if err := dg.Open(); err == nil {
				fmt.Println("Successfully launched a new Discord session.")
			} else {
				return err // Error at opening Discord Session
			}
		} else {
			return err // Error at obtaining account details
		}
	} else {
		return err // Error at creating a new Discord session
	}

	if _wait_ms, err := time.ParseDuration(wait_ms); err == nil {
		messageQueue = make(chan message, q_length)
		rateLimit = time.NewTicker(_wait_ms)
		go dispatchMessages()
	} else {
		return err
	}
	return nil
}

func PostMessage(c string, m string) {
	if len(m) > 0 {
		mq := message{ChannelID: c, Message: m}
		messageQueue <- mq
	}
}

func CloseDiscordSession() {
	dg.Close()
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	if guilds, err := s.UserGuilds(); err != nil {
		fmt.Println("ERROR: We couldn't get UserGuilds")
		log.Fatalf("error: %v", err)
	} else {
		for index, guild := range guilds {
			fmt.Printf("[%d] ------------------------------\n", index)
			fmt.Println("Guild ID: ", guild.ID)
			fmt.Println("Guild Name: ", guild.Name)
			fmt.Println("Guild Permissions: ", guild.Permissions)

		}
		fmt.Println("----------------------------------")
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go _botReactions(s, m)
}

func _botReactions(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID { // ignore the bot's own messages from processing
		return
	}

	if channel, err := s.State.Channel(m.ChannelID); err == nil {
		guildID := channel.GuildID

		if guildID == "" { // log direct messages sent to the bot
			fmt.Printf("Message received from %s: %s\n", m.Author.Username, m.Content)
			_postReactions(m.Message, &discordgo.Member{})
			return
		}

		if member, err := s.State.Member(guildID, m.Author.ID); err == nil {
			_postReactions(m.Message, member)
		} else {
			log.Printf("warning: %v", err) // Non-fatal error at s.State.Member() call
		}
	} else {
		log.Printf("warning: %v", err) // Non-fatal error at s.State.Channel() call
	}
}

func _postReactions(m *discordgo.Message, member *discordgo.Member) {
	for _, reaction := range bot_reactions.GetReactions(m, member) {
		PostMessage(m.ChannelID, reaction)
	}
}
