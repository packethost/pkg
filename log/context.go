package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxLogger struct{}

// ContextWithLogger returns a copy of parent in which the logger is embedded.
func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, logger)
}

// GetLogger returns the logger associated with this context, or no-op logger
// if no logger is embedded with key.
func GetLogger(ctx context.Context) Logger {
	if logger, ok := ctx.Value(ctxLogger{}).(Logger); ok {
		return logger
	}
	return Logger{s: zap.NewNop().Sugar()}
}
