package doc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	"github.com/ikkerens/gophbot/wizard"
)

var noButtons = make([]string, 0)

func makeDefaultHelpPage(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Color = 0x96D6FF
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: gophbot.Self.AvatarURL("48")}
	embed.Footer = &discordgo.MessageEmbedFooter{Text: "Please use the reactions below to make a selection. Press " + wizard.CloseEmoji + " to finish."}
	return embed
}
