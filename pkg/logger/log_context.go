package logger

import "context"

type logCtx struct {
    EntityId string
    Method   string
    Extra    map[string]any
}

type ctxKey struct{}

var key = ctxKey{}

type Builder struct {
    ctx    context.Context
    logCtx logCtx
}

func NewBuilder(ctx context.Context) *Builder {
    existing := logCtx{}
    if c, ok := ctx.Value(key).(logCtx); ok {
        existing = c
    }
    return &Builder{
        ctx:    ctx,
        logCtx: existing,
    }
}

func (b *Builder) WithEntityID(id string) *Builder {
    b.logCtx.EntityId = id
    return b
}

func (b *Builder) WithMethod(method string) *Builder {
    b.logCtx.Method = method
    return b
}

func (b *Builder) Build() context.Context {
    return context.WithValue(b.ctx, key, b.logCtx)
}
