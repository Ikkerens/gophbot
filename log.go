package gophbot

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the bots main logging utility
var Log *zap.Logger

func init() {
	var (
		outFilter = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			if len(sessions) > 0 && sessions[0].LogLevel == discordgo.LogDebug && lvl == zap.DebugLevel {
				return true
			}

			return lvl == zapcore.InfoLevel
		})
		errFilter = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl > zapcore.InfoLevel
		})

		stdOut = zapcore.Lock(os.Stdout)
		stdErr = zapcore.Lock(os.Stderr)

		config = zap.NewProductionEncoderConfig()
	)
	config.EncodeTime = encodeTime
	console := zapcore.NewConsoleEncoder(config)

	Log = zap.New(zapcore.NewTee(
		zapcore.NewCore(console, stdOut, outFilter),
		zapcore.NewCore(console, stdErr, errFilter),
	))
}

func encodeTime(t time.Time, e zapcore.PrimitiveArrayEncoder) {
	e.AppendString(t.Format(time.Stamp))
}
