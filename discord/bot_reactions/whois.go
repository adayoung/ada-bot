package bot_reactions

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

func (w *whois) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) string {
	rPlayer := strings.ToLower(strings.TrimSpace(m.Content[len(settings.Settings.Discord.BotPrefix)+len(w.Trigger):]))
	if gPlayer, err := ire.GetPlayer(rPlayer); err == nil {
		return fmt.Sprintf("```%s```", gPlayer)
	} else {
		log.Printf("error: %v", err) // Not a fatal error
	}
	return fmt.Sprintf("Oops, I couldn't find %s :frowning:", rPlayer)
}

func init() {
	_whois := &whois{
		Trigger: "whois",
	}
	addReaction(_whois.Trigger, "CREATE", _whois)
}
