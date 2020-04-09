package cmd

import (
	"github.com/bwmarrin/discordgo"
)

// Unregister allows a queue user to remove themselves from the queue
func Unregister(cmdInfo CommandInfo) {
	// cmdInfo.CmdOps[1:] starts after ;unregister
	if len(cmdInfo.CmdOps) != 3 {
		// Error - not enough arguments
		msg := cmdInfo.createMsgEmbed(
			"Error: Wrong Arguments", errThumbURL, "Try checking your syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"unregister event [event ID]", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}
	args := cmdInfo.CmdOps[1:]
	switch args[0] {
	case "event":
		cmdInfo.removeFromEvent(args[1], cmdInfo.Msg.Author)
	}
}

// helper func which removes queue users from an event queue
func (c CommandInfo) removeFromEvent(eventID string, user *discordgo.User) {
	if !c.Service.Event.EventExists(eventID) {
		// event does not exist
		msg := c.createMsgEmbed(
			"Error: Event does not exist", errThumbURL, "Event ID: "+eventID, errColor,
			format(
				createFields("Suggestion", "Try checking if you supplied the correct Event ID", false),
			))
		c.Ses.ChannelMessageSendEmbed(c.Msg.ChannelID, msg)
		return
	}

	// remove user
	c.Service.Event.Remove(eventID, user)
	// Remove tracking on user
	c.Service.User.RemoveQueue(eventID, user)

	// successfully removed user
	msg := c.createMsgEmbed(
		"Removed "+c.Msg.Author.String()+" from Event", checkThumbURL, "Queue ID: "+c.CmdOps[2],
		successColor, format(
			createFields("Suggestion", "Feel free to queue for any other events or create your own.", false),
		))
	c.Ses.ChannelMessageSendEmbed(c.Msg.ChannelID, msg)
}
