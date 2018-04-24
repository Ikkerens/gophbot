package commands

import (
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	commands "github.com/ikkerens/gophbot/handlers"
	"go.uber.org/zap"
)

func init() {
	commands.AddCommand("prefix", setPrefix)
	commands.AddCommand("setprefix", setPrefix)
}

func setPrefix(_ *discordgo.Session, cmd *commands.InvokedCommand) {
	if !cmd.HasPermission(discordgo.PermissionManageServer) {
		cmd.Reply("I'm sorry, you need the `MANAGE_SERVER` permission to use this feature.")
		return
	}

	if len(cmd.Args) != 1 {
		cmd.Reply("This command requires one prefix, without spaces.")
		return
	}

	if utf8.RuneCountInString(cmd.Args[0]) > 10 {
		cmd.Reply("A prefix can at most be 5 characters long.")
		return
	}

	cmd.DBGuild.Prefix = cmd.Args[0]
	if err := gophbot.DB.Save(cmd.DBGuild).Error; err != nil {
		gophbot.Log.Error("Could not update guild prefix", zap.Error(err))
		return
	}

	cmd.Reply("Success! I will now use the prefix `" + cmd.Args[0] + "` for all commands.")
}
