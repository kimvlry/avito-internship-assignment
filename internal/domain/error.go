package domain

type Error string

func (e Error) Error() string {
    return string(e)
}

const (
    ErrUserNotFound             Error = "user not found"
    ErrTeamNotFound             Error = "team not found"
    ErrTeamAlreadyExists        Error = "team already exists"
    ErrPullRequestNotFound      Error = "pull request not found"
    ErrPullRequestAlreadyExists Error = "pull request already exists"
    ErrNoReviewerCandidate      Error = "no reviewer candidate available"
    ErrPullRequestIsMerged      Error = "pull request is merged"
    ErrReviewerNotAssigned      Error = "reviewer not assigned"
)
