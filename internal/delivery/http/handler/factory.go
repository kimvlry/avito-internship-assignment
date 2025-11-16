package handler

import (
    "github.com/kimvlry/avito-internship-assignment/api"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"
)

// Handlers implements generated api.StrictServerInterface
type Handlers struct {
    *pullRequestHandler
    *teamHandler
    *userHandler
}

var _ api.StrictServerInterface = (*Handlers)(nil)

func NewHandlers(services *service.Services) *Handlers {
    return &Handlers{
        newPullRequestHandler(services.PullRequestService),
        newTeamHandler(services.TeamService),
        newUserHandler(services.UserService),
    }
}
