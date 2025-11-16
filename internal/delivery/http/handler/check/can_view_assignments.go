package check

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/middleware"
)

func CanViewUserAssignments(ctx context.Context, isActive bool) bool {
    userID := middleware.GetUserID(ctx)
    if middleware.IsAdmin(ctx) {
        return true
    }
    if isActive {
        return userID != ""
    }
    return false
}
