package handler

import (
    "context"

    "github.com/kimvlry/avito-internship-assignment/api"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/constructor"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/handler/check"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"
)

type userHandler struct {
    svc service.User
}

func NewUserHandler(service service.User) UserHandler {
    return &userHandler{svc: service}
}

func (h *userHandler) PostUsersSetIsActive(
    ctx context.Context,
    req api.PostUsersSetIsActiveRequestObject,
) (api.PostUsersSetIsActiveResponseObject, error) {

    if !check.IsAdmin(ctx) {
        return api.PostUsersSetIsActive401JSONResponse{
            Error: constructor.ErrorResponse("UNAUTHORIZED", "admin access required"),
        }, nil
    }

    user, err := h.svc.SetIsActive(ctx, req.Body.UserId, req.Body.IsActive)
    if err != nil {
        return api.PostUsersSetIsActive404JSONResponse{
            Error: constructor.ErrorResponse(api.NOTFOUND, err.Error()),
        }, nil
    }

    return api.PostUsersSetIsActive200JSONResponse{
        User: &api.User{
            UserId:   user.ID,
            Username: user.Username,
            TeamName: user.TeamName,
            IsActive: user.IsActive,
        },
    }, nil
}

func (h *userHandler) GetUsersGetReview(
    ctx context.Context,
    req api.GetUsersGetReviewRequestObject,
) (api.GetUsersGetReviewResponseObject, error) {

    if !check.IsAdminOrOwner(ctx, req.Params.UserId) {
        return api.GetUsersGetReview200JSONResponse{
            UserId:       req.Params.UserId,
            PullRequests: []api.PullRequestShort{},
        }, nil
    }

    reviews, err := h.svc.GetReviewAssignments(ctx, req.Params.UserId)
    if err != nil {
        return api.GetUsersGetReview200JSONResponse{
            UserId:       req.Params.UserId,
            PullRequests: []api.PullRequestShort{},
        }, nil
    }

    prs := make([]api.PullRequestShort, 0, len(reviews))
    for _, r := range reviews {
        prs = append(prs, api.PullRequestShort{
            PullRequestId:   r.ID,
            PullRequestName: r.Name,
            AuthorId:        r.AuthorID,
            Status:          api.PullRequestShortStatus(r.Status),
        })
    }

    return api.GetUsersGetReview200JSONResponse{
        UserId:       req.Params.UserId,
        PullRequests: prs,
    }, nil
}
