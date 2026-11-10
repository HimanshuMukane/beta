package middleware

import (
	"net/http"
	"strings"

	"core/models"
	"core/repositories"
	"core/utils"
	"github.com/google/uuid"
)

func resolveActiveTenant(r *http.Request) (uuid.UUID, bool) {
	if headerVal := strings.TrimSpace(r.Header.Get("X-Tenant-ID")); headerVal != "" {
		if id, err := uuid.Parse(headerVal); err == nil && id != uuid.Nil {
			return id, true
		}
	}
	if cookie, err := r.Cookie("active_tenant_id"); err == nil && cookie.Value != "" {
		if id, err := uuid.Parse(cookie.Value); err == nil && id != uuid.Nil {
			return id, true
		}
	}
	return uuid.Nil, false
}

func TenantContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.UserFromContext(r.Context())
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}

		db := utils.IdentityDBFromContext(r.Context(), nil)
		if db == nil {
			next.ServeHTTP(w, r)
			return
		}

		tenantRepo := repositories.NewTenantRepository(db)
		memberships, err := tenantRepo.ListUserMemberships(r.Context(), user.ID)
		if err != nil || len(memberships) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		var activeMembership *models.Membership = nil
		if reqTenantID, ok := resolveActiveTenant(r); ok {
			for index := range memberships {
				if memberships[index].TenantID == reqTenantID && memberships[index].Status == "active" {
					activeMembership = &models.Membership{
						ID:                       memberships[index].ID,
						TenantID:                 memberships[index].TenantID,
						UserID:                   memberships[index].UserID,
						Role:                     memberships[index].Role,
						Status:                   memberships[index].Status,
						TenantName:               memberships[index].TenantName,
						TenantSlug:               memberships[index].TenantSlug,
						TenantBusinessType:       memberships[index].TenantBusinessType,
						TenantStatus:             memberships[index].TenantStatus,
						TenantSubscriptionStatus: memberships[index].TenantSubscriptionStatus,
					}
					break
				}
			}
		}

		if activeMembership == nil {
			for index := range memberships {
				if memberships[index].Status == "active" {
					activeMembership = &models.Membership{
						ID:                       memberships[index].ID,
						TenantID:                 memberships[index].TenantID,
						UserID:                   memberships[index].UserID,
						Role:                     memberships[index].Role,
						Status:                   memberships[index].Status,
						TenantName:               memberships[index].TenantName,
						TenantSlug:               memberships[index].TenantSlug,
						TenantBusinessType:       memberships[index].TenantBusinessType,
						TenantStatus:             memberships[index].TenantStatus,
						TenantSubscriptionStatus: memberships[index].TenantSubscriptionStatus,
					}
					break
				}
			}
		}

		if activeMembership == nil {
			next.ServeHTTP(w, r)
			return
		}

		isReadOnly := activeMembership.TenantSubscriptionStatus == "expired"

		tenantContext := &utils.TenantContext{
			TenantID:           activeMembership.TenantID,
			TenantName:         activeMembership.TenantName,
			TenantSlug:         activeMembership.TenantSlug,
			Role:               activeMembership.Role,
			SubscriptionStatus: activeMembership.TenantSubscriptionStatus,
			IsReadOnly:         isReadOnly,
		}

		ctx := utils.WithTenantContext(r.Context(), tenantContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireTenantContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tc := utils.TenantFromContext(r.Context())
		if tc != nil && tc.TenantID != uuid.Nil {
			next.ServeHTTP(w, r)
			return
		}

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
		if err != nil || len(memberships) == 0 {
			utils.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": map[string]any{
					"code":    "FORBIDDEN",
					"message": "Valid workspace membership required",
				},
			})
			return
		}

		var activeMembership *models.Membership = nil
		if reqTenantID, ok := resolveActiveTenant(r); ok {
			for index := range memberships {
				if memberships[index].TenantID == reqTenantID && memberships[index].Status == "active" {
					activeMembership = &models.Membership{
						ID:                       memberships[index].ID,
						TenantID:                 memberships[index].TenantID,
						UserID:                   memberships[index].UserID,
						Role:                     memberships[index].Role,
						Status:                   memberships[index].Status,
						TenantName:               memberships[index].TenantName,
						TenantSlug:               memberships[index].TenantSlug,
						TenantBusinessType:       memberships[index].TenantBusinessType,
						TenantStatus:             memberships[index].TenantStatus,
						TenantSubscriptionStatus: memberships[index].TenantSubscriptionStatus,
					}
					break
				}
			}
		}

		if activeMembership == nil {
			for index := range memberships {
				if memberships[index].Status == "active" {
					activeMembership = &models.Membership{
						ID:                       memberships[index].ID,
						TenantID:                 memberships[index].TenantID,
						UserID:                   memberships[index].UserID,
						Role:                     memberships[index].Role,
						Status:                   memberships[index].Status,
						TenantName:               memberships[index].TenantName,
						TenantSlug:               memberships[index].TenantSlug,
						TenantBusinessType:       memberships[index].TenantBusinessType,
						TenantStatus:             memberships[index].TenantStatus,
						TenantSubscriptionStatus: memberships[index].TenantSubscriptionStatus,
					}
					break
				}
			}
		}

		if activeMembership == nil {
			utils.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": map[string]any{
					"code":    "FORBIDDEN",
					"message": "Valid workspace membership required",
				},
			})
			return
		}

		tenantContext := &utils.TenantContext{
			TenantID:           activeMembership.TenantID,
			TenantName:         activeMembership.TenantName,
			TenantSlug:         activeMembership.TenantSlug,
			Role:               activeMembership.Role,
			SubscriptionStatus: activeMembership.TenantSubscriptionStatus,
			IsReadOnly:         activeMembership.TenantSubscriptionStatus == "expired",
		}

		ctx := utils.WithTenantContext(r.Context(), tenantContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
