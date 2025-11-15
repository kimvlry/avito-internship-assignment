package entity

import "time"

type PullRequestStatus string

const (
    PROpen   PullRequestStatus = "OPEN"
    PRMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
    ID                string
    Name              string
    AuthorID          string
    Status            PullRequestStatus
    AssignedReviewers []string
    CreatedAt         time.Time
    MergedAt          *time.Time
}

func (p *PullRequest) SetMerged() error {
    now := time.Now()
    p.Status = PRMerged
    p.MergedAt = &now
    return nil
}

func (p *PullRequest) IsMerged() bool {
    return p.Status == PRMerged
}

func (p *PullRequest) HasReviewer(userId string) bool {
    for _, reviewer := range p.AssignedReviewers {
        if reviewer == userId {
            return true
        }
    }
    return false
}
