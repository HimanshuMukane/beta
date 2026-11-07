package core

import (
	"context"
	"testing"

	"core/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestIdentityDBFromContextIsolation(t *testing.T) {
	appDB := &sqlx.DB{}
	identityDB := &sqlx.DB{}

	ctx := context.Background()

	// 1. Empty context returns nil
	if got := utils.IdentityDBFromContext(ctx); got != nil {
		t.Fatalf("expected nil from empty context, got %v", got)
	}

	// 2. Context with ONLY app DB must NOT be returned by IdentityDBFromContext (NO silent fallback)
	ctxWithAppDB := utils.WithDB(ctx, appDB)
	if got := utils.IdentityDBFromContext(ctxWithAppDB); got != nil {
		t.Fatalf("expected nil when only app DB is in context (silent fallback prevented), got %v", got)
	}

	// 3. Context with Identity DB returns identity DB
	ctxWithBoth := utils.WithIdentityDB(ctxWithAppDB, identityDB)
	if got := utils.IdentityDBFromContext(ctxWithBoth); got != identityDB {
		t.Fatalf("expected identityDB, got %v", got)
	}

	// Verify DBFromContext returns appDB
	if got := utils.DBFromContext(ctxWithBoth); got != appDB {
		t.Fatalf("expected appDB, got %v", got)
	}

	// 4. Fallback is honored only when explicitly provided
	fallbackDB := &sqlx.DB{}
	if got := utils.IdentityDBFromContext(ctx, fallbackDB); got != fallbackDB {
		t.Fatalf("expected fallbackDB, got %v", got)
	}
}
