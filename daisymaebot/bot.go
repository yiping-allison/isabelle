package daisymaebot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/daisymae/cmd"
	"github.com/yiping-allison/daisymae/models"
)

// Bot represents a daisymae bot instance
type Bot struct {
	DS       *discordgo.Session
	Service  models.Services
	Prefix   string
	Commands map[string]Command
}

// New creates a new daisymae bot instance and loads bot commands.
//
// It will return the finished bot and nil upon success or
// empty bot and err upon failure
func New(bc string) (*Bot, error) {
	if bc == "" {
		return nil, errors.New("daisymaebot: you need to input a botKey in the .config file")
	}
	discord, err := discordgo.New("Bot " + bc)
	if err != nil {
		return nil, errors.New("daisymaebot: error connecting to discord")
	}
	// Commands Setup
	cmds := make(map[string]Command, 0)
	daisy := &Bot{
		Prefix:   "?",
		DS:       discord,
		Service:  models.Services{},
		Commands: cmds,
	}
	daisy.compileCommands()
	// Add Handlers
	daisy.DS.AddHandler(daisy.ready)
	daisy.DS.AddHandler(daisy.handleMessage)
	return daisy, nil
}

// ready will update bot status after bot receives "ready" event from
// discord
func (b *Bot) ready(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(0, b.Prefix+"list")
}

// handleMessage handles all new discord messages which the bot uses to determine
// actions
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, b.Prefix) {
		b.processCmd(s, m)
	}
}

// Command represents a discord bot command
type Command struct {
	Cmd func(cmd.CommandInfo)
}

// processCmd attemps to process any string that is prefixed with bot notifier
//
// Valid commands will be run while invalid commands will be ignored
//
// Example bot commands:
//
// ?search
//
// ?search help
func (b *Bot) processCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmds := regexp.MustCompile("\\s+").Split(m.Content[len(b.Prefix):], -1)
	trim := strings.TrimPrefix(cmds[0], b.Prefix)
	res := b.find(trim)
	if res == nil {
		return
	}
	var commands []string
	for key := range b.Commands {
		commands = append(commands, key)
	}
	ci := cmd.CommandInfo{
		Ses:     s,
		Msg:     m,
		Service: b.Service,
		Prefix:  b.Prefix,
		CmdName: trim,
		CmdOps:  cmds,
		CmdList: commands,
	}
	// Run command
	res.Cmd(ci)
}

// finds a command in the command map
//
// If it exists, it returns the Command
// If not, it returns nil
func (b *Bot) find(name string) *Command {
	if val, ok := b.Commands[name]; ok {
		return &val
	}
	return nil
}

// compileCommands contains all commands the bot should add to the bot command map
func (b *Bot) compileCommands() {
	b.addCommand("search", cmd.Search)
	b.addCommand("help", cmd.Help)
	b.addCommand("list", cmd.List)
	b.addCommand("event", cmd.Event)
	b.addCommand("queue", cmd.Queue)
	b.addCommand("close", cmd.Close)
	b.addCommand("ping", cmd.Ping)
	b.addCommand("pong", cmd.Pong)
}

// utility func to add command to bot command map
func (b *Bot) addCommand(name string, cmd func(cmd.CommandInfo)) {
	if _, ok := b.Commands[name]; ok {
		fmt.Printf("addCommand: %s already exists in the map\n", name)
		return
	}
	command := Command{
		Cmd: cmd,
	}
	b.Commands[name] = command
}

// SetPrefix sets user directed bot prefix from .config
func (b *Bot) SetPrefix(newPrefix string) {
	b.Prefix = newPrefix
}
