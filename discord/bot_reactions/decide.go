package bot_reactions

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
)

type Decide struct {
	Trigger string
}

func (d *Decide) Help() string {
	return "Let the bot decide between two or more things for you!"
}

func (d *Decide) HelpDetail(*discordgo.Message) string {
	return d.Help()
}

func (d *Decide) Reaction(m *discordgo.Message, a *discordgo.Member, u bool) string {
	choices := strings.Split(m.Content[len(settings.Settings.Discord.BotPrefix)+len(d.Trigger):], " or ")
	theAnswer := choices[rand.Intn(len(choices))]
	return fmt.Sprintf("The correct answer is **%s**", strings.TrimSpace(theAnswer))
}

func init() {
	rand.Seed(time.Now().Unix())

	decide := &Decide{
		Trigger: "decide",
	}
	addReaction(decide.Trigger, decide)
}
