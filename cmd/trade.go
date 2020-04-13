package cmd

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/isabelle/models"
)

type trade struct {
	// item user is willing to trade
	item string

	// msg (preferably) about what item they're looking for to trade
	msg string
}

// Trade will handle trade options within the server
//
// Trade is meant to be paired along with offer.go in order
// for members to trade and offer items among each other
func Trade(cmdInfo CommandInfo) {
	user := cmdInfo.Msg.Author

	if _, err := strconv.Atoi(cmdInfo.CmdOps[1]); err == nil {
		// This is a list command - print all currently offered to tradeID
		offers := cmdInfo.Service.Trade.GetAllOffers(cmdInfo.CmdOps[1])
		printTradeList(offers, cmdInfo, cmdInfo.CmdOps[1])
		return
	}

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
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	// attempt to parse command
	t := parseTradeCmd(strings.Join(cmdInfo.CmdOps[1:], " "))
	if t == nil {
		// error parse failed
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Trade", errThumbURL, "Syntax error", errColor,
			format(
				createFields("Suggestion", "Try checking if you input the command correctly.", false),
				createFields("EXAMPLE", cmdInfo.Prefix+"trade item=\"blue mountain coffee\" msg=\"looking for geisha coffee\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
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
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
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
	msg := cmdInfo.createMsgEmbed(
		"Trade", tradeThumbURL, user.Mention(), tradeColor,
		format(
			createFields("Trade ID", id, true),
			createFields("Reputation", strconv.Itoa(reps), true),
			createFields("Trade Listing", strings.Title(t.item), false),
			createFields("Message", strings.Title(t.msg), false),
		),
	)
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.ListingID, msg)
}

// printTradeList handles printing large amounts of trade offers (since trade offers has no limit)
func printTradeList(offers []models.TradeOfferer, cmdInfo CommandInfo, tradeID string) {
	var fields []*discordgo.MessageEmbedField
	for _, o := range offers {
		fields = append(fields, createFields(o.User.String(), o.Offer, true))
	}
	for i := 0; i < len(fields); i += 15 {
		j := i + 15
		if j > len(fields) {
			j = len(fields)
		}
		msg := cmdInfo.createMsgEmbed("Total Offers", tradeThumbURL, "TradeID: "+tradeID, tradeColor, fields[i:j])
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
	}
}

// parseTradeCmd will take a full command string and return a trade object
// if the command was correctly parsed
//
// else, nil
func parseTradeCmd(fullCmd string) *trade {
	f := func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r)
	}
	bytes := bytes.FieldsFunc([]byte(fullCmd), f)
	var parsedCmd []string
	for _, e := range bytes {
		parsedCmd = append(parsedCmd, strings.ToLower(strings.TrimSpace(string(e))))
	}
	if len(parsedCmd) != 4 || !validTrade(parsedCmd) {
		// error - syntax not parsed correctly
		return nil
	}
	var t trade
	for i := 0; i < len(parsedCmd); i += 2 {
		if parsedCmd[i] == "item" {
			t.item = parsedCmd[i+1]
			continue
		}
		if parsedCmd[i] == "msg" {
			t.msg = parsedCmd[i+1]
			continue
		}
	}
	return &t
}

// validTrade checks if the given command contains the two keywords:
// item & msg
func validTrade(cmds []string) bool {
	var tmp []string
	for _, i := range cmds {
		if i == "item" || i == "msg" {
			tmp = append(tmp, i)
		}
	}
	return len(tmp) == 2
}
