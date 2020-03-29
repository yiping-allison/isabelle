package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "token", "", "Set the bot's discord token using this flag.")
	flag.Parse()
}

func main() {
	bc, err := LoadConfig()
	if err != nil {
		fmt.Printf("error loading config; err = %s\n", err)
		return
	}
	discord, err := discordgo.New("Bot " + bc.BotKey)
	if err != nil {
		fmt.Printf("error connecting to discord; err = %s\n", err)
		return
	}
	discord.AddHandler(messageCreate)

	err = discord.Open()
	defer discord.Close()
	if err != nil {
		fmt.Printf("Error opening connection; err = %s", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "?") {
		switch m.Content[1:] {
		case "ping":
			emField := &discordgo.MessageEmbedField{
				Name:   "hello",
				Value:  "https://clips.twitch.tv/EsteemedHumbleStarlingGOWSkull",
				Inline: true,
			}
			emThumb := &discordgo.MessageEmbedThumbnail{
				URL:    "http://upload.wikimedia.org/wikipedia/commons/thumb/2/2a/New_Logo_AD.jpg/266px-New_Logo_AD.jpg",
				Width:  100,
				Height: 100,
			}
			emImg := &discordgo.MessageEmbedImage{
				URL:    "https://i.huffpost.com/gen/1226279/images/o-BOX-TURTLE-facebook.jpg",
				Height: 300,
				Width:  300,
			}
			emMsg := &discordgo.MessageEmbed{
				Title:       "Why are there ads?",
				Description: "Turti asks why there are ads",
				Thumbnail:   emThumb,
				Fields:      []*discordgo.MessageEmbedField{emField},
				Image:       emImg,
			}
			s.ChannelMessageSend(m.ChannelID, m.Author.Mention())
			s.ChannelMessageSendEmbed(m.ChannelID, emMsg)
		case "pong":
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		case "ad":
			s.ChannelMessageSend(m.ChannelID, "https://clips.twitch.tv/EsteemedHumbleStarlingGOWSkull")
		default:
			s.ChannelMessageSend(m.ChannelID, "???")
		}
	}
}
