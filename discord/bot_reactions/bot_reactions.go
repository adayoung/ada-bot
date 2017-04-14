package bot_reactions

import (
	"github.com/bwmarrin/discordgo"
	// "github.com/adayoung/ada-bot/settings"
)

type BotReaction interface {
	Help() string
	HelpDetail(*discordgo.Message) string
	Reaction(message *discordgo.Message, author *discordgo.Member) string
}

var _botReactions map[string][]BotReaction

func init() {
	_botReactions = make(map[string][]BotReaction)
}

func addReaction(trigger string, reaction BotReaction) {
	// FIXME: calls to addReaction should be idempotent, let's
	// not add multiple instances of the same reaction to a trigger
	_botReactions[trigger] = append(_botReactions[trigger], reaction)
}

func GetReactions(trigger string, message *discordgo.Message, author *discordgo.Member) []string {
	var reactions []string
	if _, ok := _botReactions[trigger]; ok {
		for _, reaction := range _botReactions[trigger] {
			reactions = append(reactions, reaction.Reaction(message, author))
		}
	}
	return reactions
}

func GenHelp() map[string][]string {
	help := make(map[string][]string)
	for k, v := range _botReactions {
		help[k] = []string{}
		for _, item := range v {
			help[k] = append(help[k], item.Help())
		}
	}
	return help
}

func GetHelpDetail(trigger string, message *discordgo.Message) []string {
	var help []string
	if _, ok := _botReactions[trigger]; ok {
		for _, reaction := range _botReactions[trigger] {
			help = append(help, reaction.HelpDetail(message))
		}
	}
	return help
}

func HelloWorld() string {
	return "Hello, world!"
}
