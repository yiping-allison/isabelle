package cmd

import (
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

	if !cmdInfo.Service.Trade.Exists(id) {
		// error - no trade with that id exists
		msg := cmdInfo.createMsgEmbed(
			"Error: Trade Event Does Not Exist", errThumbURL, "Trade ID: "+id,
			errColor, format(
				createFields("Suggestion", "Try checking if you supplied the correct Trade ID.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}
	offer := strings.Join(cmdInfo.CmdOps[2:], " ")
	offer = strings.Title(offer)

	// add offer to tracking
	err := cmdInfo.Service.Trade.AddOffer(id, offer, cmdInfo.Msg.Author)
	if err != nil {
		// error - user already offered
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Add "+cmdInfo.Msg.Author.String()+" To Trade", errThumbURL, strings.Title(err.Error()),
			errColor, format(
				createFields("Suggestion", "You can remove your existing offer and try again.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	host := cmdInfo.Service.Trade.GetHost(id)
	// print success msg
	embed := cmdInfo.createMsgEmbed(
		"Successfully Added "+cmdInfo.Msg.Author.String()+"'s Offer!", checkThumbURL, "Trade ID: "+id,
		successColor, format(
			createFields("Offer", offer, true),
			createFields("Suggestion", "Please Wait Until Trader Makes a Decision. Thank you!", false),
		))
	cplx := &discordgo.MessageSend{
		Content: host.Mention() + ": A new person has offered to your trade!",
		Embed:   embed,
	}
	cmdInfo.Ses.ChannelMessageSendComplex(cmdInfo.Msg.ChannelID, cplx)
}
