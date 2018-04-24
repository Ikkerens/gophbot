package handlers

import (
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	"go.uber.org/zap"
)

var commands = make(map[string]CommandHandler)

// CommandHandler describes what a command handler function should look like
type CommandHandler = func(discord *discordgo.Session, cmd *InvokedCommand)

func init() {
	gophbot.AddHandler(handleCommand)
}

// AddCommand registers a new command for the command handler
func AddCommand(command string, handler CommandHandler) {
	commands[strings.ToLower(command)] = handler
}

func handleCommand(discord *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == gophbot.Self.ID {
		return
	}

	channel, err := gophbot.GetChannel(discord, event.ChannelID)
	if err != nil {
		gophbot.Log.Error("Could not fetch channel a command is executed in", zap.Error(err))
		return
	}

	if channel.Type != discordgo.ChannelTypeGuildText {
		discord.ChannelMessageSend(event.ChannelID, "I'm sorry, but I only listen commands inside servers.")
		return
	}

	guild, err := gophbot.GetGuild(channel.GuildID)
	if err != nil {
		gophbot.Log.Error("Could not fetch guild a command is executed in", zap.Error(err))
		return
	}

	member, err := gophbot.GetGuildMember(guild.ID, event.Author.ID)
	if err != nil {
		gophbot.Log.Error("Could not get membership definition for the user that invoked a command", zap.Error(err))
		return
	}

	dbGuild := &gophbot.Guild{ID: guild.ID}
	if err = gophbot.DB.Where(dbGuild).Find(dbGuild).Error; err != nil {
		gophbot.Log.Error("Could not get guild information from database", zap.Error(err))
		return
	}

	if !strings.HasPrefix(event.Content, dbGuild.Prefix) {
		return
	}

	args := strings.Split(event.Content, " ")
	if utf8.RuneCountInString(args[0]) <= len(dbGuild.Prefix) {
		return
	}

	commandStr := strings.ToLower(args[0][len(dbGuild.Prefix):])
	command, exists := commands[commandStr]
	if exists {
		command(discord, &InvokedCommand{
			Session: discord,
			Guild:   guild,
			Channel: channel,
			User:    event.Author,
			Member:  member,

			DBGuild: dbGuild,

			Args: args[1:],
		})
	}
}

// InvokedCommand is a convenience struct that holds all the collected information before invoking a command.
type InvokedCommand struct {
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Member  *discordgo.Member
	User    *discordgo.User

	DBGuild *gophbot.Guild

	Args []string
}

func (c *InvokedCommand) HasPermission(permission int) bool {
	base, err := gophbot.ComputeBasePermissions(c.Member, c.Guild)
	if err != nil {
		gophbot.Log.Error("Could not determine base permissions", zap.Error(err))
		return false
	}

	p, err := gophbot.ComputeOverwrites(base, c.Member, c.Channel)
	if err != nil {
		gophbot.Log.Error("Could not compute overwrites for permissions", zap.Error(err))
		return false
	}

	return (p & permission) == permission
}

func (c *InvokedCommand) Reply(reply string) (msg *discordgo.Message, err error) {
	msg, err = c.Session.ChannelMessageSend(c.Channel.ID, c.User.Mention()+" "+reply)
	return
}
