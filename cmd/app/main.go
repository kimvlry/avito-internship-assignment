package main

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/app"
    "github.com/kimvlry/avito-internship-assignment/internal/infrastructure/postgres"
    "log"
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

    //services := service.NewServices(repos.Team, repos.User, repos.PullRequest, repos.Transactor)
}
