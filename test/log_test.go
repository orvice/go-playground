package test

import (
	"context"
	"io"
	"os"
	"testing"

	"go.uber.org/zap"
	"golang.org/x/exp/slog"
)

func setup() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, nil)))
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

func BenchmarkZap(b *testing.B) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	for i := 0; i < b.N; i++ {
		sugar.Infow("hello",
			"name", "Al",
		)
	}

}
