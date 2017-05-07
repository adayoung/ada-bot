package bot_reactions

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/httpclient"
)

var apiURL string = "http://api.urbandictionary.com/v0/define"

type urbanDefinition struct {
	Definition  string `json:"definition"`
	Permalink   string `json:"permalink"`
	ThumbsUp    int    `json:"thumbs_up"`
	Author      string `json:"author"`
	Word        string `json:"word"`
	DefID       int    `json:"defid"`
	CurrentVote string `json:"current_vote"`
	Example     string `json:"example"`
	ThumbsDown  int    `json:"thumbs_down"`
}

type urbanDictionaryResult struct {
	Tags       []string          `json:"tags"`
	ResultType string            `json:"result_type"`
	List       []urbanDefinition `json:"list"`
	Sounds     []string          `json:"sounds"`
}

type urbanDefinitionByVote []urbanDefinition  // Implements sort.Interface
func (u urbanDefinitionByVote) Len() int      { return len(u) }
func (u urbanDefinitionByVote) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
func (u urbanDefinitionByVote) Less(i, j int) bool {
	return u[i].ThumbsDown < u[j].ThumbsUp
}

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

	var _results urbanDictionaryResult
	if err := httpclient.GetJSON(url, &_results); err == nil {
		sort.Sort(urbanDefinitionByVote(_results.List))
		if len(_results.List) > 0 {
			response = fmt.Sprintf(
				"```%s```", _results.List[0].Definition,
			)
			if len(_results.List[0].Example) > 0 {
				response = fmt.Sprintf("%s```Example: %s```", response, _results.List[0].Example)
			}
			if len(_results.Tags) > 0 {
				response = fmt.Sprintf("%s```Tags: %s```", response, strings.Join(_results.Tags, ", "))
			}
		}
	} else {
		log.Printf("error: %v", err) // Non fatal error at httpclient.GetJSON() call
	}

	return Reaction{Text: response}
}

func init() {
	_urbanDictionary := &urbanDictionary{
		Trigger: "define",
	}
	addReaction(_urbanDictionary.Trigger, "CREATE", _urbanDictionary)
}
