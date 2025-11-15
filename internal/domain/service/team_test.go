package service

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
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
            team: &entity.Team{
                Name: "backend",
            },
            members: []entity.User{
                {ID: "u1", Username: "Alice", IsActive: true},
                {ID: "u2", Username: "Bob", IsActive: true},
            },
            mockTeamExists: false,
            mockUsersExist: map[string]bool{
                "u1": false,
                "u2": false,
            },
            expectError: false,
        },
        {
            name: "создание с обновлением существующих пользователей",
            team: &entity.Team{
                Name: "frontend",
            },
            members: []entity.User{
                {ID: "u1", Username: "Alice", IsActive: true},
                {ID: "u3", Username: "Charlie", IsActive: true},
            },
            mockTeamExists: false,
            mockUsersExist: map[string]bool{
                "u1": true,
                "u3": false,
            },
            expectError: false,
        },
        {
            name: "ошибка: команда уже существует",
            team: &entity.Team{
                Name: "backend",
            },
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

            mockTeamRepo := &MockTeamRepository{
                ExistsFunc: func(ctx context.Context, teamName string) (bool, error) {
                    return tt.mockTeamExists, nil
                },
                CreateFunc: func(ctx context.Context, team *entity.Team) error {
                    assert.Equal(t, tt.team.Name, team.Name)
                    return nil
                },
            }

            mockUserRepo := &MockUserRepository{
                CheckUsersAvailableForTeamFunc: func(ctx context.Context, userIDs []string, teamName string) error {
                    return nil
                },
                ExistsFunc: func(ctx context.Context, userID string) (bool, error) {
                    exists, ok := tt.mockUsersExist[userID]
                    if !ok {
                        return false, nil
                    }
                    return exists, nil
                },
                CreateFunc: func(ctx context.Context, user *entity.User) error {
                    createCalls[user.ID]++
                    assert.Equal(t, tt.team.Name, user.TeamName, "user должен быть добавлен в команду")
                    return nil
                },
                UpdateFunc: func(ctx context.Context, user *entity.User) error {
                    updateCalls[user.ID]++
                    assert.Equal(t, tt.team.Name, user.TeamName, "user должен быть обновлен с новой командой")
                    return nil
                },
            }

            mockTx := &MockTransactor{}

            svc := NewTeamService(mockTeamRepo, mockUserRepo, mockTx)

            team, err := svc.CreateTeam(ctx, tt.team, tt.members)

            if tt.expectError {
                require.Error(t, err)
                if tt.expectedErrType != nil {
                    assert.ErrorIs(t, err, tt.expectedErrType)
                }
                return
            }

            require.NoError(t, err)
            require.NotNil(t, team)
            assert.Equal(t, tt.team.Name, team.Name)

            for userID, shouldExist := range tt.mockUsersExist {
                if shouldExist {
                    assert.Equal(t, 1, updateCalls[userID], "user %s должен быть updated", userID)
                    assert.Equal(t, 0, createCalls[userID], "user %s не должен быть created", userID)
                } else {
                    assert.Equal(t, 1, createCalls[userID], "user %s должен быть created", userID)
                    assert.Equal(t, 0, updateCalls[userID], "user %s не должен быть updated", userID)
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
            mockTeam: &entity.Team{
                Name: "backend",
            },
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

            mockTeamRepo := &MockTeamRepository{
                GetByNameFunc: func(ctx context.Context, teamName string) (*entity.Team, error) {
                    if tt.mockTeam == nil {
                        return nil, tt.mockError
                    }
                    return tt.mockTeam, nil
                },
            }

            mockUserRepo := &MockUserRepository{
                GetByTeamFunc: func(ctx context.Context, teamName string) ([]entity.User, error) {
                    return tt.mockMembers, nil
                },
            }

            mockTx := &MockTransactor{}

            svc := NewTeamService(mockTeamRepo, mockUserRepo, mockTx)

            team, members, err := svc.GetTeamWithMembers(ctx, tt.teamName)

            if tt.expectError {
                require.Error(t, err)
                if tt.expectedErrType != nil {
                    assert.ErrorIs(t, err, tt.expectedErrType)
                }
                return
            }

            require.NoError(t, err)
            require.NotNil(t, team)
            assert.Equal(t, tt.teamName, team.Name)
            assert.Len(t, members, len(tt.mockMembers))
        })
    }
}
