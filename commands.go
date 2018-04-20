package gophbot

import (
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

var commands = make(map[string]CommandHandler)

// CommandHandler describes what a command handler function should look like
type CommandHandler = func(session *discordgo.Session, event *discordgo.MessageCreate, args []string)

func init() {
	AddHandler(handleCommand)
}

// AddCommand registers a new command for the command handler
func AddCommand(command string, handler CommandHandler) {
	commands[strings.ToLower(command)] = handler
}

func handleCommand(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == Self.ID {
		return
	}

	if !strings.HasPrefix(event.Content, "/") {
		return
	}

	args := strings.Split(event.Content, " ")
	if utf8.RuneCountInString(args[0]) == 1 {
		return
	}

	commandStr := strings.ToLower(args[0][1:])
	command, exists := commands[commandStr]
	if exists {
		command(session, event, args[1:])
	}
}
