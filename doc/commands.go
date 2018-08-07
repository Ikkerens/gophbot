package doc

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot/wizard"
)

type CommandList struct {
}

func (*CommandList) Embed(*wizard.Session) *discordgo.MessageEmbed {
	panic("implement me")
}

func (*CommandList) Buttons(*wizard.Session) []string {
	panic("implement me")
}

func (*CommandList) FireButton(*wizard.Session, string) {
	panic("implement me")
}

type CommandDetail struct {
	Command     string
	Description string
	Aliases     []string
	Usage       string
	Flags       map[string]string
}

func (c *CommandDetail) Embed(s *wizard.Session) *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, 0, len(c.Flags)+1)
	if len(c.Aliases) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Aliases",
			Value: strings.Join(c.Aliases, ", "),
		})
	}

	if c.Flags != nil {
		for f, d := range c.Flags {
			fields = append(fields, &discordgo.MessageEmbedField{
				Inline: true,
				Name:   f,
				Value:  d,
			})
		}
	}

	return makeDefaultHelpPage(&discordgo.MessageEmbed{
		Title:       "Command " + c.Command,
		Description: c.Description,
		Fields:      fields,
	})
}

func (*CommandDetail) Buttons(*wizard.Session) []string {
	return noButtons
}

func (*CommandDetail) FireButton(*wizard.Session, string) {}
