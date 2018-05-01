package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	"go.uber.org/zap"
)

func init() {
	gophbot.AddHandler(newGuild)
	gophbot.AddHandler(removeGuild)
}

func newGuild(_ *discordgo.Session, event *discordgo.GuildCreate) {
	g := &gophbot.Guild{ID: event.ID}

	if err := gophbot.DB.FirstOrCreate(g, g).Error; err != nil {
		gophbot.Log.Error("Could not create guild definition in database", zap.Error(err))
	}
}

func removeGuild(_ *discordgo.Session, event *discordgo.GuildDelete) {
	if err := gophbot.DB.Delete(gophbot.Guild{}, "id LIKE ?", event.ID).Error; err != nil {
		gophbot.Log.Error("Could not delete guild definition from database", zap.Error(err))
	}
}
