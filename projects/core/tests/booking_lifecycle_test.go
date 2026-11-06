package tests

import (
	"context"
	"testing"

	"core/models"
	"core/repositories"
	"core/services"
)

func TestProjectEventsSessionsAndBookingLifecycle(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	tenantRepo := repositories.NewTenantRepository(db)
	crmRepo := repositories.NewCRMRepository(db)
	crmService := services.NewCRMService(db)

	tenant := &models.Tenant{
		Name:         "Fine Art Studio",
		Slug:         "fine-art-studio",
		BusinessType: "wedding",
	}
	if err := tenantRepo.CreateTenant(ctx, tenant); err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	client := &repositories.CreateClientParam{
		TenantID:   tenant.ID,
		FirstName:  "Kavya",
		LastName:   "Rajendran",
		Email:      "kavya@example.com",
		Phone:      "+919876543210",
		Address:    "Chennai, Tamil Nadu",
		LeadSource: "instagram",
	}
	if err := crmRepo.CreateClient(ctx, client); err != nil {
		t.Fatalf("create client: %v", err)
	}

	project := &repositories.CreateProjectParam{
		TenantID:    tenant.ID,
		ClientID:    client.ID,
		Title:       "Kavya Rajendran Wedding",
		ProjectType: "wedding",
		Status:      "New",
		Stage:       "Booked",
		Address:     "Grand Palace Resort, Mahabalipuram",
	}
	if err := crmRepo.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}

	evalBefore, err := crmService.EvaluateBookingConditions(ctx, tenant.ID, project.ID)
	if err != nil {
		t.Fatalf("evaluate conditions error: %v", err)
	}
	if evalBefore.CanConfirm {
		t.Fatal("booking should not be confirmable before segments/dates are scheduled")
	}

	segment := &repositories.CreateSegmentParam{
		TenantID:     tenant.ID,
		ProjectID:    project.ID,
		SegmentType:  "wedding",
		Name:         "Muhurtham & Reception",
		StartDate:    "2027-02-14",
		EndDate:      "2027-02-15",
		SegmentOrder: 1,
	}
	if err := crmRepo.CreateProjectSegment(ctx, segment); err != nil {
		t.Fatalf("create segment: %v", err)
	}

	session1 := &repositories.CreateSessionParam{
		TenantID:    tenant.ID,
		ProjectID:   project.ID,
		SegmentID:   &segment.ID,
		Title:       "Morning Muhurtham Shoot",
		SessionType: "ceremony",
		StartTime:   "2027-02-14 06:00:00",
		EndTime:     "2027-02-14 10:00:00",
		Location:    "Temple Hall",
	}
	if err := crmRepo.CreateSession(ctx, session1); err != nil {
		t.Fatalf("create session 1: %v", err)
	}

	session2 := &repositories.CreateSessionParam{
		TenantID:    tenant.ID,
		ProjectID:   project.ID,
		SegmentID:   &segment.ID,
		Title:       "Evening Grand Reception",
		SessionType: "reception",
		StartTime:   "2027-02-14 18:30:00",
		EndTime:     "2027-02-14 23:00:00",
		Location:    "Beach Lawn Ballroom",
	}
	if err := crmRepo.CreateSession(ctx, session2); err != nil {
		t.Fatalf("create session 2: %v", err)
	}

	sessions, err := crmRepo.ListSessions(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}

	evalAfter, err := crmService.EvaluateBookingConditions(ctx, tenant.ID, project.ID)
	if err != nil {
		t.Fatalf("evaluate conditions error: %v", err)
	}
	if !evalAfter.CanConfirm {
		t.Fatalf("expected booking to be confirmable after segments added, reason: %s", evalAfter.Reason)
	}

	confirmedProject, err := crmService.ConfirmBooking(ctx, tenant.ID, project.ID)
	if err != nil {
		t.Fatalf("confirm booking failed: %v", err)
	}
	if confirmedProject.IsConfirmed != 1 {
		t.Errorf("expected IsConfirmed == 1, got %d", confirmedProject.IsConfirmed)
	}
	if confirmedProject.Status != "Booked" {
		t.Errorf("expected Status == 'Booked', got %s", confirmedProject.Status)
	}
	if !confirmedProject.ConfirmedAt.Valid || confirmedProject.ConfirmedAt.String == "" {
		t.Error("expected ConfirmedAt timestamp to be populated")
	}
}
