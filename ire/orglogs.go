package ire

import (
	"fmt"
	"net/url"

	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/httpclient"
)

// OrgLog is a single log event
type OrgLog struct {
	RealDate string `json:"-"`
	Date     int    `json:"date"`
	Event    string `json:"event"`
}

// OrgLogs is a collection of OrgLog(s)
type OrgLogs struct {
	LastOrgLog int
	Events     *[]OrgLog
}

// Sync gets the latest orglogs from API endpoint, returns orglogs
func (o *OrgLogs) Sync(character, password string) ([]OrgLog, error) {
	URL := fmt.Sprintf("%s/orglogs/targossas.json", APIURL)
	o.LastOrgLog = settings.Settings.IRE.LastOrgLog

	credentials := url.Values{}
	credentials.Set("character", character)
	credentials.Set("password", password)
	URL = URL + "?" + credentials.Encode()

	var orglogs []OrgLog

	if !settings.Settings.IRE.OrgLogsEnabled { // Oops, we're disabled, bail out
		return orglogs, nil
	}

	if err := httpclient.GetJSON(URL, &o.Events); err == nil {
		for _, event := range *o.Events {
			if event.Date > o.LastOrgLog {
				o.LastOrgLog = event.Date
				orglogs = append(orglogs, event)
			}
		}
	} else {
		return nil, err // Error at httpclient.GetJSON() call
	}

	settings.Settings.Lock()
	defer settings.Settings.Unlock()
	settings.Settings.IRE.LastOrgLog = o.LastOrgLog
	return orglogs, nil
}
