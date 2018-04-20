package gophbot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func setupLog() {
	var (
		outFilter = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapcore.InfoLevel || (sessions[0].LogLevel == discordgo.LogDebug && lvl == zapcore.DebugLevel)
		})
		errFilter = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl > zapcore.InfoLevel
		})

		stdOut = zapcore.Lock(os.Stdout)
		stdErr = zapcore.Lock(os.Stderr)

		console = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	)

	Log = zap.New(zapcore.NewTee(
		zapcore.NewCore(console, stdOut, outFilter),
		zapcore.NewCore(console, stdErr, errFilter),
	))
}
