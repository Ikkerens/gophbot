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
	gophbot.Log.Info("Joining server.", zap.String("name", event.Name))
	g := &gophbot.Guild{ID: event.ID}

	if err := gophbot.DB.FirstOrCreate(g, g).Error; err != nil {
		gophbot.Log.Error("Could not create guild definition in database", zap.Error(err))
	}

	gophbot.State.Lock()
	gophbot.State.Guilds[event.ID] = g
	gophbot.State.Unlock()
}

func removeGuild(_ *discordgo.Session, event *discordgo.GuildDelete) {
	gophbot.Log.Info("Leaving server.", zap.String("name", event.Name))

	if err := gophbot.DB.Delete(gophbot.Guild{}, "id LIKE ?", event.ID).Error; err != nil {
		gophbot.Log.Error("Could not delete guild definition from database", zap.Error(err))
	}

	gophbot.State.Lock()
	delete(gophbot.State.Guilds, event.ID)
	gophbot.State.Unlock()
}
