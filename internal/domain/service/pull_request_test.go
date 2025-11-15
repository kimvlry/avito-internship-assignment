package service

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
)

func TestPullRequestService_CreatePullRequestWithReviewers(t *testing.T) {
    tests := []struct {
        name              string
        prID              string
        prName            string
        authorID          string
        mockPRExists      bool
        mockAuthor        *entity.User
        mockReviewers     []entity.User
        mockError         error
        expectedReviewers int
        expectError       bool
        expectedErrType   error
    }{
        {
            name:         "успешное создание с 2 ревьюверами",
            prID:         "pr-1",
            prName:       "Feature X",
            authorID:     "u1",
            mockPRExists: false,
            mockAuthor: &entity.User{
                ID:       "u1",
                Username: "Alice",
                TeamName: "backend",
                IsActive: true,
            },
            mockReviewers: []entity.User{
                {ID: "u2", Username: "Bob", IsActive: true},
                {ID: "u3", Username: "Charlie", IsActive: true},
            },
            expectedReviewers: 2,
            expectError:       false,
        },
        {
            name:         "создание с 1 ревьювером (доступен только 1)",
            prID:         "pr-2",
            prName:       "Feature Y",
            authorID:     "u1",
            mockPRExists: false,
            mockAuthor: &entity.User{
                ID:       "u1",
                Username: "Alice",
                TeamName: "backend",
                IsActive: true,
            },
            mockReviewers: []entity.User{
                {ID: "u2", Username: "Bob", IsActive: true},
            },
            expectedReviewers: 1,
            expectError:       false,
        },
        {
            name:         "создание без ревьюверов (нет доступных)",
            prID:         "pr-3",
            prName:       "Feature Z",
            authorID:     "u1",
            mockPRExists: false,
            mockAuthor: &entity.User{
                ID:       "u1",
                Username: "Alice",
                TeamName: "backend",
                IsActive: true,
            },
            mockReviewers:     []entity.User{},
            expectedReviewers: 0,
            expectError:       false,
        },
        {
            name:            "ошибка: PR уже существует",
            prID:            "pr-4",
            prName:          "Feature W",
            authorID:        "u1",
            mockPRExists:    true,
            expectError:     true,
            expectedErrType: domain.ErrPullRequestAlreadyExists,
        },
        {
            name:            "ошибка: автор не найден",
            prID:            "pr-5",
            prName:          "Feature V",
            authorID:        "unknown",
            mockPRExists:    false,
            mockAuthor:      nil,
            mockError:       domain.ErrUserNotFound,
            expectError:     true,
            expectedErrType: domain.ErrUserNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            mockPRRepo := &MockPullRequestRepository{
                ExistsFunc: func(ctx context.Context, prID string) (bool, error) {
                    return tt.mockPRExists, nil
                },
                CreateWithReviewersFunc: func(ctx context.Context, pr *entity.PullRequest) error {
                    return nil
                },
            }

            mockUserRepo := &MockUserRepository{
                GetByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
                    if tt.mockAuthor == nil {
                        return nil, tt.mockError
                    }
                    return tt.mockAuthor, nil
                },
                GetRandomActiveTeamUsersFunc: func(ctx context.Context, teamName string, excludeIDs []string, maxCount int) ([]entity.User, error) {
                    assert.Contains(t, excludeIDs, tt.authorID)
                    assert.Equal(t, 2, maxCount)
                    return tt.mockReviewers, nil
                },
            }

            mockTx := &MockTransactor{}
            svc := NewPullRequestService(mockPRRepo, mockUserRepo, mockTx)

            pr, err := svc.CreatePullRequestWithReviewers(ctx, tt.prID, tt.prName, tt.authorID)

            if tt.expectError {
                require.Error(t, err)
                if tt.expectedErrType != nil {
                    assert.ErrorIs(t, err, tt.expectedErrType)
                }
                return
            }

            require.NoError(t, err)
            require.NotNil(t, pr)
            assert.Equal(t, tt.prID, pr.ID)
            assert.Equal(t, tt.prName, pr.Name)
            assert.Equal(t, tt.authorID, pr.AuthorID)
            assert.Equal(t, entity.PROpen, pr.Status)
            assert.Len(t, pr.AssignedReviewers, tt.expectedReviewers)

            for _, reviewerID := range pr.AssignedReviewers {
                assert.NotEqual(t, tt.authorID, reviewerID, "автор не должен быть ревьювером")
            }
        })
    }
}

func TestPullRequestService_Merge(t *testing.T) {
    tests := []struct {
        name          string
        prID          string
        mockPR        *entity.PullRequest
        mockError     error
        expectError   bool
        expectedCalls int
    }{
        {
            name: "успешный merge открытого PR",
            prID: "pr-1",
            mockPR: &entity.PullRequest{
                ID:        "pr-1",
                Name:      "Feature X",
                AuthorID:  "u1",
                Status:    entity.PROpen,
                CreatedAt: time.Now(),
            },
            expectError:   false,
            expectedCalls: 1,
        },
        {
            name: "идемпотентность: повторный merge уже merged PR",
            prID: "pr-2",
            mockPR: &entity.PullRequest{
                ID:        "pr-2",
                Name:      "Feature Y",
                AuthorID:  "u1",
                Status:    entity.PRMerged,
                CreatedAt: time.Now(),
                MergedAt:  func() *time.Time { t := time.Now(); return &t }(),
            },
            expectError:   false,
            expectedCalls: 0,
        },
        {
            name:        "ошибка: PR не найден",
            prID:        "pr-999",
            mockPR:      nil,
            mockError:   domain.ErrPullRequestNotFound,
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            updateStatusCalls := 0

            mockPRRepo := &MockPullRequestRepository{
                GetByIDFunc: func(ctx context.Context, prID string) (*entity.PullRequest, error) {
                    if tt.mockPR == nil {
                        return nil, tt.mockError
                    }
                    return tt.mockPR, nil
                },
                UpdateStatusFunc: func(ctx context.Context, prID string, status entity.PullRequestStatus) error {
                    updateStatusCalls++
                    assert.Equal(t, entity.PRMerged, status)
                    return nil
                },
            }

            mockUserRepo := &MockUserRepository{}
            mockTx := &MockTransactor{}

            svc := NewPullRequestService(mockPRRepo, mockUserRepo, mockTx)

            pr, err := svc.Merge(ctx, tt.prID)

            if tt.expectError {
                require.Error(t, err)
                if tt.mockError != nil {
                    assert.ErrorIs(t, err, tt.mockError)
                }
                return
            }

            require.NoError(t, err)
            require.NotNil(t, pr)
            assert.Equal(t, entity.PRMerged, pr.Status)
            assert.NotNil(t, pr.MergedAt)

            assert.Equal(t, tt.expectedCalls, updateStatusCalls, "неправильное количество вызовов UpdateStatus")
        })
    }
}

func TestPullRequestService_ReassignReviewer(t *testing.T) {
    tests := []struct {
        name             string
        prID             string
        oldUserID        string
        mockPR           *entity.PullRequest
        mockOldUser      *entity.User
        mockReplacements []entity.User
        expectError      bool
        expectedErrType  error
    }{
        {
            name:      "успешная замена ревьювера",
            prID:      "pr-1",
            oldUserID: "u2",
            mockPR: &entity.PullRequest{
                ID:                "pr-1",
                AuthorID:          "u1",
                Status:            entity.PROpen,
                AssignedReviewers: []string{"u2", "u3"},
            },
            mockOldUser: &entity.User{
                ID:       "u2",
                TeamName: "backend",
                IsActive: true,
            },
            mockReplacements: []entity.User{
                {ID: "u4", Username: "Dave", IsActive: true},
            },
            expectError: false,
        },
        {
            name:      "ошибка: PR уже merged",
            prID:      "pr-2",
            oldUserID: "u2",
            mockPR: &entity.PullRequest{
                ID:                "pr-2",
                AuthorID:          "u1",
                Status:            entity.PRMerged,
                AssignedReviewers: []string{"u2", "u3"},
            },
            expectError:     true,
            expectedErrType: domain.ErrPullRequestIsMerged,
        },
        {
            name:      "ошибка: ревьювер не назначен",
            prID:      "pr-3",
            oldUserID: "u999",
            mockPR: &entity.PullRequest{
                ID:                "pr-3",
                AuthorID:          "u1",
                Status:            entity.PROpen,
                AssignedReviewers: []string{"u2", "u3"},
            },
            expectError:     true,
            expectedErrType: domain.ErrReviewerNotAssigned,
        },
        {
            name:      "ошибка: нет доступных кандидатов",
            prID:      "pr-4",
            oldUserID: "u2",
            mockPR: &entity.PullRequest{
                ID:                "pr-4",
                AuthorID:          "u1",
                Status:            entity.PROpen,
                AssignedReviewers: []string{"u2", "u3"},
            },
            mockOldUser: &entity.User{
                ID:       "u2",
                TeamName: "backend",
                IsActive: true,
            },
            mockReplacements: []entity.User{},
            expectError:      true,
            expectedErrType:  domain.ErrNoReviewerCandidate,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            mockPRRepo := &MockPullRequestRepository{
                GetByIDFunc: func(ctx context.Context, prID string) (*entity.PullRequest, error) {
                    return tt.mockPR, nil
                },
                ReplaceReviewerFunc: func(ctx context.Context, prID, oldUserID, newUserID string) error {
                    assert.Equal(t, tt.oldUserID, oldUserID)
                    assert.NotEmpty(t, newUserID)
                    return nil
                },
            }

            mockUserRepo := &MockUserRepository{
                GetByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
                    if tt.mockOldUser != nil && id == tt.oldUserID {
                        return tt.mockOldUser, nil
                    }
                    return nil, domain.ErrUserNotFound
                },
                GetRandomActiveTeamUsersFunc: func(ctx context.Context, teamName string, excludeIDs []string, maxCount int) ([]entity.User, error) {
                    assert.Contains(t, excludeIDs, tt.mockPR.AuthorID, "автор должен быть исключен")
                    for _, reviewerID := range tt.mockPR.AssignedReviewers {
                        assert.Contains(t, excludeIDs, reviewerID, "текущий ревьювер должен быть исключен")
                    }
                    assert.Equal(t, 1, maxCount)
                    return tt.mockReplacements, nil
                },
            }

            mockTx := &MockTransactor{}

            svc := NewPullRequestService(mockPRRepo, mockUserRepo, mockTx)

            pr, newUserID, err := svc.ReassignReviewer(ctx, tt.prID, tt.oldUserID)

            if tt.expectError {
                require.Error(t, err)
                if tt.expectedErrType != nil {
                    assert.ErrorIs(t, err, tt.expectedErrType)
                }
                return
            }

            require.NoError(t, err)
            require.NotNil(t, pr)
            assert.NotEmpty(t, newUserID)
            assert.Equal(t, tt.mockReplacements[0].ID, newUserID)
        })
    }
}
