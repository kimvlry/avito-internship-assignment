package handler

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/api"
)

type Server struct {
    PullRequestHandler
    TeamHandler
    UserHandler
}

func NewServer() *Server {
    return &Server{}
}

type PullRequestHandler interface {
    PostPullRequestCreate(ctx context.Context, request api.PostPullRequestCreateRequestObject) (api.PostPullRequestCreateResponseObject, error)
    PostPullRequestMerge(ctx context.Context, request api.PostPullRequestMergeRequestObject) (api.PostPullRequestMergeResponseObject, error)
    PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (
        api.PostPullRequestReassignResponseObject,
        error,
    )
}

type TeamHandler interface {
    PostTeamAdd(ctx context.Context, request api.PostTeamAddRequestObject) (api.PostTeamAddResponseObject, error)
    GetTeamGet(ctx context.Context, request api.GetTeamGetRequestObject) (api.GetTeamGetResponseObject, error)
}

type UserHandler interface {
    PostUsersSetIsActive(ctx context.Context, request api.PostUsersSetIsActiveRequestObject) (api.PostUsersSetIsActiveResponseObject, error)
    GetUsersGetReview(ctx context.Context, request api.GetUsersGetReviewRequestObject) (api.GetUsersGetReviewResponseObject, error)
}
