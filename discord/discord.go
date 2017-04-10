package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strings"
	"time"
)

var BotID string
var dg *discordgo.Session

func init() {
	rand.Seed(time.Now().Unix())
}

func InitDiscordSession(token string) error {
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

	return nil
}

func PostMessage(c string, m string) {
	_, _ = dg.ChannelMessageSend(c, m)
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
	if m.Author.ID == BotID { // ignore the bot's own messages from processing
		return
	}

	if c, err := s.State.Channel(m.ChannelID); err != nil {
		fmt.Println("Oops, error at getting session.State.Channel,", err)
		return // Not a fatal error
	} else {
		if c.GuildID == "" {
			fmt.Printf("Message received from %s: %s\n", m.Author.Username, m.Content)
			if strings.ToLower(m.Content) == "ping" {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
			}
		}
	}

	if strings.ToLower(m.Content) == "!ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.ToLower(m.Content) == "!pink" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "I love pink!")
	}

	if strings.HasPrefix(strings.ToLower(m.Content), "!decide") {
		choices := strings.Split(m.Content[8:], " or ")
		the_answer := choices[rand.Intn(len(choices))]
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The correct answer is **%s**", the_answer))
	}
}
