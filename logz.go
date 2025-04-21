package logz

import (
	"log/slog"
	"os"
)

// TODO: context to log

// Init initializes the logger with the provided options and set it to slog's default logger.
// Default writer is set to [os.Stdout] if not provided.
// Default log level is set to Info if not provided.
//
// **You should use the gloabal instance from slog package to log messages.**
func Init(serviceName string, opts ...option) {
	cfg := config{}
	for _, fn := range opts {
		fn(&cfg)
	}
	if cfg.writer == nil {
		cfg.writer = os.Stdout
	}
	if !cfg.replacerEnabled {
		cfg.replacer = nil
	}

	logger := slog.New(slog.NewJSONHandler(cfg.writer, &slog.HandlerOptions{
		AddSource:   cfg.callerEnabled,
		Level:       cfg.level,
		ReplaceAttr: cfg.replacer,
	}))
	logger = logger.With(slog.String("service", serviceName))
	slog.SetDefault(logger)

	slog.Info("logz initialized")
}
