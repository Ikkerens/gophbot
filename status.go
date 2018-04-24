package gophbot

func statusLoop() {
	for _, discord := range sessions {
		discord.UpdateStatus(0, "with Development")
	}
}
