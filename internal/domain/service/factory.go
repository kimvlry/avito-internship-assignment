package service

import "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"

type Services struct {
    TeamService        *TeamService
    UserService        *UserService
    PullRequestService *PullRequestService
    Transactor         Transactor
}

func NewServices(
    teamRepository repository.TeamRepository,
    userRepository repository.UserRepository,
    pullRequestRepository repository.PullRequestRepository,
    tx Transactor,
) *Services {
    return &Services{
        TeamService:        NewTeamService(teamRepository, userRepository, tx),
        UserService:        NewUserService(userRepository, pullRequestRepository),
        PullRequestService: NewPullRequestService(pullRequestRepository, userRepository, tx),
    }
}
