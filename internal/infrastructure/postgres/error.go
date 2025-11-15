package postgres

import (
    "errors"

    "github.com/jackc/pgx/v5/pgconn"
)

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
    ErrReviewerNotAssigned      Error = "reviewer not assigned"
    ErrUserHasActiveAssignments Error = "user has active PR assignments in other team"
    ErrTooManyReviewers         Error = "too many reviewers assigned"
    ErrUserAlreadyExists        Error = "user already exists"
)

const (
    pgUniqueViolation     = "23505"
    pgForeignKeyViolation = "23503"
    pgCheckViolation      = "23514"
)

func isPgUniqueViolation(err error) bool {
    var pgErr *pgconn.PgError
    return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}

func isPgForeignKeyViolation(err error) bool {
    var pgErr *pgconn.PgError
    return errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation
}

func isPgCheckViolation(err error) bool {
    var pgErr *pgconn.PgError
    return errors.As(err, &pgErr) && pgErr.Code == pgCheckViolation
}
