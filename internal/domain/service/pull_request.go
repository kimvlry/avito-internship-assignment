package service

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
    "time"
)

type PullRequestService struct {
    prRepository   repository.PullRequestRepository
    userRepository repository.UserRepository
    tx             Transactor
}

func NewPullRequestService(
    prRepo repository.PullRequestRepository,
    userRepo repository.UserRepository,
    tx Transactor,
) *PullRequestService {
    return &PullRequestService{
        prRepository:   prRepo,
        userRepository: userRepo,
        tx:             tx,
    }
}

func (s *PullRequestService) CreatePullRequestWithReviewers(
    ctx context.Context,
    prId,
    prName,
    authorId string,
) (*entity.PullRequest, error) {

    ok, err := s.prRepository.Exists(ctx, prId)
    if err != nil {
        return nil, fmt.Errorf("check pr exists: %w", err)
    }
    if ok {
        return nil, domain.ErrPullRequestAlreadyExists
    }

    author, err := s.userRepository.GetByID(ctx, authorId)
    if err != nil {
        return nil, fmt.Errorf("get pr author: %w", err)
    }

    var createdPr *entity.PullRequest
    err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        reviewers, err := s.userRepository.GetRandomActiveTeamUsers(txCtx, author.TeamName, []string{authorId}, 2)
        if err != nil {
            return fmt.Errorf("get reviewers: %w", err)
        }
        reviewersIds := make([]string, len(reviewers))
        for i, reviewer := range reviewers {
            reviewersIds[i] = reviewer.ID
        }

        pr := &entity.PullRequest{
            ID:                prId,
            Name:              prName,
            AuthorID:          authorId,
            Status:            entity.PROpen,
            AssignedReviewers: reviewersIds,
            CreatedAt:         time.Now(),
            MergedAt:          nil,
        }

        if err := s.prRepository.CreateWithReviewers(txCtx, pr); err != nil {
            return fmt.Errorf("create pr: %w", err)
        }
        createdPr = pr
        return nil
    })

    if err != nil {
        return nil, err
    }
    return createdPr, nil
}

func (s *PullRequestService) ReassignReviewer(
    ctx context.Context,
    prId,
    oldUserId string,
) (*entity.PullRequest, string, error) {

    var updatedPr *entity.PullRequest
    var newUserId string

    err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        pr, err := s.prRepository.GetByID(txCtx, prId)
        if err != nil {
            return fmt.Errorf("get pr: %w", err)
        }
        if pr.IsMerged() {
            return domain.ErrPullRequestIsMerged
        }
        if !pr.HasReviewer(oldUserId) {
            return domain.ErrReviewerNotAssigned
        }

        oldUser, err := s.userRepository.GetByID(txCtx, oldUserId)
        if err != nil {
            return fmt.Errorf("get old user: %w", err)
        }

        excludedIds := append(pr.AssignedReviewers, pr.AuthorID)

        replacements, err := s.userRepository.GetRandomActiveTeamUsers(txCtx, oldUser.TeamName, excludedIds, 1)
        if err != nil {
            return fmt.Errorf("get replacements: %w", err)
        }
        if len(replacements) == 0 {
            return domain.ErrNoReviewerCandidate
        }

        newReviewer := replacements[0]
        newUserId = newReviewer.ID

        if err := s.prRepository.ReplaceReviewer(txCtx, prId, oldUserId, newUserId); err != nil {
            return fmt.Errorf("replace reviewer: %w", err)
        }
        updatedPr, err = s.prRepository.GetByID(txCtx, prId)
        if err != nil {
            return fmt.Errorf("get updated pr: %w", err)
        }
        return nil
    })

    if err != nil {
        return nil, "", err
    }
    return updatedPr, newUserId, nil
}

func (s *PullRequestService) Merge(ctx context.Context, prId string) (*entity.PullRequest, error) {
    pr, err := s.prRepository.GetByID(ctx, prId)
    if err != nil {
        return nil, fmt.Errorf("get pr: %w", err)
    }
    if pr.IsMerged() {
        return pr, nil
    }
    if err := pr.SetMerged(); err != nil {
        return nil, fmt.Errorf("merge pr: %w", err)
    }
    if err = s.prRepository.UpdateStatus(ctx, prId, entity.PRMerged); err != nil {
        return nil, fmt.Errorf("update pr status: %w", err)
    }
    return pr, nil
}
