package auth

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"student_service_app/backend/internal/authctx"
	domain "student_service_app/backend/internal/domain/auth"
	"student_service_app/backend/internal/errs"
	authrepo "student_service_app/backend/internal/repository/auth"
	settingsrepo "student_service_app/backend/internal/repository/settings"
)

var tenantSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type SignUpInput struct {
	SchoolName    string
	TenantSlug    string
	AdminName     string
	Email         string
	Password      string
	SchoolPhone   string
	SchoolAddress string
}

type LoginInput struct {
	TenantSlug string
	Email      string
	Password   string
}

type Service interface {
	SignUp(ctx context.Context, input SignUpInput) (*domain.Session, error)
	Login(ctx context.Context, input LoginInput) (*domain.Session, error)
	Me(ctx context.Context) (*domain.Session, error)
}

type service struct {
	db           *sql.DB
	repo         authrepo.Repository
	settingsRepo settingsrepo.Repository
	tokenManager *TokenManager
}

func NewService(db *sql.DB, repo authrepo.Repository, settingsRepo settingsrepo.Repository, tokenManager *TokenManager) Service {
	return &service{
		db:           db,
		repo:         repo,
		settingsRepo: settingsRepo,
		tokenManager: tokenManager,
	}
}

func (s *service) SignUp(ctx context.Context, input SignUpInput) (*domain.Session, error) {
	tenantSlug := normalizeTenantSlug(input.TenantSlug)
	email := normalizeEmail(input.Email)
	if err := validateSignUpInput(input, tenantSlug, email); err != nil {
		return nil, err
	}

	existingTenant, err := s.repo.GetTenantBySlug(ctx, tenantSlug)
	if err != nil {
		return nil, err
	}
	if existingTenant != nil {
		userCount, err := s.repo.CountUsersByTenantID(ctx, existingTenant.ID)
		if err != nil {
			return nil, err
		}
		if userCount > 0 {
			return nil, errs.Conflict("workspace slug already exists")
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	tenant := existingTenant
	if tenant == nil {
		tenant = &domain.Tenant{
			Slug:          tenantSlug,
			SchoolName:    strings.TrimSpace(input.SchoolName),
			SchoolAddress: strings.TrimSpace(input.SchoolAddress),
			SchoolPhone:   strings.TrimSpace(input.SchoolPhone),
			IsActive:      true,
		}
		if err := s.repo.CreateTenant(ctx, tx, tenant); err != nil {
			return nil, mapCreateConflict(err)
		}
	} else {
		tenant.SchoolName = strings.TrimSpace(input.SchoolName)
		tenant.SchoolAddress = strings.TrimSpace(input.SchoolAddress)
		tenant.SchoolPhone = strings.TrimSpace(input.SchoolPhone)
		tenant.IsActive = true
		if err := s.repo.UpdateTenant(ctx, tx, tenant); err != nil {
			return nil, err
		}
	}

	if err := s.settingsRepo.EnsureForTenant(ctx, tx, settingsrepo.TenantDefaults{
		TenantID:      tenant.ID,
		SchoolName:    tenant.SchoolName,
		SchoolAddress: tenant.SchoolAddress,
		SchoolPhone:   tenant.SchoolPhone,
		ReceiptPrefix: deriveReceiptPrefix(tenant.SchoolName, tenant.Slug),
	}); err != nil {
		return nil, err
	}

	loginAt := time.Now().UTC()
	user := &domain.User{
		TenantID:     tenant.ID,
		FullName:     strings.TrimSpace(input.AdminName),
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         "owner",
		IsActive:     true,
		LastLoginAt:  &loginAt,
	}
	if err := s.repo.CreateUser(ctx, tx, user); err != nil {
		return nil, mapCreateConflict(err)
	}

	if err := s.repo.UpdateLastLogin(ctx, tx, user.ID, loginAt); err != nil {
		return nil, err
	}

	session := buildSession(domain.Account{
		UserID:        user.ID,
		TenantID:      tenant.ID,
		FullName:      user.FullName,
		Email:         user.Email,
		Role:          user.Role,
		UserActive:    user.IsActive,
		TenantActive:  tenant.IsActive,
		TenantSlug:    tenant.Slug,
		SchoolName:    tenant.SchoolName,
		SchoolAddress: tenant.SchoolAddress,
		SchoolPhone:   tenant.SchoolPhone,
	})
	token, err := s.tokenManager.Sign(session)
	if err != nil {
		return nil, err
	}
	session.AccessToken = token

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *service) Login(ctx context.Context, input LoginInput) (*domain.Session, error) {
	tenantSlug := normalizeTenantSlug(input.TenantSlug)
	email := normalizeEmail(input.Email)
	if tenantSlug == "" || email == "" || strings.TrimSpace(input.Password) == "" {
		return nil, errs.BadRequest("tenant_slug, email, and password are required")
	}

	account, err := s.repo.GetAccountByTenantAndEmail(ctx, tenantSlug, email)
	if err != nil {
		return nil, err
	}
	if account == nil || !account.UserActive || !account.TenantActive {
		return nil, errs.Unauthorized("invalid login credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errs.Unauthorized("invalid login credentials")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.repo.UpdateLastLogin(ctx, tx, account.UserID, time.Now().UTC()); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	session := buildSession(*account)
	token, err := s.tokenManager.Sign(session)
	if err != nil {
		return nil, err
	}
	session.AccessToken = token
	return &session, nil
}

func (s *service) Me(ctx context.Context) (*domain.Session, error) {
	principal, ok := authctx.PrincipalFromContext(ctx)
	if !ok || principal.UserID <= 0 {
		return nil, errs.Unauthorized("authentication required")
	}

	account, err := s.repo.GetAccountByUserID(ctx, principal.UserID)
	if err != nil {
		return nil, err
	}
	if account == nil || account.TenantID != principal.TenantID || !account.UserActive || !account.TenantActive {
		return nil, errs.Unauthorized("authentication required")
	}

	session := buildSession(*account)
	return &session, nil
}

func buildSession(account domain.Account) domain.Session {
	return domain.Session{
		User: domain.SessionUser{
			ID:       account.UserID,
			FullName: account.FullName,
			Email:    account.Email,
			Role:     account.Role,
		},
		Tenant: domain.SessionTenant{
			ID:            account.TenantID,
			Slug:          account.TenantSlug,
			SchoolName:    account.SchoolName,
			SchoolAddress: account.SchoolAddress,
			SchoolPhone:   account.SchoolPhone,
		},
	}
}

func validateSignUpInput(input SignUpInput, tenantSlug, email string) error {
	if strings.TrimSpace(input.SchoolName) == "" || strings.TrimSpace(input.AdminName) == "" {
		return errs.BadRequest("school_name and admin_name are required")
	}
	if email == "" || strings.TrimSpace(input.Password) == "" {
		return errs.BadRequest("email and password are required")
	}
	if !tenantSlugPattern.MatchString(tenantSlug) {
		return errs.BadRequest("tenant_slug must use lowercase letters, numbers, and hyphens only")
	}
	return nil
}

func normalizeTenantSlug(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func deriveReceiptPrefix(schoolName, slug string) string {
	source := schoolName
	if source == "" {
		source = slug
	}

	parts := strings.FieldsFunc(strings.ToUpper(source), func(r rune) bool {
		return (r < 'A' || r > 'Z') && (r < '0' || r > '9')
	})

	prefix := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		prefix += string(part[0])
		if len(prefix) >= 3 {
			break
		}
	}
	if len(prefix) < 2 {
		clean := strings.Map(func(r rune) rune {
			if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			return -1
		}, strings.ToUpper(source))
		if len(clean) >= 3 {
			prefix = clean[:3]
		} else if clean != "" {
			prefix = clean
		}
	}
	if prefix == "" {
		return "SCH"
	}
	if len(prefix) > 6 {
		return prefix[:6]
	}
	return prefix
}

func mapCreateConflict(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
		return errs.Conflict("workspace or user already exists")
	}
	return err
}
