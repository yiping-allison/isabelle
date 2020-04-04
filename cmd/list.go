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
	fields := format(
		createFields("search", "?search [item]", true),
		createFields("help", "?help [command_name]", false),
		createFields("ping", "?ping", true),
		createFields("pong", "?pong", true),
	)
	emMsg := &discordgo.MessageEmbed{
		Title:     "Commands",
		Thumbnail: emThumb,
		Fields:    fields,
	}
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
}
