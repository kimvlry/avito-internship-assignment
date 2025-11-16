package logger

import (
    "context"
    "log/slog"
)

func Info(ctx context.Context, msg string, args ...any) {
    slog.InfoContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
    slog.ErrorContext(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
    slog.DebugContext(ctx, msg, args...)
}
