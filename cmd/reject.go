package cmd

// Reject allows admins to reject reputation requests
func Reject(cmdInfo CommandInfo) {
	if len(cmdInfo.CmdOps) != 2 || !isAdmin(cmdInfo.Msg.Member.Roles, cmdInfo.AdminRole) {
		// command length must be 2 and must be admin
		return
	}
	repID := cmdInfo.CmdOps[1]
	if !cmdInfo.Service.Rep.RepIDExists(repID) {
		// repID doesn't exist
		return
	}
	userID := cmdInfo.Service.Rep.GetUser(repID)
	cmdInfo.Service.Rep.Clean(repID)
	// print rejection msg
	embed := cmdInfo.createMsgEmbed(
		"Rejected Reputation Request", errThumbURL, "App ID: "+repID,
		errColor, format(
			createFields("Sorry, we think there's something wrong.", mentionUser(userID), true),
			createFields("Suggestion", "If this is a mistake, please PM the mods, thanks!", true),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, embed)
}
