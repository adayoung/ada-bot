package bot_reactions

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	eliza "github.com/necrophonic/go-eliza"

	"github.com/adayoung/ada-bot/settings"
)

type Eliza struct {
	Trigger string
}

func (e *Eliza) Help() string {
	return "Talk to Eliza!"
}

func (e *Eliza) HelpDetail(m *discordgo.Message) string {
	return e.Help()
}

var requestRegexp *regexp.Regexp = regexp.MustCompile("(?i)[^a-z!',. ]+")

func (e *Eliza) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	if strings.HasPrefix(m.Content, settings.Settings.Discord.BotPrefix) {
		return "" // It's an explicit bot reaction-request, bail out
	}

	if a.GuildID != "" {
		return "" // Let's not talk on a channel unless it's a DM
	}

	request := requestRegexp.ReplaceAllString(m.Content, "")
	response, err := eliza.AnalyseString(request)
	if err != nil {
		log.Printf("error: %v", err) // Error with eliza.AnalyseString() call
	}
	return response
}

func init() {
	_eliza := &Eliza{
		Trigger: "*",
	}
	addReaction(_eliza.Trigger, _eliza)
}
