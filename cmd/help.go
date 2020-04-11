package cmd

import (
	"strings"
)

// Help defines the bot's help command
//
// It prints the help message of a specific bot command using
// Discord's message embedding
func Help(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) == 1 {
		// When user only writes: ?help
		// No valid command input
		msg := cmdInfo.createMsgEmbed("Error", errThumbURL, "You must enter a valid command",
			errColor, format(createFields("EXAMPLE", cmdInfo.Prefix+"help search", true)))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	full := strings.Join(cmdInfo.CmdOps[1:], " ")
	if !find(full, cmdInfo) {
		// Command not found
		msg := cmdInfo.createMsgEmbed(full, errThumbURL, "Command Not Found", errColor, format(
			createFields("To List All Commands: ", cmdInfo.Prefix+"list", true),
		))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	// Valid commands
	switch full {
	case "search":
		msg := cmdInfo.createMsgEmbed("Search", helpThumbURL, "Looks up an item from bug and fish database.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"search emperor butterfly", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"search north bug", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "event":
		msg := cmdInfo.createMsgEmbed("Event", helpThumbURL, "Creates visitation events.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event celeste limit=\"2\" msg=\"Come on over for shooting stars\"", false),
				createFields("EXAMPLE", cmdInfo.Prefix+"event 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "queue":
		msg := cmdInfo.createMsgEmbed("Queue", helpThumbURL, "Join a queue for visitation events.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "close":
		msg := cmdInfo.createMsgEmbed("Close", helpThumbURL, "Ends events.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"close event 1234", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"close trade 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "unregister":
		msg := cmdInfo.createMsgEmbed("Unregister", helpThumbURL, "Removes yourself from listings.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"unregister event 1234", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"unregister trade 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "accept":
		msg := cmdInfo.createMsgEmbed("Accept", helpThumbURL, "Accepts reputation applications.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"accept 1234", true),
				createFields("NOTE", "This command is only available to moderators.", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "reject":
		msg := cmdInfo.createMsgEmbed("Reject", helpThumbURL, "Rejects reputation applications.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"reject 1234", true),
				createFields("NOTE", "This command is only available to moderators.", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "rep":
		msg := cmdInfo.createMsgEmbed("Rep", helpThumbURL, "Creates a new reputation application.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"rep @awesome-person successfully traded coffee beans", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "trade":
		msg := cmdInfo.createMsgEmbed("Trade", helpThumbURL, "Creates a new trade event.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"trade blue mountain coffee", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "offer":
		msg := cmdInfo.createMsgEmbed("Offer", helpThumbURL, "Provide an offer to a trade event.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"offer 1234 geisha coffee beans", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)

	case "list":
		msg := cmdInfo.createMsgEmbed("List", helpThumbURL, "Displays all bot commands.", helpColor,
			format(createFields("EXAMPLE", cmdInfo.Prefix+"list", true)))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
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
