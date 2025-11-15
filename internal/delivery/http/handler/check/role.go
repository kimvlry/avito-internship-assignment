package check

import (
    "context"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/middleware"
)

func IsAdmin(ctx context.Context) bool {
    return middleware.IsAdmin(ctx)
}

func IsAdminOrOwner(ctx context.Context, ownerID string) bool {
    return middleware.IsAdmin(ctx) || middleware.GetUserID(ctx) == ownerID
}
