package main

import (
	"github.com/ikkerens/gophbot"
	_ "github.com/ikkerens/gophbot/commands"
)

func main() {
	defer gophbot.Start()
}
