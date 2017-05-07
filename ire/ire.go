package ire

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/httpclient"
)

// APIURL is read and set from config.yaml by main.init()
var APIURL string

// Event represents a single event per IRE's gamefeed
type Event struct {
	ID          int    `json:"id"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

// Gamefeed is a collection of Events, ahead of LastID
type Gamefeed struct {
	LastID int
	Events *[]Event
}

type eventsByDate []Event            // Implements sort.Interface
func (d eventsByDate) Len() int      { return len(d) }
func (d eventsByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d eventsByDate) Less(i, j int) bool {
	return d[i].Date < d[j].Date
}

// Sync gets the latest events from API endpoint, returns deathsights
func (g *Gamefeed) Sync() ([]Event, error) {
	url := fmt.Sprintf("%s/gamefeed.json", APIURL)
	g.LastID = settings.Settings.IRE.LastID
	if g.LastID > 0 {
		url = fmt.Sprintf("%s?id=%d", url, g.LastID)
	}

	var deathsights []Event

	if !settings.Settings.IRE.DeathsightEnabled { // Oops, we're disabled, bail out
		return deathsights, nil
	}

	if err := httpclient.GetJSON(url, &g.Events); err == nil {
		for _, event := range *g.Events {
			go logEvent(event)
			if event.ID > g.LastID {
				g.LastID = event.ID
			}

			if event.Type == "DEA" {
				deathsights = append(deathsights, event)
			}
		}
	} else {
		return nil, err // Error at httpclient.GetJSON() call
	}

	settings.Settings.IRE.LastID = g.LastID
	sort.Sort(eventsByDate(deathsights))
	return deathsights, nil
}

// Player represents a player per IRE's API
type Player struct {
	Name         string `json:"name"`
	Fullname     string `json:"fullname"`
	City         string `json:"city"`
	House        string `json:"house"`
	Level        string `json:"level"`
	Class        string `json:"class"`
	MobKills     string `json:"mob_kills"`
	PlayerKills  string `json:"player_kills"`
	XPRank       string `json:"xp_rank"`
	ExplorerRank string `json:"explorer_rank"`
}

func (s *Player) String() string {
	player := fmt.Sprintf(`
           Name: %s
          Class: %s (Level %s)
           City: %s
          House: %s
    Kills (Mob): %s
Kills (Players): %s
      Rank (XP): %s
Rank (Explorer): %s`,
		s.Fullname, strings.Title(s.Class), s.Level,
		strings.Title(s.City), strings.Title(s.House),
		s.MobKills, s.PlayerKills, s.XPRank, s.ExplorerRank)
	return player
}

// GetPlayer performs a lookup and retrieve a player from IRE's API
func GetPlayer(player string) (*Player, error) {
	if !(len(player) > 0) {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid player name supplied: %s", player))
	}
	if match, err := regexp.MatchString("(?i)[^a-z]+", player); err == nil {
		if !match {
			url := fmt.Sprintf("%s/characters/%s.json", APIURL, player)
			_player := &Player{}
			if err := httpclient.GetJSON(url, &_player); err == nil {
				return _player, nil
			} else {
				return nil, err // Error at httpclient.GetJson() call
			}
		} else {
			return nil, fmt.Errorf(fmt.Sprintf("Invalid player name supplied: %s", player))
		}
	} else {
		return nil, err // Error at regexp.MatchString() call
	}
}
