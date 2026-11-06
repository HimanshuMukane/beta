package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/thecodearcher/limen"
	sqladapter "github.com/thecodearcher/limen/adapters/sql"
	credentialpassword "github.com/thecodearcher/limen/plugins/credential-password"
)

type UUIDGenerator struct{}

func (g *UUIDGenerator) Generate(ctx context.Context) (any, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (g *UUIDGenerator) GetColumnType() limen.ColumnType {
	return limen.ColumnTypeUUID
}

func NewLimenAuth(db *sql.DB, secret string) (*limen.Limen, error) {
	if len(secret) < 32 {
		secret = "especa_limen_default_session_secret_32bytes"
	}

	adapter := sqladapter.NewMySQL(db)
	pwPlugin := credentialpassword.New()

	schemaConfig := limen.NewDefaultSchemaConfig(
		limen.WithSchemaIDGenerator(&UUIDGenerator{}),
		limen.WithSchemaUser(limen.WithUserTableName("limen_users")),
		limen.WithSchemaSession(limen.WithSessionTableName("limen_sessions")),
		limen.WithSchemaVerification(limen.WithVerificationTableName("limen_verifications")),
		limen.WithSchemaAccount(limen.WithAccountTableName("limen_accounts")),
		limen.WithSchemaRateLimit(limen.WithRateLimitTableName("limen_rate_limits")),
	)

	auth, err := limen.New(&limen.Config{
		Database: adapter,
		Plugins:  []limen.Plugin{pwPlugin},
		Schema:   schemaConfig,
		Secret:   []byte(secret[:32]),
	})
	if err != nil {
		return nil, fmt.Errorf("initialize limen auth: %w", err)
	}

	return auth, nil
}

func LimenSessionUser(auth *limen.Limen, r *http.Request) (*limen.ValidatedSession, error) {
	if auth == nil {
		return nil, fmt.Errorf("limen auth not configured")
	}
	return auth.GetSession(r)
}
