package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"core/models"
	"core/repositories"
	"core/security"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		return nil, "", fmt.Errorf("create admin session: %w", err)
	}

	admin.PasswordHash = ""
	return admin, token, nil
}

func (s *PlatformAdminService) SubmitRegistration(ctx context.Context, req *models.TenantRegistrationRequest, plainPassword string) error {
	if strings.TrimSpace(req.BusinessName) == "" {
		return errors.New("business name is required")
	}
	if strings.TrimSpace(req.OwnerEmail) == "" {
		return errors.New("owner email is required")
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

func (s *PlatformAdminService) ApproveRegistration(ctx context.Context, requestID uuid.UUID, reviewerID uuid.UUID) (*models.Tenant, error) {
	req, err := s.adminRepo.GetRegistrationRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req.Status != "pending" {
		return nil, ErrRequestAlreadyProcessed
	}

	slug := slugify(req.BusinessName)
	tenantRepo := repositories.NewTenantRepository(s.db)
	userRepo := repositories.NewUserRepository(s.db)
	roleRepo := repositories.NewRoleRepository(s.db)

	ownerRole, err := roleRepo.GetRoleByName(ctx, "owner")
	if err != nil {
		ownerRole, err = roleRepo.CreateRole(ctx, "owner")
		if err != nil {
			return nil, fmt.Errorf("ensure owner role: %w", err)
		}
	}

	tenant := &models.Tenant{
		Name:         req.BusinessName,
		Slug:         slug,
		BusinessType: req.BusinessType,
		Status:       "active",
	}
	if err := tenantRepo.CreateTenant(ctx, tenant); err != nil {
		tenant.Slug = fmt.Sprintf("%s-%s", slug, uuid.New().String()[:6])
		if err := tenantRepo.CreateTenant(ctx, tenant); err != nil {
			return nil, fmt.Errorf("create tenant: %w", err)
		}
	}

	ownerUser, err := userRepo.GetUserByEmail(ctx, req.OwnerEmail)
	if err != nil {
		ownerUser = &models.User{
			FirstName:    req.OwnerFirstName,
			LastName:     req.OwnerLastName,
			Email:        req.OwnerEmail,
			PasswordHash: req.PasswordHash,
			RoleID:       ownerRole.ID,
		}
		if err := userRepo.CreateUser(ctx, ownerUser); err != nil {
			return nil, fmt.Errorf("create owner user: %w", err)
		}
	}

	membership := &models.Membership{
		TenantID: tenant.ID,
		UserID:   ownerUser.ID,
		Role:     "owner",
		Status:   "active",
	}
	if err := tenantRepo.CreateMembership(ctx, membership); err != nil {
		return nil, fmt.Errorf("create owner membership: %w", err)
	}

	initTenantWorkflows(ctx, s.db, tenant.ID)

	if err := s.adminRepo.UpdateRegistrationStatus(ctx, requestID, "approved", reviewerID); err != nil {
		return nil, fmt.Errorf("update registration status: %w", err)
	}

	return tenant, nil
}

func (s *PlatformAdminService) RejectRegistration(ctx context.Context, requestID uuid.UUID, reviewerID uuid.UUID) error {
	req, err := s.adminRepo.GetRegistrationRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if req.Status != "pending" {
		return ErrRequestAlreadyProcessed
	}
	return s.adminRepo.UpdateRegistrationStatus(ctx, requestID, "rejected", reviewerID)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func initTenantWorkflows(ctx context.Context, db *sqlx.DB, tenantID uuid.UUID) {
	wfID := uuid.New()
	_, _ = db.ExecContext(ctx, "INSERT INTO lead_workflows (id, tenant_id, name, is_default) VALUES (?, ?, ?, 1)", wfID.String(), tenantID.String(), "Standard Photography Pipeline")

	stages := []struct {
		name  string
		color string
		order int
		won   int
		lost  int
	}{
		{"New Inquiry", "#0984e3", 1, 0, 0},
		{"Contacted", "#fdcb6e", 2, 0, 0},
		{"Qualified", "#7c5dfa", 3, 0, 0},
		{"Negotiation", "#e17055", 4, 0, 0},
		{"Ready to Book", "#00b894", 5, 0, 0},
		{"Converted", "#00cec9", 6, 1, 0},
		{"Lost", "#d63031", 7, 0, 1},
	}

	for _, st := range stages {
		_, _ = db.ExecContext(ctx, `
			INSERT INTO lead_stages (id, tenant_id, workflow_id, name, color, stage_order, is_won, is_lost)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, uuid.New().String(), tenantID.String(), wfID.String(), st.name, st.color, st.order, st.won, st.lost)
	}
}

