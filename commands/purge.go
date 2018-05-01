package commands

import (
	"flag"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

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
	channel := cmd.Channel
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

	if count > 100 {
		if err := discord.MessageReactionAdd(cmd.Channel.ID, cmd.Message.ID, gophbot.LoadingReaction); err != nil {
			gophbot.Log.Error("Could not send progress reaction", zap.Error(err))
		}
	}

	success := false
	defer sendPurgeFailed(&success, cmd)

	// Build a list of all messages
	before := cmd.Message.ID
	oldMessageIDs := make([]gophbot.Snowflake, 0, count)
	newMessageIDs := make([]gophbot.Snowflake, 0, count)
	for count > 0 {
		messages, err := discord.ChannelMessages(channel.ID, 100, before, "", "")
		if err != nil {
			gophbot.Log.Error("Could not request discord channel messages", zap.Error(err))
			return
		}

		//before = messages[0].ID

		for _, message := range messages {
			before = message.ID

			if *bots && !message.Author.Bot {
				continue
			}

			if *images && (len(message.Attachments) == 0 || message.Attachments[0].Height == 0) {
				continue
			}

			if authorID != "" && authorID != message.Author.ID {
				continue
			}

			t, err := message.Timestamp.Parse()
			if err != nil {
				gophbot.Log.Error("Discord provided an unparseable date", zap.Error(err))
				return
			}

			if time.Since(t) > (2 * 7 * 24 * time.Hour) {
				oldMessageIDs = append(oldMessageIDs, message.ID)
			} else {
				newMessageIDs = append(newMessageIDs, message.ID)
			}
			count--

			if count == 0 {
				break
			}
		}

		if len(messages) != 100 {
			break // end of channel
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(newMessageIDs)))
	sort.Sort(sort.Reverse(sort.StringSlice(oldMessageIDs)))

	for i := 0; i < len(newMessageIDs); i += 100 {
		c := int(math.Min(float64(len(newMessageIDs)), 100))
		if err = discord.ChannelMessagesBulkDelete(channel.ID, newMessageIDs[i:i+c]); err != nil {
			gophbot.Log.Error("Could not bulk purge messages", zap.Error(err))
			return
		}
	}
	for _, id := range oldMessageIDs {
		if err = discord.ChannelMessageDelete(channel.ID, id); err != nil {
			gophbot.Log.Error("Could not purge messages", zap.Error(err))
			return
		}
	}

	discord.ChannelMessageDelete(cmd.Channel.ID, cmd.Message.ID)
	success = true

	message, err := cmd.Reply(gophbot.OkHand)
	if err == nil {
		time.Sleep(2 * time.Second)
		discord.ChannelMessageDelete(message.ChannelID, message.ID)
	}
}

func sendPurgeFailed(success *bool, cmd *commands.InvokedCommand) {
	if !*success {
		// We don't care about errors here, it's only a status report, which would be nice if it works
		cmd.Session.MessageReactionRemove(cmd.Channel.ID, cmd.Message.ID, gophbot.LoadingReaction, gophbot.Self.ID)
		cmd.Session.MessageReactionAdd(cmd.Channel.ID, cmd.Message.ID, gophbot.Error)
	}
}
