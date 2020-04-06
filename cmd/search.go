package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/daisymae/models"
)

// Search will look up a possible insect or fish in the database and display to the user
func Search(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) == 1 {
		return
	}
	if strings.ToLower(cmdInfo.CmdOps[1]) == "north" || strings.ToLower(cmdInfo.CmdOps[1]) == "south" {
		// ByMonth search
		byMonth(cmdInfo.CmdOps[1:], cmdInfo)
		return
	}
	formatStr := toLowerAndFormat(cmdInfo.CmdOps[1:])
	entry, err := cmdInfo.Service.Entry.ByName(formatStr, "bug_and_fish")
	searchItem := formatName(cmdInfo.CmdOps[1:])
	if err != nil {
		// If entry was not found in database
		word := strings.Split(searchItem, " ")
		vals := cmdInfo.Service.Entry.FindLike(toLowerAndFormat(word), "bug_and_fish")
		var fields []*discordgo.MessageEmbedField
		for _, val := range vals {
			fields = append(fields, createFields(strings.Title(val.Type), strings.Title(removeUnderscore(val.Name)), true))
		}
		if len(fields) == 0 {
			// If no similar entries were found
			msg := cmdInfo.createMsgEmbed(searchItem, errThumbURL, "No similar entries found.", errColor, fields)
			cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
			return
		}
		msg := cmdInfo.createMsgEmbed(searchItem, errThumbURL, "Entry Not Found in Database... Perhaps you meant...?", errColor, fields)
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
		return
	}
	nHemi, sHemi := parseHemi(entry.NorthSt, entry.NorthEnd, entry.SouthSt, entry.SouthEnd)
	fields := format(
		createFields("Price", strconv.Itoa(entry.SellPrice)+" Bells", true),
		createFields("Location", removeUnderscore(entry.Location), true),
		createFields("Time", removeUnderscore(entry.Time), false),
		createFields("Northern Hemisphere", nHemi, false),
		createFields("Southern Hemisphere", sHemi, false),
	)
	msg := cmdInfo.createMsgEmbed(searchItem, entry.Image, strings.Title(entry.Type), searchColor, fields)
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
}

// Utility func which returns true when the given slice
// contains "bug" or "fish"
func checkValid(args []string) bool {
	for _, item := range args {
		if item == "bug" || item == "fish" {
			return true
		}
	}
	return false
}

// byMonth prints out glob searches based on hemisphere, current month, and entry type
//
// Default month gathered by Go will depend on local time settings on server.
//
// E.g. if your server is on Pacific time and it's the last day of April, there will
// be a mismatch for someone querying on Australian time
func byMonth(cmds []string, cmdInfo CommandInfo) {
	if !checkValid(cmds) {
		cmdInfo.createMsgEmbed("Error", errThumbURL, "Please specify 'bug' or 'fish'", errColor,
			format(createFields("EXAMPLE", cmdInfo.Prefix+"search south fish", true)))
		return
	}
	var entries []models.Entry
	if cmds[0] == "north" {
		if strings.ToLower(cmds[1]) == "bug" {
			entries = cmdInfo.Service.Entry.ByMonth("north_hemi_months", time.Now().Month().String(), "bug")
		} else {
			entries = cmdInfo.Service.Entry.ByMonth("north_hemi_months", time.Now().Month().String(), "fish")
		}
	} else {
		if strings.ToLower(cmds[1]) == "bug" {
			entries = cmdInfo.Service.Entry.ByMonth("south_hemi_months", time.Now().Month().String(), "bug")
		} else {
			entries = cmdInfo.Service.Entry.ByMonth("south_hemi_months", time.Now().Month().String(), "fish")
		}
	}
	var fields []*discordgo.MessageEmbedField
	for _, val := range entries {
		fields = append(fields, createFields(strings.Title(removeUnderscore(val.Location)), strings.Title(removeUnderscore(val.Name)), true))
	}
	massPrint(fields, "Search By Hemisphere & Current Month", strings.Title(cmds[0]), cmdInfo)
}

// massPrints splits large slices into separate discord embed print statements specifically
// for search results
func massPrint(fields []*discordgo.MessageEmbedField, title, desc string, cmdInfo CommandInfo) {
	for i := 0; i < len(fields); i += 21 {
		j := i + 21
		if j > len(fields) {
			j = len(fields)
		}
		msg := cmdInfo.createMsgEmbed(title, blobSThumbURL, desc, searchColor, fields[i:j])
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
	}
}

// utility func to parse hemisphere data and return as
// useful text users can read
func parseHemi(ns, ne, st, se string) (string, string) {
	northSt := strings.Split(ns, "|")
	northEnd := strings.Split(ne, "|")
	southSt := strings.Split(st, "|")
	southEnd := strings.Split(se, "|")
	if len(northSt) == 1 && len(southSt) == 1 {
		return formatDate(northSt[0] + " " + northEnd[0]), formatDate(southSt[0] + " " + southEnd[0])
	}
	var northMonths []string
	var southMonths []string
	if len(northSt) > 1 && len(northEnd) > 1 {
		for i := 0; i < len(northSt); i++ {
			northMonths = append(northMonths, northSt[i]+" "+northEnd[i])
		}
	}
	if len(southSt) > 1 && len(southEnd) > 1 {
		for i := 0; i < len(southSt); i++ {
			southMonths = append(southMonths, southSt[i]+" "+southEnd[i])
		}
	}
	return wrapDate(northMonths, southMonths)
}

// utility func which wraps entries where there are multiple
// location months
//
// E.g. May to June, September to November
func wrapDate(north, south []string) (string, string) {
	var n []string
	var s []string
	for i := 0; i < len(north); i++ {
		n = append(n, formatDate(north[i]))
	}
	for i := 0; i < len(south); i++ {
		s = append(s, formatDate(south[i]))
	}
	return strings.Join(n, ", "), strings.Join(s, ", ")
}

// utility func to format dates
//
// This function parses date strings that have been turned into
// month names and appends the month with 'to'
func formatDate(date string) string {
	d := strings.Split(date, " ")
	if len(d) == 1 {
		return ""
	}
	var dateS []string
	for i := 0; i < len(d); i++ {
		if contains(dateS, getMonth(d[i])) {
			continue
		}
		dateS = append(dateS, getMonth(d[i]))
	}
	return strings.Join(dateS, " to ")
}

// Helper func
//
// contains returns true if a string is found inside a string slice
func contains(c []string, want string) bool {
	for _, e := range c {
		if e == want {
			return true
		}
	}
	return false
}

// helper func to return a month name in string format
// from integer format
func getMonth(monthInt string) string {
	monthNames := []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}
	m, err := strconv.Atoi(monthInt)
	if err != nil {
		return ""
	}
	return monthNames[m-1]
}

// Helper func which formats argument list to match database keys
//
// Returns a valid database key format
func toLowerAndFormat(args []string) string {
	var endStr []string
	for _, word := range args {
		endStr = append(endStr, strings.ToLower(word))
	}
	if len(args) > 1 {
		return strings.Join(endStr, "_")
	}
	return endStr[0]
}

// helper func to format search names for pretty printing
//
// Example: []{"tHiS", "iS", "A", "WORD"}
//
// will become: This Is A Word
func formatName(str []string) string {
	// BUG: Go's string package has a bug in which unicode punctuation aren't
	// accounted for - will fail "with apostrophes" test
	// FIXME: Either find a workaround for this or wait til new golang patch
	var endStr []string
	for _, word := range str {
		tmp := strings.Title(strings.ToLower(word))
		endStr = append(endStr, tmp)
	}
	return strings.Join(endStr, " ")
}

// utility func to replace all underscores with a space
//
// normalization for users
func removeUnderscore(str string) string {
	return strings.ReplaceAll(str, "_", " ")
}
