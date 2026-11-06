package handlers

import (
	"net/http"
	"time"

	"core/repositories"
	"core/utils"
	"github.com/google/uuid"
)

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

	db := utils.DBFromContext(r.Context(), nil)
	if db != nil {
		tenantRepo := repositories.NewTenantRepository(db)
		_, err := tenantRepo.GetMembership(r.Context(), tenantID, user.ID)
		if err != nil {
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
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func WorkspacesAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.UserFromContext(r.Context())
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
			return
		}

		tenantRepo := repositories.NewTenantRepository(db)
		memberships, err := tenantRepo.ListUserMemberships(r.Context(), user.ID)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		activeCtx := utils.TenantFromContext(r.Context())
		var activeTenantID string
		if activeCtx != nil {
			activeTenantID = activeCtx.TenantID.String()
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"active_tenant_id": activeTenantID,
			"workspaces":       memberships,
		})
	}
}

