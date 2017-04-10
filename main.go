package main

import (
		"fmt"
		"os"
		"os/signal"
		"syscall"
		"time"
)

func main() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
				<-c
				fmt.Println("cleanup")
				os.Exit(0)
		}()

		ticker := time.NewTicker(time.Millisecond * 3000) // 3 second ticker
		go func() {
				for t := range ticker.C {
					fmt.Println("Tick at", t)
				}
		}()

		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		<-make(chan int) // block forever till SIGINT / SIGTERM
}
