package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/daisymae/models"
)

const (
	listThumbURL  string = "https://cdn2.iconfinder.com/data/icons/business-office-icons/256/To-do_List-512.png"
	helpThumbURL  string = "https://www.bbqguru.com/content/images/manual-bbq-icon.png"
	errThumbURL   string = "http://static2.wikia.nocookie.net/__cb20131020025649/fantendo/images/b/b2/Sad_Face.png"
	blobSThumbURL string = "https://gerhardinger.org/wp-content/uploads/2017/05/icon-world.png"
	checkThumbURL string = "http://www.providesupport.com/blog/wp-content/uploads/2013/08/green-check-mark.png"
)

const (
	listColor    int = 13473141
	helpColor    int = 9410425
	errColor     int = 14886454
	searchColor  int = 9526403
	eventColor   int = 3108709
	successColor int = 3764015
)

// CommandInfo represents all metadata discord and bot needs to
// execute certain API callbacks and commands
//
// Ses: discord session (discord)
//
// Msg: discord message (discord)
//
// Service: contains all services needed by bot (bot)
//
// Prefix: the prefix the bot recognizes set in .config
//
// CmdName: contains the command name
//
// CmdOps: the full slice of commands (unparsed)
//
// CmdList: contains the names of all commands
type CommandInfo struct {
	Ses     *discordgo.Session
	Msg     *discordgo.MessageCreate
	Service models.Services
	Prefix  string
	CmdName string
	CmdOps  []string
	CmdList []string
}

// format is a utility func which takes in a variadic parameter of discord message embed field
// types and returns them as a slice
func format(f ...*discordgo.MessageEmbedField) []*discordgo.MessageEmbedField { return f }

// createFields is a utility func which creates individual discord Message Embed Field types
func createFields(text, val string, inline bool) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   text,
		Value:  val,
		Inline: inline,
	}
}

// createMsgEmbed is a utility function to be used by all command types to print messages using
// discord message embed
func (c CommandInfo) createMsgEmbed(title, tURL, desc string, color int, fields []*discordgo.MessageEmbedField) {
	emThumb := &discordgo.MessageEmbedThumbnail{
		URL:    tURL,
		Width:  200,
		Height: 200,
	}
	emMsg := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Thumbnail:   emThumb,
		Fields:      fields,
		Color:       color,
	}
	c.Ses.ChannelMessageSendEmbed(c.Msg.ChannelID, emMsg)
}
