package check

import (
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/api"
    "strings"
)

type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidPullRequestCreate(req api.PostPullRequestCreateRequestObject) error {
    if strings.TrimSpace(req.Body.PullRequestId) == "" {
        return ValidationError{"pull_request_id", "empty"}
    }
    if strings.TrimSpace(req.Body.PullRequestName) == "" {
        return ValidationError{"pull_request_name", "empty"}
    }
    if strings.TrimSpace(req.Body.AuthorId) == "" {
        return ValidationError{"author_id", "empty"}
    }
    return nil
}

func ValidTeamCreate(req api.PostTeamAddRequestObject) error {
    if strings.TrimSpace(req.Body.TeamName) == "" {
        return ValidationError{"team_name", "cannot be empty"}
    }
    return nil
}

func ValidUserID(userID string) error {
    if strings.TrimSpace(userID) == "" {
        return ValidationError{"user_id", "cannot be empty"}
    }
    return nil
}
