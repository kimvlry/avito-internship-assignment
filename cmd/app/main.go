package main

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/app"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/handler"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"
    "github.com/kimvlry/avito-internship-assignment/internal/infrastructure/postgres"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    ctx := context.Background()

    cfg, err := app.LoadConfig()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    repos, db, err := postgres.NewRepositories(ctx, cfg.Postgres.GetConnString())
    if err != nil {
        log.Fatalf("failed to create repos and connect to db: %v", err)
    }
    defer repos.Close(db)

    services := service.NewServices(repos.Team, repos.User, repos.PullRequest, repos.Transactor)

    handlers := handler.NewHandlers(services)
    server := http.NewServer(cfg.Http, handlers)

    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("failed to start server: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
    defer shutdownCancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("failed to shutdown server: %v", err)
    }
    log.Println("Server gracefully stopped")
}
