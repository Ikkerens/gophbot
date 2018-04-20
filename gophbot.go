package gophbot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Self is the bot user itself
var Self *discordgo.User
var sessions []*discordgo.Session

// Snowflake is a convenience typealias depicting the format used to store snowflakes
type Snowflake = string

func ensureShardsSetup() {
	if Self != nil || os.Getenv("TOKEN") == "" {
		return
	}

	setupLog()
	Log.Info("Setting up shards")

	// Set up the first session so we can request the amount of required shards
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	Self, err = discord.User("@me")
	if err != nil {
		panic(err)
	}

	// Request the gateway from Discord, and the shards.
	gateway, err := discord.GatewayBot()
	if err != nil {
		panic(err)
	}

	// Make a list of all handlers
	sessions = make([]*discordgo.Session, gateway.Shards)
	sessions[0] = discord

	// Set up all the sessions with a token and shard identifier
	for i := range sessions {
		if i != 0 {
			sessions[i], err = discordgo.New("Bot " + os.Getenv("TOKEN"))
			if err != nil {
				panic(err)
			}
		}

		// sessions[i].LogLevel = discordgo.LogDebug
		sessions[i].ShardID = i
		sessions[i].ShardCount = gateway.Shards
	}
}

// Start is the main entry point for the bot, outside of the main package to allow external packages to call back
func Start() {
	Log.Info("Connecting all shards")

	// Connect all the shards
	for _, shard := range sessions {
		if err := shard.Open(); err != nil {
			panic(err)
		}
		//noinspection GoDeferInLoop
		defer shard.Close()
	}

	Log.Info("Bot is now running. Press Ctrl+C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// AddHandler adds an event handler to all shards
func AddHandler(handler interface{}) {
	ensureShardsSetup()

	for _, session := range sessions {
		session.AddHandler(handler)
	}
}
