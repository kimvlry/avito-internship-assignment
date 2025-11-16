package logger

import (
    "context"
    "log/slog"
)

// HandlerMiddleware is a decorator for slog.Handler that injects
// contextual fields (user_id, request_id, etc.) into every log record.
// It implements the slog.Handler interface.
type HandlerMiddleware struct {
    next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
    return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
    return h.next.Enabled(ctx, level)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
    if c, ok := ctx.Value(key).(logCtx); ok {
        if c.EntityId != "" {
            rec.Add("entity_id", c.EntityId)
        }
        if c.Method != "" {
            rec.Add("method", c.Method)
        }
        for k, v := range c.Extra {
            rec.Add(k, v)
        }
    }
    return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &HandlerMiddleware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
    return &HandlerMiddleware{next: h.next.WithGroup(name)}
}
