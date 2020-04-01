package cmd

import "github.com/bwmarrin/discordgo"

// Help defines the bot's help command
//
// It prints the help message of a specific bot command using
// Discord's message embedding
func Help(cmdInfo CommandInfo) {
	// TODO: Prettier formatting?
	emThumb := &discordgo.MessageEmbedThumbnail{
		URL:    "https://www.bbqguru.com/content/images/manual-bbq-icon.png",
		Width:  100,
		Height: 100,
	}
	emMsg := &discordgo.MessageEmbed{
		Title:       cmdInfo.CmdName,
		Description: cmdInfo.CmdHlp,
		Thumbnail:   emThumb,
	}
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
}
