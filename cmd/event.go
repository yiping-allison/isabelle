package cmd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/isabelle/models"
)

type newEvent struct {
	// Name of the event
	Name string

	// Link to the event image
	Img string

	// Limit specifies max amount of queue individuals for event host
	Limit string

	// Custom message set by event hosts - must be within character ranges
	Msg string
}

const (
	diyURL      string = "https://cdn.discordapp.com/attachments/693564368423616562/696635733368111144/DIY.png"
	saharahURL  string = "https://vignette.wikia.nocookie.net/animalcrossing/images/d/d7/Acnl-saharah.png/revision/latest/scale-to-width-down/344?cb=20130707101048"
	celesteURL  string = "https://vignette.wikia.nocookie.net/animalcrossing/images/a/a5/Acnl-celeste.png/revision/latest/scale-to-width-down/350?cb=20130703203412"
	daisymaeURL string = "https://vignette.wikia.nocookie.net/animalcrossing/images/8/85/Daisy_Mae.png/revision/latest?cb=20200220213944"
	meteorURL   string = "https://static0.srcdn.com/wordpress/wp-content/uploads/2020/04/animal-crossing-new-horizon-meteor-shower.jpg"
	kicksURL    string = "https://vignette.wikia.nocookie.net/animalcrossing/images/2/29/200px-Kicks_3DS.png/revision/latest/scale-to-width-down/350?cb=20140718172000"
)

// Event will parse through event commands and display embed with
// role ping
func Event(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) == 1 {
		// No arguments supplied - error
		msg := cmdInfo.createMsgEmbed(
			"Error: No Arguments", errThumbURL, "You must supply arguments to the command.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"2\" msg=\"bonsai tree recipe\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
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
			cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
			return
		}
		queue := cmdInfo.Service.Event.GetQueue(cmdInfo.CmdOps[1])
		fields := queueToFields(queue)
		msg := cmdInfo.createMsgEmbed(
			"Current Queue", queueThumbURL, "Queue ID: "+cmdInfo.CmdOps[1], eventColor, fields)
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
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
	case "turnip":
		event = parseCmd(cmd, "Turnip - High Sell Price", daisymaeURL)
	case "kicks":
		event = parseCmd(cmd, "Kicks", kicksURL)
	default:
		return
	}

	if event == nil {
		// Couldn't create an event - error
		fmt.Println(event)
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Try checking your command's syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"5\" msg=\"ironwood bed\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	limit, err := strconv.Atoi(event.Limit)
	if err != nil || limit > 20 || limit < 1 {
		// Error - couldn't convert limit value into a number (limit MUST be a number)
		// Limit must be within bounds 1 - 20
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Your limit must be a valid number", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"5\" msg=\"ironwood bed\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	if !validMsg(event.Msg, eventName) {
		// Error - message must be within 50 or 100 characters
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Your message must be within valid length", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"5\" msg=\"ironwood bed\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	if cmdInfo.Service.User.LimitEvent(cmdInfo.Msg.Author) {
		// Error - Cannot create anymore events
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "You already have the max amount of events.", errColor,
			format(
				createFields("Suggestion", "Either end one of your events or wait until your events are finished before creating another.", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	user := cmdInfo.Msg.Author
	// if user doesn't exist in rep database, create a new one
	if !cmdInfo.Service.Rep.Exists(user.ID) {
		cmdInfo.newRep(user.ID)
	}

	// if the user doesn't currently exist in tracking, create a new one
	if !cmdInfo.Service.User.UserExists(user) {
		cmdInfo.Service.User.AddUser(user)
	}

	// generate a random id with at least 4 digits
	id := generateID(1000, 9999)

	if cmdInfo.Service.Event.EventExists(id) {
		// There's an event with the same ID
		msg := cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "There is already an event with this ID; Please try again.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"5\" msg=\"ironwood bed\"", false),
			))
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.BotChID, msg)
		return
	}

	// Add the event to tracking
	cmdInfo.Service.Event.AddEvent(user, id, limit)
	// record expiration time
	expire := cmdInfo.Service.Event.GetExpiration(id)
	cmdInfo.Service.User.AddEvent(user, id, expire)
	// retrieve reputation
	rep := cmdInfo.Service.Rep.GetRep(user.ID)

	msg := cmdInfo.createMsgEmbed(
		"Event: "+event.Name,
		event.Img,
		"Queue ID: "+id,
		eventColor,
		format(
			createFields("Hosted By", user.Mention(), true),
			createFields("Reputation", strconv.Itoa(rep), true),
			createFields("Limit", event.Limit, false),
			createFields("Message", event.Msg, false),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.ListingID, msg)
	cmdInfo.Ses.ChannelMessageSend(cmdInfo.BotChID, "Listing Posted!")
}

// validMsg checks if a message is within text length
func validMsg(msg, event string) bool {
	max := 50
	if event == "saharah" {
		max = 100
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
