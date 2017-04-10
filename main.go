package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"gopkg.in/yaml.v2"

	"ada-bot/discord"
)

type config struct {
	Discord struct {
		BotKey  string
		Channel string
	}

	IronRealms struct {
		GameFeed string
	}
}

var _config config

func init() {
	if data, err := ioutil.ReadFile("env.yaml"); err == nil {
		if err := yaml.Unmarshal([]byte(data), &_config); err != nil {
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

	// ticker := time.NewTicker(time.Millisecond * 3000) // 3 second ticker
	// go func() {
	// 	for t := range ticker.C {
	// 		fmt.Println("Tick at", t)
	// 	}
	// }()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan int) // block forever till SIGINT / SIGTERM
}
