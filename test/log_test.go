package test

import (
	"context"
	"io"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type discardWriteSyncer struct{}

func (discardWriteSyncer) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discardWriteSyncer) Sync() error {
	return nil
}

func newZapCore() zapcore.Core {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{}), discardWriteSyncer{}, zap.DebugLevel)
	return core
}

func newZapLogger() *zap.Logger {
	logger, _ := zap.NewProduction(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return newZapCore()
	}))
	return logger
}

func BenchmarkZapStruct(b *testing.B) {
	logger := newZapLogger()
	defer logger.Sync()
	for i := 0; i < b.N; i++ {
		logger.Info("hello", zap.String("name", "Al"))
	}

}

func BenchmarkZap(b *testing.B) {
	logger := newZapLogger()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	for i := 0; i < b.N; i++ {
		sugar.Infow("hello",
			"name", "Al",
		)
	}

}
