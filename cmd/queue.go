package cmd

import "strings"

// Queue handles the queue-ing system; queue ids are retrieved from
// event ids defined in event.go
func Queue(cmdInfo CommandInfo) {
	// TODO: Add NAMES to queue print messages otherwise no one knows who was rejected
	// or accepted to queue
	if len(cmdInfo.CmdOps) == 1 {
		// Not enough arguments: only ;queue
		cmdInfo.createMsgEmbed(
			"Error: Syntax", errThumbURL, "Not enough arguments supplied.",
			errColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue [queue_ID]", true),
			))
		return
	}
	if !cmdInfo.Service.Event.EventExists(cmdInfo.CmdOps[1]) {
		// Error - event doesn't exists
		cmdInfo.createMsgEmbed(
			"Error: Event Not Found", errThumbURL, cmdInfo.CmdOps[1],
			errColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue [queue_ID]", true),
			))
		return
	}
	// Add user to queue
	err := cmdInfo.Service.Event.AddToQueue(cmdInfo.Msg.Author, cmdInfo.CmdOps[1])
	if err != nil {
		// Check error
		// TODO: Event Hosts cannot queue for their own event
		cmdInfo.createMsgEmbed(
			"Error: Couldn't Add You To Queue", errThumbURL, strings.Title(err.Error()),
			errColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue [queue_ID]", true),
			))
		return
	}
	cmdInfo.createMsgEmbed(
		"Successfully Added You to Queue!", checkThumbURL, cmdInfo.CmdOps[1],
		successColor, format(
			createFields("Please Wait Until You're Pinged", "Thank you!", true),
		))
}
