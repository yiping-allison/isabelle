package cmd

import "github.com/bwmarrin/discordgo"

// CommandInfo represents all metadata discord needs to
// execute certain API callbacks
type CommandInfo struct {
	Ses *discordgo.Session
	Msg *discordgo.MessageCreate
}
