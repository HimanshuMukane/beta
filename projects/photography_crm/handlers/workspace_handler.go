package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"core/models"
	"core/repositories"
	"core/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SwitchWorkspaceRequest struct {
	TenantID string `json:"tenantId"`
}

type SwitchWorkspaceSuccessResponse struct {
	Message        string             `json:"message"`
	ActiveTenantID string             `json:"activeTenantId"`
	Workspace      *models.Membership `json:"workspace"`
}

type WorkspaceItemDTO struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenantId"`
	TenantName         string `json:"tenantName"`
	TenantSlug         string `json:"tenantSlug"`
	TenantBusinessType string `json:"tenantBusinessType"`
	Role               string `json:"role"`
	Status             string `json:"status"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}

type WorkspacesListResponse struct {
	ActiveTenantID string             `json:"activeTenantId"`
	Workspaces     []WorkspaceItemDTO `json:"workspaces"`
}

func SwitchTenantHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var tenantIDString string = request.URL.Query().Get("tenant_id")
	if tenantIDString == "" {
		http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
		return
	}

	var parsedTenantID uuid.UUID
	var parseErr error
	parsedTenantID, parseErr = uuid.Parse(tenantIDString)
	if parseErr != nil {
		http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
		return
	}

	var authenticatedUser *models.User = utils.UserFromContext(request.Context())
	if authenticatedUser == nil {
		http.Redirect(responseWriter, request, "/signin", http.StatusSeeOther)
		return
	}

	var identityDatabase *sqlx.DB = utils.IdentityDBFromContext(request.Context(), nil)
	if identityDatabase != nil {
		var tenantRepo *repositories.TenantRepository = repositories.NewTenantRepository(identityDatabase)
		var userMembership *models.Membership
		var membershipErr error
		userMembership, membershipErr = tenantRepo.GetMembership(request.Context(), parsedTenantID, authenticatedUser.ID)
		if membershipErr != nil || userMembership == nil || userMembership.Status != "active" {
			http.Error(responseWriter, "Unauthorized workspace access", http.StatusForbidden)
			return
		}
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "active_tenant_id",
		Value:    parsedTenantID.String(),
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   utils.AppConfig.SessionSecure,
	})

	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}

func SwitchWorkspaceAPIHandler(configuration utils.Config) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var authenticatedUser *models.User = utils.UserFromContext(request.Context())
		if authenticatedUser == nil {
			utils.WriteJSON(responseWriter, http.StatusUnauthorized, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "AUTHENTICATION_ERROR",
					Message: "Authentication required",
				},
			})
			return
		}

		var switchPayload SwitchWorkspaceRequest
		var decodeErr error = json.NewDecoder(request.Body).Decode(&switchPayload)
		if decodeErr != nil || switchPayload.TenantID == "" {
			utils.WriteJSON(responseWriter, http.StatusBadRequest, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "VALIDATION_ERROR",
					Message: "Valid tenantId is required",
				},
			})
			return
		}

		var parsedTenantID uuid.UUID
		var parseErr error
		parsedTenantID, parseErr = uuid.Parse(switchPayload.TenantID)
		if parseErr != nil || parsedTenantID == uuid.Nil {
			utils.WriteJSON(responseWriter, http.StatusBadRequest, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "VALIDATION_ERROR",
					Message: "Invalid tenantId format",
				},
			})
			return
		}

		var identityDatabase *sqlx.DB = utils.IdentityDBFromContext(request.Context(), nil)
		if identityDatabase == nil {
			utils.WriteJSON(responseWriter, http.StatusServiceUnavailable, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "DATABASE_ERROR",
					Message: "Database unavailable",
				},
			})
			return
		}

		var tenantRepo *repositories.TenantRepository = repositories.NewTenantRepository(identityDatabase)
		var userMembership *models.Membership
		var membershipErr error
		userMembership, membershipErr = tenantRepo.GetMembership(request.Context(), parsedTenantID, authenticatedUser.ID)
		if membershipErr != nil || userMembership == nil || userMembership.Status != "active" {
			utils.WriteJSON(responseWriter, http.StatusForbidden, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "FORBIDDEN",
					Message: "You do not have active access to this workspace",
				},
			})
			return
		}

		http.SetCookie(responseWriter, &http.Cookie{
			Name:     "active_tenant_id",
			Value:    parsedTenantID.String(),
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   configuration.SessionSecure,
		})

		utils.WriteJSON(responseWriter, http.StatusOK, SwitchWorkspaceSuccessResponse{
			Message:        "Workspace switched successfully",
			ActiveTenantID: parsedTenantID.String(),
			Workspace:      userMembership,
		})
	}
}

func WorkspacesAPIHandler() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var authenticatedUser *models.User = utils.UserFromContext(request.Context())
		if authenticatedUser == nil {
			utils.WriteJSON(responseWriter, http.StatusUnauthorized, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "AUTHENTICATION_ERROR",
					Message: "Authentication required",
				},
			})
			return
		}

		var identityDatabase *sqlx.DB = utils.IdentityDBFromContext(request.Context(), nil)
		if identityDatabase == nil {
			utils.WriteJSON(responseWriter, http.StatusServiceUnavailable, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "DATABASE_ERROR",
					Message: "Database unavailable",
				},
			})
			return
		}

		var tenantRepo *repositories.TenantRepository = repositories.NewTenantRepository(identityDatabase)
		var userMemberships []models.Membership
		var membershipsErr error
		userMemberships, membershipsErr = tenantRepo.ListUserMemberships(request.Context(), authenticatedUser.ID)
		if membershipsErr != nil {
			utils.WriteJSON(responseWriter, http.StatusInternalServerError, ErrorResponseEnvelope{
				Error: DetailedError{
					Code:    "INTERNAL_ERROR",
					Message: "Failed to list workspaces",
				},
			})
			return
		}

		var activeContextTenant *utils.TenantContext = utils.TenantFromContext(request.Context())
		var activeTenantID string = ""
		if activeContextTenant != nil && activeContextTenant.TenantID != uuid.Nil {
			activeTenantID = activeContextTenant.TenantID.String()
		}

		var activeList []WorkspaceItemDTO = make([]WorkspaceItemDTO, 0, len(userMemberships))
		var memberIndex int
		for memberIndex = 0; memberIndex < len(userMemberships); memberIndex++ {
			var memberItem models.Membership = userMemberships[memberIndex]
			if memberItem.Status == "active" {
				activeList = append(activeList, WorkspaceItemDTO{
					ID:                 memberItem.ID.String(),
					TenantID:           memberItem.TenantID.String(),
					TenantName:         memberItem.TenantName,
					TenantSlug:         memberItem.TenantSlug,
					TenantBusinessType: memberItem.TenantBusinessType,
					Role:               memberItem.Role,
					Status:             memberItem.Status,
					SubscriptionStatus: memberItem.TenantSubscriptionStatus,
				})
			}
		}

		if activeTenantID == "" && len(activeList) > 0 {
			activeTenantID = activeList[0].TenantID
		}

		utils.WriteJSON(responseWriter, http.StatusOK, WorkspacesListResponse{
			ActiveTenantID: activeTenantID,
			Workspaces:     activeList,
		})
	}
}
