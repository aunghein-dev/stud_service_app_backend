package settings

import (
	"context"
	"database/sql"

	"student_service_app/backend/internal/domain/settings"
)

type Repository interface {
	Get(ctx context.Context) (*settings.Setting, error)
	Update(ctx context.Context, s *settings.Setting) error
	EnsureForTenant(ctx context.Context, tx *sql.Tx, defaults TenantDefaults) error
	AllocateReceiptNo(ctx context.Context, tx *sql.Tx) (string, error)
}

type TenantDefaults struct {
	TenantID      int64
	SchoolName    string
	SchoolAddress string
	SchoolPhone   string
	ReceiptPrefix string
}
