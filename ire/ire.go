package ire

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

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

func (g *Gamefeed) Sync(url string) ([]Event, error) {
	if g.LastID > 0 {
		url = fmt.Sprintf("%s?id=%d", url, g.LastID)
	}

	var deathsights []Event
	if response, err := http.Get(url); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			if err := json.Unmarshal([]byte(data), &g.Events); err == nil {
				for _, event := range *g.Events {
					if event.ID > g.LastID {
						g.LastID = event.ID
					}

					if event.Type == "DEA" {
						deathsights = append(deathsights, event)
					}
				}
			} else {
				return nil, err
			} // Error at json.Unmarshal() call
		} else {
			return nil, err // Error at ioutil.ReadAll() call
		}
	} else {
		return nil, err // Error at http.Get() call
	}

	sort.Sort(EventsByDate(deathsights))
	return deathsights, nil
}
