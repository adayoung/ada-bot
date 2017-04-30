package bot_reactions

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var magic []string = []string{
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

type eightBall struct {
	Trigger string
}

func (p *eightBall) Help() string {
	return "Let the magic 8-ball guide your destiny (Y/N questions only)."
}

func (p *eightBall) HelpDetail() string {
	return p.Help()
}

func (p *eightBall) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	theAnswer := magic[rand.Intn(len(magic))]
	response := fmt.Sprintf("```%s```", theAnswer)
	return Reaction{Text: response}
}

func init() {
	rand.Seed(time.Now().Unix())

	_eightball := &eightBall{
		Trigger: "8ball",
	}
	addReaction(_eightball.Trigger, "CREATE", _eightball)
}
