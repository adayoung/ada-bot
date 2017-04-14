package bot_reactions

import (
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

func (p *Ping) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	return "Pong!"
}

func init() {
	ping := &Ping{
		Trigger: "ping",
	}
	addReaction(ping.Trigger, ping)
}
