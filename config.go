package logz

import (
	"io"
	"log/slog"
	"strings"
)

type config struct {
	writer          io.Writer
	sourceEnabled   bool
	level           slog.Leveler
	replacer        func(groups []string, a slog.Attr) slog.Attr
	replacerEnabled bool
	serviceVersion  string
	env             string
}

type option func(*config)

func WithWriter(w io.Writer) option {
	return option(func(cfg *config) {
		if w != nil {
			cfg.writer = w
		}
	})
}

func WithSourceEnabled(enabled bool) option {
	return option(func(cfg *config) {
		cfg.sourceEnabled = enabled
	})
}

func WithLevel(level string) option {
	return option(func(cfg *config) {
		cfg.level = parseLogLevel(level)
	})
}

func WithReplacer(replacer func(groups []string, a slog.Attr) slog.Attr) option {
	return option(func(cfg *config) {
		if replacer != nil {
			cfg.replacer = replacer
		}
	})
}

func WithReplacerEnabled(enabled bool) option {
	return option(func(cfg *config) {
		cfg.replacerEnabled = enabled
	})
}

func WithServiceVersion(version string) option {
	return option(func(cfg *config) {
		cfg.serviceVersion = version
	})
}

func WithEnv(env string) option {
	return option(func(cfg *config) {
		cfg.env = env
	})
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
