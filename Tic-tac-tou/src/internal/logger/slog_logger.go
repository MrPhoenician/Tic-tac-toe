package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	*slog.Logger
}

func NewSlogLogger() Logger {
	return &SlogLogger{
		Logger: slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}
}

func (sl *SlogLogger) Debug(msg string, args ...any) {
	sl.Logger.Debug(msg, args...)
}

func (sl *SlogLogger) Info(msg string, args ...any) {
	sl.Logger.Info(msg, args...)
}

func (sl *SlogLogger) Warn(msg string, args ...any) {
	sl.Logger.Warn(msg, args...)
}

func (sl *SlogLogger) Error(msg string, args ...any) {
	sl.Logger.Error(msg, args...)
}
