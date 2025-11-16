package handler

import (
    "context"
    "errors"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/handler/check"

    "github.com/kimvlry/avito-internship-assignment/api"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/constructor"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"
)

type teamHandler struct {
    svc *service.Team
}

func newTeamHandler(s *service.Team) *teamHandler {
    return &teamHandler{svc: s}
}

func (h *teamHandler) PostTeamAdd(
    ctx context.Context,
    req api.PostTeamAddRequestObject,
) (api.PostTeamAddResponseObject, error) {

    if err := check.ValidTeamCreate(req); err != nil {
        return api.PostTeamAdd400JSONResponse{
            Error: constructor.ErrorResponse(api.BADREQUEST, err.Error()),
        }, nil
    }

    for _, m := range req.Body.Members {
        if err := check.ValidUserID(m.UserId); err != nil {
            return api.PostTeamAdd400JSONResponse{
                Error: constructor.ErrorResponse(api.BADREQUEST, err.Error()),
            }, nil
        }
    }

    team := &entity.Team{Name: req.Body.TeamName}
    members := make([]entity.User, 0, len(req.Body.Members))
    for _, m := range req.Body.Members {
        members = append(members, entity.User{
            ID:       m.UserId,
            Username: m.Username,
            IsActive: m.IsActive,
        })
    }

    createdTeam, err := h.svc.CreateTeam(ctx, team, members)
    if err != nil {
        if errors.Is(err, domain.ErrTeamAlreadyExists) {
            return api.PostTeamAdd400JSONResponse{
                Error: constructor.ErrorResponse(api.TEAMEXISTS, err.Error()),
            }, nil
        }
        return api.PostTeamAdd500JSONResponse{
            Error: constructor.ErrorResponse("INTERNAL_SERVER_ERROR", err.Error()),
        }, nil
    }

    responseMembers := make([]api.TeamMember, 0, len(createdTeam.Members))
    for _, m := range createdTeam.Members {
        responseMembers = append(responseMembers, api.TeamMember{
            UserId:   m.ID,
            Username: m.Username,
            IsActive: m.IsActive,
        })
    }

    return api.PostTeamAdd201JSONResponse{
        Team: &api.Team{
            TeamName: createdTeam.Name,
            Members:  responseMembers,
        },
    }, nil
}

func (h *teamHandler) GetTeamGet(
    ctx context.Context,
    req api.GetTeamGetRequestObject,
) (api.GetTeamGetResponseObject, error) {

    team, users, err := h.svc.GetTeamWithMembers(ctx, req.Params.TeamName)
    if err != nil {
        if errors.Is(err, domain.ErrTeamNotFound) {
            return api.GetTeamGet404JSONResponse{
                Error: constructor.ErrorResponse(api.NOTFOUND, err.Error()),
            }, nil
        }
        return api.GetTeamGet500JSONResponse{
            Error: constructor.ErrorResponse("INTERNAL_SERVER_ERROR", err.Error()),
        }, nil
    }

    apiMembers := make([]api.TeamMember, 0, len(users))
    for _, m := range users {
        apiMembers = append(apiMembers, api.TeamMember{
            UserId:   m.ID,
            Username: m.Username,
            IsActive: m.IsActive,
        })
    }

    return api.GetTeamGet200JSONResponse{
        TeamName: team.Name,
        Members:  apiMembers,
    }, nil
}
