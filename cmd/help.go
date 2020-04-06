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
		msg := cmdInfo.createMsgEmbed("Error", errThumbURL, "You must query a valid command",
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
		msg := cmdInfo.createMsgEmbed("Search", helpThumbURL, "Search will look up an item from New Horizon's bug and fish database.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"search emperor butterfly", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"search north bug", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
	case "event":
		msg := cmdInfo.createMsgEmbed("Event", helpThumbURL, "Event handles the event creation service (paired with queue-ing) to host visitation events.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event celeste limit=\"2\" msg=\"Come on over for shooting stars!\"", false),
				createFields("EXAMPLE", cmdInfo.Prefix+"event 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
	case "queue":
		msg := cmdInfo.createMsgEmbed("Queue", helpThumbURL, "Queue handles the queue creation service (paired with events) to join visitation events.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
	case "close":
		msg := cmdInfo.createMsgEmbed("Close", helpThumbURL, "Ends events.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"close event 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
	case "list":
		msg := cmdInfo.createMsgEmbed("List", helpThumbURL, "List will show all commands the bot understands.", helpColor,
			format(createFields("EXAMPLE", cmdInfo.Prefix+"list", true)))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
	case "pong":
		msg := cmdInfo.createMsgEmbed("Pong", helpThumbURL, "Playing with pong.", helpColor, format(
			createFields("EXAMPLE", cmdInfo.Prefix+"pong", true),
		))
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
