package cmd

import (
	"github.com/bwmarrin/discordgo"
)

// Ping test func to play with rich embedding
func Ping(cmdInfo CommandInfo) {
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
	cmdInfo.Ses.ChannelMessageSend(cmdInfo.Msg.ChannelID, cmdInfo.Msg.Author.Mention())
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
}
