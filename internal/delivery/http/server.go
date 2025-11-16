package http

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/app"
    "github.com/kimvlry/avito-internship-assignment/pkg/logger"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    chimiddleware "github.com/go-chi/chi/v5/middleware"

    "github.com/kimvlry/avito-internship-assignment/api"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/handler"
    "github.com/kimvlry/avito-internship-assignment/internal/delivery/http/middleware"
)

type Server struct {
    srv *http.Server
}

func NewServer(cfg app.HttpConfig, handlers *handler.Handlers) *Server {
    router := setupRouter(cfg.JwtSecret, handlers)

    return &Server{
        srv: &http.Server{
            Addr:         cfg.Addr(),
            Handler:      router,
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            IdleTimeout:  cfg.IdleTimeout,
        },
    }
}

func (s *Server) Start() error {
    logger.Info(context.Background(), fmt.Sprintf("Starting HTTP server on %s", s.srv.Addr))
    if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        return fmt.Errorf("http server error: %w", err)
    }
    return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
    logger.Info(context.Background(), "Shutting down HTTP server...")
    if err := s.srv.Shutdown(ctx); err != nil {
        return fmt.Errorf("http server shutdown failed: %w", err)
    }
    logger.Info(context.Background(), "HTTP server stopped")
    return nil
}

func setupRouter(jwtSecret string, handlers *handler.Handlers) http.Handler {
    r := chi.NewRouter()

    r.Use(chimiddleware.RequestID)
    r.Use(chimiddleware.RealIP)
    r.Use(chimiddleware.Logger)
    r.Use(chimiddleware.Recoverer)
    r.Use(chimiddleware.Timeout(60 * time.Second))

    jwtAuth := middleware.NewJWTMiddleware(jwtSecret)

    strictHandler := api.NewStrictHandler(handlers, nil)

    r.Group(func(r chi.Router) {
        r.Post("/team/add", strictHandler.PostTeamAdd)
    })

    r.Group(func(r chi.Router) {
        r.Use(jwtAuth)

        r.Post("/pullRequest/create", strictHandler.PostPullRequestCreate)
        r.Post("/pullRequest/merge", strictHandler.PostPullRequestMerge)
        r.Post("/pullRequest/reassign", strictHandler.PostPullRequestReassign)

        r.Get("/team/get", handleGetWithQuery(
            "team_name",
            func(ctx context.Context, teamName string) (api.GetTeamGetResponseObject, error) {
                return handlers.GetTeamGet(ctx, api.GetTeamGetRequestObject{
                    Params: api.GetTeamGetParams{TeamName: teamName},
                })
            },
        ))

        r.Post("/users/setIsActive", strictHandler.PostUsersSetIsActive)
        r.Get("/users/getReview", handleGetWithQuery(
            "user_id",
            func(ctx context.Context, userID string) (api.GetUsersGetReviewResponseObject, error) {
                return handlers.GetUsersGetReview(ctx, api.GetUsersGetReviewRequestObject{
                    Params: api.GetUsersGetReviewParams{UserId: userID},
                })
            },
        ))
        r.Get("/stats/assignments", strictHandler.GetStatsAssignments)
    })
    return r
}

type ResponseVisitor interface {
    VisitGetTeamGetResponse(w http.ResponseWriter) error
    VisitGetUsersGetReviewResponse(w http.ResponseWriter) error
}

func handleGetWithQuery[T any](
    paramName string,
    handlerFunc func(ctx context.Context, paramValue string) (T, error),
) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        paramValue := r.URL.Query().Get(paramName)
        if paramValue == "" {
            writeError(w, http.StatusBadRequest, "BAD_REQUEST", fmt.Sprintf("%s is required", paramName))
            return
        }

        resp, err := handlerFunc(r.Context(), paramValue)
        if err != nil {
            writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
            return
        }

        if visitor, ok := any(resp).(interface {
            VisitGetTeamGetResponse(w http.ResponseWriter) error
        }); ok {
            if err := visitor.VisitGetTeamGetResponse(w); err != nil {
                writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to write response")
            }
            return
        }

        if visitor, ok := any(resp).(interface {
            VisitGetUsersGetReviewResponse(w http.ResponseWriter) error
        }); ok {
            if err := visitor.VisitGetUsersGetReviewResponse(w); err != nil {
                writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to write response")
            }
            return
        }
        writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "unknown response type")
    }
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    err := json.NewEncoder(w).Encode(map[string]interface{}{
        "error": map[string]string{
            "code":    code,
            "message": message,
        },
    })
    if err != nil {
        return
    }
}
