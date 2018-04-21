package main

import (
	"github.com/ikkerens/gophbot"
	_ "github.com/ikkerens/gophbot/commands"
	_ "github.com/ikkerens/gophbot/handlers"
)

func main() {
	defer gophbot.Start()
}
