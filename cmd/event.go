package cmd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type newEvent struct {
	Name  string
	Img   string
	Limit string
	Msg   string
}

const (
	diyURL      string = "https://static1.srcdn.com/wordpress/wp-content/uploads/2020/02/Animal-Crossing-New-Horizons-Crafting-Guide.jpg"
	saharahURL  string = "https://www.imore.com/sites/imore.com/files/styles/large/public/field/image/2020/03/animal-crossing-new-horizons-switch-confirmed-characters-saharah.jpg?itok=iFqeqqRc"
	celesteURL  string = "https://www.imore.com/sites/imore.com/files/styles/large/public/field/image/2020/03/animal-crossing-new-horizons-switch-confirmed-characters-celeste.jpg?itok=-2ib_Yfs"
	daisymaeURL string = "https://www.imore.com/sites/imore.com/files/styles/large/public/field/image/2020/03/animal-crossing-new-horizons-switch-confirmed-characters-daisy-mae.jpg?itok=yPB96-2n"
)

// Event will parse through event commands and display embed with
// role ping
func Event(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) == 1 {
		// No arguments supplied - error
		cmdInfo.createMsgEmbed(
			"Error: No Arguments", errThumbURL, "You must supply arguments to the command.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy time=\"[Your time]\" msg=\"[Your message]\"", false),
			))
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
	default:
		return
	}
	if event == nil {
		// Couldn't create an event - error
		cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Try checking your command's syntax.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit]\" msg=\"[Your message]\"", false),
			))
		return
	}
	if cmdInfo.Service.Event.EventExists(cmdInfo.Msg.ID) {
		// There's an event with the same ID
		cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "There is already an event with this ID; Please try again.", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit]\" msg=\"[Your message]\"", false),
			))
		return
	}
	limit, err := strconv.Atoi(event.Limit)
	if err != nil {
		// Error - couldn't convert limit value into a number (limit MUST be a number)
		cmdInfo.createMsgEmbed(
			"Error: Couldn't Create Event", errThumbURL, "Your limit must be a valid number", errColor,
			format(
				createFields("EXAMPLE", cmdInfo.Prefix+"event diy limit=\"[Your limit]\" msg=\"[Your message]\"", false),
			))
		return
	}
	cmdInfo.Service.Event.AddEvent(cmdInfo.Msg.Author, cmdInfo.Msg.ID[10:14], limit)
	cmdInfo.createMsgEmbed(
		"Event: "+event.Name,
		event.Img,
		"To Queue: "+cmdInfo.Msg.ID[10:14],
		eventColor,
		format(
			createFields("Limit", event.Limit, true),
			createFields("Message", event.Msg, false),
		))
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
	fmt.Println(fullCmd)
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
// i.e. a validEvent must have two keywords (ONLY), time and msg
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
