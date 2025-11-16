package logger

import (
    "log/slog"
    "os"
    "time"
)

func Init(appMode string) {
    levels := map[string]slog.Leveler{
        "local": slog.LevelDebug,
        "dev":   slog.LevelInfo,
    }

    timeFormat := time.RFC3339
    base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     levels[appMode],
        AddSource: appMode == "local",
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.TimeKey {
                return slog.String(a.Key, a.Value.Time().Format(timeFormat))
            }
            return a
        },
    })

    handler := NewHandlerMiddleware(base).WithAttrs([]slog.Attr{
        slog.String("env", appMode),
    })
    slog.SetDefault(slog.New(handler))
}
