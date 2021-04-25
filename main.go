package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/adayoung/ada-bot/discord"
	"github.com/adayoung/ada-bot/ire"
	"github.com/adayoung/ada-bot/settings"
	"github.com/adayoung/ada-bot/utils/storage"
)

type config struct {
	Discord struct {
		BotKey  string
		Channel string
		QLength int
		WaitMS  string
	}

	IronRealms struct {
		APIURL          string
		Character       string
		Password        string
		DefenceChannels []string
		OffenceChannels []string
		CityChannel     string
	}

	Database struct {
		Connection string
	}

	SettingsFile string
}

var _config config

func init() {
	if data, err := ioutil.ReadFile("config.yaml"); err == nil {
		if err = yaml.Unmarshal([]byte(data), &_config); err == nil {
			ire.APIURL = _config.IronRealms.APIURL
		} else {
			log.Println("ERROR: Error with parsing config.yaml.")
			log.Fatalf("error: %v", err)
		}
	} else {
		log.Println("ERROR: The file 'config.yaml' could not be read.")
		log.Fatalf("error: %v", err)
	}

	if err := storage.InitDB(_config.Database.Connection); err != nil {
		log.Println("The database could not be initialized, DB will not unavailable.")
		log.Printf("error: %v", err)
	}

	if err := settings.Init(_config.SettingsFile); err != nil {
		fmt.Printf("ERROR: The settings file, '%s' could not be processed.\n", _config.SettingsFile)
		log.Fatalf("error: %v", err)
	}
}

func main() {
	err := discord.InitDiscordSession(
		_config.Discord.BotKey,
		_config.Discord.QLength,
		_config.Discord.WaitMS,
	)
	if err != nil {
		fmt.Println("ERROR: We couldn't initialize a Discord session.")
		log.Fatalf("error: %v", err)
	}

	// Setup signals handling here
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Interrupt received, cleaning up...")
		discord.CloseDiscordSession()
		if err := settings.Settings.Save(); err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Exit(0)
	}()

	discord.PostMessage(_config.Discord.Channel, "```---------- ada-bot restarting ... ----------```")

	IRE := ire.Gamefeed{}
	ticker := time.NewTicker(time.Second * 30) // 30 second ticker
	go func() {
		for range ticker.C {
			if deathsights, err := IRE.Sync(); err == nil {
				for _, event := range deathsights {
					discord.PostMessage(_config.Discord.Channel, fmt.Sprintf("```%s - %s```", event.Date, event.Description))
				}
			} else {
				fmt.Println("ERROR: We couldn't get new deathsights.")
				log.Printf("warning: %v", err) // Not a fatal error
			}
		}
	}()

	OrgLogs := ire.OrgLogs{}
	orglogTicker := time.NewTicker(time.Second * 90) // 90 second ticker
	go func() {
		for range orglogTicker.C {
			if orglogs, err := OrgLogs.Sync(_config.IronRealms.Character, _config.IronRealms.Password); err == nil {
				for _, event := range orglogs {
					event.RealDate = time.Unix(int64(event.Date), 0).Format(time.Kitchen)

					if strings.HasSuffix(event.Event, "has declared a sanctioned raid against the City") ||
						strings.Contains(event.Event, " slew ") ||
						strings.Contains(event.Event, " has disarmed a hostile tank.") ||
						strings.HasSuffix(event.Event, "has withdrawn from the City, ending the sanctioned raid") ||
						strings.HasSuffix(event.Event, "have declared their intent to retaliate against our previous conquests") {
						for _, channel := range _config.IronRealms.DefenceChannels {
							discord.PostMessage(channel, fmt.Sprintf("```%s - %s```", event.RealDate, event.Event))
						}
					}

					if match, err := regexp.MatchString(`^The forces of \S+ have destroyed .+`, event.Event); err == nil {
						if match {
							for _, channel := range _config.IronRealms.DefenceChannels {
								discord.PostMessage(channel, fmt.Sprintf("```%s - %s```", event.RealDate, event.Event))
							}
						}
					} else {
						log.Printf("warning: %v", err) // Not a fatal error
					}

					if strings.Contains(event.Event, "has sanctioned a raid against") ||
						strings.HasPrefix(event.Event, "The sanctioned raid in ") {
						for _, channel := range _config.IronRealms.OffenceChannels {
							discord.PostMessage(channel, fmt.Sprintf("```%s - %s```", event.RealDate, event.Event))
						}
					}

					if match, err := regexp.MatchString(`, and \S+ have destroyed .+`, event.Event); err == nil {
						if match {
							for _, channel := range _config.IronRealms.OffenceChannels {
								discord.PostMessage(channel, fmt.Sprintf("```%s - %s```", event.RealDate, event.Event))
							}
						}
					} else {
						log.Printf("warning: %v", err) // Not a fatal error
					}

					if match, err := regexp.MatchString(`^\S+ completed the bounty on \S+ and has received 10000 gold sovereigns$`, event.Event); err == nil {
						if match {
							for _, channel := range _config.IronRealms.OffenceChannels {
								discord.PostMessage(channel, fmt.Sprintf("```%s - %s```", event.RealDate, event.Event))
							}
						}
					} else {
						log.Printf("warning: %v", err) // Not a fatal error
					}

					discord.PostMessage(_config.IronRealms.CityChannel, fmt.Sprintf("```%s - %s```", event.RealDate, event.Event))
				}
			} else {
				fmt.Println("ERROR: We couldn't get new orglogs.")
				log.Printf("warning: %v", err) // Not a fatal error
			}
		}
	}()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan int) // block forever till SIGINT / SIGTERM
}
