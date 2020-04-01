package cmd

import (
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
		emThumb := &discordgo.MessageEmbedThumbnail{
			URL:    entry.Image,
			Width:  200,
			Height: 200,
		}
		emMsg := &discordgo.MessageEmbed{
			Title:       searchItem,
			Description: entry.Location,
			Thumbnail:   emThumb,
		}
		cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, emMsg)
	}
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
