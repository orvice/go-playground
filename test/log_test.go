package test

import (
	"context"
	"os"
	"testing"

	"golang.org/x/exp/slog"
)

func setup() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func BenchmarkSlogKV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slog.Info("hello", "name", "Al")
	}
}

func BenchmarkSlogAttr(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		slog.LogAttrs(ctx, slog.LevelInfo, "hello", slog.String("name", "Al"))
	}
}
