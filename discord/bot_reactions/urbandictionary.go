package bot_reactions

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/httpclient"
)

var apiURL string = "http://api.urbandictionary.com/v0/define"

type urbanDictionary struct {
	Trigger string
}

func (u *urbanDictionary) Help() string {
	return "Ask Urban Dictionary for a meaning!"
}

func (u *urbanDictionary) HelpDetail() string {
	return u.Help()
}

func (u *urbanDictionary) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	request := strings.ToLower(strings.TrimSpace(m.Content[len(settings.Settings.Discord.BotPrefix)+len(u.Trigger):]))
	response := fmt.Sprintf("I couldn't find the meaning of %s :frowning:", request)
	urlArgs := url.Values{"term": []string{request}}
	url := fmt.Sprintf("%s?%s", apiURL, urlArgs.Encode())

	var _results interface{}
	// FIXME: What follows looks like unholy witchcraft, I think I got carried away with type assertion XD
	if err := httpclient.GetJSON(url, &_results); err == nil {
		if results, ok := _results.(map[string]interface{}); ok {
			if rType, ok := results["result_type"]; ok {
				if _, ok := rType.(string); ok {
					if rType == "exact" {
						if _definitions, ok := results["list"]; ok {
							if definitions, ok := _definitions.([]interface{}); ok {
								if len(definitions) > 0 {
									firstDefinition := definitions[0]
									if tehDefinition, ok := firstDefinition.(map[string]interface{}); ok {
										if _definition, ok := tehDefinition["definition"]; ok {
											if definition, ok := _definition.(string); ok {
												response = fmt.Sprintf("```%s```", definition)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		log.Printf("error: %v", err) // Not a fatal error
	}

	return Reaction{Text: response}
}

func init() {
	_urbanDictionary := &urbanDictionary{
		Trigger: "define",
	}
	addReaction(_urbanDictionary.Trigger, "CREATE", _urbanDictionary)
}
