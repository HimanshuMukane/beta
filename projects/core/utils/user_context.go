package utils

import (
	"context"

	"core/models"
	"github.com/google/uuid"
)

type userContextKey struct{}
type tenantContextKey struct{}

var currentUserKey = userContextKey{}
var currentTenantKey = tenantContextKey{}

type TenantContext struct {
	TenantID           uuid.UUID
	TenantName         string
	TenantSlug         string
	Role               string
	SubscriptionStatus string
	IsReadOnly         bool
}

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

func UserFromContext(ctx context.Context) *models.User {
	if ctx == nil {
		return nil
	}
	if u, ok := ctx.Value(currentUserKey).(*models.User); ok {
		return u
	}
	return nil
}

func WithTenantContext(ctx context.Context, tc *TenantContext) context.Context {
	return context.WithValue(ctx, currentTenantKey, tc)
}

func TenantFromContext(ctx context.Context) *TenantContext {
	if ctx == nil {
		return nil
	}
	if tc, ok := ctx.Value(currentTenantKey).(*TenantContext); ok {
		return tc
	}
	return nil
}
