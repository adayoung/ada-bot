package bot_reactions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Ping struct {
	Trigger string
}

func (p *Ping) Help() string {
	return "Pong!"
}

func (p *Ping) HelpDetail(m *discordgo.Message) string {
	return p.Help()
}

func (p *Ping) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) string {
	if strings.Contains(strings.ToLower(m.Content), "pong") {
		return "Ping!"
	}
	return "Pong!"
}

func init() {
	ping := &Ping{
		Trigger: "ping",
	}
	addReaction(ping.Trigger, "CREATE", ping)

	pong := &Ping{
		Trigger: "pong",
	}
	addReaction(pong.Trigger, "CREATE", pong)
}
