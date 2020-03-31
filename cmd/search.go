package cmd

import "strings"

// Search will look up a possible insect or fish in the database and display to the user
func Search(cmdInfo CommandInfo) {
	// TODO: Finish this
	// NOTE: Ignore cmdInfo.Cmd[0] - start parsing and database search from [1:]
	cmdInfo.Ses.ChannelMessageSend(cmdInfo.Msg.ChannelID, "Testing search!")
	cmdInfo.Ses.ChannelMessageSend(cmdInfo.Msg.ChannelID, printStr(cmdInfo.CmdOps))
}

func printStr(str []string) string {
	final := strings.Join(str[1:], " ")
	return final
}
