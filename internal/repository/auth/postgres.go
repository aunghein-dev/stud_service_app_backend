package auth

import (
	"context"
	"database/sql"
	"time"

	domain "student_service_app/backend/internal/domain/auth"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetTenantBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	query := `SELECT id, slug, school_name, school_address, school_phone, is_active, created_at, updated_at
	FROM tenants WHERE slug=$1`
	var tenant domain.Tenant
	if err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&tenant.ID,
		&tenant.Slug,
		&tenant.SchoolName,
		&tenant.SchoolAddress,
		&tenant.SchoolPhone,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *postgresRepository) CountUsersByTenantID(ctx context.Context, tenantID int64) (int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenant_users WHERE tenant_id=$1`, tenantID).Scan(&total)
	return total, err
}

func (r *postgresRepository) CreateTenant(ctx context.Context, tx *sql.Tx, tenant *domain.Tenant) error {
	query := `INSERT INTO tenants (slug, school_name, school_address, school_phone, is_active)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, created_at, updated_at`
	return tx.QueryRowContext(ctx, query,
		tenant.Slug,
		tenant.SchoolName,
		tenant.SchoolAddress,
		tenant.SchoolPhone,
		tenant.IsActive,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

func (r *postgresRepository) UpdateTenant(ctx context.Context, tx *sql.Tx, tenant *domain.Tenant) error {
	query := `UPDATE tenants
	SET school_name=$2, school_address=$3, school_phone=$4, is_active=$5, updated_at=NOW()
	WHERE id=$1
	RETURNING updated_at`
	return tx.QueryRowContext(ctx, query,
		tenant.ID,
		tenant.SchoolName,
		tenant.SchoolAddress,
		tenant.SchoolPhone,
		tenant.IsActive,
	).Scan(&tenant.UpdatedAt)
}

func (r *postgresRepository) CreateUser(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	query := `INSERT INTO tenant_users (tenant_id, full_name, email, password_hash, role, is_active, last_login_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, created_at, updated_at`
	return tx.QueryRowContext(ctx, query,
		user.TenantID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		user.LastLoginAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *postgresRepository) GetAccountByTenantAndEmail(ctx context.Context, tenantSlug, email string) (*domain.Account, error) {
	query := `SELECT
		u.id,
		u.tenant_id,
		u.full_name,
		u.email,
		u.password_hash,
		u.role,
		u.is_active,
		t.is_active,
		t.slug,
		t.school_name,
		t.school_address,
		t.school_phone
	FROM tenant_users u
	JOIN tenants t ON t.id = u.tenant_id
	WHERE t.slug=$1 AND LOWER(u.email)=LOWER($2)`
	var account domain.Account
	if err := r.db.QueryRowContext(ctx, query, tenantSlug, email).Scan(
		&account.UserID,
		&account.TenantID,
		&account.FullName,
		&account.Email,
		&account.PasswordHash,
		&account.Role,
		&account.UserActive,
		&account.TenantActive,
		&account.TenantSlug,
		&account.SchoolName,
		&account.SchoolAddress,
		&account.SchoolPhone,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *postgresRepository) GetAccountByUserID(ctx context.Context, userID int64) (*domain.Account, error) {
	query := `SELECT
		u.id,
		u.tenant_id,
		u.full_name,
		u.email,
		u.password_hash,
		u.role,
		u.is_active,
		t.is_active,
		t.slug,
		t.school_name,
		t.school_address,
		t.school_phone
	FROM tenant_users u
	JOIN tenants t ON t.id = u.tenant_id
	WHERE u.id=$1`
	var account domain.Account
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&account.UserID,
		&account.TenantID,
		&account.FullName,
		&account.Email,
		&account.PasswordHash,
		&account.Role,
		&account.UserActive,
		&account.TenantActive,
		&account.TenantSlug,
		&account.SchoolName,
		&account.SchoolAddress,
		&account.SchoolPhone,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *postgresRepository) UpdateLastLogin(ctx context.Context, tx *sql.Tx, userID int64, loginAt time.Time) error {
	_, err := tx.ExecContext(ctx, `UPDATE tenant_users SET last_login_at=$2, updated_at=NOW() WHERE id=$1`, userID, loginAt)
	return err
}
