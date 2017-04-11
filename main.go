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
)

type config struct {
	Discord struct {
		BotKey  string
		Channel string
	}

	IronRealms struct {
		APIURL string
	}
}

var _config config

func init() {
	if data, err := ioutil.ReadFile("env.yaml"); err == nil {
		if err := yaml.Unmarshal([]byte(data), &_config); err == nil {
			ire.APIURL = _config.IronRealms.APIURL
		} else {
			fmt.Println("ERROR: Error with parsing env.yaml.")
			log.Fatalf("error: %v", err)
		}
	} else {
		fmt.Println("ERROR: The file 'env.yaml' could not be read.")
		log.Fatalf("error: %v", err)
	}
}

func main() {
	err := discord.InitDiscordSession(_config.Discord.BotKey)
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
		os.Exit(0)
	}()

	IRE := ire.Gamefeed{}

	discord.PostMessage(_config.Discord.Channel, "```---------- ada-bot restarting ... ----------```")
	ticker := time.NewTicker(time.Millisecond * 30000) // 30 second ticker
	go func() {
		for _ = range ticker.C {
			if deathsights, err := IRE.Sync(); err == nil {
				for _, event := range deathsights {
					time.Sleep(time.Millisecond * 500) // Wait half a second FIXME: this way is awkward
					discord.PostMessage(_config.Discord.Channel, fmt.Sprintf("```%s - %s```", event.Date, event.Description))
				}
			} else {
				fmt.Println("ERROR: We couldn't get new deathsights.") // Not a fatal error
				log.Printf("error: %v", err)
			}
		}
	}()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan int) // block forever till SIGINT / SIGTERM
}
