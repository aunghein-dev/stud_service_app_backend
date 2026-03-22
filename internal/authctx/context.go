package authctx

import "context"

type Principal struct {
	UserID     int64
	TenantID   int64
	Email      string
	FullName   string
	Role       string
	TenantSlug string
	SchoolName string
}

type principalKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok
}
