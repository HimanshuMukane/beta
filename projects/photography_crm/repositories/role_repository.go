package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"chi-mysql-server/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// CreateRole inserts a new role.
func (r *RoleRepository) CreateRole(
	ctx context.Context,
	name string,
) (*models.Role, error) {
	name = strings.TrimSpace(name)

	if name == "" || len(name) > 50 {
		return nil, errors.New("role name must be between 1 and 50 bytes")
	}

	role := &models.Role{
		ID:   uuid.New(),
		Name: name,
	}

	query := `
		INSERT INTO roles (id, name)
		VALUES (?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		models.UUIDToString(role.ID),
		role.Name,
	)
	if err != nil {
		// SQLite unique constraint violations use extended error code 2067.
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return nil, ErrRoleAlreadyExists
		}

		return nil, fmt.Errorf("create role: %w", err)
	}

	err = r.db.GetContext(
		ctx,
		role,
		`
			SELECT created_at, updated_at
			FROM roles
			WHERE id = ?
			LIMIT 1
		`,
		models.UUIDToString(role.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("load created role: %w", err)
	}

	return role, nil
}

// GetRoleByName retrieves an active role by name.
func (r *RoleRepository) GetRoleByName(
	ctx context.Context,
	name string,
) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNotFound
	}

	var row struct {
		ID        string       `db:"id"`
		Name      string       `db:"name"`
		CreatedAt time.Time    `db:"created_at"`
		UpdatedAt time.Time    `db:"updated_at"`
		DeletedAt sql.NullTime `db:"deleted_at"`
	}

	err := r.db.GetContext(
		ctx,
		&row,
		`
			SELECT id, name, created_at, updated_at, deleted_at
			FROM roles
			WHERE name = ?
			  AND deleted_at IS NULL
			LIMIT 1
		`,
		name,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	id, err := models.UUIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("decode role ID: %w", err)
	}

	role := &models.Role{
		ID:        id,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.DeletedAt.Valid {
		deletedAt := row.DeletedAt.Time
		role.DeletedAt = &deletedAt
	}

	return role, nil
}

// GetRoleByID retrieves an active role by UUID.
func (r *RoleRepository) GetRoleByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Role, error) {
	if id == uuid.Nil {
		return nil, ErrRoleNotFound
	}

	var row struct {
		ID        string       `db:"id"`
		Name      string       `db:"name"`
		CreatedAt time.Time    `db:"created_at"`
		UpdatedAt time.Time    `db:"updated_at"`
		DeletedAt sql.NullTime `db:"deleted_at"`
	}

	err := r.db.GetContext(
		ctx,
		&row,
		`
			SELECT id, name, created_at, updated_at, deleted_at
			FROM roles
			WHERE id = ?
			  AND deleted_at IS NULL
			LIMIT 1
		`,
		models.UUIDToString(id),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by ID: %w", err)
	}

	roleID, err := models.UUIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("decode role ID: %w", err)
	}

	role := &models.Role{
		ID:        roleID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.DeletedAt.Valid {
		deletedAt := row.DeletedAt.Time
		role.DeletedAt = &deletedAt
	}

	return role, nil
}
