package commands

import (
	"github.com/bwmarrin/discordgo"
	commands "github.com/ikkerens/gophbot/handlers"
)

func init() {
	commands.AddCommand("purge", purge)
	commands.AddCommand("clean", purge)
}

func purge(discord *discordgo.Session, cmd *commands.InvokedCommand) {
	if !cmd.HasPermission(discordgo.PermissionManageMessages) {
		cmd.Reply("I'm sorry, you need the `MANAGE_MESSAGES` permission to use this feature.")
		return
	}

	discord.ChannelMessageSend(cmd.Channel.ID, "Purrrrrrrge.")
}
