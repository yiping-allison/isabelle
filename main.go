package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/yiping-allison/daisymae/daisymaebot"
	"github.com/yiping-allison/daisymae/service"
)

func main() {
	bc, err := LoadConfig()
	if err != nil {
		fmt.Printf("error loading config; err = %s\n", err)
		return
	}
	// TODO: Implement database connection
	// TODO: Test database connection
	dbCfg := bc.Database
	services, err := service.NewServices(
		service.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		service.WithLogMode(true),
	)
	if err != nil {
		fmt.Println(err)
	}
	daisy, err := daisymaebot.New(bc.BotKey)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	daisy.Service = *services
	// FIXME: Uncomment after integrating sql types
	// daisy.Service.AutoMigrate()
	defer daisy.Service.Close()
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
