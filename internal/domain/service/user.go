package service

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
    "github.com/kimvlry/avito-internship-assignment/pkg/logger"
)

type User struct {
    userRepo repository.UserRepository
    prRepo   repository.PullRequestRepository
}

func NewUser(userRepo repository.UserRepository, prRepo repository.PullRequestRepository) *User {
    return &User{
        userRepo: userRepo,
        prRepo:   prRepo,
    }
}

func (s *User) SetIsActive(ctx context.Context, userId string, isActive bool) (*entity.User, error) {
    user, err := s.userRepo.SetIsActive(ctx, userId, isActive)
    if err != nil {
        return nil, fmt.Errorf("set user active status: %w", err)
    }
    return user, nil
}

func (s *User) GetReviewAssignments(ctx context.Context, userId string) ([]*entity.PullRequest, error) {
    pullRequests, err := s.prRepo.GetByReviewer(ctx, userId)
    if err != nil {
        logger.Error(ctx, fmt.Sprintf("get pull requests: %s", userId), err)
        return nil, err
    }
    return pullRequests, nil
}

func (s *User) GetByID(ctx context.Context, userID string) (*entity.User, error) {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("get user by id: %w", err)
    }
    return user, nil
}
