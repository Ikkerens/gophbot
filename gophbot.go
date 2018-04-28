package gophbot

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var (
	// DB is the database connection pool that we use for obtaining guild/channel-specific settings.
	DB *gorm.DB

	// Self is the bot user itself
	Self     *discordgo.User
	sessions []*discordgo.Session
)

// Snowflake is a convenience typealias depicting the format used to store snowflakes
type Snowflake = string

func ensureShardsSetup() {
	// Don't run this in test mode
	if Self != nil || flag.Lookup("test.v") != nil {
		return
	}

	Log.Info("Setting up shards")

	// Set up the first session so we can request the amount of required shards
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		Log.Panic("Could not create initial discord session", zap.Error(err))
	}

	Self, err = discord.User("@me")
	if err != nil {
		Log.Fatal("Could not request bot user information", zap.Error(err))
	}

	// Request the gateway from Discord, and the shards.
	gateway, err := discord.GatewayBot()
	if err != nil {
		Log.Fatal("Could not request bot gateway information", zap.Error(err))
	}

	// Make a list of all handlers
	sessions = make([]*discordgo.Session, gateway.Shards)
	sessions[0] = discord

	// Set up all the sessions with a token and shard identifier
	for i := range sessions {
		if i != 0 {
			sessions[i], err = discordgo.New("Bot " + os.Getenv("TOKEN"))
			if err != nil {
				Log.Panic("Could not create shard session", zap.Error(err))
			}
		}

		// sessions[i].LogLevel = discordgo.LogDebug
		sessions[i].ShardID = i
		sessions[i].ShardCount = gateway.Shards
	}
}

// Start is the main entry point for the bot, outside of the main package to allow external packages to call back
func Start() {
	Log.Info("Connecting to database")
	db, err := setupDB()
	if err != nil {
		Log.Fatal("Could not connect to database", zap.Error(err))
	}
	defer db.Close()
	DB = db

	Log.Info("Connecting all shards")

	// Connect all the shards
	for _, shard := range sessions {
		if err := shard.Open(); err != nil {
			Log.Fatal("Could not open shard websocket", zap.Error(err))
		}
		//noinspection GoDeferInLoop
		defer shard.Close()
	}

	go statusLoop()

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
