package test

import (
	"context"
	"testing"

	"golang.org/x/exp/slog"
)

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
