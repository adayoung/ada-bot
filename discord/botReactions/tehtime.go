package botReactions

import (
	"regexp"
	"strings"
	"time"

	"github.com/adayoung/ada-bot/settings"
	"github.com/bwmarrin/discordgo"
)

var argsPattern = regexp.MustCompile(`(\S+)( (?:\+|-)\w+)?`)
var IsDirectory = map[string]string{
	"Africa": "Johannesburg",
	"America": "Chicago",
	"Asia": "Shanghai",
	"Australia": "Sydney",
	"Europe": "Zurich",
}

var commonZones = map[string]string {
	"BST": "Europe/London",
	"CET": "Europe/Zurich",
	"CEST": "Europe/Zurich",
	"CDT": "America/Chicago",
	"EDT": "America/New_York",
	"PDT": "America/Los_Angeles",
}

var underscoreExceptions = map[string]string {
	"dar_es_salaam": "Africa/Dar_es_Salaam",
	"port_of_spain": "America/Port_of_Spain",
	"ho_chi_minh": "Asia/Ho_Chi_Minh",
	"isle_of_man": "Europe/Isle_of_Man",
}

type tehtime struct {
	Trigger string
}

func (t *tehtime) Help() string {
	return "What's teh time?! Try something like " +
		settings.Settings.Discord.BotPrefix +
		"time Europe/Paris +6h or " +
		settings.Settings.Discord.BotPrefix +
		"time EDT"
}

func (t *tehtime) HelpDetail() string {
	return t.Help()
}

func (t *tehtime) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	timeNow := time.Now()
	timezoneExtra := ""

	request := m.Content[len(settings.Settings.Discord.BotPrefix)+len(t.Trigger):]
	request = strings.TrimSpace(request)

	args := argsPattern.FindStringSubmatch(request)
	if args != nil {
		timezone := ""
		duration := ""

		if len(args[1]) > 0 && len(args[2]) > 0 {
			timezone = args[1]
			duration = args[2]
		} else if len(args[1]) > 0 && len(args[2]) == 0 {
			if match, _ := regexp.MatchString(`\+|-`, args[1]); match {
				duration = args[1]
			} else {
				timezone = args[1]
			}
		}

		if cz, ok := commonZones[strings.ToUpper(timezone)]; ok {
			timezone = cz
			timezoneExtra = " " + cz
		}

		timezone = strings.Title(strings.ToLower(timezone))

		if strings.Contains(timezone, "_") {
			timeZoneSplit := strings.SplitAfterN(timezone, "_", 2)
			timezone = timeZoneSplit[0] + strings.Title(strings.ToLower(timeZoneSplit[1]))
		}

		timezoneSplit := strings.Split(timezone, "/")
		timezoneSplitLast := timezoneSplit[len(timezoneSplit) - 1]
		if cz, ok := underscoreExceptions[strings.ToLower(timezoneSplitLast)]; ok {
			timezone = cz
			timezoneExtra = " " + cz
		}

		if len(timezone) > 0 {
			if location, err := time.LoadLocation(timezone); err == nil {
				timeNow = timeNow.In(location)
			} else {
				if err.Error() == "is a directory" {
					if _, ok := IsDirectory[timezone]; ok {
						if location, err := time.LoadLocation(timezone + "/" + IsDirectory[timezone]); err == nil {
							timeNow = timeNow.In(location)
							timezoneExtra = " (" + timezone + "/" + IsDirectory[timezone] + ")"
						} else {
							return Reaction{Text: "Oop, invalid time zone: " + err.Error()}
						}
					} else {
						return Reaction{Text: "Oop, invalid time zone: " + err.Error()}
					}
				} else {
					if !strings.Contains(timezone, "/") {
						for key := range IsDirectory {
							if location, err := time.LoadLocation(key + "/" + timezone); err == nil {
								timeNow = timeNow.In(location)
								timezoneExtra = " (" + key + "/" + timezone + ")"
							}
						}

						if len(timezoneExtra) == 0 {
							return Reaction{Text: "Oop, invalid time zone: " + err.Error()}
						}
					} else {
						return Reaction{Text: "Oop, invalid time zone: " + err.Error()}
					}
				}
			}
		}

		if len(duration) > 0 {
			if d, err := time.ParseDuration(strings.TrimSpace(duration)); err == nil {
				timeNow = timeNow.Add(d)
			} else {
				return Reaction{Text: "Oop, invalid duration: " + err.Error()}
			}
		}
	}

	return Reaction{Text: timeNow.Format("2006-01-02 15:04:05 MST" + timezoneExtra)}
}

func init() {
	_tehtime := &tehtime{
		Trigger: "time",
	}
	addReaction(_tehtime.Trigger, "CREATE", _tehtime)
}
