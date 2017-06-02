package botReactions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
)

type readNews struct {
	Trigger string
}

func (r *readNews) Help() string {
	return "Generate link to Achaea NEWS article!"
}

func (r *readNews) HelpDetail() string {
	return r.Help()
}

var readNewsRegexp = regexp.MustCompile(`(?i)([a-z]+)\s([0-9]+)`)

func (r *readNews) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	request := strings.TrimSpace(m.Content[len(settings.Settings.Discord.BotPrefix)+len(r.Trigger):])
	rMatch := readNewsRegexp.FindStringSubmatch(request)
	if len(rMatch) > 0 {
		section, number := rMatch[1], rMatch[2]
		response := fmt.Sprintf("https://www.achaea.com/news/?game=Achaea&section=%s&number=%s", section, number)
		return Reaction{Text: response}
	}
	return Reaction{}
}

func init() {
	_readNews := &readNews{
		Trigger: "readnews",
	}
	addReaction(_readNews.Trigger, "CREATE", _readNews)
}
