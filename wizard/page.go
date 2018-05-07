package wizard

import "github.com/bwmarrin/discordgo"

// CloseEmoji is the ❌ emoji in Discord, used for the close-wizard button
const CloseEmoji = "\U0000274c"

// Page defines what functionality a wizard page should have.
type Page interface {
	Embed(*Session) *discordgo.MessageEmbed
	Buttons(*Session) []string
	FireButton(session *Session, emoji string)
}

// SwitchToPage causes this Session to switch to another page, and replace the buttons.
func (s *Session) SwitchToPage(page Page) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.switchToPage(page)
}

func (s *Session) switchToPage(page Page) error {
	var err error
	if s.container == nil {
		s.container, err = s.discord.ChannelMessageSendEmbed(s.channel, page.Embed(s))
	} else {
		_, err = s.discord.ChannelMessageEditEmbed(s.channel, s.container.ID, page.Embed(s))
		if err != nil {
			return err
		}

		err = s.discord.MessageReactionsRemoveAll(s.channel, s.container.ID)
	}

	if err != nil {
		return err
	}

	go func() {
		for _, reaction := range page.Buttons(s) {
			s.discord.MessageReactionAdd(s.channel, s.container.ID, reaction)
		}
		s.discord.MessageReactionAdd(s.channel, s.container.ID, CloseEmoji)
	}()

	return nil
}
