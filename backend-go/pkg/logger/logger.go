package logger

import (
	"log/slog"
	"os"
)

var L *slog.Logger

func Init(env string) {
	var handler slog.Handler

	if env == "production" {
		// JSON logs for production (CloudWatch, ELK, etc.)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Readable text logs for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	L = slog.New(handler)
	slog.SetDefault(L)
}

func Info(msg string, args ...any) {
	L.Info(msg, args...)
}

func Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("error", err))
	}
	L.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	L.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	L.Warn(msg, args...)
}
