package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/daisymae/service"
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
// Cmd: contains full valid command input by user in discord (bot)
type CommandInfo struct {
	Ses     *discordgo.Session
	Msg     *discordgo.MessageCreate
	Service service.Services
	CmdName string
	CmdOps  []string
	CmdHlp  string
}
