package postgres

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type Repositories struct {
    Team        repository.TeamRepository
    User        repository.UserRepository
    PullRequest repository.PullRequestRepository
    Transactor  repository.Transactor
}

func NewRepositories(ctx context.Context, connString string) (*Repositories, *DB, error) {
    db, err := NewDB(ctx, connString)
    if err != nil {
        return nil, nil, fmt.Errorf("create db connection: %w", err)
    }

    return &Repositories{
        Team:        NewTeamRepository(db),
        User:        NewUserRepository(db),
        PullRequest: NewPullRequestRepository(db),
        Transactor:  NewTransactor(db.Pool),
    }, db, nil
}

func (r *Repositories) Close(db *DB) {
    if db != nil {
        db.Close()
    }
}
