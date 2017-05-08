package botReactions

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/utils/httpclient"
)

type character struct {
	Name string
	URI  string
}

type qwho struct {
	Count      string      `json:"count"`
	Characters []character `json:"characters"`
}

type charactersByName []character        // Implements sort.Interface
func (q charactersByName) Len() int      { return len(q) }
func (q charactersByName) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q charactersByName) Less(i, j int) bool {
	return q[i].Name < q[j].Name
}

type qwhoTrigger struct {
	Trigger string
}

func (q *qwhoTrigger) Help() string {
	return "Check for online players visible on Achaea at the moment."
}

func (q *qwhoTrigger) HelpDetail() string {
	return q.Help()
}

func (q *qwhoTrigger) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	response := "```"
	url := "http://api.achaea.com/characters.json"

	var _results qwho
	var characters []string
	if err := httpclient.GetJSON(url, &_results); err == nil {
		sort.Sort(charactersByName(_results.Characters))
		for _, character := range _results.Characters {
			characters = append(characters, character.Name)
		}
		response = fmt.Sprintf("%sPlayers: %s", response, strings.Join(characters, ", "))
		response = fmt.Sprintf("%s\nTotal: %s", response, _results.Count)
	} else {
		log.Printf("error: %v", err) // Non fatal error at httpclient.GetJSON() call
	}

	response = fmt.Sprintf("%s```", response)
	return Reaction{Text: response}
}

func init() {
	_qwho := &qwhoTrigger{
		Trigger: "qwho",
	}
	addReaction(_qwho.Trigger, "CREATE", _qwho)
}
