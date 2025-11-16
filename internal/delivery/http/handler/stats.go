package handler

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/api"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"
)

type statsHandler struct {
    statsSvc *service.StatsService
}

func newStatsHandler(statsSvc *service.StatsService) *statsHandler {
    return &statsHandler{statsSvc: statsSvc}
}

func (h *statsHandler) GetStatsAssignments(ctx context.Context, request api.GetStatsAssignmentsRequestObject) (api.GetStatsAssignmentsResponseObject, error) {
    userStats, err := h.statsSvc.GetUserAssignmentStats(ctx)
    if err != nil {
        return nil, err
    }

    resp := api.GetStatsAssignments200JSONResponse{
        ByUser: func() *[]api.AssignmentCountPerUser {
            res := make([]api.AssignmentCountPerUser, 0, len(userStats))
            for _, stat := range userStats {
                res = append(res, api.AssignmentCountPerUser{
                    UserId:        stat.UserID,
                    AssignedCount: stat.AssignedPRs,
                })
            }
            return &res
        }(),
    }
    return resp, nil
}
