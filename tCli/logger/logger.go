package logger

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func Init(debug bool) {
	level := zerolog.ErrorLevel
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	logPath := filepath.Join(homeDir, "Library", "Logs", "Halo", "halo.log")

	_ = os.MkdirAll(filepath.Dir(logPath), 0700)

	rotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	writer := zerolog.MultiLevelWriter(os.Stderr, rotator)

	Logger = zerolog.New(writer).With().Timestamp().Logger()
}
