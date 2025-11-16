package service

import (
    "context"
    "fmt"

    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type StatsService struct {
    prRepo repository.PullRequestRepository
}

func NewStatsService(prRepo repository.PullRequestRepository) *StatsService {
    return &StatsService{prRepo: prRepo}
}

type UserAssignmentStat struct {
    UserID      string `json:"user_id"`
    AssignedPRs int    `json:"assigned_prs"`
}

type PRReviewerStat struct {
    PullRequestID string   `json:"pull_request_id"`
    ReviewerCount int      `json:"reviewer_count"`
    Reviewers     []string `json:"reviewers"`
}

func (s *StatsService) GetUserAssignmentStats(ctx context.Context) ([]UserAssignmentStat, error) {
    prs, err := s.prRepo.GetAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("get all PRs: %w", err)
    }

    counts := make(map[string]int)
    for _, pr := range prs {
        for _, r := range pr.AssignedReviewers {
            counts[r]++
        }
    }

    stats := make([]UserAssignmentStat, 0, len(counts))
    for userID, c := range counts {
        stats = append(stats, UserAssignmentStat{
            UserID:      userID,
            AssignedPRs: c,
        })
    }
    return stats, nil
}
