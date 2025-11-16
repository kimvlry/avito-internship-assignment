package postgres

import (
    "context"
    "errors"
    "fmt"
    "github.com/jackc/pgx/v5"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/pkg/logger"

    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type pullRequestRepository struct {
    db *DB
}

func NewPullRequestRepository(db *DB) repository.PullRequestRepository {
    return &pullRequestRepository{db: db}
}

func (r *pullRequestRepository) CreateWithReviewers(
    ctx context.Context,
    pr *entity.PullRequest,
) error {
    query := `
		WITH inserted_pr AS (
			INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING pull_request_id
		)
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		SELECT inserted_pr.pull_request_id, unnest($6::text[])
		FROM inserted_pr
	`

    querier := r.db.GetQuerier(ctx)

    _, err := querier.Exec(ctx, query,
        pr.ID,
        pr.Name,
        pr.AuthorID,
        pr.Status,
        pr.CreatedAt,
        pr.AssignedReviewers,
    )

    if err != nil {
        if isPgUniqueViolation(err) {
            return domain.ErrPullRequestAlreadyExists
        }
        if isPgForeignKeyViolation(err) {
            return domain.ErrUserNotFound
        }
        return fmt.Errorf("exec create pr with reviewers: %w", err)
    }
    return nil
}

func (r *pullRequestRepository) GetByID(
    ctx context.Context,
    id string,
) (*entity.PullRequest, error) {
    query := `
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status,
			pr.created_at,
			pr.merged_at,
			COALESCE(
				array_agg(prr.reviewer_id) 
				FILTER (WHERE prr.reviewer_id IS NOT NULL), 
				'{}'
			) as reviewers
		FROM pull_requests pr
		LEFT JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		WHERE pr.pull_request_id = $1
		GROUP BY pr.pull_request_id
	`

    querier := r.db.GetQuerier(ctx)

    var pr entity.PullRequest
    err := querier.QueryRow(ctx, query, id).Scan(
        &pr.ID,
        &pr.Name,
        &pr.AuthorID,
        &pr.Status,
        &pr.CreatedAt,
        &pr.MergedAt,
        &pr.AssignedReviewers,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrPullRequestNotFound
        }
        return nil, fmt.Errorf("query pr by id: %w", err)
    }
    return &pr, nil
}

func (r *pullRequestRepository) Exists(ctx context.Context, id string) (bool, error) {
    query := `
		SELECT EXISTS(
			SELECT 1 FROM pull_requests WHERE pull_request_id = $1
		)
	`

    querier := r.db.GetQuerier(ctx)

    var exists bool
    err := querier.QueryRow(ctx, query, id).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check pr exists: %w", err)
    }

    return exists, nil
}

func (r *pullRequestRepository) UpdateStatus(
    ctx context.Context,
    prID string,
    status entity.PullRequestStatus,
) error {
    query := `
        UPDATE pull_requests
        SET status = $2::varchar, 
            merged_at = CASE WHEN $2::varchar = 'MERGED' THEN NOW() ELSE merged_at END
        WHERE pull_request_id = $1
    `

    querier := r.db.GetQuerier(ctx)

    result, err := querier.Exec(ctx, query, prID, string(status))
    if err != nil {
        return fmt.Errorf("exec update pr status: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrPullRequestNotFound
    }

    return nil
}

func (r *pullRequestRepository) ReplaceReviewer(
    ctx context.Context,
    prID, oldUserID, newUserID string,
) error {
    query := `
		WITH deleted AS (
			DELETE FROM pull_request_reviewers
			WHERE pull_request_id = $1 AND reviewer_id = $2
			RETURNING pull_request_id
		)
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		SELECT pull_request_id, $3
		FROM deleted
	`

    querier := r.db.GetQuerier(ctx)

    result, err := querier.Exec(ctx, query, prID, oldUserID, newUserID)
    if err != nil {
        if isPgForeignKeyViolation(err) {
            return domain.ErrUserNotFound
        }
        return fmt.Errorf("exec replace reviewer: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrReviewerNotAssigned
    }

    return nil
}

func (r *pullRequestRepository) GetByReviewer(
    ctx context.Context,
    userID string,
) ([]*entity.PullRequest, error) {
    query := `
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status,
			pr.created_at,
			pr.merged_at,
			array_agg(prr.reviewer_id) as reviewers
		FROM pull_requests pr
		JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		WHERE EXISTS (
			SELECT 1
			FROM pull_request_reviewers prr2
			WHERE prr2.pull_request_id = pr.pull_request_id
			  AND prr2.reviewer_id = $1
		)
		GROUP BY pr.pull_request_id
	`
    logger.Debug(ctx, "executing GetByReview query")

    querier := r.db.GetQuerier(ctx)

    rows, err := querier.Query(ctx, query, userID)
    if err != nil {
        logger.Error(ctx, fmt.Sprintf("query prs by reviewer: %s", userID), err)
        return nil, err
    }
    defer rows.Close()

    prs, err := scanPullRequests(rows)
    logger.Debug(ctx, fmt.Sprintf("rows affected by GetByReview query: %v", len(prs)))
    return prs, err
}

func scanPullRequests(rows pgx.Rows) ([]*entity.PullRequest, error) {
    var prs []*entity.PullRequest

    for rows.Next() {
        var pr entity.PullRequest
        err := rows.Scan(
            &pr.ID,
            &pr.Name,
            &pr.AuthorID,
            &pr.Status,
            &pr.CreatedAt,
            &pr.MergedAt,
            &pr.AssignedReviewers,
        )
        if err != nil {
            return nil, fmt.Errorf("scan pr: %w", err)
        }
        prs = append(prs, &pr)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return prs, nil
}

func (r *pullRequestRepository) GetAll(ctx context.Context) ([]*entity.PullRequest, error) {
    query := `
        SELECT 
            pr.pull_request_id,
            pr.pull_request_name,
            pr.author_id,
            pr.status,
            pr.created_at,
            pr.merged_at,
            array_agg(prr.reviewer_id) as reviewers
        FROM pull_requests pr
        LEFT JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        GROUP BY pr.pull_request_id
    `

    rows, err := r.db.GetQuerier(ctx).Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    return scanPullRequests(rows)
}
