package auth

import domain "student_service_app/backend/internal/domain/auth"

type SignUpRequest struct {
	SchoolName    string `json:"school_name" validate:"required,min=2,max=150"`
	TenantSlug    string `json:"tenant_slug" validate:"required,min=3,max=50"`
	AdminName     string `json:"admin_name" validate:"required,min=2,max=150"`
	Email         string `json:"email" validate:"required,email,max=150"`
	Password      string `json:"password" validate:"required,min=8,max=72"`
	SchoolPhone   string `json:"school_phone" validate:"max=30"`
	SchoolAddress string `json:"school_address" validate:"max=255"`
}

type LoginRequest struct {
	TenantSlug string `json:"tenant_slug" validate:"required,min=3,max=50"`
	Email      string `json:"email" validate:"required,email,max=150"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
}

type SessionResponse struct {
	AccessToken string         `json:"access_token,omitempty"`
	User        UserResponse   `json:"user"`
	Tenant      TenantResponse `json:"tenant"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type TenantResponse struct {
	ID            int64  `json:"id"`
	Slug          string `json:"slug"`
	SchoolName    string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
	SchoolPhone   string `json:"school_phone"`
}

func FromSession(session domain.Session) SessionResponse {
	return SessionResponse{
		AccessToken: session.AccessToken,
		User: UserResponse{
			ID:       session.User.ID,
			FullName: session.User.FullName,
			Email:    session.User.Email,
			Role:     session.User.Role,
		},
		Tenant: TenantResponse{
			ID:            session.Tenant.ID,
			Slug:          session.Tenant.Slug,
			SchoolName:    session.Tenant.SchoolName,
			SchoolAddress: session.Tenant.SchoolAddress,
			SchoolPhone:   session.Tenant.SchoolPhone,
		},
	}
}
