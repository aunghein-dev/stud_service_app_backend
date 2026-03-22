package auth

import "time"

type Tenant struct {
	ID            int64
	Slug          string
	SchoolName    string
	SchoolAddress string
	SchoolPhone   string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	ID           int64
	TenantID     int64
	FullName     string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Account struct {
	UserID        int64
	TenantID      int64
	FullName      string
	Email         string
	PasswordHash  string
	Role          string
	UserActive    bool
	TenantActive  bool
	TenantSlug    string
	SchoolName    string
	SchoolAddress string
	SchoolPhone   string
}

type Session struct {
	AccessToken string
	User        SessionUser
	Tenant      SessionTenant
}

type SessionUser struct {
	ID       int64
	FullName string
	Email    string
	Role     string
}

type SessionTenant struct {
	ID            int64
	Slug          string
	SchoolName    string
	SchoolAddress string
	SchoolPhone   string
}
