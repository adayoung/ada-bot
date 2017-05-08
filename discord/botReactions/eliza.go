package botReactions

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	_eliza "github.com/necrophonic/go-eliza"

	"github.com/adayoung/ada-bot/settings"
)

type eliza struct {
	Trigger string
}

func (e *eliza) Help() string {
	return "Talk to Eliza!"
}

func (e *eliza) HelpDetail() string {
	return e.Help()
}

var requestRegexp *regexp.Regexp = regexp.MustCompile("(?i)[^a-z!',. ]+")

func (e *eliza) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	if strings.HasPrefix(m.Content, settings.Settings.Discord.BotPrefix) {
		return Reaction{Text: ""} // It's an explicit bot reaction-request, bail out
	}

	if a.GuildID != "" {
		return Reaction{Text: ""} // Let's not talk on a channel unless it's a DM
	}

	request := requestRegexp.ReplaceAllString(m.Content, "")
	response, err := _eliza.AnalyseString(request)
	if err != nil {
		log.Printf("error: %v", err) // Error with eliza.AnalyseString() call
	}
	return Reaction{Text: response}
}

func init() {
	_eliza := &eliza{
		Trigger: "*",
	}
	addReaction(_eliza.Trigger, "CREATE", _eliza)
}
