package wizard

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
)

const timeout = 30 * time.Second

// Session is a wizard session, using embeds to display a menu in Discord and reactions for buttons.
type Session struct {
	discord   *discordgo.Session
	page      Page
	container *discordgo.Message
	owner     gophbot.Snowflake
	channel   gophbot.Snowflake

	/* Runtime */
	mutex                sync.Mutex
	closeReactionHandler func()
	reactions            chan *discordgo.MessageReactionAdd
}

// New creates a new Wizard session.
func New(discord *discordgo.Session, channel, user gophbot.Snowflake, page Page) (*Session, error) {
	var (
		err error

		session = &Session{
			discord: discord,
			owner:   user,
			channel: channel,

			page:      page,
			reactions: make(chan *discordgo.MessageReactionAdd),
		}
	)

	session.mutex.Lock()
	session.closeReactionHandler = discord.AddHandler(session.handleReaction)

	err = session.switchToPage(page)
	if err != nil {
		return nil, err
	}
	session.mutex.Unlock()

	go session.loop()
	return session, nil
}

func (s *Session) handleReaction(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if event.UserID == gophbot.Self.ID {
		return
	}

	s.mutex.Lock()
	if s.container == nil || event.MessageID != s.container.ID {
		s.mutex.Unlock()
		return
	}
	s.mutex.Unlock()

	discord.MessageReactionRemove(event.ChannelID, event.MessageID, event.Emoji.APIName(), event.UserID)

	if event.UserID != s.owner {
		return
	}

	s.reactions <- event
}

func (s *Session) loop() {
	defer s.closeReactions()
	defer s.cleanupContainer()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event := <-s.reactions:
			if !timer.Reset(timeout) {
				<-timer.C
			}

			switch event.Emoji.APIName() {
			case CloseEmoji:
				return
			default:
				s.page.FireButton(s, event.Emoji.APIName())
			}
		case <-timer.C:
			return
		}
	}
}

func (s *Session) closeReactions() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.closeReactionHandler() // Prevent future calls
	s.container = nil        // Prevent pending calls

	close(s.reactions)
	for range s.reactions {
		// Discarding all reactions here, as the session is already over
	}
}

func (s *Session) cleanupContainer() {
	embed := s.container.Embeds[0]
	embed.Color = 0
	embed.Footer = &discordgo.MessageEmbedFooter{Text: "This dialog has either been ended or has expired. " +
		"Please use the same command again to start a new session."}

	s.discord.ChannelMessageEditEmbed(s.container.ChannelID, s.container.ID, embed)
	s.discord.MessageReactionsRemoveAll(s.container.ChannelID, s.container.ID)
}
