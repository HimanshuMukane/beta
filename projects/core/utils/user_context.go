package utils

import (
	"context"

	"core/models"
)

type userContextKey struct{}

var currentUserKey = userContextKey{}

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

