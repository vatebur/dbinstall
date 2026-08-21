package logging

import (
	"io"
	"log/slog"
)

func New(writer io.Writer, jsonOutput bool, level slog.Level) *slog.Logger {
	options := &slog.HandlerOptions{Level: level}
	if jsonOutput {
		return slog.New(slog.NewJSONHandler(writer, options))
	}
	return slog.New(slog.NewTextHandler(writer, options))
}
