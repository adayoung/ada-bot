package bot_reactions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ping struct {
	Trigger string
}

func (p *ping) Help() string {
	return "Pong!"
}

func (p *ping) HelpDetail(m *discordgo.Message) string {
	return p.Help()
}

func (p *ping) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) string {
	if strings.Contains(strings.ToLower(m.Content), "pong") {
		return "Ping!"
	}
	return "Pong!"
}

func init() {
	_ping := &ping{
		Trigger: "ping",
	}
	addReaction(_ping.Trigger, "CREATE", _ping)

	_pong := &ping{
		Trigger: "pong",
	}
	addReaction(_pong.Trigger, "CREATE", _pong)
}
