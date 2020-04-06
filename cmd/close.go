package cmd

import (
	"strings"
)

// Close will attempt to parse event or trade types
// and close the event by ID through host request
func Close(cmdInfo CommandInfo) {
	// Arg at cmdInfo.CmdOps[2] should be specifiying event or trade
	if len(cmdInfo.CmdOps) != 3 {
		// wrong arguments - error
		msg := cmdInfo.createMsgEmbed(
			"Error: Incorrect Arguments Supplied", errThumbURL, "Please check your syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"close event 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}
	t := strings.ToLower(cmdInfo.CmdOps[1])
	switch t {
	case "event":
		// close an event
		closeEvent(cmdInfo.CmdOps[2], cmdInfo)
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
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}
	// close the event
	cmdInfo.Service.Event.Close(eventID)
	embed := cmdInfo.createMsgEmbed(
		"Successfully Removed Event "+eventID+" from listings!", checkThumbURL, "Thank you for hosting!",
		successColor, format(
			createFields("Suggestion", "If you are planning on opening another event, it is safe to do so now.", false),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, embed)
}
