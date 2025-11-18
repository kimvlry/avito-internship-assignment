package main

import (
	"context"
	"github.com/kimvlry/avito-internship-assignment/internal/app"
	"github.com/kimvlry/avito-internship-assignment/internal/delivery/http"
	"github.com/kimvlry/avito-internship-assignment/internal/delivery/http/handler"
	"github.com/kimvlry/avito-internship-assignment/internal/domain/service"
	"github.com/kimvlry/avito-internship-assignment/internal/infrastructure/postgres"
	"github.com/kimvlry/avito-internship-assignment/pkg/logger"
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

	logger.Init(cfg.AppMode)
	logger.Info(ctx, "Starting application",
		"mode", cfg.AppMode,
		"port", app.HttpPort,
	)

	repos, db, err := postgres.NewRepositories(ctx, cfg.Postgres.GetConnString())
	if err != nil {
		logger.Error(ctx, "failed to create repos and connect to db: %v", err)
		os.Exit(1)
	}
	defer repos.Close(db)

	services := service.NewServices(repos.Team, repos.User, repos.PullRequest, repos.Transactor)

	handlers := handler.NewHandlers(services)
	server := http.NewServer(cfg.Http, handlers)

	go func() {
		if err := server.Start(); err != nil {
			logger.Error(ctx, "failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "failed to shutdown server: %v", err)
		os.Exit(1)
	}
	logger.Info(ctx, "Server gracefully stopped")
}
