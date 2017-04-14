package bot_reactions

import (
	"fmt"
	"strings"

	"github.com/adayoung/ada-bot/settings"
	"github.com/bwmarrin/discordgo"
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

func GetReactions(message *discordgo.Message, author *discordgo.Member) []string {
	if _, ok := _botReactions["*"]; ok { // Run wildcard triggers first
		for _, reaction := range _botReactions["*"] {
			reaction.Reaction(message, author)
		}
	}

	var reactions []string
	if !strings.HasPrefix(message.Content, settings.Settings.Discord.BotPrefix) {
		return reactions // The message is irrelevant, bail out with no reactions
	}

	if strings.TrimSpace(message.Content) == fmt.Sprintf("%shelp", settings.Settings.Discord.BotPrefix) {
		reactions = append(reactions, GenHelp())
		return reactions
	}

	for trigger, _reactions := range _botReactions {
		if strings.HasPrefix(message.Content[len(settings.Settings.Discord.BotPrefix):], trigger) {
			for _, reaction := range _reactions {
				reactions = append(reactions, reaction.Reaction(message, author))
			}
		}
	}

	return reactions
}

func GenHelp() string { // FIXME: Use padding to correctly align the help text
	help := "I have the following commands available:"
	for k, v := range _botReactions {
		for _, item := range v {
			help = fmt.Sprintf("%s\n%s%s - %s", help, settings.Settings.Discord.BotPrefix, k, item.Help())
		}
	}
	return fmt.Sprintf("```%s```", help)
}

func GetHelpDetail(trigger string, message *discordgo.Message) string {
	return "" // TODO: Not implemented yet
	// var help []string
	// if _, ok := _botReactions[trigger]; ok {
	// 	for _, reaction := range _botReactions[trigger] {
	// 		help = append(help, reaction.HelpDetail(message))
	// 	}
	// }
	// return help
}
