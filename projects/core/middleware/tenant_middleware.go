package middleware

import (
	"net/http"

	"core/models"
	"core/repositories"
	"core/utils"
	"github.com/google/uuid"
)

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
		tenantCookie, err := r.Cookie("active_tenant_id")
		if err == nil && tenantCookie.Value != "" {
			reqTenantID, parseErr := uuid.Parse(tenantCookie.Value)
			if parseErr == nil {
				for index := range memberships {
					if memberships[index].TenantID == reqTenantID {
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
		}

		if activeMembership == nil {
			activeMembership = &models.Membership{
				ID:                       memberships[0].ID,
				TenantID:                 memberships[0].TenantID,
				UserID:                   memberships[0].UserID,
				Role:                     memberships[0].Role,
				Status:                   memberships[0].Status,
				TenantName:               memberships[0].TenantName,
				TenantSlug:               memberships[0].TenantSlug,
				TenantBusinessType:       memberships[0].TenantBusinessType,
				TenantStatus:             memberships[0].TenantStatus,
				TenantSubscriptionStatus: memberships[0].TenantSubscriptionStatus,
			}
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
