package cmd

// Pong test func to play around with multiple commands
func Pong(cmdInfo CommandInfo) {
	cmdInfo.Ses.ChannelMessageSend(cmdInfo.Msg.ChannelID, "ping!")
}
