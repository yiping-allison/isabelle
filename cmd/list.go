package cmd

// List will list all bot commands to the user
func List(cmdInfo CommandInfo) {
	fields := format(
		createFields("NOTE", "Use the help command for more detailed usage examples.", false),
		createFields("search", cmdInfo.Prefix+"search ...", true),
		createFields("help", cmdInfo.Prefix+"help ...", true),
		createFields("event", cmdInfo.Prefix+"event ...", true),
		createFields("queue", cmdInfo.Prefix+"queue ...", true),
		createFields("close", cmdInfo.Prefix+"close ...", true),
		createFields("unregister", cmdInfo.Prefix+"unregister ...", true),
		createFields("trade", cmdInfo.Prefix+"trade ...", true),
		createFields("offer", cmdInfo.Prefix+"offer ...", true),
		createFields("accept", cmdInfo.Prefix+"accept ...", true),
		createFields("reject", cmdInfo.Prefix+"reject ...", true),
		createFields("rep", cmdInfo.Prefix+"rep ...", true),
	)
	msg := cmdInfo.createMsgEmbed("Commands", listThumbURL, "", listColor, fields)
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, msg)
}
