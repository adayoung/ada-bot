package botReactions

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
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

var qWhoLast = time.Now().AddDate(0, 0, -1)

func (q *qwhoTrigger) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	/* begin rate limit qwho */
	timeNow := time.Now()
	if timeNow.Sub(qWhoLast) < time.Second*60 {
		return Reaction{Text: fmt.Sprintf("Oops, %s%s is rate limited to once per minute only :shrug:",
			settings.Settings.Discord.BotPrefix,
			q.Trigger,
		)}
	} else {
		qWhoLast = timeNow
	}
	/* end rate limit qwho */

	url := "http://api.achaea.com/characters.json"
	response := "```"

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
