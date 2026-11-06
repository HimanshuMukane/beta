package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type TenantContextKey struct{}

var ActiveTenantKey = TenantContextKey{}

func WithTenantTx(ctx context.Context, db *sql.DB, tenantID uuid.UUID, fn func(*sql.Tx) error) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tenant transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID.String())
	if err != nil {
		return fmt.Errorf("set transaction tenant context: %w", err)
	}

	txCtx := context.WithValue(ctx, ActiveTenantKey, tenantID)
	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tenant transaction: %w", err)
	}
	_ = txCtx
	return nil
}

