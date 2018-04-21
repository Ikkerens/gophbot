package gophbot

import "github.com/bwmarrin/discordgo"

// ComputeBasePermissions calculates the permissions a guild member has, outside the scope of a channel
func ComputeBasePermissions(member *discordgo.Member, guild *discordgo.Guild) (int, error) {
	if guild.OwnerID == member.User.ID {
		return discordgo.PermissionAll, nil
	}

	everyone, err := GetRole(guild.ID, guild.ID)
	if err != nil {
		return 0, err
	}

	permissions := everyone.Permissions
	for _, roleID := range member.Roles {
		role, err := GetRole(guild.ID, roleID)
		if err != nil {
			return 0, err
		}

		permissions |= role.Permissions
	}

	if permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		return discordgo.PermissionAll, nil
	}

	return permissions, nil
}

// ComputeOverwrites calculates the permissions a channel member has, given the guilds base permissions
func ComputeOverwrites(basePermissions int, member *discordgo.Member, channel *discordgo.Channel) (int, error) {
	if basePermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		return discordgo.PermissionAll, nil
	}

	everyone, err := GetRole(channel.GuildID, channel.GuildID)
	if err != nil {
		return 0, err
	}

	everyoneOverwrite := GetOverwrite(channel, everyone)
	if everyoneOverwrite != nil {
		basePermissions &= ^everyoneOverwrite.Deny
		basePermissions |= everyoneOverwrite.Allow
	}

	var allow, deny int
	for _, roleID := range member.Roles {
		overwrite := GetOverwriteByID(channel, roleID, "role")
		if overwrite != nil {
			allow |= overwrite.Allow
			deny |= overwrite.Deny
		}
	}

	basePermissions &= ^deny
	basePermissions |= allow

	memberOverwrite := GetOverwriteByID(channel, member.User.ID, "user")
	if memberOverwrite != nil {
		basePermissions &= ^memberOverwrite.Deny
		basePermissions |= memberOverwrite.Allow
	}

	return basePermissions, nil
}

// ComputePermissions calculates the permissions a channel member has, using the guilds base permissions
func ComputePermissions(member *discordgo.Member, channel *discordgo.Channel) (int, error) {
	guild, err := GetGuild(channel.GuildID)
	if err != nil {
		return 0, err
	}

	base, err := ComputeBasePermissions(member, guild)
	if err != nil {
		return 0, err
	}

	return ComputeOverwrites(base, member, channel)
}
