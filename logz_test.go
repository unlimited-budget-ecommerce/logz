package logz

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitWithDefaultOption(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w // redirect stdout to a pipe
	Init("service-name")
	slog.Info("info")
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = stdout // restore stdout
	strs := strings.Split(string(out), "\n")
	m := map[string]any{}

	err := json.Unmarshal([]byte(strs[0]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "[LOGZ] logz initialized", "service-name", "", "SIT")

	err = json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "service-name", "", "SIT")
}

func TestInitWithWriter(t *testing.T) {
	b := bytes.Buffer{}
	Init("", WithWriter(&b))
	slog.Info("info")
	m := map[string]any{}
	strs := strings.Split(b.String(), "\n")

	err := json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "", "", "SIT")
}

func TestInitWithCaller(t *testing.T) {
	b := bytes.Buffer{}
	Init("", WithWriter(&b), WithSourceEnabled(true))
	slog.Info("info")
	m := map[string]any{}
	strs := strings.Split(b.String(), "\n")

	err := json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "", "", "SIT")
	assert.NotEmpty(t, m["source"])

	src, ok := m["source"].(map[string]any)

	assert.True(t, ok)
	assert.NotZero(t, src["file"])
	assert.NotZero(t, src["function"])
	assert.NotZero(t, src["line"])
}

func TestInitWithLevel(t *testing.T) {
	// test level debug
	b := bytes.Buffer{}
	Init("", WithWriter(&b), WithLevel("debug"))
	slog.Debug("debug")
	slog.Info("info")
	slog.Warn("warn")
	slog.Error("error")
	m := map[string]any{}

	strs := strings.Split(b.String(), "\n")

	assert.Len(t, strs, 6) // include newline at the end

	err := json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "DEBUG", "debug", "", "", "SIT")

	// test level info
	b.Reset()
	Init("", WithWriter(&b), WithLevel("info"))
	slog.Debug("debug")
	slog.Info("info")
	slog.Warn("warn")
	slog.Error("error")

	strs = strings.Split(b.String(), "\n")

	assert.Len(t, strs, 5)

	err = json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "", "", "SIT")

	// test level warn
	b.Reset()
	Init("", WithWriter(&b), WithLevel("warn"))
	slog.Debug("debug")
	slog.Info("info")
	slog.Warn("warn")
	slog.Error("error")

	strs = strings.Split(b.String(), "\n")

	assert.Len(t, strs, 3) // init msg won't be logged

	err = json.Unmarshal([]byte(strs[0]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "WARN", "warn", "", "", "SIT")

	// test level error
	b.Reset()
	Init("", WithWriter(&b), WithLevel("error"))
	slog.Debug("debug")
	slog.Info("info")
	slog.Warn("warn")
	slog.Error("error")

	strs = strings.Split(b.String(), "\n")

	assert.Len(t, strs, 2)

	err = json.Unmarshal([]byte(strs[0]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "ERROR", "error", "", "", "SIT")
}

func TestInitWithReplacer(t *testing.T) {
	b := bytes.Buffer{}
	Init(
		"",
		WithWriter(&b),
		WithReplacer(func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "name" {
				a.Value = slog.StringValue(MaskName(a.Value.String()))
			} else if a.Key == "email" {
				a.Value = slog.StringValue(MaskEmail(a.Value.String()))
			}
			return a
		}),
		WithReplacerEnabled(true),
	)
	slog.With("name", "john doe", "email", "john@doe.com").Info("info")
	m := map[string]any{}
	strs := strings.Split(b.String(), "\n")

	err := json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "", "", "SIT")
	assert.Equal(t, "j**n d*e", m["name"])
	assert.Equal(t, "j**n@doe.com", m["email"])
}

func TestLogContextValue(t *testing.T) {
	b := bytes.Buffer{}
	Init("", WithWriter(&b))
	ctx := SetContextAttrs(context.Background(), slog.String("uid", "123"), slog.String("trace_id", "456"))
	slog.InfoContext(ctx, "info")
	m := map[string]any{}
	strs := strings.Split(b.String(), "\n")

	err := json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "", "", "SIT")
	assert.Equal(t, "123", m["uid"])
	assert.Equal(t, "456", m["trace_id"])
}

func TestLogContextWithReplacer(t *testing.T) {
	b := bytes.Buffer{}
	Init(
		"",
		WithWriter(&b),
		WithReplacer(func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "name" {
				a.Value = slog.StringValue(MaskName(a.Value.String()))
			}
			return a
		}),
		WithReplacerEnabled(true),
	)
	ctx := SetContextAttrs(context.Background(), slog.String("name", "john doe"))
	slog.InfoContext(ctx, "info")
	m := map[string]any{}
	strs := strings.Split(b.String(), "\n")

	err := json.Unmarshal([]byte(strs[1]), &m)

	assert.NoError(t, err)
	assertBaseFields(t, m, "INFO", "info", "", "", "SIT")
	assert.Equal(t, "j**n d*e", m["name"])
}

func TestLogContextConcurrently(t *testing.T) {
	b := bytes.Buffer{}
	Init("", WithWriter(&b))
	const numGoroutines = 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			ctx := SetContextAttrs(context.Background(), slog.Int("request_id", id))
			slog.InfoContext(ctx, "")
		}(i)
	}
	wg.Wait()

	strs := strings.Split(strings.TrimSpace(b.String()), "\n")

	assert.Len(t, strs, numGoroutines+1)

	seenID := make([]bool, numGoroutines)
	for _, s := range strs[1:] { // skip init msg
		var m map[string]any

		err := json.Unmarshal([]byte(s), &m)

		assert.NoError(t, err)
		assertBaseFields(t, m, "INFO", "", "", "", "SIT")

		reqID, ok := m["request_id"].(float64)

		assert.True(t, ok)
		assert.False(t, seenID[int(reqID)])

		seenID[int(reqID)] = true
	}
}

func assertBaseFields(
	t *testing.T,
	m map[string]any,
	level, msg, serviceName, serviceVersion, env string,
) {
	assert.NotZero(t, m["time"])
	assert.Equal(t, level, m["level"])
	assert.Equal(t, msg, m["msg"])
	assert.Equal(t, serviceName, m["service.name"])
	assert.Equal(t, serviceVersion, m["service.version"])
	assert.Equal(t, env, m["deployment.environment.name"])
}
