package cmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Rep will allow server members to update reputation points
// on another member
func Rep(cmdInfo CommandInfo) {
	userID := stripPing(cmdInfo.CmdOps[1])
	if !cmdInfo.Service.Rep.Exists(userID) {
		// if the user doesn't exist in rep database, create a new one
		cmdInfo.newRep(userID)
	}
	// generate random 4 digit ID for acception event
	id := generateID(1000, 9999)
	cmdInfo.Service.Rep.AddRep(userID, id)
	userMsg := strings.Join(cmdInfo.CmdOps[2:], " ")
	// print rep msg
	msg := cmdInfo.createMsgEmbed(
		"Reputation Application", thumbThumbURL, "App ID: "+id,
		appColor, format(
			createFields("Nominee", mentionUser(userID), true),
			createFields("Message", userMsg, true),
			createFields("Note", "The mods will try to process this app ASAP. Thank you for submitting!", false),
		))
	cplx := &discordgo.MessageSend{
		Content: mentionRole(cmdInfo.AdminRole) + ": New Reputation App!",
		Embed:   msg,
	}
	cmdInfo.Ses.ChannelMessageSendComplex(cmdInfo.AppID, cplx)
}

// mentionUser is a helper func which mentions a user by ID
func mentionUser(user string) string {
	return "<@!" + user + ">"
}

// mention is a helper func to mention role IDs
func mentionRole(role string) string {
	return "<@&" + role + ">"
}

// stripPing is a helper func which turns discord pings into a regular user id
// string
func stripPing(ping string) string {
	id := strings.TrimPrefix(ping, "<@!")
	id = strings.TrimSuffix(id, ">")
	return id
}
