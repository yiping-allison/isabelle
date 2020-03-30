package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yiping-allison/daisymae/daisymaebot"
)

func main() {
	bc, err := LoadConfig()
	if err != nil {
		fmt.Printf("error loading config; err = %s\n", err)
		return
	}
	daisy, err := daisymaebot.New(bc.BotKey)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	daisy.SetPrefix(bc.BotPrefix)
	err = daisy.DS.Open()
	if err != nil {
		fmt.Printf("Error opening connection; err = %v\n", err)
		return
	}
	defer daisy.DS.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
