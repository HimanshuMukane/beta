package utils

import (
	"context"

	"core/models"
	"github.com/google/uuid"
)

type userContextKey struct{}
type tenantContextKey struct{}
type projectContextKey struct{}

var currentUserKey = userContextKey{}
var currentTenantKey = tenantContextKey{}
var currentProjectKey = projectContextKey{}

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

func WithProjectName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, currentProjectKey, name)
}

func ProjectNameFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if name, ok := ctx.Value(currentProjectKey).(string); ok {
		return name
	}
	return ""
}
