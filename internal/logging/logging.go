package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level  string
	Format string
}

func New(c Config) (*slog.Logger, error) { return NewWithWriter(c, os.Stdout) }

func NewWithWriter(c Config, w io.Writer) (*slog.Logger, error) {
	l, e := parseLevel(c.Level)
	if e != nil {
		return nil, e
	}
	o := &slog.HandlerOptions{Level: l}
	var h slog.Handler
	switch strings.ToLower(strings.TrimSpace(c.Format)) {
	case "json", "":
		h = slog.NewJSONHandler(w, o)
	case "text":
		h = slog.NewTextHandler(w, o)
	default:
		return nil, fmt.Errorf("unsupported log format %q: expected json or text", c.Format)
	}
	return slog.New(h), nil
}

func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported log level %q: expected debug, info, warn, or error", s)
	}
}
