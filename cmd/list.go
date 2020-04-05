package cmd

// List will list all bot commands to the user
func List(cmdInfo CommandInfo) {
	fields := format(
		createFields("search", cmdInfo.Prefix+"search [item]", true),
		createFields("help", cmdInfo.Prefix+"help [command_name]", false),
		createFields("ping", cmdInfo.Prefix+"ping", true),
		createFields("pong", cmdInfo.Prefix+"pong", true),
	)
	cmdInfo.createMsgEmbed("Commands", listThumbURL, "", listColor, fields)
}
