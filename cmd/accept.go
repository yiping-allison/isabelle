package cmd

import "fmt"

// Accept will allow moderators (or bot controllers) to accept
// reputation application requests
//
// The command usage should look like: ;accept 1234
func Accept(cmdInfo CommandInfo) {
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
	if !cmdInfo.Service.Rep.Exists(userID) {
		// if someone has been repped, they need to have an event or trade already...
		// this case shouldn't be possible so we automatically reject
		return
	}
	err := cmdInfo.Service.Rep.Increase(userID)
	if err != nil {
		fmt.Println("error updating individual...")
		return
	}
	// clean temp map
	cmdInfo.Service.Rep.Clean(repID)
	// print success msg
	embed := cmdInfo.createMsgEmbed(
		"Accepted Reputation Request", checkThumbURL, "App ID: "+repID,
		successColor, format(
			createFields("Congratulations, your rep has increased!", mentionUser(userID), true),
		))
	cmdInfo.Ses.ChannelMessageSendEmbed(cmdInfo.Msg.ChannelID, embed)
}

// isAdmin is a helper func which returns true if role container
// has an admin role
func isAdmin(roles []string, adminID string) bool {
	for _, r := range roles {
		if r == adminID {
			return true
		}
	}
	return false
}
