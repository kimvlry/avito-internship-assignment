package repository

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Exists(ctx context.Context, id string) (bool, error)
    GetByID(ctx context.Context, id string) (*entity.User, error)
    GetByTeam(ctx context.Context, teamName string) ([]entity.User, error)
    SetIsActive(ctx context.Context, id string, isActive bool) (*entity.User, error)
    GetRandomActiveTeamUsers(
        ctx context.Context,
        teamName string,
        excludeUserIDs []string,
        maxCount int,
    ) ([]entity.User, error)
    CheckUsersAvailableForTeam(ctx context.Context, userIDs []string, teamName string) error
}
