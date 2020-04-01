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
// CmdHlp: contains the command's help string
type CommandInfo struct {
	Ses     *discordgo.Session
	Msg     *discordgo.MessageCreate
	Service models.Services
	CmdName string
	CmdOps  []string
	CmdHlp  string
}
