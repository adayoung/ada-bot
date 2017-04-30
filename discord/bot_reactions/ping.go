package bot_reactions

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ping struct {
	Trigger string
}

func (p *ping) Help() string {
	return "Pong!"
}

func (p *ping) HelpDetail() string {
	return p.Help()
}

func (p *ping) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	if strings.Contains(strings.ToLower(m.Content), "pong") {
		// return "Ping!"
		return Reaction{Text: "Ping!"}
	}
	// return "Pong!"
	return Reaction{Text: "Pong!"}
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
