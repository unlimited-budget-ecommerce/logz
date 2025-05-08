# logz - structured logging

`logz` simplifies the initialization of `log/slog` package with customisable options.

## Base Fields

- time
- level (`DEBUG`, `INFO`, `WARN`, `ERROR`)
- msg
- service.name
- service.version
- deployment.environment.name (`SIT`, `UAT`, `PRD`)

## Features

- JSON output
- Configurable
- Replace field value
- Context logging
- Mask sensitive data

## Installation

```sh
go get github.com/unlimited-budget-ecommerce/logz
```

## Usage

### Initializing logz

```go
// This function also sets the global instance logger of `log/slog` package.
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
    logz.WithServiceVersion(""),    // default: ""
    logz.WithEnv(""),               // default: "SIT"
)
```

### Basic logging

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

### Context logging

```go
ctx := context.Background()
ctx = logz.SetContextAttrs(ctx, slog.String("trace_id", "123"))

slog.InfoContext(ctx, "info") // trace_id is included in log.
```

### Masking sensitive data

use `logz.MaskXxx` with `logz.WithReplacer` option to masks matched field's data or attach `LogValue` method to a struct to control the log output of that struct.

```go
logz.WithReplacer(func(_ []string, a slog.Attr) slog.Attr {
    if a.Key == "name" {
        a.Value = slog.StringValue(logz.MaskName(a.Value.String()))
    } else if a.Key == "email" {
        a.Value = slog.StringValue(logz.MaskEmail(a.Value.String()))
    }
    return a
}),
```

use `logz.SetReplacerMap` to set replacer map for `logz.MaskMap` to masks those keys's value based on replacer map.

```go
// This function is unsafe for concurrent calls.
logz.SetReplacerMap(map[string]func(string) string{
    "name":  logz.MaskName,
    "email": logz.MaskEmail,
})

maskedMap := logz.MaskMap(m)
```
