package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/Oloruntobi1/grey/pkg/otel"
)

var (
	onceSlog   sync.Once
	slogLogger *slog.Logger
)

func NewSlog(ctx context.Context) *slog.Logger {
	onceSlog.Do(func() {
		opts := slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		jh := slog.NewJSONHandler(os.Stdout, &opts)

		h := otel.OtelHandler{H: jh}
		l := slog.New(h).With("app", "grey-wallet-app")
		slogLogger = l
	})

	return slogLogger
}
