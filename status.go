package gophbot

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
)

var startup = time.Now()
var statusBuilders = []func() string{
	getGuildCount,
	getMemberCount,
	timeSinceStart,
}

func statusLoop() {
	next := 0
	for {
		status := statusBuilders[next]()
		next = (next + 1) % len(statusBuilders)

		for _, discord := range sessions {
			discord.UpdateStatus(0, status)
		}
		time.Sleep(30 * time.Second)
	}
}

func getGuildCount() string {
	guilds := 0
	for _, discord := range sessions {
		guilds += len(discord.State.Guilds)
	}

	return fmt.Sprintf("in %d servers", guilds)
}

func getMemberCount() string {
	members := 0
	for _, discord := range sessions {
		discord.State.Lock()
		for _, guild := range discord.State.Guilds {

			members += guild.MemberCount
		}
		discord.State.Unlock()
	}

	return fmt.Sprintf("with %d users", members)
}

func timeSinceStart() string {
	return "for " + durafmt.ParseShort(time.Since(startup)).String()
}
