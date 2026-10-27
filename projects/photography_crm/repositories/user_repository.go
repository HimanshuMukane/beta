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
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a user into SQLite.
//
// The caller must hash the plaintext password before calling this method.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = normalizeEmail(user.Email)

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
		user.ID = uuid.New()
	}

	query := `
		INSERT INTO users (
			id,
			first_name,
			last_name,
			email,
			password_hash,
			role_id
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		models.UUIDToString(user.ID),
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		models.UUIDToString(user.RoleID),
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return ErrEmailAlreadyExists
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves an active user and their active role.
func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	query := `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.password_hash,
			u.role_id,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			r.name AS role_name
		FROM users AS u
		INNER JOIN roles AS r
			ON r.id = u.role_id
		WHERE u.email = ?
			AND u.deleted_at IS NULL
			AND r.deleted_at IS NULL
		LIMIT 1
	`

	return r.getUser(ctx, query, normalizeEmail(email))
}

// GetUserByID retrieves an active user and their active role.
func (r *UserRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	query := `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.password_hash,
			u.role_id,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			r.name AS role_name
		FROM users AS u
		INNER JOIN roles AS r
			ON r.id = u.role_id
		WHERE u.id = ?
			AND u.deleted_at IS NULL
			AND r.deleted_at IS NULL
		LIMIT 1
	`

	return r.getUser(ctx, query, models.UUIDToString(id))
}

func (r *UserRepository) getUser(
	ctx context.Context,
	query string,
	args ...any,
) (*models.User, error) {
	var row struct {
		ID           string         `db:"id"`
		FirstName    string         `db:"first_name"`
		LastName     string         `db:"last_name"`
		Email        string         `db:"email"`
		PasswordHash string         `db:"password_hash"`
		RoleID       string         `db:"role_id"`
		CreatedAt    string         `db:"created_at"`
		UpdatedAt    string         `db:"updated_at"`
		DeletedAt    sql.NullString `db:"deleted_at"`
		RoleName     string         `db:"role_name"`
	}

	err := r.db.GetContext(ctx, &row, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	id, err := models.UUIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("decode user ID: %w", err)
	}

	roleID, err := models.UUIDFromString(row.RoleID)
	if err != nil {
		return nil, fmt.Errorf("decode user role ID: %w", err)
	}

	user := &models.User{
		ID:           id,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleID:       roleID,
		RoleName:     row.RoleName,
	}

	user.CreatedAt, err = parseSQLiteTime(row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse user creation time: %w", err)
	}

	user.UpdatedAt, err = parseSQLiteTime(row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse user update time: %w", err)
	}

	if row.DeletedAt.Valid && row.DeletedAt.String != "" {
		deletedAt, err := parseSQLiteTime(row.DeletedAt.String)
		if err != nil {
			return nil, fmt.Errorf("parse user deletion time: %w", err)
		}
		user.DeletedAt = &deletedAt
	}

	return user, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
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
