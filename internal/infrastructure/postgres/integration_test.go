package postgres_test

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "os"
    "testing"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    testcontainerPg "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"

    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/infrastructure/postgres"
)

type TestDB struct {
    DB        *postgres.DB
    ConnStr   string
    Container testcontainers.Container
}

type NoopLogger struct{}

func (n *NoopLogger) Printf(format string, v ...interface{}) {}

func SetupTestDB(t *testing.T) *TestDB {
    ctx := context.Background()

    opts := []testcontainers.ContainerCustomizer{
        testcontainerPg.WithDatabase("integration_test"),
        testcontainerPg.WithUsername("integration_test"),
        testcontainerPg.WithPassword("integration_test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2),
        ),
    }

    if os.Getenv("INTEGRATION_LOGS") != "1" {
        opts = append(opts, testcontainers.WithLogger(&NoopLogger{}))
    }

    container, err := testcontainerPg.Run(ctx, "postgres:15-alpine", opts...)
    require.NoError(t, err)

    connStr, err := container.ConnectionString(ctx, "sslmode=disable", "connect_timeout=5")
    require.NoError(t, err)

    db, err := postgres.NewDB(ctx, connStr)
    require.NoError(t, err)

    m, err := migrate.New("file://../../../migrations", connStr)
    require.NoError(t, err)

    err = m.Up()
    require.NoError(t, err)

    return &TestDB{
        DB:        db,
        ConnStr:   connStr,
        Container: container,
    }
}

func (tdb *TestDB) Teardown() {
    ctx := context.Background()
    if tdb.DB != nil {
        tdb.DB.Close()
    }
    if tdb.Container != nil {
        _ = tdb.Container.Terminate(ctx)
    }
}

func (tdb *TestDB) CleanDatabase(t *testing.T) {
    ctx := context.Background()
    query := `
        TRUNCATE TABLE 
            pull_request_reviewers, 
            pull_requests, 
            users, 
            teams 
        CASCADE
    `
    _, err := tdb.DB.Exec(ctx, query)
    require.NoError(t, err)
}

func TestRepositories(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Teardown()
    testDB.CleanDatabase(t)

    teamRepo := postgres.NewTeamRepository(testDB.DB)
    userRepo := postgres.NewUserRepository(testDB.DB)
    prRepo := postgres.NewPullRequestRepository(testDB.DB)
    transactor := postgres.NewTransactor(testDB.DB.Pool)

    ctx := context.Background()

    t.Run("TeamRepository", func(t *testing.T) {
        testDB.CleanDatabase(t)

        team := &entity.Team{Name: "backend"}
        err := teamRepo.Create(ctx, team)
        require.NoError(t, err)

        fetched, err := teamRepo.GetByName(ctx, "backend")
        require.NoError(t, err)
        assert.Equal(t, "backend", fetched.Name)

        exists, err := teamRepo.Exists(ctx, "backend")
        require.NoError(t, err)
        assert.True(t, exists)

        err = teamRepo.Create(ctx, team)
        assert.ErrorIs(t, err, domain.ErrTeamAlreadyExists)

        _, err = teamRepo.GetByName(ctx, "nonexistent")
        assert.ErrorIs(t, err, domain.ErrTeamNotFound)
    })

    t.Run("UserRepository", func(t *testing.T) {
        testDB.CleanDatabase(t)

        team := &entity.Team{Name: "team1"}
        err := teamRepo.Create(ctx, team)
        require.NoError(t, err)

        user := &entity.User{
            ID:       "user1",
            Username: "alice",
            TeamName: "team1",
            IsActive: true,
        }
        err = userRepo.Create(ctx, user)
        require.NoError(t, err)

        fetched, err := userRepo.GetByID(ctx, "user1")
        require.NoError(t, err)
        assert.Equal(t, "alice", fetched.Username)
        assert.Equal(t, "team1", fetched.TeamName)
        assert.True(t, fetched.IsActive)

        exists, err := userRepo.Exists(ctx, "user1")
        require.NoError(t, err)
        assert.True(t, exists)

        user.Username = "alice_updated"
        err = userRepo.Update(ctx, user)
        require.NoError(t, err)

        updated, err := userRepo.GetByID(ctx, "user1")
        require.NoError(t, err)
        assert.Equal(t, "alice_updated", updated.Username)

        users, err := userRepo.GetByTeam(ctx, "team1")
        require.NoError(t, err)
        assert.Len(t, users, 1)
        assert.Equal(t, "user1", users[0].ID)
    })

    t.Run("PullRequestRepository", func(t *testing.T) {
        testDB.CleanDatabase(t)

        team := &entity.Team{Name: "dev-team"}
        err := teamRepo.Create(ctx, team)
        require.NoError(t, err)

        author := &entity.User{ID: "author1", Username: "author", TeamName: "dev-team", IsActive: true}
        err = userRepo.Create(ctx, author)
        require.NoError(t, err)

        reviewer := &entity.User{ID: "reviewer1", Username: "reviewer", TeamName: "dev-team", IsActive: true}
        err = userRepo.Create(ctx, reviewer)
        require.NoError(t, err)

        pr := &entity.PullRequest{
            ID:                "pr1",
            Name:              "Add feature",
            AuthorID:          "author1",
            Status:            entity.PROpen,
            AssignedReviewers: []string{"reviewer1"},
        }
        err = prRepo.CreateWithReviewers(ctx, pr)
        require.NoError(t, err)

        fetched, err := prRepo.GetByID(ctx, "pr1")
        require.NoError(t, err)
        assert.Equal(t, "Add feature", fetched.Name)
        assert.Equal(t, "author1", fetched.AuthorID)
        assert.Equal(t, entity.PROpen, fetched.Status)
        assert.Equal(t, []string{"reviewer1"}, fetched.AssignedReviewers)

        exists, err := prRepo.Exists(ctx, "pr1")
        require.NoError(t, err)
        assert.True(t, exists)

        err = prRepo.UpdateStatus(ctx, "pr1", entity.PRMerged)
        require.NoError(t, err)

        merged, err := prRepo.GetByID(ctx, "pr1")
        require.NoError(t, err)
        assert.Equal(t, entity.PRMerged, merged.Status)
        assert.NotNil(t, merged.MergedAt)

        prs, err := prRepo.GetByReviewer(ctx, "reviewer1")
        require.NoError(t, err)
        assert.Len(t, prs, 1)
        assert.Equal(t, "pr1", prs[0].ID)
    })

    t.Run("Transactor", func(t *testing.T) {
        testDB.CleanDatabase(t)

        err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
            team := &entity.Team{Name: "tx-team"}
            err := teamRepo.Create(ctx, team)
            if err != nil {
                return err
            }

            user := &entity.User{ID: "tx-user", Username: "tx-user", TeamName: "tx-team"}
            return userRepo.Create(ctx, user)
        })
        require.NoError(t, err)

        team, err := teamRepo.GetByName(ctx, "tx-team")
        require.NoError(t, err)
        assert.Equal(t, "tx-team", team.Name)

        user, err := userRepo.GetByID(ctx, "tx-user")
        require.NoError(t, err)
        assert.Equal(t, "tx-user", user.Username)
    })
}
