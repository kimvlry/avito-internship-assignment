package service

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
)

type MockUserRepository struct {
    GetByIDFunc                    func(ctx context.Context, id string) (*entity.User, error)
    GetByTeamFunc                  func(ctx context.Context, teamName string) ([]entity.User, error)
    GetRandomActiveTeamUsersFunc   func(ctx context.Context, teamName string, excludeIDs []string, maxCount int) ([]entity.User, error)
    SetIsActiveFunc                func(ctx context.Context, userID string, isActive bool) error
    ExistsFunc                     func(ctx context.Context, userID string) (bool, error)
    CreateFunc                     func(ctx context.Context, user *entity.User) error
    UpdateFunc                     func(ctx context.Context, user *entity.User) error
    CheckUsersAvailableForTeamFunc func(ctx context.Context, userIDs []string, teamName string) error
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(ctx, id)
    }
    return nil, nil
}

func (m *MockUserRepository) GetByTeam(ctx context.Context, teamName string) ([]entity.User, error) {
    if m.GetByTeamFunc != nil {
        return m.GetByTeamFunc(ctx, teamName)
    }
    return nil, nil
}

func (m *MockUserRepository) GetRandomActiveTeamUsers(ctx context.Context, teamName string, excludeIDs []string, maxCount int) ([]entity.User, error) {
    if m.GetRandomActiveTeamUsersFunc != nil {
        return m.GetRandomActiveTeamUsersFunc(ctx, teamName, excludeIDs, maxCount)
    }
    return nil, nil
}

func (m *MockUserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) error {
    if m.SetIsActiveFunc != nil {
        return m.SetIsActiveFunc(ctx, userID, isActive)
    }
    return nil
}

func (m *MockUserRepository) Exists(ctx context.Context, userID string) (bool, error) {
    if m.ExistsFunc != nil {
        return m.ExistsFunc(ctx, userID)
    }
    return false, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(ctx, user)
    }
    return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
    if m.UpdateFunc != nil {
        return m.UpdateFunc(ctx, user)
    }
    return nil
}

func (m *MockUserRepository) CheckUsersAvailableForTeam(ctx context.Context, userIDs []string, teamName string) error {
    if m.CheckUsersAvailableForTeamFunc != nil {
        return m.CheckUsersAvailableForTeamFunc(ctx, userIDs, teamName)
    }
    return nil
}

type MockPullRequestRepository struct {
    ExistsFunc              func(ctx context.Context, prID string) (bool, error)
    GetByIDFunc             func(ctx context.Context, prID string) (*entity.PullRequest, error)
    CreateWithReviewersFunc func(ctx context.Context, pr *entity.PullRequest) error
    ReplaceReviewerFunc     func(ctx context.Context, prID, oldUserID, newUserID string) error
    UpdateStatusFunc        func(ctx context.Context, prID string, status entity.PullRequestStatus) error
    GetByReviewerFunc       func(ctx context.Context, userID string) ([]*entity.PullRequest, error)
}

func (m *MockPullRequestRepository) Exists(ctx context.Context, prID string) (bool, error) {
    if m.ExistsFunc != nil {
        return m.ExistsFunc(ctx, prID)
    }
    return false, nil
}

func (m *MockPullRequestRepository) GetByID(ctx context.Context, prID string) (*entity.PullRequest, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(ctx, prID)
    }
    return nil, nil
}

func (m *MockPullRequestRepository) CreateWithReviewers(ctx context.Context, pr *entity.PullRequest) error {
    if m.CreateWithReviewersFunc != nil {
        return m.CreateWithReviewersFunc(ctx, pr)
    }
    return nil
}

func (m *MockPullRequestRepository) ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error {
    if m.ReplaceReviewerFunc != nil {
        return m.ReplaceReviewerFunc(ctx, prID, oldUserID, newUserID)
    }
    return nil
}

func (m *MockPullRequestRepository) UpdateStatus(ctx context.Context, prID string, status entity.PullRequestStatus) error {
    if m.UpdateStatusFunc != nil {
        return m.UpdateStatusFunc(ctx, prID, status)
    }
    return nil
}

func (m *MockPullRequestRepository) GetByReviewer(ctx context.Context, userID string) ([]*entity.PullRequest, error) {
    if m.GetByReviewerFunc != nil {
        return m.GetByReviewerFunc(ctx, userID)
    }
    return nil, nil
}

type MockTeamRepository struct {
    ExistsFunc    func(ctx context.Context, teamName string) (bool, error)
    GetByNameFunc func(ctx context.Context, teamName string) (*entity.Team, error)
    CreateFunc    func(ctx context.Context, team *entity.Team) error
}

func (m *MockTeamRepository) Exists(ctx context.Context, teamName string) (bool, error) {
    if m.ExistsFunc != nil {
        return m.ExistsFunc(ctx, teamName)
    }
    return false, nil
}

func (m *MockTeamRepository) GetByName(ctx context.Context, teamName string) (*entity.Team, error) {
    if m.GetByNameFunc != nil {
        return m.GetByNameFunc(ctx, teamName)
    }
    return nil, nil
}

func (m *MockTeamRepository) Create(ctx context.Context, team *entity.Team) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(ctx, team)
    }
    return nil
}

type MockTransactor struct{}

func (m *MockTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return fn(ctx)
}
