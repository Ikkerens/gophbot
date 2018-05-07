package main

import (
	"os"

	"github.com/ikkerens/gophbot"
	_ "github.com/ikkerens/gophbot/commands"
	_ "github.com/ikkerens/gophbot/handlers"
)

func main() {
	if os.Getenv("TOKEN") == "" {
		gophbot.Log.Fatal("Environment variable TOKEN, which is required, is not set.")
	}
	if os.Getenv("SQL_DSN") == "" {
		gophbot.Log.Fatal("Environment variable SQL_DSN, which is required, is not set.")
	}

	gophbot.Start()
}
