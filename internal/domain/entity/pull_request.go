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

func (p *PullRequest) Merge() {
	if p.Status == PRMerged {
		return
	}
	now := time.Now()
	p.Status = PRMerged
	p.MergedAt = &now
}

func (p *PullRequest) CanBeModified() bool {
	return p.Status == PROpen
}
