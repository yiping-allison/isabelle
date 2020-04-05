package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/daisymae/models"
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
