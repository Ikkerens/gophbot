package gophbot

import (
	"os"
	"os/signal"
	"syscall"
)

// Start is the main entry point for the bot, outside of the main package to allow external packages to call back
func Start() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
