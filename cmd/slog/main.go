package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	slogwebhhok "github.com/samber/slog-webhook"
	"golang.org/x/exp/slog"
)

func setHandler() {
	slog.SetDefault(newLogger())
}

func webhookLogger() *slog.Logger {
	url := os.Getenv("SLOG_ENDPOINT")

	logger := slog.New(slogwebhhok.Option{Level: slog.LevelDebug, Endpoint: url}.NewWebhookHandler())
	logger = logger.With("release", "v1.0.0")
	return logger
}

func newLogger() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(newHttpWriter(), &slog.HandlerOptions{
		AddSource: false,
	})
	var handler = jsonHandler.WithAttrs([]slog.Attr{
		slog.String("sub_type", "sub_type_value"),
		slog.String("system", "optimus"),
	}).WithGroup("message")
	return slog.New(handler)
}

func main() {
	setHandler()
	slog.Info("hello", "name", "Al")
	slog.Error("oops", "status", 500)
	ctx := context.Background()
	slog.LogAttrs(ctx, slog.LevelError, "oops",
		slog.Int("status", 500), slog.Any("err", net.ErrClosed))

	logger := newLogger()
	logger.LogAttrs(ctx, slog.LevelError, "oops", slog.Int("status", http.StatusAccepted))
	logger.LogAttrs(ctx, slog.LevelInfo, "", slog.Group("group", "key", "value"))
	logger.Info("Usage Statistics",
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
		slog.Int("cpu", 10),
		slog.String("version", "v0.0.1"),
	)

	webhookLogger := webhookLogger()
	webhookLogger.Info("Usage Statistics",
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
		slog.Int("cpu", 10),
		slog.String("version", "v0.0.1"),
	)
}

func newHttpWriter() io.Writer {
	return &httpWriter{
		client:   http.DefaultClient,
		endpoint: os.Getenv("SLOG_ENDPOINT"),
	}
}

type httpWriter struct {
	client   *http.Client
	endpoint string
}

func (h *httpWriter) Write(p []byte) (int, error) {
	req, err := http.NewRequest(http.MethodPost, h.endpoint, bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	fmt.Print(string(p))
	return len(p), nil
}
