package cmd

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Search will look up a possible insect or fish in the database and display to the user
func Search(cmdInfo CommandInfo) {
	// TODO: Refactor and Prettier Formatting?
	// TODO: List out top 3 Like matches in not found option
	if len(cmdInfo.CmdOps) == 1 {
		return
	}
	formatStr := toLowerAndFormat(cmdInfo.CmdOps[1:])
	entry, err := cmdInfo.Service.Entry.ByName(formatStr, "bug_and_fish")
	searchItem := formatName(cmdInfo.CmdOps[1:])
	if err != nil {
		emThumb := &discordgo.MessageEmbedThumbnail{
			URL:    "http://static2.wikia.nocookie.net/__cb20131020025649/fantendo/images/b/b2/Sad_Face.png",
			Width:  100,
			Height: 100,
		}
		emMsg := &discordgo.MessageEmbed{
			Title:       searchItem,
			Description: "Entry Not Found in Database",
			Thumbnail:   emThumb,
		}
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
	} else {
		nHemi, sHemi := parseHemi(entry.NorthSt, entry.NorthEnd, entry.SouthSt, entry.SouthEnd)
		emThumb := &discordgo.MessageEmbedThumbnail{
			URL:    entry.Image,
			Width:  200,
			Height: 200,
		}
		emPrice := &discordgo.MessageEmbedField{
			Name:  "Price",
			Value: strconv.Itoa(entry.SellPrice) + " Bells",
		}
		emHemiNorth := &discordgo.MessageEmbedField{
			Name:  "Northern Hemisphere Months",
			Value: nHemi,
		}
		emHemiSouth := &discordgo.MessageEmbedField{
			Name:  "Southern Hemisphere Months",
			Value: sHemi,
		}
		emTime := &discordgo.MessageEmbedField{
			Name:  "Time",
			Value: removeUnderscore(entry.Time),
		}
		emLocation := &discordgo.MessageEmbedField{
			Name:  "Location",
			Value: removeUnderscore(entry.Location),
		}
		emMsg := &discordgo.MessageEmbed{
			Title:       searchItem,
			Description: strings.Title(entry.Type),
			Thumbnail:   emThumb,
			Fields:      []*discordgo.MessageEmbedField{emLocation, emPrice, emTime, emHemiNorth, emHemiSouth},
		}
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
	}
}

// utility func to parse hemisphere data and return as
// useful text users can read
func parseHemi(ns, ne, st, se string) (string, string) {
	northSt := strings.Split(ns, "|")
	northEnd := strings.Split(ne, "|")
	southSt := strings.Split(st, "|")
	southEnd := strings.Split(se, "|")
	if len(northSt) == 1 {
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
// E.g. May to June AND September to November
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
