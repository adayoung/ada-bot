package bot_reactions

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/ire"
)

type Whois struct {
	Trigger string
}

func (w *Whois) Help() string {
	return "Lookup <name> in game and report findings."
}

func (w *Whois) HelpDetail(m *discordgo.Message) string {
	return w.Help()
}

func (w *Whois) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	r_player := strings.ToLower(strings.TrimSpace(m.Content[len(w.Trigger)+1:]))
	if g_player, err := ire.GetPlayer(r_player); err == nil {
		return fmt.Sprintf("```%s```", g_player)
	} else {
		log.Printf("error: %v", err) // Not a fatal error
	}
	return fmt.Sprintf("Oops, I couldn't find %s :frowning:", r_player)
}

func init() {
	whois := &Whois{
		Trigger: "whois",
	}
	addReaction(whois.Trigger, whois)
}
