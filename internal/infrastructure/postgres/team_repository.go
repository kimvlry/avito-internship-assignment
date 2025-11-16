package postgres

import (
    "context"
    "errors"
    "fmt"
    "github.com/jackc/pgx/v5"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type teamRepository struct {
    db *DB
}

func NewTeamRepository(db *DB) repository.TeamRepository {
    return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *entity.Team) error {
    query := `
		INSERT INTO teams (name)
		VALUES ($1)
	`

    querier := r.db.GetQuerier(ctx)

    _, err := querier.Exec(ctx, query, team.Name)
    if err != nil {
        if isPgUniqueViolation(err) {
            return domain.ErrTeamAlreadyExists
        }
        return fmt.Errorf("exec create team: %w", err)
    }
    return nil
}

func (r *teamRepository) GetByName(ctx context.Context, name string) (*entity.Team, error) {
    query := `
		SELECT name
		FROM teams
		WHERE name = $1
	`

    querier := r.db.GetQuerier(ctx)

    var team entity.Team
    err := querier.QueryRow(ctx, query, name).Scan(
        &team.Name,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrTeamNotFound
        }
        return nil, fmt.Errorf("query team by name: %w", err)
    }

    return &team, nil
}

func (r *teamRepository) Exists(ctx context.Context, name string) (bool, error) {
    query := `
		SELECT EXISTS(
			SELECT 1 FROM teams WHERE name = $1
		)
	`

    querier := r.db.GetQuerier(ctx)

    var exists bool
    err := querier.QueryRow(ctx, query, name).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check team exists: %w", err)
    }

    return exists, nil
}
