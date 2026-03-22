package auth

import (
	"context"
	"database/sql"
	"time"

	domain "student_service_app/backend/internal/domain/auth"
)

type Repository interface {
	GetTenantBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	CountUsersByTenantID(ctx context.Context, tenantID int64) (int, error)
	CreateTenant(ctx context.Context, tx *sql.Tx, tenant *domain.Tenant) error
	UpdateTenant(ctx context.Context, tx *sql.Tx, tenant *domain.Tenant) error
	CreateUser(ctx context.Context, tx *sql.Tx, user *domain.User) error
	GetAccountByTenantAndEmail(ctx context.Context, tenantSlug, email string) (*domain.Account, error)
	GetAccountByUserID(ctx context.Context, userID int64) (*domain.Account, error)
	UpdateLastLogin(ctx context.Context, tx *sql.Tx, userID int64, loginAt time.Time) error
}
