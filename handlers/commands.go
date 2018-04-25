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

	cmd, err := newInvokedCommand(discord, event)
	if err != nil {
		gophbot.Log.Error("Could not fetch command metadata", zap.Error(err))
		return
	}

	if cmd.Channel.Type != discordgo.ChannelTypeGuildText {
		discord.ChannelMessageSend(event.ChannelID, "I'm sorry, but I only listen commands inside servers.")
		return
	}

	if !strings.HasPrefix(event.Content, cmd.DBGuild.Prefix) {
		return
	}

	args := strings.Split(event.Content, " ")
	if utf8.RuneCountInString(args[0]) <= len(cmd.DBGuild.Prefix) {
		return
	}

	commandStr := strings.ToLower(args[0][len(cmd.DBGuild.Prefix):])
	cmd.Args = args[1:]
	command, exists := commands[commandStr]
	if exists {
		command(discord, cmd)
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

func newInvokedCommand(discord *discordgo.Session, event *discordgo.MessageCreate) (*InvokedCommand, error) {
	channel, err := gophbot.GetChannel(discord, event.ChannelID)
	if err != nil {
		return nil, err
	}

	guild, err := gophbot.GetGuild(channel.GuildID)
	if err != nil {
		return nil, err
	}

	member, err := gophbot.GetGuildMember(guild.ID, event.Author.ID)
	if err != nil {
		return nil, err
	}

	dbGuild := &gophbot.Guild{ID: guild.ID}
	if err = gophbot.DB.Where(dbGuild).Find(dbGuild).Error; err != nil {
		return nil, err
	}

	return &InvokedCommand{
		Session: discord,
		Guild:   guild,
		Channel: channel,
		User:    event.Author,
		Member:  member,

		DBGuild: dbGuild,
	}, nil
}

// HasPermission is a convenience command that allows for quick & simple permission checking.
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

// Reply is a convenience method that allows for quick replies to a command.
func (c *InvokedCommand) Reply(reply string) (msg *discordgo.Message, err error) {
	msg, err = c.Session.ChannelMessageSend(c.Channel.ID, c.User.Mention()+" "+reply)
	return
}
