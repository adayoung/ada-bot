package botReactions

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/ire"
	"github.com/adayoung/ada-bot/settings"
)

type whois struct {
	Trigger string
}

func (w *whois) Help() string {
	return "Lookup <name> in game and report findings."
}

func (w *whois) HelpDetail() string {
	return w.Help()
}

func (w *whois) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	var response string
	rPlayer := strings.ToLower(strings.TrimSpace(m.Content[len(settings.Settings.Discord.BotPrefix)+len(w.Trigger):]))
	if gPlayer, err := ire.GetPlayer(rPlayer); err == nil {
		response = fmt.Sprintf("```%s```", gPlayer)
		return Reaction{Text: response}
	} else {
		log.Printf("error: %v", err) // Not a fatal error
	}
	response = fmt.Sprintf("Oops, I couldn't find %s :frowning:", rPlayer)
	return Reaction{Text: response}
}

func init() {
	_whois := &whois{
		Trigger: "whois",
	}
	addReaction(_whois.Trigger, "CREATE", _whois)
}
