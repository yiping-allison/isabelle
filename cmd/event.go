package cmd

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/isabelle/models"
)

type newEvent struct {
	Name  string
	Img   string
	Limit string
	Msg   string
}

const (
	diyURL      string = "https://cdn.discordapp.com/attachments/693564368423616562/696635733368111144/DIY.png"
	saharahURL  string = "https://vignette.wikia.nocookie.net/animalcrossing/images/d/d7/Acnl-saharah.png/revision/latest/scale-to-width-down/344?cb=20130707101048"
	celesteURL  string = "https://vignette.wikia.nocookie.net/animalcrossing/images/a/a5/Acnl-celeste.png/revision/latest/scale-to-width-down/350?cb=20130703203412"
	daisymaeURL string = "https://vignette.wikia.nocookie.net/animalcrossing/images/8/85/Daisy_Mae.png/revision/latest?cb=20200220213944"
	meteorURL   string = "https://static0.srcdn.com/wordpress/wp-content/uploads/2020/04/animal-crossing-new-horizon-meteor-shower.jpg"
)

// Event will parse through event commands and display embed with
// role ping
func Event(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) == 1 {
		// No arguments supplied - error
		msg := cmdInfo.createMsgEmbed(
			"Error: No Arguments", errThumbURL, "You must supply arguments to the command.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy time=\"[Your time]\" msg=\"[Your message]\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	if _, err := strconv.Atoi(cmdInfo.CmdOps[1]); err == nil {
		// This is a list command - print all users in queue
		if !cmdInfo.Service.Event.EventExists(cmdInfo.CmdOps[1]) {
			// event doesn't exists - print error
			msg := cmdInfo.createMsgEmbed(
				"Error: "+cmdInfo.CmdOps[1]+" Event Does Not Exist", errThumbURL, "Please check your event ID.", errColor,
				format(
					createFields("EXAMPLE", cmdInfo.Prefix+"event 1234", false),
				))
			cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
			return
		}
		queue := cmdInfo.Service.Event.GetQueue(cmdInfo.CmdOps[1])
		fields := queueToFields(queue)
		msg := cmdInfo.createMsgEmbed(
			"Current Queue", queueThumbURL, "Queue ID: "+cmdInfo.CmdOps[1], eventColor, fields)
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	eventName := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(cmdInfo.CmdOps[1])), " ", "")
	cmd := strings.Join(cmdInfo.CmdOps[2:], " ")
	var event *newEvent
	switch eventName {
	case "celeste":
		event = parseCmd(cmd, "Celeste", celesteURL)
	case "daisymae":
		event = parseCmd(cmd, "Daisy Mae", daisymaeURL)
	case "saharah":
		event = parseCmd(cmd, "Saharah", saharahURL)
	case "diy":
		event = parseCmd(cmd, "DIY", diyURL)
	case "meteor":
		event = parseCmd(cmd, "Meteor Shower", meteorURL)
	default:
		return
	}

	if event == nil {
		// Couldn't create an event - error
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Try checking your command's syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit (number within 1-20)]\" msg=\"[Your message]\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	if cmdInfo.Service.Event.EventExists(cmdInfo.Msg.ID) {
		// There's an event with the same ID
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "There is already an event with this ID; Please try again.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit (number within 1-20)]\" msg=\"[Your message]\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	limit, err := strconv.Atoi(event.Limit)
	if err != nil || limit > 20 || limit < 1 {
		// Error - couldn't convert limit value into a number (limit MUST be a number)
		// Limit must be within bounds 1 - 20
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Your limit must be a valid number", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit (number within 1-20)]\" msg=\"[Your message]\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	if !validMsg(event.Msg, eventName) {
		// Error - message must be within 50 or 100 characters
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Your message must be within valid length", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit (number within 1-20)]\" msg=\"[Your message (50 chars for all events except Saharah at 100 chars)]\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}

	cmdInfo.Service.Event.AddEvent(cmdInfo.Msg.Author, cmdInfo.Msg.ID[10:14], limit)
	msg := cmdInfo.createMsgEmbed(
		"Event: "+event.Name,
		event.Img,
		"Queue ID: "+cmdInfo.Msg.ID[10:14],
		eventColor,
		format(
			createFields("Hosted By", cmdInfo.Msg.Author.String(), true),
			createFields("Limit", event.Limit, true),
			createFields("Message", event.Msg, false),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
}

// validMsg checks if a message is within text length
func validMsg(msg, event string) bool {
	var max int
	if event == "saharah" {
		max = 100
	} else {
		max = 50
	}
	if len([]rune(msg)) > max {
		return false
	}
	return true
}

// queueToFields is specifically made to create field embeds based on variable
// number to type QueueUser
func queueToFields(user *[]models.QueueUser) []*discordgo.MessageEmbedField {
	var f []*discordgo.MessageEmbedField
	for i, u := range *user {
		tmp := createFields(
			"Queue Number: "+strconv.Itoa(i+1),
			u.DiscordUser.String(),
			true,
		)
		f = append(f, tmp)
	}
	return f
}

// parseCmd will attempt to parse a user's set event command.
//
// If successful, it will return a pointer to the new event
//
// else, it will return nil
func parseCmd(fullCmd, name, imgURL string) *newEvent {
	var event newEvent
	event.Name = name
	event.Img = imgURL
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsSpace(c) && !unicode.IsNumber(c)
	}
	commandBytes := bytes.FieldsFunc([]byte(fullCmd), f)
	var parsedCmd []string
	for _, e := range commandBytes {
		parsedCmd = append(parsedCmd, strings.ToLower(strings.TrimSpace(string(e))))
	}
	if !validEvent(parsedCmd) || len(parsedCmd) != 4 {
		return nil
	}
	for i := 0; i < len(parsedCmd); i += 2 {
		if strings.ToLower(parsedCmd[i]) == "limit" {
			event.Limit = strings.ToUpper(parsedCmd[i+1])
		}
		if strings.ToLower(parsedCmd[i]) == "msg" {
			event.Msg = strings.Title(parsedCmd[i+1])
		}
	}
	return &event
}

// validEvent will check if the user set enough fields with valid naming.
//
// i.e. a validEvent must have two keywords (ONLY), limit and msg
func validEvent(args []string) bool {
	var set []string
	for _, e := range args {
		if e == "limit" || e == "msg" {
			set = append(set, e)
		}
	}
	if len(set) != 2 {
		return false
	}
	return true
}
