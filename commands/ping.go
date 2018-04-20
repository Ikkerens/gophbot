package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	"go.uber.org/zap"
)

func init() {
	gophbot.AddCommand("ping", ping)
}

func ping(session *discordgo.Session, event *discordgo.MessageCreate, _ []string) {
	var (
		lock     sync.RWMutex
		message  *discordgo.Message
		callback = make(chan struct{})
	)

	// Set up an event handler now, in case we get the WS event before we have the HTTP response ready
	defer session.AddHandler(func(_ *discordgo.Session, event *discordgo.MessageCreate) {
		lock.RLock()
		defer lock.RUnlock()

		if message != nil && message.ID == event.ID {
			close(callback)
		}
	})()

	// Prepare the initial message
	pong := event.Author.Mention() + " Pong!"

	// Lock the message variable for writing & send the message
	lock.Lock()
	start := time.Now()
	msg, err := session.ChannelMessageSend(event.ChannelID, pong)
	if err != nil {
		lock.Unlock()
		gophbot.Log.Error("Could not send ping message.", zap.Error(err))
		return
	}

	sent := time.Now()
	message = msg
	lock.Unlock()

	pong += fmt.Sprintf(" (REST: %dms)", sent.Sub(start).Nanoseconds()/1e6)

	select {
	case <-callback:
		pong += fmt.Sprintf(" (WS: %dms)", time.Since(sent).Nanoseconds()/1e6)
	case <-time.After(5 * time.Second):
		pong += " (WS: timed out, >5s lag)"
	}

	_, err = session.ChannelMessageEdit(msg.ChannelID, msg.ID, pong)
	if err != nil {
		gophbot.Log.Error("Could not edit ping message.", zap.Error(err))
	}
}
