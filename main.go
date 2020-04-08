package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/yiping-allison/isabelle/isabellebot"
	"github.com/yiping-allison/isabelle/models"
)

func main() {
	bc, err := LoadConfig()
	if err != nil {
		fmt.Printf("error loading config; err = %s\n", err)
		return
	}
	dbCfg := bc.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(true),
		models.WithEntries(),
		models.WithEvents(),
	)
	if err != nil {
		fmt.Println(err)
	}
	isa, err := isabellebot.New(bc.BotKey, bc.AdminRole)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	isa.Service = *services
	// FIXME: Uncomment after integrating sql types
	// daisy.Service.AutoMigrate()
	defer isa.Service.Close()
	isa.SetPrefix(bc.BotPrefix)

	go func() {
		// check map and remove anything that expired
		// every 15 minutes (can be changed if you want frequent or longer time interval cleaning)
		cleanTicker := time.NewTicker(15 * time.Minute)
		for {
			select {
			case <-cleanTicker.C:
				isa.Service.Event.Clean()
			}
		}
	}()

	err = isa.DS.Open()
	if err != nil {
		fmt.Printf("Error opening connection; err = %v\n", err)
		return
	}
	defer isa.DS.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
