package cmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Help defines the bot's help command
//
// It prints the help message of a specific bot command using
// Discord's message embedding
func Help(cmdInfo CommandInfo) {
	title := ""
	desc := ""
	extraTitle := ""
	extraText := ""
	if len(cmdInfo.CmdOps) == 1 {
		// When user only writes: ?help
		title = "Error"
		desc = "You must enter a valid command."
		extraTitle = "EXAMPLE"
		extraText = "?help search"
		prettyPrint(title, desc, extraTitle, extraText, cmdInfo)
		return
	}
	full := strings.Join(cmdInfo.CmdOps[1:], " ")
	if !find(full, cmdInfo) {
		title = full
		desc = "Command not found"
		extraTitle = "To list out all commands:"
		extraText = "?list"
		prettyPrint(title, desc, extraTitle, extraText, cmdInfo)
		return
	}
	switch full {
	case "search":
		title = "Search"
		desc = "Search will look up an item from New Horizon's bug and fish database."
		extraTitle = "EXAMPLE"
		extraText = "?search emperor butterfly or ?search north bugs"
	case "pong":
		title = "Pong"
		desc = "Playing with pong."
		extraTitle = "EXAMPLE"
		extraText = "?pong"
	}
	prettyPrint(title, desc, extraTitle, extraText, cmdInfo)
}

// find returns true if a specific command name
// was found in the command name list
func find(ops string, cmdInfo CommandInfo) bool {
	for _, cmd := range cmdInfo.CmdList {
		if ops == cmd {
			return true
		}
	}
	return false
}

// prettyPrint uses discord message embedding to print help messages
func prettyPrint(title, desc, innerT, innerTxt string, cmdInfo CommandInfo) {
	emThumb := &discordgo.MessageEmbedThumbnail{
		URL:    "https://www.bbqguru.com/content/images/manual-bbq-icon.png",
		Width:  100,
		Height: 100,
	}
	emDes := &discordgo.MessageEmbedField{
		Name:  innerT,
		Value: innerTxt,
	}
	emMsg := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Thumbnail:   emThumb,
		Fields:      []*discordgo.MessageEmbedField{emDes},
	}
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
}
