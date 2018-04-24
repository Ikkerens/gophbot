package handlers

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/ikkerens/gophbot"
)

func TestMain(m *testing.M) {
	gophbot.Self = &discordgo.User{}
	os.Exit(m.Run())
}

func TestCase(t *testing.T) {
	AddCommand("TeSt", nil)
	_, exists := commands["test"]
	if !exists {
		t.Fail()
	}
}

func TestCommandParser(t *testing.T) {
	callback := make(chan bool, 1)
	AddCommand("parsetest", func(session *discordgo.Session, cmd *InvokedCommand) {
		if len(cmd.Args) != 4 {
			callback <- false
			return
		}

		callback <- true
	})

	handleCommand(nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:  &discordgo.User{ID: "user"},
			Content: "/parsetest a b c d",
		},
	})

	select {
	case result := <-callback:
		if !result {
			t.Fail()
		}
	default:
		t.Fail()
	}
}
