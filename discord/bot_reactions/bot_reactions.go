package bot_reactions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
)

type BotReaction interface {
	Help() string
	HelpDetail(*discordgo.Message) string
	Reaction(message *discordgo.Message, author *discordgo.Member, update bool) string
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

func GetReactions(message *discordgo.Message, author *discordgo.Member, update bool) []string {
	var reactions []string
	if _, ok := _botReactions["*"]; ok { // Run wildcard triggers first
		for _, reaction := range _botReactions["*"] {
			if author.GuildID == "" {
				reactions = append(reactions, reaction.Reaction(message, author, update))
			} else {
				_ = reaction.Reaction(message, author, update) // Wildcard triggers should not respond on channels
			}
		}
	}

	if strings.HasPrefix(message.Content, fmt.Sprintf("%s*", settings.Settings.Discord.BotPrefix)) {
		return reactions // Attempted wildcard trigger! Abort abort!
	}

	if !strings.HasPrefix(message.Content, settings.Settings.Discord.BotPrefix) {
		return reactions // The message is irrelevant, bail out with no reactions
	}

	if strings.TrimSpace(strings.ToLower(message.Content)) == fmt.Sprintf("%shelp", settings.Settings.Discord.BotPrefix) {
		reactions = append(reactions, GenHelp())
		return reactions
	}

	for trigger, _reactions := range _botReactions {
		if strings.HasPrefix(strings.ToLower(message.Content[len(settings.Settings.Discord.BotPrefix):]), strings.ToLower(trigger)) {
			for _, reaction := range _reactions {
				reactions = append(reactions, reaction.Reaction(message, author, update))
			}
		}
	}

	return reactions
}

func GenHelp() string {
	_longestTrigger := 0
	triggers := []string{}
	for trigger := range _botReactions {
		if trigger != "*" {
			triggers = append(triggers, trigger)
			if len(trigger) > _longestTrigger {
				_longestTrigger = len(trigger)
			}
		}
	}
	sort.Strings(triggers)

	help := "I have the following commands available:"
	for _, trigger := range triggers {
		for _, item := range _botReactions[trigger] {
			help = fmt.Sprintf(
				"%s\n%s%s - %s", help, settings.Settings.Discord.BotPrefix,
				fmt.Sprintf("%s%s", trigger, strings.Repeat(" ", _longestTrigger-len(trigger))),
				item.Help(),
			)
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
