package middleware

import (
	"net/http"
	"strings"
	"time"

	"student_service_app/backend/internal/authctx"
	"student_service_app/backend/internal/errs"
	"student_service_app/backend/internal/response"
	authservice "student_service_app/backend/internal/service/auth"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Recoverer() func(http.Handler) http.Handler {
	return chimiddleware.Recoverer
}

func RequestID() func(http.Handler) http.Handler {
	return chimiddleware.RequestID
}

func Logger(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func RequireAuth(tokenManager *authservice.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				response.Error(w, errs.Unauthorized("missing bearer token"))
				return
			}

			token := strings.TrimSpace(authHeader[len("Bearer "):])
			claims, err := tokenManager.Parse(token)
			if err != nil {
				response.Error(w, errs.Unauthorized("invalid or expired token"))
				return
			}

			ctx := authctx.WithPrincipal(r.Context(), authctx.Principal{
				UserID:     claims.UserID,
				TenantID:   claims.TenantID,
				Email:      claims.Email,
				FullName:   claims.FullName,
				Role:       claims.Role,
				TenantSlug: claims.TenantSlug,
				SchoolName: claims.SchoolName,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
