package cmd

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Offer handles the offer capabilities to trade
// options
func Offer(cmdInfo CommandInfo) {
	// only argument to offer must be trade id
	if len(cmdInfo.CmdOps) < 3 {
		// wrong arguments to command
		return
	}
	id := cmdInfo.CmdOps[1]
	user := cmdInfo.Msg.Author
	// if user doesn't exist in rep database, create a new one
	if !cmdInfo.Service.Rep.Exists(user.ID) {
		cmdInfo.newRep(user.ID)
	}

	if !cmdInfo.Service.Trade.Exists(id) {
		// error - no trade with that id exists
		msg := cmdInfo.createMsgEmbed(
			"Error: Trade Event Does Not Exist", errThumbURL, "Trade ID: "+id,
			errColor, format(
				createFields("Suggestion", "Try checking if you supplied the correct Trade ID.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}
	offer := strings.Join(cmdInfo.CmdOps[2:], " ")
	offer = strings.Title(offer)

	// add offer to tracking
	err := cmdInfo.Service.Trade.AddOffer(id, offer, user)
	if err != nil {
		// error - user already offered
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Add To Trade", errThumbURL, strings.Title(err.Error()),
			errColor, format(
				createFields("User", user.Mention(), true),
				createFields("Suggestion", "You can remove your existing offer and try again.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	// if user doesn't exist, create new user
	if !cmdInfo.Service.User.UserExists(user) {
		// create new user tracking
		cmdInfo.Service.User.AddUser(user)
	}

	// add user to user tracking
	expire := cmdInfo.Service.Trade.GetExpiration(id)
	cmdInfo.Service.User.AddOffer(id, user, expire)
	// get original trade host info
	host := cmdInfo.Service.Trade.GetHost(id)
	rep := cmdInfo.Service.Rep.GetRep(user.ID)
	// print success msg
	embed := cmdInfo.createMsgEmbed(
		"Successfully Added Offer!", checkThumbURL, "Trade ID: "+id,
		successColor, format(
			createFields("Offerer", user.Mention(), true),
			createFields("Offer Item", offer, true),
			createFields("Reputation", strconv.Itoa(rep), true),
			createFields("Suggestion", "Please Wait Until Trader Makes a Decision. Thank you!", false),
		))
	cplx := &discordgo.MessageSend{
		Content: host.Mention() + ": A new person has offered to your trade!",
		Embed:   embed,
	}
	cmdInfo.Ses.ChannelMessageSendComplex(cmdInfo.BotChID, cplx)
}
