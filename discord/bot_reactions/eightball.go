package bot_reactions

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Magic []string = []string{
	"By the Logos, it is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it, yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
	"No -- Lorielan, the Jade Empress",
}

type EightBall struct {
	Trigger string
}

func (p *EightBall) Help() string {
	return "Let the magic 8-ball guide your destiny (Y/N questions only)."
}

func (p *EightBall) HelpDetail(m *discordgo.Message) string {
	return p.Help()
}

func (p *EightBall) Reaction(m *discordgo.Message, a *discordgo.Member, u bool) string {
	theAnswer := Magic[rand.Intn(len(Magic))]
	return fmt.Sprintf("```%s```", theAnswer)
}

func init() {
	rand.Seed(time.Now().Unix())

	eightball := &EightBall{
		Trigger: "8ball",
	}
	addReaction(eightball.Trigger, eightball)
}
