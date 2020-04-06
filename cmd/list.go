package cmd

// List will list all bot commands to the user
func List(cmdInfo CommandInfo) {
	fields := format(
		createFields("search", cmdInfo.Prefix+"search [item]", true),
		createFields("help", cmdInfo.Prefix+"help [command_name]", false),
		createFields("event", cmdInfo.Prefix+"event [arguments]", true),
		createFields("queue", cmdInfo.Prefix+"queue [event ID]", true),
		createFields("close", cmdInfo.Prefix+"close event [event ID]", true),
		createFields("ping", cmdInfo.Prefix+"ping", true),
		createFields("pong", cmdInfo.Prefix+"pong", true),
	)
	msg := cmdInfo.createMsgEmbed("Commands", listThumbURL, "", listColor, fields)
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
}
