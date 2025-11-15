package service

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
)

func TestUserService_SetIsActive(t *testing.T) {
    tests := []struct {
        name            string
        userID          string
        isActive        bool
        mockError       error
        expectError     bool
        expectedErrType error
    }{
        {
            name:        "успешная активация пользователя",
            userID:      "u1",
            isActive:    true,
            expectError: false,
        },
        {
            name:        "успешная деактивация пользователя",
            userID:      "u2",
            isActive:    false,
            expectError: false,
        },
        {
            name:            "ошибка: пользователь не найден",
            userID:          "u999",
            isActive:        true,
            mockError:       domain.ErrUserNotFound,
            expectError:     true,
            expectedErrType: domain.ErrUserNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            setIsActiveCalled := false
            var capturedIsActive bool

            mockUserRepo := &MockUserRepository{
                SetIsActiveFunc: func(ctx context.Context, userID string, isActive bool) error {
                    setIsActiveCalled = true
                    capturedIsActive = isActive
                    assert.Equal(t, tt.userID, userID)
                    return tt.mockError
                },
            }

            mockPRRepo := &MockPullRequestRepository{}

            svc := NewUserService(mockUserRepo, mockPRRepo)

            err := svc.SetIsActive(ctx, tt.userID, tt.isActive)

            if tt.expectError {
                require.Error(t, err)
                if tt.expectedErrType != nil {
                    assert.ErrorIs(t, err, tt.expectedErrType)
                }
                return
            }

            require.NoError(t, err)
            assert.True(t, setIsActiveCalled, "SetIsActive должен быть вызван")
            assert.Equal(t, tt.isActive, capturedIsActive, "передан правильный is_active")
        })
    }
}

func TestUserService_GetReviewAssignments(t *testing.T) {
    tests := []struct {
        name            string
        userID          string
        mockPRs         []*entity.PullRequest
        mockError       error
        expectError     bool
        expectedErrType error
    }{
        {
            name:   "успешное получение назначений пользователя",
            userID: "u1",
            mockPRs: []*entity.PullRequest{
                {
                    ID:                "pr-1",
                    Name:              "Feature X",
                    AuthorID:          "u2",
                    Status:            entity.PROpen,
                    AssignedReviewers: []string{"u1", "u3"},
                },
                {
                    ID:                "pr-2",
                    Name:              "Feature Y",
                    AuthorID:          "u3",
                    Status:            entity.PROpen,
                    AssignedReviewers: []string{"u1"},
                },
            },
            expectError: false,
        },
        {
            name:        "пустой список для пользователя без назначений",
            userID:      "u2",
            mockPRs:     []*entity.PullRequest{},
            expectError: false,
        },
        {
            name:   "получение назначений включая merged PR (администратор видит)",
            userID: "u1",
            mockPRs: []*entity.PullRequest{
                {
                    ID:                "pr-1",
                    Name:              "Feature X",
                    AuthorID:          "u2",
                    Status:            entity.PROpen,
                    AssignedReviewers: []string{"u1"},
                },
                {
                    ID:                "pr-2",
                    Name:              "Feature Y",
                    AuthorID:          "u3",
                    Status:            entity.PRMerged,
                    AssignedReviewers: []string{"u1"},
                },
            },
            expectError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            getByReviewerCalled := false

            mockUserRepo := &MockUserRepository{}

            mockPRRepo := &MockPullRequestRepository{
                GetByReviewerFunc: func(ctx context.Context, userID string) ([]*entity.PullRequest, error) {
                    getByReviewerCalled = true
                    assert.Equal(t, tt.userID, userID)
                    if tt.mockError != nil {
                        return nil, tt.mockError
                    }
                    return tt.mockPRs, nil
                },
            }

            svc := NewUserService(mockUserRepo, mockPRRepo)

            prs, err := svc.GetReviewAssignments(ctx, tt.userID)

            if tt.expectError {
                require.Error(t, err)
                if tt.expectedErrType != nil {
                    assert.ErrorIs(t, err, tt.expectedErrType)
                }
                return
            }

            require.NoError(t, err)
            assert.True(t, getByReviewerCalled, "GetByReviewer должен быть вызван")
            assert.Len(t, prs, len(tt.mockPRs))

            for _, pr := range prs {
                assert.Contains(t, pr.AssignedReviewers, tt.userID,
                    "PR %s должен содержать пользователя %s в ревьюверах", pr.ID, tt.userID)
            }
        })
    }
}
