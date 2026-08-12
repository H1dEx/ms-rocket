package logger

import (
	"context"

	"go.uber.org/zap"
)

type NoopLogger struct{}

func (l *NoopLogger) Info(_ context.Context, _ string, _ ...zap.Field)  {}
func (l *NoopLogger) Error(_ context.Context, _ string, _ ...zap.Field) {}
