package service

import "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"

type Services struct {
    TeamService        *Team
    UserService        *User
    PullRequestService *PullRequest
    StatsService       *StatsService
    Transactor         repository.Transactor
}

func NewServices(
    teamRepository repository.TeamRepository,
    userRepository repository.UserRepository,
    pullRequestRepository repository.PullRequestRepository,
    tx repository.Transactor,
) *Services {
    return &Services{
        TeamService:        NewTeam(teamRepository, userRepository, tx),
        UserService:        NewUser(userRepository, pullRequestRepository),
        PullRequestService: NewPullRequest(pullRequestRepository, userRepository, tx),
        StatsService:       NewStatsService(pullRequestRepository),
    }
}
