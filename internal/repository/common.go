package repository

import (
	"context"
	"database/sql"

	"student_service_app/backend/internal/authctx"
	"student_service_app/backend/internal/errs"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func TenantID(ctx context.Context) (int64, error) {
	principal, ok := authctx.PrincipalFromContext(ctx)
	if !ok || principal.TenantID <= 0 {
		return 0, errs.Unauthorized("authentication required")
	}
	return principal.TenantID, nil
}
