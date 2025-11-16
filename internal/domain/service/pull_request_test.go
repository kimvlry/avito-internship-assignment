package service

import (
    "context"
    "github.com/stretchr/testify/mock"
    "testing"
    "time"

    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
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
            name:            "ошибка: PR уже существует",
            prID:            "pr-2",
            prName:          "Feature Y",
            authorID:        "u1",
            mockPRExists:    true,
            expectError:     true,
            expectedErrType: domain.ErrPullRequestAlreadyExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            mockPRRepo := mocks.NewPullRequestRepository(t)
            mockUserRepo := mocks.NewUserRepository(t)
            mockTx := mocks.NewTransactor(t)

            mockPRRepo.On("Exists", ctx, tt.prID).Return(tt.mockPRExists, nil)
            if !tt.mockPRExists && tt.mockAuthor != nil {
                mockUserRepo.On("GetByID", ctx, tt.authorID).Return(tt.mockAuthor, nil)
                mockUserRepo.On("GetRandomActiveTeamUsers", ctx, tt.mockAuthor.TeamName, mock.Anything, 2).
                    Return(tt.mockReviewers, nil)
                mockPRRepo.On("CreateWithReviewers", ctx, mock.AnythingOfType("*entity.PullRequest")).Return(nil)

                mockTx.On(
                    "WithinTransaction",
                    mock.Anything,
                    mock.AnythingOfType("func(context.Context) error"),
                ).Return(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
                    return fn(ctx)
                })
            }
            svc := NewPullRequest(mockPRRepo, mockUserRepo, mockTx)

            pr, err := svc.CreatePullRequestWithReviewers(ctx, tt.prID, tt.prName, tt.authorID)

            if tt.expectError {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.expectedErrType)
                return
            }

            require.NoError(t, err)
            require.NotNil(t, pr)
            assert.Equal(t, tt.prID, pr.ID)
            assert.Equal(t, tt.prName, pr.Name)
            assert.Equal(t, tt.authorID, pr.AuthorID)
            assert.Equal(t, entity.PROpen, pr.Status)
            assert.Len(t, pr.AssignedReviewers, tt.expectedReviewers)
        })
    }
}

func TestPullRequestService_Merge(t *testing.T) {
    t.Run("успешный merge", func(t *testing.T) {
        ctx := context.Background()
        pr := &entity.PullRequest{
            ID:        "pr-1",
            Name:      "Feature X",
            AuthorID:  "u1",
            Status:    entity.PROpen,
            CreatedAt: time.Now(),
        }

        mockPRRepo := mocks.NewPullRequestRepository(t)
        mockUserRepo := mocks.NewUserRepository(t)
        mockTx := mocks.NewTransactor(t)

        mockPRRepo.On("GetByID", ctx, pr.ID).Return(pr, nil)
        mockPRRepo.On("UpdateStatus", ctx, pr.ID, entity.PRMerged).Return(nil)

        svc := NewPullRequest(mockPRRepo, mockUserRepo, mockTx)
        gotPr, err := svc.Merge(ctx, pr.ID)
        require.NoError(t, err)
        assert.Equal(t, entity.PRMerged, gotPr.Status)
        assert.NotNil(t, gotPr.MergedAt)
    })
}

func TestPullRequestService_ReassignReviewer(t *testing.T) {
    t.Run("успешная замена ревьювера", func(t *testing.T) {
        ctx := context.Background()
        pr := &entity.PullRequest{
            ID:                "pr-1",
            AuthorID:          "u1",
            Status:            entity.PROpen,
            AssignedReviewers: []string{"u2", "u3"},
        }
        oldUser := &entity.User{
            ID:       "u2",
            TeamName: "backend",
            IsActive: true,
        }
        newReviewer := entity.User{ID: "u4"}

        mockPRRepo := mocks.NewPullRequestRepository(t)
        mockUserRepo := mocks.NewUserRepository(t)
        mockTx := mocks.NewTransactor(t)

        mockPRRepo.On("GetByID", ctx, pr.ID).Return(pr, nil)
        mockPRRepo.On("ReplaceReviewer", ctx, pr.ID, oldUser.ID, newReviewer.ID).Return(nil)
        mockUserRepo.On("GetByID", ctx, oldUser.ID).Return(oldUser, nil)
        mockUserRepo.On("GetRandomActiveTeamUsers", ctx, oldUser.TeamName, mock.Anything, 1).
            Return([]entity.User{newReviewer}, nil)

        mockTx.On(
            "WithinTransaction",
            mock.Anything,
            mock.AnythingOfType("func(context.Context) error"),
        ).Return(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
            return fn(ctx)
        })

        svc := NewPullRequest(mockPRRepo, mockUserRepo, mockTx)
        gotPr, gotNewID, err := svc.ReassignReviewer(ctx, pr.ID, oldUser.ID)
        require.NoError(t, err)
        assert.Equal(t, newReviewer.ID, gotNewID)
        assert.Equal(t, pr, gotPr)
    })
}
