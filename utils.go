package gophbot

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func getShardID(guildID Snowflake) (int, error) {
	snow, err := strconv.ParseInt(guildID, 10, 64)
	if err != nil {
		return 0, err
	}

	shard := int((snow >> 22) % int64(sessions[0].ShardCount))

	if shard < 0 || shard >= len(sessions) {
		return 0, errors.New("shard id out of bounds")
	}

	return shard, nil
}

func getSession(guildID Snowflake) (*discordgo.Session, error) {
	shard, err := getShardID(guildID)
	if err != nil {
		return nil, err
	}

	return sessions[shard], nil
}

// GetGuild attempts to get a guild instance from the shard state cache, and if none exists, attempts to obtain it
// from the Discord API. Will err if the guild does not exist or if this guild is unreachable for the bot.
func GetGuild(guildID Snowflake) (*discordgo.Guild, error) {
	discord, err := getSession(guildID)
	if err != nil {
		return nil, err
	}

	guild, err := discord.State.Guild(guildID)
	if err == nil {
		return guild, nil
	}

	guild, err = discord.Guild(guildID)
	if err != nil {
		return nil, err
	}
	discord.State.GuildAdd(guild)

	return guild, nil
}

func GetChannel(discord *discordgo.Session, channelID Snowflake) (*discordgo.Channel, error) {
	channel, err := discord.State.Channel(channelID)
	if err == nil {
		return channel, nil
	}

	channel, err = discord.Channel(channelID)
	if err != nil {
		return nil, err
	}
	discord.State.ChannelAdd(channel)

	return channel, err
}

// GetGuildMember will attempt to obtain a member instance for this guild member.
// Will err if this user is not a member of this guild
func GetGuildMember(guildID, userID Snowflake) (*discordgo.Member, error) {
	discord, err := getSession(guildID)
	if err != nil {
		return nil, err
	}

	member, err := discord.State.Member(guildID, userID)
	if err == nil {
		return member, nil
	}

	member, err = discord.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}
	discord.State.MemberAdd(member)

	return member, nil
}

// GetRole will attempt to obtain a role instance for this role
// Will err if the given ID is not a role in the guild.
func GetRole(guildID, roleID Snowflake) (*discordgo.Role, error) {
	discord, err := getSession(guildID)
	if err != nil {
		return nil, err
	}

	role, err := discord.State.Role(guildID, roleID)
	if err == nil {
		return role, nil
	}

	roles, err := discord.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.ID == roleID {
			discord.State.RoleAdd(guildID, role)
			return role, nil
		}
	}

	return nil, errors.New("role does not exist or is not part of that guild")
}

// GetOverwrite will attempt to obtain the PermissionOverwrite instance for the target (Role, Member or User)
// Will return nil if no such overwrite exists
func GetOverwrite(channel *discordgo.Channel, target interface{}) *discordgo.PermissionOverwrite {
	var id Snowflake
	var typ string

	switch t := target.(type) {
	case *discordgo.Role:
		id = t.ID
		typ = "role"
	case *discordgo.Member:
		id = t.User.ID
		typ = "user"
	case *discordgo.User:
		id = t.ID
		typ = "user"
	default:
		panic(errors.New("invalid overwrite target type"))
	}

	return GetOverwriteByID(channel, id, typ)
}

// GetOverwriteByID will attempt to obtain the PermissionOverwrite instance for the given ID and type ("role" or "user")
// Will return nil if no such overwrite exists
func GetOverwriteByID(channel *discordgo.Channel, id Snowflake, typ string) *discordgo.PermissionOverwrite {
	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == typ && overwrite.ID == id {
			return overwrite
		}
	}

	return nil
}
