package cmd

import (
	"strconv"
	"strings"
)

// Trade will handle trade options within the server
//
// Trade is meant to be paired along with offer.go in order
// for members to trade and offer items among each other
func Trade(cmdInfo CommandInfo) {
	user := cmdInfo.Msg.Author
	// if user doesn't exist in rep database, create a new one
	if !cmdInfo.Service.Rep.Exists(user.ID) {
		cmdInfo.newRep(user.ID)
	}

	// If the user currently doesn't exist in server tracking, make a new one
	if !cmdInfo.Service.User.UserExists(user) {
		// Create a user
		cmdInfo.Service.User.AddUser(user)
	}

	if cmdInfo.Service.User.LimitTrade(user.ID) {
		// Max trades created - can't make anymore
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Trade", errThumbURL, "You already have the max trade events.", errColor,
			format(
				createFields("Suggestion", "Either end one of your trades or wait until they are finished before creating another.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	// generate trade id
	id := generateID(1000, 9999)

	if cmdInfo.Service.Trade.Exists(id) {
		// error - id exists
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Trade", errThumbURL, "This ID already exists.", errColor,
			format(
				createFields("Suggestion", "Try re-creating the trade event.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	// Add trade event
	cmdInfo.Service.Trade.AddTrade(id, user)
	// Add trade tracking to user
	expire := cmdInfo.Service.Trade.GetExpiration(id)
	cmdInfo.Service.User.AddTrade(user, id, expire)

	// retrieve reps from database
	reps := cmdInfo.Service.Rep.GetRep(user.ID)

	// Print Trade Offer
	offer := strings.Join(cmdInfo.CmdOps[1:], " ")
	msg := cmdInfo.createMsgEmbed(
		"Trade Offer", tradeThumbURL, user.String(), tradeColor,
		format(
			createFields("Trade ID", id, true),
			createFields("Reputation", strconv.Itoa(reps), true),
			createFields("Trade Listing", offer, false),
		),
	)
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
}
