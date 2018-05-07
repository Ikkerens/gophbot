package doc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot/wizard"
)

// IndexPage is the main entry point for the help command
type IndexPage struct{}

// Embed ...
func (IndexPage) Embed(*wizard.Session) *discordgo.MessageEmbed {
	return makeDefaultHelpPage(&discordgo.MessageEmbed{
		Title:       "GophBot help page",
		Description: "This is a very helpful help page.",
	})
}

// Buttons ...
func (IndexPage) Buttons(*wizard.Session) []string {
	return noButtons
}

// FireButton ...
func (IndexPage) FireButton(session *wizard.Session, emoji string) {}
