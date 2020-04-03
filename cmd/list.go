package cmd

import (
	"github.com/bwmarrin/discordgo"
)

// List will list all bot commands to the user
func List(cmdInfo CommandInfo) {
	emThumb := &discordgo.MessageEmbedThumbnail{
		URL:    "https://cdn2.iconfinder.com/data/icons/business-office-icons/256/To-do_List-512.png",
		Width:  100,
		Height: 100,
	}
	fields := createFields(cmdInfo)
	emMsg := &discordgo.MessageEmbed{
		Title:     "Commands",
		Thumbnail: emThumb,
		Fields:    fields,
	}
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
}

// createFields is a helper func which combines
func createFields(cmdInfo CommandInfo) []*discordgo.MessageEmbedField {
	format := func(f ...*discordgo.MessageEmbedField) []*discordgo.MessageEmbedField { return f }
	search := &discordgo.MessageEmbedField{
		Name:  "search",
		Value: "?search [item]",
	}
	help := &discordgo.MessageEmbedField{
		Name:  "help",
		Value: "?help [command_name]",
	}
	ping := &discordgo.MessageEmbedField{
		Name:  "ping",
		Value: "?ping",
	}
	pong := &discordgo.MessageEmbedField{
		Name:  "pong",
		Value: "?pong",
	}
	return format(search, help, ping, pong)
}
