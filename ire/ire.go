package ire

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/adayoung/ada-bot/settings"
)

var APIURL string

type Event struct {
	ID          int    `json:"id"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

type Gamefeed struct {
	LastID int
	Events *[]Event
}

type EventsByDate []Event            // Implements sort.Interface
func (d EventsByDate) Len() int      { return len(d) }
func (d EventsByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d EventsByDate) Less(i, j int) bool {
	return d[i].Date < d[j].Date
}

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

	if err := getJSON(url, &g.Events); err == nil {
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
		return nil, err // Error at http.Get() call
	}

	settings.Settings.IRE.LastID = g.LastID
	sort.Sort(EventsByDate(deathsights))
	return deathsights, nil
}

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

func GetPlayer(player string) (*Player, error) {
	if !(len(player) > 0) {
		return nil, errors.New(fmt.Sprintf("Invalid player name supplied: %s", player))
	}
	if match, err := regexp.MatchString("(?i)[^a-z]+", player); err == nil {
		if !match {
			url := fmt.Sprintf("%s/characters/%s.json", APIURL, player)
			_player := &Player{}
			if err := getJSON(url, &_player); err == nil {
				return _player, nil
			} else {
				return nil, err // Error at getJson() call
			}
			return _player, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid player name supplied: %s", player))
		}
	} else {
		return nil, err // Error at regexp.MatchString() call
	}
	return nil, nil
}
