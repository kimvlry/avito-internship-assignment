package service

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
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
        return nil, fmt.Errorf("get pull requests: %w", err)
    }
    return pullRequests, nil
}
