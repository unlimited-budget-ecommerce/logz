package logz

import (
	"context"
	"log/slog"
	"maps"
	"os"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

type (
	logzHandler struct{ slog.Handler }
	ctxKey      struct{}
)

func (h *logzHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(ctxKey{}).(map[string]slog.Value); ok {
		for k, v := range attrs {
			r.AddAttrs(slog.Attr{Key: k, Value: v})
		}
	}

	return h.Handler.Handle(ctx, r)
}

func (h *logzHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *logzHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logzHandler{h.Handler.WithAttrs(attrs)}
}

func (h *logzHandler) WithGroup(name string) slog.Handler {
	return &logzHandler{h.Handler.WithGroup(name)}
}

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
	if cfg.env == "" {
		cfg.env = "SIT"
	}

	logger := slog.New(&logzHandler{slog.NewJSONHandler(
		cfg.writer,
		&slog.HandlerOptions{
			AddSource:   cfg.sourceEnabled,
			Level:       cfg.level,
			ReplaceAttr: cfg.replacer,
		}),
	})
	logger = logger.With(
		slog.String(string(semconv.ServiceNameKey), serviceName),
		slog.String(string(semconv.ServiceVersionKey), cfg.serviceVersion),
		slog.String(string(semconv.DeploymentEnvironmentNameKey), cfg.env),
	)
	slog.SetDefault(logger)

	slog.Info("[LOGZ] logz initialized")
}

// AddContexts adds attributes to the context for logging only.
func AddContexts(parent context.Context, attrs ...slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	var oldAttrs map[string]slog.Value
	if m, ok := parent.Value(ctxKey{}).(map[string]slog.Value); ok {
		oldAttrs = m
	}

	newAttrs := make(map[string]slog.Value, len(oldAttrs)+len(attrs))
	maps.Copy(newAttrs, oldAttrs)
	for _, a := range attrs {
		newAttrs[a.Key] = a.Value
	}

	return context.WithValue(parent, ctxKey{}, newAttrs)
}
