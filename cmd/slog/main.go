package main

import (
	"context"
	"net"

	"golang.org/x/exp/slog"
)

func main() {
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr)))
	slog.Info("hello", "name", "Al")
	slog.Error("oops", "status", 500)
	ctx := context.Background()
	slog.LogAttrs(ctx, slog.LevelError, "oops",
		slog.Int("status", 500), slog.Any("err", net.ErrClosed))
}
