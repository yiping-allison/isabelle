package cmd

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yiping-allison/isabelle/models"
)

const (
	listThumbURL  string = "https://cdn.discordapp.com/attachments/693564368423616562/696628637163847730/Commands2.png"
	helpThumbURL  string = "https://cdn.discordapp.com/attachments/693564368423616562/696624102101876746/Search.png"
	errThumbURL   string = "http://static2.wikia.nocookie.net/__cb20131020025649/fantendo/images/b/b2/Sad_Face.png"
	blobSThumbURL string = "https://gerhardinger.org/wp-content/uploads/2017/05/icon-world.png"
	checkThumbURL string = "https://cdn.discordapp.com/attachments/693564368423616562/696816758945742868/Success.png"
	queueThumbURL string = "https://cdn.discordapp.com/attachments/693564368423616562/696802401490960505/Queue.png"
	tradeThumbURL string = "https://cdn.discordapp.com/attachments/693564368423616562/698023785164439562/trade-icon.png"
	thumbThumbURL string = "http://cdn.onlinewebfonts.com/svg/img_504758.png"
)

const (
	listColor    int = 4772300
	helpColor    int = 9410425
	errColor     int = 14886454
	searchColor  int = 9526403
	eventColor   int = 3108709
	successColor int = 3764015
	tradeColor   int = 15893760
	appColor     int = 4617611
)

// CommandInfo represents all metadata discord and bot needs to
// execute certain API callbacks and commands
type CommandInfo struct {
	// AdminRole: ID of the role which can control bot
	AdminRole string

	// Ses: discord session (discord)
	Ses *discordgo.Session

	// Msg: discord message (discord)
	Msg *discordgo.MessageCreate

	// Service: contains all services needed by bot (bot)
	Service models.Services

	// Channel ID to post event notices
	ListingID string

	// Channel ID to post general bot commands
	BotChID string

	// Channel ID of rep applications
	AppID string

	// Prefix: the prefix the bot recognizes set in .config
	Prefix string

	// CmdName: contains the command name
	CmdName string

	// CmdOps: the full slice of commands (unparsed)
	CmdOps []string

	// CmdList: contains the names of all commands
	CmdList []string
}

// newRep creates a new rep database objects and inserts it into
// postgreSQL
func (c CommandInfo) newRep(userID string) {
	// user does not exist in rep database
	// initialize user to 0
	rep := models.Rep{
		DiscordID: userID,
		RepNum:    0,
	}
	err := c.Service.Rep.Create(&rep)
	if err != nil {
		// error creating new user in database
		return
	}
}

// generateID will come up with a pseudo random number with 4 digits
// and return it in string format
//
// This is used to generate different event IDs
func generateID(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	id := min + rand.Intn(max-min)
	return strconv.Itoa(id)
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
func (c CommandInfo) createMsgEmbed(title, tURL, desc string, color int, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	emThumb := &discordgo.MessageEmbedThumbnail{
		URL:    tURL,
		Width:  400,
		Height: 400,
	}
	emMsg := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Thumbnail:   emThumb,
		Fields:      fields,
		Color:       color,
	}
	return emMsg
}
