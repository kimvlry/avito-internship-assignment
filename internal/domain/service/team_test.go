package service

import (
    "context"
    "testing"

    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestTeamService_CreateTeam(t *testing.T) {
    tests := []struct {
        name            string
        team            *entity.Team
        members         []entity.User
        mockTeamExists  bool
        mockUsersExist  map[string]bool
        expectError     bool
        expectedErrType error
    }{
        {
            name: "успешное создание команды с новыми пользователями",
            team: &entity.Team{Name: "backend"},
            members: []entity.User{
                {ID: "u1", Username: "Alice", IsActive: true},
                {ID: "u2", Username: "Bob", IsActive: true},
            },
            mockTeamExists: false,
            mockUsersExist: map[string]bool{"u1": false, "u2": false},
            expectError:    false,
        },
        {
            name: "ошибка: команда уже существует",
            team: &entity.Team{Name: "backend"},
            members: []entity.User{
                {ID: "u1", Username: "Alice", IsActive: true},
            },
            mockTeamExists:  true,
            expectError:     true,
            expectedErrType: domain.ErrTeamAlreadyExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()
            createCalls := make(map[string]int)
            updateCalls := make(map[string]int)

            mockTeamRepo := mocks.NewTeamRepository(t)
            mockUserRepo := mocks.NewUserRepository(t)
            mockTx := mocks.NewTransactor(t)

            mockTx.On(
                "WithinTransaction",
                mock.Anything,
                mock.AnythingOfType("func(context.Context) error"),
            ).Return(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
                return fn(ctx)
            })

            mockTeamRepo.On("Exists", mock.Anything, tt.team.Name).Return(tt.mockTeamExists, nil)
            if !tt.mockTeamExists {
                mockTeamRepo.On("Create", mock.Anything, tt.team).Return(nil)

                for _, user := range tt.members {
                    if exists, ok := tt.mockUsersExist[user.ID]; ok && exists {
                        mockUserRepo.On("Exists", mock.Anything, user.ID).Return(true, nil)
                        mockUserRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u.ID == user.ID })).Return(nil).Run(func(args mock.Arguments) {
                            updateCalls[user.ID]++
                        })
                    } else {
                        mockUserRepo.On("Exists", mock.Anything, user.ID).Return(false, nil)
                        mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u.ID == user.ID })).Return(nil).Run(func(args mock.Arguments) {
                            createCalls[user.ID]++
                        })
                    }
                }
                mockUserRepo.On("GetByTeam", mock.Anything, tt.team.Name).Return([]entity.User{}, nil)

                userIDs := make([]string, len(tt.members))
                for i, u := range tt.members {
                    userIDs[i] = u.ID
                }
                mockUserRepo.On("CheckUsersAvailableForTeam", mock.Anything, userIDs, tt.team.Name).Return(nil)
            }

            svc := NewTeam(mockTeamRepo, mockUserRepo, mockTx)
            team, err := svc.CreateTeam(ctx, tt.team, tt.members)

            if tt.expectError {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.expectedErrType)
                return
            }

            require.NoError(t, err)
            require.NotNil(t, team)
            assert.Equal(t, tt.team.Name, team.Name)

            for userID, shouldExist := range tt.mockUsersExist {
                if shouldExist {
                    assert.Equal(t, 1, updateCalls[userID], "user %s должен быть обновлен", userID)
                    assert.Equal(t, 0, createCalls[userID], "user %s не должен быть создан", userID)
                } else {
                    assert.Equal(t, 1, createCalls[userID], "user %s должен быть создан", userID)
                    assert.Equal(t, 0, updateCalls[userID], "user %s не должен быть обновлен", userID)
                }
            }
        })
    }
}

func TestTeamService_GetTeamWithMembers(t *testing.T) {
    tests := []struct {
        name            string
        teamName        string
        mockTeam        *entity.Team
        mockMembers     []entity.User
        mockError       error
        expectError     bool
        expectedErrType error
    }{
        {
            name:     "успешное получение команды с участниками",
            teamName: "backend",
            mockTeam: &entity.Team{Name: "backend"},
            mockMembers: []entity.User{
                {ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
                {ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
            },
            expectError: false,
        },
        {
            name:            "ошибка: команда не найдена",
            teamName:        "nonexistent",
            mockTeam:        nil,
            mockError:       domain.ErrTeamNotFound,
            expectError:     true,
            expectedErrType: domain.ErrTeamNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            mockTeamRepo := mocks.NewTeamRepository(t)
            mockUserRepo := mocks.NewUserRepository(t)
            mockTx := mocks.NewTransactor(t)

            if tt.mockTeam != nil {
                mockTeamRepo.On("GetByName", mock.Anything, tt.teamName).Return(tt.mockTeam, nil)
                mockUserRepo.On("GetByTeam", mock.Anything, tt.teamName).Return(tt.mockMembers, nil)
            } else {
                mockTeamRepo.On("GetByName", mock.Anything, tt.teamName).Return(nil, tt.mockError)
            }

            svc := NewTeam(mockTeamRepo, mockUserRepo, mockTx)
            team, members, err := svc.GetTeamWithMembers(ctx, tt.teamName)

            if tt.expectError {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.expectedErrType)
                return
            }

            require.NoError(t, err)
            require.NotNil(t, team)
            assert.Equal(t, tt.teamName, team.Name)
            assert.Len(t, members, len(tt.mockMembers))
        })
    }
}
