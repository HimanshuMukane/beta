package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"core/repositories"
	"core/utils"
	"github.com/google/uuid"
)

type SwitchWorkspaceRequest struct {
	TenantID string `json:"tenantId"`
}

func SwitchTenantHandler(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	if tenantIDStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := utils.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	db := utils.IdentityDBFromContext(r.Context(), nil)
	if db != nil {
		tenantRepo := repositories.NewTenantRepository(db)
		m, err := tenantRepo.GetMembership(r.Context(), tenantID, user.ID)
		if err != nil || m == nil || m.Status != "active" {
			http.Error(w, "Unauthorized workspace access", http.StatusForbidden)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "active_tenant_id",
		Value:    tenantID.String(),
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   utils.AppConfig.SessionSecure,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SwitchWorkspaceAPIHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.UserFromContext(r.Context())
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": map[string]any{
					"code":    "AUTHENTICATION_ERROR",
					"message": "Authentication required",
				},
			})
			return
		}

		var req SwitchWorkspaceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TenantID == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": map[string]any{
					"code":    "VALIDATION_ERROR",
					"message": "Valid tenantId is required",
				},
			})
			return
		}

		tenantID, err := uuid.Parse(req.TenantID)
		if err != nil || tenantID == uuid.Nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": map[string]any{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid tenantId format",
				},
			})
			return
		}

		db := utils.IdentityDBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
				"error": map[string]any{
					"code":    "DATABASE_ERROR",
					"message": "Database unavailable",
				},
			})
			return
		}

		tenantRepo := repositories.NewTenantRepository(db)
		membership, err := tenantRepo.GetMembership(r.Context(), tenantID, user.ID)
		if err != nil || membership == nil || membership.Status != "active" {
			utils.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": map[string]any{
					"code":    "FORBIDDEN",
					"message": "You do not have active access to this workspace",
				},
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "active_tenant_id",
			Value:    tenantID.String(),
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   cfg.SessionSecure,
		})

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message":        "Workspace switched successfully",
			"activeTenantId": tenantID.String(),
			"workspace":      membership,
		})
	}
}

func WorkspacesAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.UserFromContext(r.Context())
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": map[string]any{
					"code":    "AUTHENTICATION_ERROR",
					"message": "Authentication required",
				},
			})
			return
		}

		db := utils.IdentityDBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
				"error": map[string]any{
					"code":    "DATABASE_ERROR",
					"message": "Database unavailable",
				},
			})
			return
		}

		tenantRepo := repositories.NewTenantRepository(db)
		memberships, err := tenantRepo.ListUserMemberships(r.Context(), user.ID)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to list workspaces",
				},
			})
			return
		}

		activeCtx := utils.TenantFromContext(r.Context())
		var activeTenantID string
		if activeCtx != nil && activeCtx.TenantID != uuid.Nil {
			activeTenantID = activeCtx.TenantID.String()
		}

		type WorkspaceItem struct {
			ID                 string `json:"id"`
			TenantID           string `json:"tenantId"`
			TenantName         string `json:"tenantName"`
			TenantSlug         string `json:"tenantSlug"`
			TenantBusinessType string `json:"tenantBusinessType"`
			Role               string `json:"role"`
			Status             string `json:"status"`
			SubscriptionStatus string `json:"subscriptionStatus"`
		}

		activeList := make([]WorkspaceItem, 0, len(memberships))
		for _, m := range memberships {
			if m.Status == "active" {
				activeList = append(activeList, WorkspaceItem{
					ID:                 m.ID.String(),
					TenantID:           m.TenantID.String(),
					TenantName:         m.TenantName,
					TenantSlug:         m.TenantSlug,
					TenantBusinessType: m.TenantBusinessType,
					Role:               m.Role,
					Status:             m.Status,
					SubscriptionStatus: m.TenantSubscriptionStatus,
				})
			}
		}

		if activeTenantID == "" && len(activeList) > 0 {
			activeTenantID = activeList[0].TenantID
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"activeTenantId": activeTenantID,
			"workspaces":     activeList,
		})
	}
}
