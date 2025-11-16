package service

import (
    "context"
    "testing"

    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
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
        {"успешная активация пользователя", "u1", true, nil, false, nil},
        {"успешная деактивация пользователя", "u2", false, nil, false, nil},
        {"ошибка: пользователь не найден", "u999", true, domain.ErrUserNotFound, true, domain.ErrUserNotFound},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            mockUserRepo := mocks.NewUserRepository(t)
            mockPRRepo := mocks.NewPullRequestRepository(t)

            mockUserRepo.On("SetIsActive", ctx, tt.userID, tt.isActive).
                Return(nil, tt.mockError)

            svc := NewUser(mockUserRepo, mockPRRepo)

            _, err := svc.SetIsActive(ctx, tt.userID, tt.isActive)

            if tt.expectError {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.expectedErrType)
                return
            }

            require.NoError(t, err)
            mockUserRepo.AssertCalled(t, "SetIsActive", ctx, tt.userID, tt.isActive)
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
            "успешное получение назначений пользователя",
            "u1",
            []*entity.PullRequest{
                {ID: "pr-1", Name: "Feature X", AuthorID: "u2", Status: entity.PROpen, AssignedReviewers: []string{"u1", "u3"}},
                {ID: "pr-2", Name: "Feature Y", AuthorID: "u3", Status: entity.PROpen, AssignedReviewers: []string{"u1"}},
            },
            nil, false, nil,
        },
        {
            "пустой список для пользователя без назначений",
            "u2", []*entity.PullRequest{}, nil, false, nil,
        },
        {
            "получение назначений включая merged PR (администратор видит)",
            "u1",
            []*entity.PullRequest{
                {ID: "pr-1", Name: "Feature X", AuthorID: "u2", Status: entity.PROpen, AssignedReviewers: []string{"u1"}},
                {ID: "pr-2", Name: "Feature Y", AuthorID: "u3", Status: entity.PRMerged, AssignedReviewers: []string{"u1"}},
            },
            nil, false, nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            mockUserRepo := mocks.NewUserRepository(t)
            mockPRRepo := mocks.NewPullRequestRepository(t)

            mockPRRepo.On("GetByReviewer", ctx, tt.userID).Return(tt.mockPRs, tt.mockError)

            svc := NewUser(mockUserRepo, mockPRRepo)

            prs, err := svc.GetReviewAssignments(ctx, tt.userID)

            if tt.expectError {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.expectedErrType)
                return
            }

            require.NoError(t, err)
            mockPRRepo.AssertCalled(t, "GetByReviewer", ctx, tt.userID)
            assert.Len(t, prs, len(tt.mockPRs))

            for _, pr := range prs {
                assert.Contains(t, pr.AssignedReviewers, tt.userID,
                    "PR %s должен содержать пользователя %s в ревьюверах", pr.ID, tt.userID)
            }
        })
    }
}
