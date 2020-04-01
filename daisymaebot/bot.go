package daisymaebot

import (
	"errors"
	"fmt"
	"reflect"
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
		return &Bot{}, errors.New("daisymaebot: you need to input a botKey in the .config file")
	}
	discord, err := discordgo.New("Bot " + bc)
	if err != nil {
		return &Bot{}, errors.New("daisymaebot: error connecting to discord")
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
	daisy.DS.AddHandler(daisy.messageCreate)
	return daisy, nil
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	Cmd  func(cmd.CommandInfo)
	Help string
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
	// TODO: Be able to parse search commands: ?search monarch butterfly brings up
	//       all info of insect in New Horizons game
	cmds := strings.Split(m.Content[len(b.Prefix):], " ")
	trim := strings.TrimPrefix(cmds[0], b.Prefix)
	res := b.find(trim)
	if reflect.DeepEqual(res, Command{}) {
		return
	}
	ci := cmd.CommandInfo{
		Ses:     s,
		Msg:     m,
		Service: b.Service,
		CmdName: trim,
		CmdOps:  cmds,
		CmdHlp:  res.Help,
	}
	if trim == "help" && len(cmds) > 1 {
		// if command given is a help command
		ci.CmdName = strings.Join(cmds[1:], " ")
		tmp := b.find(ci.CmdName)
		if reflect.DeepEqual(tmp, Command{}) {
			ci.CmdHlp = "This command does not exist."
		} else {
			ci.CmdHlp = b.find(ci.CmdName).Help
		}
	}
	// Run command
	res.Cmd(ci)
}

// finds a command in the command map
//
// If it exists, it returns the Command
// If not, it returns an empty Command
func (b *Bot) find(name string) Command {
	if val, ok := b.Commands[name]; ok {
		return val
	}
	return Command{}
}

// compileCommands contains all commands the bot should add to the bot command map
func (b *Bot) compileCommands() {
	b.addCommand("search", `Search for an insect or fish in the database
		Example:
		?search Tarantula`, cmd.Search)
	b.addCommand("help", `Search how to use a command
		Example:
		?help search`, cmd.Help)
	b.addCommand("ping", "Tells the bot to ping", cmd.Ping)
	b.addCommand("pong", "Tells bot to respond with ping", cmd.Pong)
}

// utility func to add command to bot command map
func (b *Bot) addCommand(name, help string, cmd func(cmd.CommandInfo)) {
	if _, ok := b.Commands[name]; ok {
		fmt.Printf("addCommand: %s already exists in the map\n", name)
		return
	}
	command := Command{
		Cmd:  cmd,
		Help: help,
	}
	b.Commands[name] = command
}

// SetPrefix sets user directed bot prefix from .config
func (b *Bot) SetPrefix(newPrefix string) {
	b.Prefix = newPrefix
}
