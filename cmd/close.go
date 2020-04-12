package cmd

import (
	"strings"
)

// Close will attempt to parse event or trade types
// and close the event by ID through host request
func Close(cmdInfo CommandInfo) {
	// Arg at cmdInfo.CmdOps[1] should be specifiying event or trade
	if len(cmdInfo.CmdOps) != 3 {
		// wrong arguments - error
		msg := cmdInfo.createMsgEmbed(
			"Error: Incorrect Arguments Supplied", errThumbURL, "Please check your syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"close event 1234", true),
				createFields("EXAMPLE", cmdInfo.Prefix+"close trade 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	t := strings.ToLower(cmdInfo.CmdOps[1])
	switch t {
	case "event":
		closeEvent(cmdInfo.CmdOps[2], cmdInfo)
	case "trade":
		closeTrade(cmdInfo.CmdOps[2], cmdInfo)
	}
}

// closeEvent will handle closing events by host
func closeEvent(eventID string, cmdInfo CommandInfo) {
	if !cmdInfo.Service.Event.EventExists(eventID) {
		// Error - event not found; can't close
		msg := cmdInfo.createMsgEmbed(
			"Error: Event Not Found", errThumbURL, "Event ID: "+eventID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied a valid Event ID.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	// store original host before removing event
	host := cmdInfo.Service.Event.GetHost(eventID)

	// attempt to close the event
	err := cmdInfo.Service.Event.Close(eventID, cmdInfo.AdminRole, cmdInfo.Msg.Author, cmdInfo.Msg.Member.Roles)
	if err != nil {
		// Error - Permission denied
		msg := cmdInfo.createMsgEmbed(
			"Error: You do not have permission to delete this event", errThumbURL, "Event ID: "+eventID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied the right Event ID.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	// remove all people tracking event
	cmdInfo.Service.User.RemoveAllQueue(eventID)
	// Remove user from tracking
	cmdInfo.Service.User.RemoveEvent(host, eventID)

	// print msg
	embed := cmdInfo.createMsgEmbed(
		"Successfully Removed Event "+eventID+" from listings!", checkThumbURL, "Thank you for hosting!",
		successColor, format(
			createFields("Host", host.Mention(), true),
			createFields("Suggestion", "If you are planning on opening another event, it is safe to do so now.", false),
			createFields("Suggestion", "If your event was deleted by a moderator, please make sure to follow event guidelines.", false),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.ListingID, embed)
}

// closeTrade is a helper func which closes a trade event and
// removes all trade tracking from the original user
func closeTrade(tradeID string, cmdInfo CommandInfo) {
	if !cmdInfo.Service.Trade.Exists(tradeID) {
		// error - trade event doesn't exist
		msg := cmdInfo.createMsgEmbed(
			"Error: Trade Not Found", errThumbURL, "Trade ID: "+tradeID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied a valid Trade ID.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	// get the original creator of the trade
	host := cmdInfo.Service.Trade.GetHost(tradeID)
	// attempt to close the trade
	err := cmdInfo.Service.Trade.Close(tradeID, cmdInfo.Msg.Author, cmdInfo.Msg.Member.Roles, cmdInfo.AdminRole)
	if err != nil {
		// error - user does not have permission to close event
		msg := cmdInfo.createMsgEmbed(
			"Error: You do not have permission to delete this trade", errThumbURL, "Trade ID: "+tradeID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied the right Trade ID.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}
	// remove host from user tracking
	cmdInfo.Service.User.RemoveTrade(tradeID, host)

	// print msg
	embed := cmdInfo.createMsgEmbed(
		"Successfully Removed Trade "+tradeID+" from listings!", checkThumbURL, "Thank you for hosting!",
		successColor, format(
			createFields("Host", host.Mention(), true),
			createFields("Suggestion", "If you are planning on opening another trade, it is safe to do so now.", false),
			createFields("Suggestion", "If your trade was deleted by a moderator, please make sure to follow trade guidelines.", false),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.ListingID, embed)
}
