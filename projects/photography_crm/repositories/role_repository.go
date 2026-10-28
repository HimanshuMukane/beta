package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

// CreateRole inserts a new role and returns the stored record.
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

	_, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO roles (id, name)
			VALUES (?, ?)
		`,
		models.UUIDToString(role.ID),
		role.Name,
	)
	if err != nil {
		if strings.Contains(
			strings.ToLower(err.Error()),
			"unique constraint failed",
		) {
			return nil, ErrRoleAlreadyExists
		}

		return nil, fmt.Errorf("create role: %w", err)
	}

	createdRole, err := r.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("retrieve created role: %w", err)
	}

	return createdRole, nil
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
		ID        string         `db:"id"`
		Name      string         `db:"name"`
		CreatedAt string         `db:"created_at"`
		UpdatedAt string         `db:"updated_at"`
		DeletedAt sql.NullString `db:"deleted_at"`
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

	return mapRoleRow(
		row.ID,
		row.Name,
		row.CreatedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)
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
		ID        string         `db:"id"`
		Name      string         `db:"name"`
		CreatedAt string         `db:"created_at"`
		UpdatedAt string         `db:"updated_at"`
		DeletedAt sql.NullString `db:"deleted_at"`
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

	return mapRoleRow(
		row.ID,
		row.Name,
		row.CreatedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)
}

func mapRoleRow(
	idValue string,
	name string,
	createdAtValue string,
	updatedAtValue string,
	deletedAtValue sql.NullString,
) (*models.Role, error) {
	id, err := models.UUIDFromString(idValue)
	if err != nil {
		return nil, fmt.Errorf("decode role ID: %w", err)
	}

	createdAt, err := parseSQLiteTime(createdAtValue)
	if err != nil {
		return nil, fmt.Errorf("parse role creation time: %w", err)
	}

	updatedAt, err := parseSQLiteTime(updatedAtValue)
	if err != nil {
		return nil, fmt.Errorf("parse role update time: %w", err)
	}

	role := &models.Role{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if deletedAtValue.Valid && deletedAtValue.String != "" {
		deletedAt, err := parseSQLiteTime(deletedAtValue.String)
		if err != nil {
			return nil, fmt.Errorf("parse role deletion time: %w", err)
		}

		role.DeletedAt = &deletedAt
	}

	return role, nil
}
