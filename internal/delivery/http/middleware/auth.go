package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
    ContextUserID  contextKey = "user_id"
    ContextIsAdmin contextKey = "is_admin"
)

type JWTConfig struct {
    Secret string
}

func NewJWTMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenString := r.Header.Get("Authorization")
            if tokenString == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            tokenString = strings.TrimPrefix(tokenString, "Bearer ")

            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, jwt.ErrTokenUnverifiable
                }
                return []byte(secret), nil
            })
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok || !token.Valid {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            userID, ok := claims["user_id"].(string)
            if !ok {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            isAdmin, _ := claims["is_admin"].(bool)

            ctx := context.WithValue(r.Context(), ContextUserID, userID)
            ctx = context.WithValue(ctx, ContextIsAdmin, isAdmin)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserID(ctx context.Context) string {
    if userID, ok := ctx.Value(ContextUserID).(string); ok {
        return userID
    }
    return ""
}

func IsAdmin(ctx context.Context) bool {
    if isAdmin, ok := ctx.Value(ContextIsAdmin).(bool); ok {
        return isAdmin
    }
    return false
}
