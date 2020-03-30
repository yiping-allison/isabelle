package daisymaebot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/daisymae/cmd"
)

// Bot represents a daisymae bot instance
type Bot struct {
	DS       *discordgo.Session
	Commands map[string]Command
}

// New creates a new daisymae bot instance and loads bot commands.
//
// It will return the finished bot and nil upon success or
// empty bot and err upon failure
func New(bc string) (*Bot, error) {
	discord, err := discordgo.New("Bot " + bc)
	if err != nil {
		return &Bot{}, errors.New("daisymaebot: error connecting to discord")
	}
	// Commands Setup
	cmds := make(map[string]Command, 0)
	daisy := &Bot{
		DS:       discord,
		Commands: cmds,
	}
	daisy.compileCommands()
	// Add Handlers
	daisy.DS.AddHandler(daisy.messageCreate)
	return daisy, nil
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "?") {
		b.processCmd(s, m)
	}
}

// Command represents a discord bot command
type Command struct {
	Cmd  func(cmd.CommandInfo)
	Help string
}

// processCmd attemps to process any string that is prefixed with bot notifier
//
// Valid commands will be run while invalid commands will be ignored
func (b *Bot) processCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Give option to customize prefix
	// TODO: Be able to parse multi-word commands (E.g. ?ping help brings up command description of ping)
	if val, ok := b.Commands[m.Content[1:]]; ok {
		ci := cmd.CommandInfo{
			Ses: s,
			Msg: m,
		}
		val.Cmd(ci)
	}
}

// utility func for all commands bot should add to command map
func (b *Bot) compileCommands() {
	b.addCommand("ping", "tells the bot to ping", cmd.Ping)
	b.addCommand("pong", "tells bot to respond with ping", cmd.Pong)
}

// utility func to add command to command map
func (b *Bot) addCommand(name, help string, cmd func(cmd.CommandInfo)) {
	if _, ok := b.Commands[name]; ok {
		fmt.Printf("addCommand err: %s already exists in the map\n", name)
		return
	}
	command := Command{
		Cmd:  cmd,
		Help: help,
	}
	b.Commands[name] = command
}
