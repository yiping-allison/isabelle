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
	if len(cmdInfo.CmdOps) == 1 {
		// When user only writes: ?help
		prettyPrintHelp(
			"Error",
			"You must query a valid command.",
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"help search", true),
			),
			cmdInfo,
			14886454,
		)
		return
	}
	full := strings.Join(cmdInfo.CmdOps[1:], " ")
	if !find(full, cmdInfo) {
		prettyPrintHelp(
			full,
			"Command Not Found",
			format(
				createFields("To List All Commands:", cmdInfo.Prefix+"list", true),
			),
			cmdInfo,
			14886454,
		)
		return
	}
	// Valid commands
	switch full {
	case "search":
		prettyPrintHelp(
			"Search",
			"Search will look up an item from New Horizon's bug and fish database.",
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"search emperor butterfly", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"search north bug", true),
			),
			cmdInfo,
			9410425,
		)
	case "list":
		prettyPrintHelp(
			"List",
			"List will show all commands the bot understands.",
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"list", true),
			),
			cmdInfo,
			9410425,
		)
	case "pong":
		prettyPrintHelp(
			"Pong",
			"Playing with pong.",
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"pong", true),
			),
			cmdInfo,
			9410425,
		)
	}
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
func prettyPrintHelp(title, desc string, fields []*discordgo.MessageEmbedField, cmdInfo CommandInfo, color int) {
	emThumb := &discordgo.MessageEmbedThumbnail{
		URL:    "https://www.bbqguru.com/content/images/manual-bbq-icon.png",
		Width:  100,
		Height: 100,
	}
	emMsg := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Thumbnail:   emThumb,
		Fields:      fields,
		Color:       color,
	}
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
}
