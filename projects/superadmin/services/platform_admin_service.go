package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"core/security"
	"github.com/jmoiron/sqlx"
	"superadmin/models"
	"superadmin/repositories"
)

var (
	ErrRequestAlreadyProcessed = errors.New("request has already been reviewed")
)

type PlatformAdminService struct {
	db        *sqlx.DB
	adminRepo *repositories.PlatformAdminRepository
	secret    []byte
	ttl       time.Duration
}

func NewPlatformAdminService(db *sqlx.DB, secret string, ttl time.Duration) *PlatformAdminService {
	return &PlatformAdminService{
		db:        db,
		adminRepo: repositories.NewPlatformAdminRepository(db),
		secret:    []byte(secret),
		ttl:       ttl,
	}
}

func (s *PlatformAdminService) Login(ctx context.Context, email, password string) (*models.PlatformAdmin, string, error) {
	admin, err := s.adminRepo.GetPlatformAdminByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	match, err := security.VerifyPassword(password, admin.PasswordHash)
	if err != nil || !match {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := security.CreateSessionToken(admin.ID, s.secret, s.ttl)
	if err != nil {
		return nil, "", fmt.Errorf("create session: %w", err)
	}

	return admin, token, nil
}

func (s *PlatformAdminService) SubmitRegistration(ctx context.Context, req *models.TenantRegistrationRequest, plainPassword string) error {
	if req == nil {
		return errors.New("registration data is required")
	}

	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.OwnerFirstName = strings.TrimSpace(req.OwnerFirstName)
	req.OwnerLastName = strings.TrimSpace(req.OwnerLastName)
	req.OwnerEmail = strings.ToLower(strings.TrimSpace(req.OwnerEmail))

	if req.BusinessName == "" {
		return errors.New("studio / business name is required")
	}
	if req.OwnerFirstName == "" || req.OwnerLastName == "" {
		return errors.New("owner name is required")
	}
	if !strings.Contains(req.OwnerEmail, "@") {
		return errors.New("valid owner email is required")
	}
	if len(plainPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hash, err := security.HashPassword(plainPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	req.PasswordHash = hash

	return s.adminRepo.CreateRegistrationRequest(ctx, req)
}

func slugify(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := re.ReplaceAllString(lower, "-")
	return strings.Trim(slug, "-")
}
