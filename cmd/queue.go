package cmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Queue handles the queue-ing system; queue ids are retrieved from
// event ids defined in event.go
func Queue(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) == 1 {
		// Not enough arguments: only ;queue
		msg := cmdInfo.createMsgEmbed(
			"Error: Syntax", errThumbURL, "Not enough arguments supplied.",
			errColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue [queue_ID]", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	if !cmdInfo.Service.Event.EventExists(cmdInfo.CmdOps[1]) {
		// Error - event doesn't exists
		msg := cmdInfo.createMsgEmbed(
			"Error: Event Not Found", errThumbURL, "ID: "+cmdInfo.CmdOps[1],
			errColor, format(
				createFields("EXAMPLE", cmdInfo.Prefix+"queue [queue_ID]", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	user := cmdInfo.Msg.Author
	if cmdInfo.Service.User.LimitQueue(user) {
		// Error - max queue reached
		msg := cmdInfo.createMsgEmbed(
			"Error: Max Queue reached", errThumbURL, "You already reached the max queue.",
			errColor, format(
				createFields("Suggestion", "Remove a queue before trying to add another queue.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	// Add user to queue
	host, err := cmdInfo.Service.Event.AddToQueue(user, cmdInfo.CmdOps[1])
	if err != nil {
		// Check error
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Add To Queue", errThumbURL, strings.Title(err.Error()),
			errColor, format(
				createFields("User", user.Mention(), true),
				createFields("EXAMPLE", cmdInfo.Prefix+"queue 1234", true),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	// if user doesn't exist in user tracking, create a new one
	if !cmdInfo.Service.User.UserExists(user) {
		cmdInfo.Service.User.AddUser(user)
	}

	// Add queue to user tracking
	cmdInfo.Service.User.AddQueue(cmdInfo.CmdOps[1], user, cmdInfo.Service.Event.GetExpiration(cmdInfo.CmdOps[1]))

	embed := cmdInfo.createMsgEmbed(
		"Successfully Added to Queue!", checkThumbURL, "Queue ID: "+cmdInfo.CmdOps[1],
		successColor, format(
			createFields("User", user.Mention(), true),
			createFields("Please Wait Until You're Pinged or Messaged!", "Thank you!", true),
		))
	cplx := &discordgo.MessageSend{
		Content: host.Mention() + ": A new person has joined your queue!",
		Embed:   embed,
	}
	cmdInfo.Ses.ChannelMessageSendComplex(cmdInfo.Msg.ChannelID, cplx)
}
