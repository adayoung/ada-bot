package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
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
		Channel []string
		QLength int
		WaitMS  string
	}

	IronRealms struct {
		APIURL string
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

	IRE := ire.Gamefeed{}

	for _, channel := range _config.Discord.Channel {
		discord.PostMessage(channel, "```---------- ada-bot restarting ... ----------```")
	}
	ticker := time.NewTicker(time.Millisecond * 30000) // 30 second ticker
	go func() {
		for range ticker.C {
			if deathsights, err := IRE.Sync(); err == nil {
				for _, event := range deathsights {
					for _, channel := range _config.Discord.Channel {
						discord.PostMessage(channel, fmt.Sprintf("```%s - %s```", event.Date, event.Description))
					}
				}
			} else {
				fmt.Println("ERROR: We couldn't get new deathsights.")
				log.Printf("warning: %v", err) // Not a fatal error
			}
		}
	}()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan int) // block forever till SIGINT / SIGTERM
}
