package commands

import (
	"flag"
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
	commands "github.com/ikkerens/gophbot/handlers"
	"go.uber.org/zap"
)

func init() {
	commands.AddCommand("purge", purge)
	commands.AddCommand("prune", purge)
	commands.AddCommand("clean", purge)
}

func purge(discord *discordgo.Session, cmd *commands.InvokedCommand) {
	if !cmd.HasPermission(discordgo.PermissionManageMessages) {
		cmd.Reply("I'm sorry, you need the `MANAGE_MESSAGES` permission to use this feature.")
		return
	}

	// This command has several flags that modify behaviour
	flags := flag.NewFlagSet("purgefs", flag.ContinueOnError)
	bots := flags.Bool("bots", false, "Setting this will cause this command to only delete bot messages")
	images := flags.Bool("images", false, "Setting this will delete only messages containing an image")
	author := flags.String("author", "", "Setting this will only delete messages by the specified author")
	channelID := flags.String("channel", "", "Setting this will cause the command to be executed in a specific channel instead of the current")

	if err := flags.Parse(cmd.Args); err != nil {
		cmd.Reply("I'm sorry, I did not understand what you are trying to do there.")
		gophbot.Log.Error("Could not parse purge flags", zap.Error(err))
		return
	}

	authorID := strings.Trim(*author, "<@!>")

	// Get the amount of to-delete messages
	countStr := flags.Arg(0)
	if countStr == "" {
		cmd.Reply("You have to specify an amount of messages that you wish to purge.")
		return
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		cmd.Reply(countStr + " is not a valid number.")
		return
	}
	if count > 1000 {
		cmd.Reply("You can not delete more than 1000 messages in 1 operation.")
		return
	}

	// Get the active channel instance, verifying it exists in the process
	var channel *discordgo.Channel
	if *channelID != "" {
		channel, err = gophbot.GetChannel(discord, strings.Trim(*channelID, "<#>"))
		if err != nil {
			cmd.Reply("I could not find  that channel!")
			gophbot.Log.Error("Could not locate specified channel", zap.String("channel", *channelID), zap.Error(err))
			return
		}

		if channel.GuildID != cmd.Guild.ID {
			cmd.Reply("That channel is not in this guild!")
			return
		}
	}

	// Build a list of all messages
	before := cmd.Message.ID
	messageIDs := make([]gophbot.Snowflake, 0, count)
	for count > 0 {
		messages, err := discord.ChannelMessages(channel.ID, int(math.Min(100, float64(count))), before, "", "")
		if err != nil {
			gophbot.Log.Error("Could not request discord channel messages", zap.Error(err))
			return
		}

		for _, message := range messages {
			if *bots && !message.Author.Bot {
				continue
			}

			if *images && (len(message.Attachments) == 0 || message.Attachments[0].Height == 0) {
				continue
			}

			if authorID == "" || authorID == message.Author.ID {
				continue
			}

			messageIDs = append(messageIDs, message.ID)
			count--
			before = messageIDs[len(messageIDs)-1]
		}
	}

	for _, id := range messageIDs {
		gophbot.Log.Info("Deleting", zap.String("msg", id))
	}
	cmd.Reply("Puuuuurge!")
	// TODO delete messages
}
