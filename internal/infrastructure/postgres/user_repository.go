package postgres

import (
    "context"
    "errors"
    "fmt"
    "github.com/Masterminds/squirrel"
    "github.com/jackc/pgx/v5"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type userRepository struct {
    db *DB
}

func NewUserRepository(db *DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    query := `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`

    querier := r.db.GetQuerier(ctx)

    _, err := querier.Exec(ctx, query,
        user.ID,
        user.Username,
        user.TeamName,
        user.IsActive,
    )

    if err != nil {
        if isPgUniqueViolation(err) {
            return ErrUserAlreadyExists
        }
        if isPgForeignKeyViolation(err) {
            return domain.ErrTeamNotFound
        }
        return fmt.Errorf("exec create user: %w", err)
    }

    return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
    query := `
		UPDATE users
		SET username = $2, team_name = $3, is_active = $4
		WHERE user_id = $1
	`

    querier := r.db.GetQuerier(ctx)

    result, err := querier.Exec(ctx, query,
        user.ID,
        user.Username,
        user.TeamName,
        user.IsActive,
    )

    if err != nil {
        if isPgForeignKeyViolation(err) {
            return domain.ErrTeamNotFound
        }
        return fmt.Errorf("exec update user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrUserNotFound
    }

    return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
    query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`

    querier := r.db.GetQuerier(ctx)

    var user entity.User
    err := querier.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.TeamName,
        &user.IsActive,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("query user by id: %w", err)
    }

    return &user, nil
}

func (r *userRepository) GetByTeam(ctx context.Context, teamName string) ([]entity.User, error) {
    query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

    querier := r.db.GetQuerier(ctx)

    rows, err := querier.Query(ctx, query, teamName)
    if err != nil {
        return nil, fmt.Errorf("query users by team: %w", err)
    }
    defer rows.Close()
    return scanUsers(rows)
}

func (r *userRepository) SetIsActive(ctx context.Context, id string, isActive bool) (*entity.User, error) {
    query := `
		UPDATE users
		SET is_active = $2
		WHERE user_id = $1
		RETURNING user_id, username, team_name, is_active
	`

    var u entity.User
    querier := r.db.GetQuerier(ctx)

    err := querier.QueryRow(ctx, query, id, isActive).Scan(
        &u.ID,
        &u.Username,
        &u.TeamName,
        &u.IsActive,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("exec set is_active: %w", err)
    }
    return &u, nil
}

func (r *userRepository) Exists(ctx context.Context, id string) (bool, error) {
    query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE user_id = $1
		)
	`

    querier := r.db.GetQuerier(ctx)

    var exists bool
    err := querier.QueryRow(ctx, query, id).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check user exists: %w", err)
    }

    return exists, nil
}

func (r *userRepository) GetRandomActiveTeamUsers(
    ctx context.Context,
    teamName string,
    excludeUserIDs []string,
    maxCount int,
) ([]entity.User, error) {
    qb := r.db.QueryBuilder().
        Select("user_id", "username", "team_name", "is_active").
        From("users").
        Where(squirrel.Eq{
            "team_name": teamName,
            "is_active": true,
        }).
        OrderBy("RANDOM()").
        Limit(uint64(maxCount))

    if len(excludeUserIDs) > 0 {
        qb = qb.Where(squirrel.NotEq{"user_id": excludeUserIDs})
    }

    query, args, err := qb.ToSql()
    if err != nil {
        return nil, fmt.Errorf("build query: %w", err)
    }

    querier := r.db.GetQuerier(ctx)

    rows, err := querier.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("query random active users: %w", err)
    }
    defer rows.Close()

    return scanUsers(rows)
}

func (r *userRepository) CheckUsersAvailableForTeam(
    ctx context.Context,
    userIDs []string,
    teamName string,
) error {
    if len(userIDs) == 0 {
        return nil
    }

    query := `
		SELECT u.user_id
		FROM users u
		WHERE u.user_id = ANY($1)
		  AND u.team_name != $2
		  AND EXISTS (
			SELECT 1
			FROM pull_request_reviewers prr
			JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id
			WHERE prr.reviewer_id = u.user_id
			  AND pr.status = 'OPEN'
		  )
	`

    querier := r.db.GetQuerier(ctx)

    rows, err := querier.Query(ctx, query, userIDs, teamName)
    if err != nil {
        return fmt.Errorf("check users available: %w", err)
    }
    defer rows.Close()

    var blockedUsers []string
    for rows.Next() {
        var userID string
        if err := rows.Scan(&userID); err != nil {
            return fmt.Errorf("scan blocked user: %w", err)
        }
        blockedUsers = append(blockedUsers, userID)
    }

    if len(blockedUsers) > 0 {
        return fmt.Errorf("%w: users %v have active PR assignments in other teams",
            ErrUserHasActiveAssignments, blockedUsers)
    }

    return nil
}

func scanUsers(rows pgx.Rows) ([]entity.User, error) {
    var users []entity.User

    for rows.Next() {
        var user entity.User
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.TeamName,
            &user.IsActive,
        )
        if err != nil {
            return nil, fmt.Errorf("scan user: %w", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }
    return users, nil
}
