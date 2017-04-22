package discord

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/discord/bot_reactions"
	"github.com/adayoung/ada-bot/settings"
)

var BotID string // Set by InitDiscordSession
var dg *discordgo.Session

func init() {
	rand.Seed(time.Now().Unix())
}

// Create a new Discord session using the provided login information.
func InitDiscordSession(token string, qLength int, waitMs string) error {
	var err error
	if dg, err = discordgo.New(fmt.Sprintf("Bot %s", token)); err == nil {
		if u, err := dg.User("@me"); err == nil {
			BotID = u.ID

			dg.AddHandler(ready)
			// Add handlers for messages received
			dg.AddHandler(messageCreate)
			dg.AddHandler(messageUpdate)
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

	if _waitMs, err := time.ParseDuration(waitMs); err == nil {
		messageQueue = make(chan message, qLength)
		rateLimit = time.NewTicker(_waitMs)
		go dispatchMessages()
	} else {
		return err
	}
	return nil
}

// Queue a message for posting via Discord API, takes channelID and message
func PostMessage(c string, m string) {
	if len(m) > 0 {
		mq := message{ChannelID: c, Message: m}
		messageQueue <- mq
	}
}

// Close Discord Gateway websocket on exit
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
	if err := s.UpdateStatus(0, "play.achaea.com"); err != nil {
		log.Printf("warning: %v", err) // Not a fatal error
	}
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	go _botReactions(s, m.Message, true)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go _botReactions(s, m.Message, false)
}

func _botReactions(s *discordgo.Session, m *discordgo.Message, update bool) {
	if m.Author == nil { // iopred - Yeah, Author isn't guaranteed to be non nil
		return // iopred - @Ada, in _botReactions you really should just do <-- that
	}

	if m.Author.ID == BotID { // ignore the bot's own messages from processing
		return
	}

	if channel, err := s.State.Channel(m.ChannelID); err == nil {
		guildID := channel.GuildID

		if guildID == "" { // log direct messages sent to the bot
			fmt.Printf("Message received from %s: %s\n", m.Author.Username, m.Content)
			_postReactions(m, &discordgo.Member{}, update)
			return
		}

		if member, err := s.State.Member(guildID, m.Author.ID); err == nil {

			if m.Author.ID == settings.Settings.Discord.BotAdmin {
				if strings.HasPrefix(m.Content, "!join") {
					vid := strings.TrimSpace(m.Content[5:])
					if len(vid) > 0 {
						JoinVoice(guildID, vid)
					}
				} else if strings.HasPrefix(m.Content, "!leave") {
					LeaveVoice()
				}
			}

			_postReactions(m, member, update)
		} else {
			log.Printf("warning: %v", err) // Non-fatal error at s.State.Member() call
		}
	} else {
		log.Printf("warning: %v", err) // Non-fatal error at s.State.Channel() call
	}
}

func _postReactions(m *discordgo.Message, member *discordgo.Member, update bool) {
	for _, reaction := range bot_reactions.GetReactions(m, member, update) {
		PostMessage(m.ChannelID, reaction)
	}
}
