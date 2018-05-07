package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	"github.com/ikkerens/gophbot/doc"
	commands "github.com/ikkerens/gophbot/handlers"
	"github.com/ikkerens/gophbot/wizard"
	"go.uber.org/zap"
)

func init() {
	commands.AddCommand("help", help)
}

func help(discord *discordgo.Session, cmd *commands.InvokedCommand) {
	pageName := "index"
	if len(cmd.Args) > 0 {
		pageName = strings.TrimSpace(strings.Join(cmd.Args, " "))
	}

	page, ok := topics[pageName]
	if !ok {
		cmd.Reply("The help topic `" + pageName + "` does not exist!")
		return
	}

	_, err := wizard.New(discord, cmd.Channel.ID, cmd.User.ID, page())
	if err != nil {
		gophbot.Log.Error("Could not start help wizard", zap.Error(err))
	}
}

var topics = map[string]func() wizard.Page{
	"index": func() wizard.Page { return new(doc.IndexPage) },
}
