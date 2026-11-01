package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"core/db/sqlc"
	"core/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleAlreadyExists  = errors.New("role already exists")
)

type RoleRepository struct {
	q *sqlc.Queries
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{q: sqlc.New(db.DB)}
}

func (r *RoleRepository) CreateRole(ctx context.Context, name string) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 50 {
		return nil, errors.New("role name must be between 1 and 50 bytes")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate role id: %w", err)
	}

	err = r.q.CreateRole(ctx, sqlc.CreateRoleParams{
		ID:   id.String(),
		Name: name,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return nil, ErrRoleAlreadyExists
		}
		return nil, fmt.Errorf("create role: %w", err)
	}

	return r.GetRoleByID(ctx, id)
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	if id == uuid.Nil {
		return nil, ErrRoleNotFound
	}

	row, err := r.q.GetRoleByID(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by id: %w", err)
	}

	return mapRole(row.ID, row.Name, row.CreatedAt, row.UpdatedAt, row.DeletedAt)
}

func (r *RoleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNotFound
	}

	row, err := r.q.GetRoleByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	return mapRole(row.ID, row.Name, row.CreatedAt, row.UpdatedAt, row.DeletedAt)
}

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{q: sqlc.New(db.DB)}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if user.FirstName == "" || len(user.FirstName) > 50 {
		return errors.New("first name must be between 1 and 50 bytes")
	}
	if user.LastName == "" || len(user.LastName) > 50 {
		return errors.New("last name must be between 1 and 50 bytes")
	}
	if user.Email == "" || len(user.Email) > 254 {
		return errors.New("email must be between 1 and 254 bytes")
	}
	if user.PasswordHash == "" || len(user.PasswordHash) > 255 {
		return errors.New("password hash must be between 1 and 255 bytes")
	}
	if user.RoleID == uuid.Nil {
		return errors.New("role ID is required")
	}

	if user.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate user id: %w", err)
		}
		user.ID = id
	}

	err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           user.ID.String(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		RoleID:       user.RoleID.String(),
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	row, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return mapUser(row.ID, row.FirstName, row.LastName, row.Email, row.PasswordHash, row.RoleID, row.RoleName, row.CreatedAt, row.UpdatedAt, row.DeletedAt)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	row, err := r.q.GetUserByID(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return mapUser(row.ID, row.FirstName, row.LastName, row.Email, row.PasswordHash, row.RoleID, row.RoleName, row.CreatedAt, row.UpdatedAt, row.DeletedAt)
}

func mapRole(idStr, name, createdStr, updatedStr string, deletedStr sql.NullString) (*models.Role, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse role id: %w", err)
	}
	createdAt, err := parseSQLiteTime(createdStr)
	if err != nil {
		return nil, fmt.Errorf("parse role created_at: %w", err)
	}
	updatedAt, err := parseSQLiteTime(updatedStr)
	if err != nil {
		return nil, fmt.Errorf("parse role updated_at: %w", err)
	}

	role := &models.Role{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if deletedStr.Valid && deletedStr.String != "" {
		t, err := parseSQLiteTime(deletedStr.String)
		if err != nil {
			return nil, fmt.Errorf("parse role deleted_at: %w", err)
		}
		role.DeletedAt = &t
	}
	return role, nil
}

func mapUser(idStr, firstName, lastName, email, pwdHash, roleIDStr, roleName, createdStr, updatedStr string, deletedStr sql.NullString) (*models.User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse user role_id: %w", err)
	}
	createdAt, err := parseSQLiteTime(createdStr)
	if err != nil {
		return nil, fmt.Errorf("parse user created_at: %w", err)
	}
	updatedAt, err := parseSQLiteTime(updatedStr)
	if err != nil {
		return nil, fmt.Errorf("parse user updated_at: %w", err)
	}

	user := &models.User{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: pwdHash,
		RoleID:       roleID,
		RoleName:     roleName,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
	if deletedStr.Valid && deletedStr.String != "" {
		t, err := parseSQLiteTime(deletedStr.String)
		if err != nil {
			return nil, fmt.Errorf("parse user deleted_at: %w", err)
		}
		user.DeletedAt = &t
	}
	return user, nil
}

func parseSQLiteTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported SQLite timestamp %q", value)
}

