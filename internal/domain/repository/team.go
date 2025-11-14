package repository

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
)

type TeamRepository interface {
    Create(ctx context.Context, team *entity.Team) error
    GetByName(ctx context.Context, name string) (*entity.Team, error)
    Exists(ctx context.Context, teamName string) (bool, error)
}
