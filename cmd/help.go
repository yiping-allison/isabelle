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
		cmdInfo.createMsgEmbed("Error", errThumbURL, "You must query a valid command",
			errColor, format(createFields("EXAMPLE", cmdInfo.Prefix+"help search", true)))
		return
	}
	full := strings.Join(cmdInfo.CmdOps[1:], " ")
	if !find(full, cmdInfo) {
		cmdInfo.createMsgEmbed(full, errThumbURL, "Command Not Found", errColor, format(
			createFields("To List All Commands: ", cmdInfo.Prefix+"list", true),
		))
		return
	}
	// Valid commands
	switch full {
	case "search":
		cmdInfo.createMsgEmbed("Search", helpThumbURL, "Search will look up an item from New Horizon's bug and fish database.",
			helpColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"search emperor butterfly", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"search north bug", true),
			))
	case "list":
		cmdInfo.createMsgEmbed("List", helpThumbURL, "List will show all commands the bot understands.", helpColor,
			format(createFields("EXAMPLE", cmdInfo.Prefix+"list", true)))
	case "pong":
		cmdInfo.createMsgEmbed("Pong", helpThumbURL, "Playing with pong.", helpColor, format(
			createFields("EXAMPLE", cmdInfo.Prefix+"pong", true),
		))
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
