package service

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type UserService struct {
    userRepo repository.UserRepository
    prRepo   repository.PullRequestRepository
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PullRequestRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
        prRepo:   prRepo,
    }
}

func (s *UserService) SetIsActive(ctx context.Context, userId string, isActive bool) error {
    if err := s.userRepo.SetIsActive(ctx, userId, isActive); err != nil {
        return fmt.Errorf("set user active status: %w", err)
    }
    return nil
}

func (s *UserService) GetReviewAssignments(ctx context.Context, userId string) ([]*entity.PullRequest, error) {
    pullRequests, err := s.prRepo.GetByReviewer(ctx, userId)
    if err != nil {
        return nil, fmt.Errorf("get pull requests: %w", err)
    }
    return pullRequests, nil
}
