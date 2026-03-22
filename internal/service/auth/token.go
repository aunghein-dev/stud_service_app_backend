package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"student_service_app/backend/internal/config"
	domain "student_service_app/backend/internal/domain/auth"
)

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

type tokenClaims struct {
	UserID     int64  `json:"user_id"`
	TenantID   int64  `json:"tenant_id"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	TenantSlug string `json:"tenant_slug"`
	SchoolName string `json:"school_name"`
	IssuedAt   int64  `json:"iat"`
	ExpiresAt  int64  `json:"exp"`
}

func NewTokenManager(cfg *config.Config) *TokenManager {
	return NewTokenManagerWithValues(cfg.Auth.Secret, time.Duration(cfg.Auth.AccessTokenTTLHours)*time.Hour)
}

func NewTokenManagerWithValues(secret string, ttl time.Duration) *TokenManager {
	if ttl <= 0 {
		ttl = 72 * time.Hour
	}
	return &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *TokenManager) Sign(session domain.Session) (string, error) {
	now := time.Now().UTC()
	claims := tokenClaims{
		UserID:     session.User.ID,
		TenantID:   session.Tenant.ID,
		Email:      session.User.Email,
		FullName:   session.User.FullName,
		Role:       session.User.Role,
		TenantSlug: session.Tenant.Slug,
		SchoolName: session.Tenant.SchoolName,
		IssuedAt:   now.Unix(),
		ExpiresAt:  now.Add(m.ttl).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := m.sign(encodedPayload)
	return encodedPayload + "." + signature, nil
}

func (m *TokenManager) Parse(token string) (*tokenClaims, error) {
	parts := splitToken(token)
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}

	if !hmac.Equal([]byte(parts[1]), []byte(m.sign(parts[0]))) {
		return nil, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims tokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid token payload")
	}
	if claims.ExpiresAt <= time.Now().UTC().Unix() {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func (m *TokenManager) sign(payload string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func splitToken(token string) []string {
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			return []string{token[:i], token[i+1:]}
		}
	}
	return nil
}
