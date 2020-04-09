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
		models.WithUsers(),
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

	cleaning := scheduleClean(clean, 15*time.Minute, isa)
	defer cleaning.Stop()

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

// clean will call the routine cleans for event and user tracking
func clean(isa *isabellebot.Bot) {
	isa.Service.Event.Clean()
	isa.Service.User.Clean()
}

// scheduleClean will run routine cleaning after specified time duration
func scheduleClean(f func(*isabellebot.Bot), interval time.Duration, isa *isabellebot.Bot) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			f(isa)
		}
	}()
	return ticker
}
