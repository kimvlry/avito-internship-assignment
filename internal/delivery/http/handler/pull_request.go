package handler

import (
    "context"
    "errors"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/handler/check"

    "github.com/kimvlry/avito-internship-assignment/api"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/constructor"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"
)

type pullRequestHandler struct {
    svc service.PullRequest
}

func NewPullRequestHandler(svc service.PullRequest) PullRequestHandler {
    return &pullRequestHandler{svc: svc}
}

func (h *pullRequestHandler) PostPullRequestCreate(
    ctx context.Context,
    req api.PostPullRequestCreateRequestObject,
) (api.PostPullRequestCreateResponseObject, error) {
    if !check.IsAdmin(ctx) {
        return api.PostPullRequestCreate404JSONResponse{
            Error: constructor.ErrorResponse(api.NOTFOUND, "resource not found"),
        }, nil
    }

    if err := check.ValidPullRequestCreate(req); err != nil {
        return api.PostPullRequestCreate400JSONResponse{
            Error: constructor.ErrorResponse(api.BADREQUEST, err.Error()),
        }, nil
    }

    pr, err := h.svc.CreatePullRequestWithReviewers(
        ctx,
        req.Body.PullRequestId,
        req.Body.PullRequestName,
        req.Body.AuthorId,
    )
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrTeamNotFound):
            return api.PostPullRequestCreate404JSONResponse{
                Error: constructor.ErrorResponse(api.NOTFOUND, err.Error()),
            }, nil
        case errors.Is(err, domain.ErrPullRequestAlreadyExists):
            return api.PostPullRequestCreate409JSONResponse{
                Error: constructor.ErrorResponse(api.PREXISTS, err.Error()),
            }, nil
        default:
            return api.PostPullRequestCreate500JSONResponse{
                Error: constructor.ErrorResponse("INTERNAL_SERVER_ERROR", err.Error()),
            }, nil
        }
    }

    return api.PostPullRequestCreate201JSONResponse{
        Pr: &api.PullRequest{
            PullRequestId:     pr.ID,
            PullRequestName:   pr.Name,
            AuthorId:          pr.AuthorID,
            Status:            api.PullRequestStatus(pr.Status),
            AssignedReviewers: pr.AssignedReviewers,
            CreatedAt:         &pr.CreatedAt,
            MergedAt:          pr.MergedAt,
        },
    }, nil
}

func (h *pullRequestHandler) PostPullRequestMerge(
    ctx context.Context,
    req api.PostPullRequestMergeRequestObject,
) (api.PostPullRequestMergeResponseObject, error) {
    if !check.IsAdmin(ctx) {
        return api.PostPullRequestMerge404JSONResponse{
            Error: constructor.ErrorResponse(api.NOTFOUND, "resource not found"),
        }, nil
    }

    pr, err := h.svc.Merge(ctx, req.Body.PullRequestId)
    if err != nil {
        if errors.Is(err, domain.ErrPullRequestNotFound) {
            return api.PostPullRequestMerge404JSONResponse{
                Error: constructor.ErrorResponse(api.NOTFOUND, err.Error()),
            }, nil
        }
        return api.PostPullRequestMerge500JSONResponse{
            Error: constructor.ErrorResponse("INTERNAL_SERVER_ERROR", err.Error()),
        }, nil
    }

    return api.PostPullRequestMerge200JSONResponse{
        Pr: &api.PullRequest{
            PullRequestId:     pr.ID,
            PullRequestName:   pr.Name,
            AuthorId:          pr.AuthorID,
            Status:            api.PullRequestStatus(pr.Status),
            AssignedReviewers: pr.AssignedReviewers,
            CreatedAt:         &pr.CreatedAt,
            MergedAt:          pr.MergedAt,
        },
    }, nil
}

func (h *pullRequestHandler) PostPullRequestReassign(
    ctx context.Context,
    req api.PostPullRequestReassignRequestObject,
) (api.PostPullRequestReassignResponseObject, error) {
    if !check.IsAdmin(ctx) {
        return api.PostPullRequestReassign404JSONResponse{
            Error: constructor.ErrorResponse(api.NOTFOUND, "resource not found"),
        }, nil
    }

    pr, newID, err := h.svc.ReassignReviewer(ctx, req.Body.PullRequestId, req.Body.OldUserId)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrPullRequestNotFound), errors.Is(err, domain.ErrUserNotFound):
            return api.PostPullRequestReassign404JSONResponse{
                Error: constructor.ErrorResponse(api.NOTFOUND, err.Error()),
            }, nil
        case errors.Is(err, domain.ErrPullRequestIsMerged):
            return api.PostPullRequestReassign409JSONResponse{
                Error: constructor.ErrorResponse(api.PRMERGED, err.Error()),
            }, nil
        case errors.Is(err, domain.ErrReviewerNotAssigned):
            return api.PostPullRequestReassign409JSONResponse{
                Error: constructor.ErrorResponse(api.NOTASSIGNED, err.Error()),
            }, nil
        case errors.Is(err, domain.ErrNoReviewerCandidate):
            return api.PostPullRequestReassign409JSONResponse{
                Error: constructor.ErrorResponse(api.NOCANDIDATE, err.Error()),
            }, nil
        default:
            return api.PostPullRequestReassign500JSONResponse{
                Error: constructor.ErrorResponse("INTERNAL_SERVER_ERROR", err.Error()),
            }, nil
        }
    }

    return api.PostPullRequestReassign200JSONResponse{
        Pr: api.PullRequest{
            PullRequestId:     pr.ID,
            PullRequestName:   pr.Name,
            AuthorId:          pr.AuthorID,
            Status:            api.PullRequestStatus(pr.Status),
            AssignedReviewers: pr.AssignedReviewers,
            CreatedAt:         &pr.CreatedAt,
            MergedAt:          pr.MergedAt,
        },
        ReplacedBy: newID,
    }, nil
}
