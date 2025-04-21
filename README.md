# logz - structured logging

`logz` simplifies the initialization of `log/slog` package with customisable options.

## Base Fields

- time
- level (`DEBUG`, `INFO`, `WARN`, `ERROR`)
- msg
- service (service name)

## Features

- JSON output
- Configurable
- Replace field value
- Context logging

## Installation

```sh
go get github.com/unlimited-budget-ecommerce/logz
```

## Usage

Initialize logz. It also sets the global instance of `log/slog` package.

```go
logz.Init(
    "service-name",
    logz.WithWriter(os.Stdout),     // default: [os.Stdout]
    logz.WithCallerEnabled(true),   // default: false
    logz.WithLevel("info"),         // default: "info"
    logz.WithReplacer(func(groups []string, a slog.Attr) slog.Attr {
        if a.Key == "name" {
            a.Value = slog.StringValue(logz.MaskName(a.Value.String()))
        } else if a.Key == "email" {
            a.Value = slog.StringValue(logz.MaskEmail(a.Value.String()))
        }
        return a
    }),                             // default: nil
    logz.WithReplacerEnabled(true), // default: false
)
```

Basic logging.

```go
slog.Debug("debug")

// add key-value pairs to log using [slog.With].
// it creates a copy of logger with inputed fields included.
// unmatched pair results in "!BADKEY".
slog.With("name", "john doe", "email", "john@doe.com").Info("info")

// wraps key-value pair with [slog.TYPE] to avoid "!BADKEY"
// and improve performance.
slog.Warn("warn", slog.String("name", "john doe"))

// groups key-value pairs using [slog.Group].
slog.Error("error", slog.Group("req",
    slog.String("method", "POST"),
    slog.String("path", "/v1/users"),
))
```

Context logging.

```go
ctx := context.Background()
ctx = logz.AddContext(ctx, slog.String("traceID", "123"))

slog.InfoContext(ctx, "info") // traceID is included in log.
```
