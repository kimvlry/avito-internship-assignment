package repository

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
)

type PullRequestRepository interface {
    CreateWithReviewers(ctx context.Context, pr *entity.PullRequest) error
    GetByID(ctx context.Context, id string) (*entity.PullRequest, error)
    Exists(ctx context.Context, id string) (bool, error)
    UpdateStatus(ctx context.Context, prId string, status entity.PullRequestStatus) error
    GetByReviewer(ctx context.Context, userId string) ([]*entity.PullRequest, error)
    ReplaceReviewer(ctx context.Context, prId, oldUserId, newUserId string) error
}
