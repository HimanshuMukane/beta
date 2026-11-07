package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"cms/models"
	"cms/repositories"
	"core/utils"
	"github.com/go-chi/chi/v5"
)

type APIContextKey string

const (
	APIKeyTenantIDKey APIContextKey = "api_tenant_id"
	APIKeyScopesKey   APIContextKey = "api_scopes"
)

// DeveloperAPIKeyMiddleware authenticates Developer API requests via Bearer token
func DeveloperAPIKeyMiddleware(repo *repositories.CMSRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := ""
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else if k := r.Header.Get("X-CMS-Key"); k != "" {
				token = k
			}

			if token == "" {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing api key: provide via Authorization: Bearer <token> or X-CMS-Key header"})
				return
			}

			apiKey, err := repo.VerifyAPIKey(r.Context(), token)
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or revoked api key"})
				return
			}

			ctx := context.WithValue(r.Context(), APIKeyTenantIDKey, apiKey.TenantID)
			ctx = context.WithValue(ctx, APIKeyScopesKey, apiKey.Scopes)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireFeatureMiddleware checks whether the tenant has the given feature licensed by Superadmin
func RequireFeatureMiddleware(repo *repositories.CMSRepository, featureKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID, ok := r.Context().Value(APIKeyTenantIDKey).(string)
			if !ok || tenantID == "" {
				// Try from tenant context
				if tc := utils.TenantFromContext(r.Context()); tc != nil {
					tenantID = tc.TenantID.String()
				}
			}

			if tenantID == "" {
				utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "tenant context missing"})
				return
			}

			enabled, err := repo.IsFeatureEnabled(r.Context(), tenantID, featureKey)
			if err != nil || !enabled {
				utils.WriteJSON(w, http.StatusForbidden, map[string]string{
					"error":   "feature not licensed for this tenant",
					"feature": featureKey,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// =========================================================================
// DEVELOPER CONTENT DELIVERY API HANDLERS
// =========================================================================

func GetContentEntriesHandler(repo *repositories.CMSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value(APIKeyTenantIDKey).(string)
		typeSlug := chi.URLParam(r, "type_slug")

		ct, err := repo.GetContentTypeBySlug(r.Context(), tenantID, typeSlug)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "content type not found"})
			return
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		entries, err := repo.ListPublishedEntries(r.Context(), tenantID, ct.ID, limit, offset)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list entries"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"type":    ct.Name,
			"slug":    ct.Slug,
			"page":    page,
			"limit":   limit,
			"entries": entries,
		})
	}
}

func GetContentEntryBySlugHandler(repo *repositories.CMSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value(APIKeyTenantIDKey).(string)
		typeSlug := chi.URLParam(r, "type_slug")
		entrySlug := chi.URLParam(r, "entry_slug")

		ct, err := repo.GetContentTypeBySlug(r.Context(), tenantID, typeSlug)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "content type not found"})
			return
		}

		entry, err := repo.GetEntryBySlug(r.Context(), tenantID, ct.ID, entrySlug)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "entry not found"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, entry)
	}
}

func GetMenuHandler(repo *repositories.CMSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value(APIKeyTenantIDKey).(string)
		menuKey := chi.URLParam(r, "menu_key")

		menu, err := repo.GetMenuByKey(r.Context(), tenantID, menuKey)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "menu not found"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, menu)
	}
}

func SubmitFormHandler(repo *repositories.CMSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value(APIKeyTenantIDKey).(string)
		formKey := chi.URLParam(r, "form_key")

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
			return
		}

		rawJSON, _ := json.Marshal(payload)
		sub := &models.FormSubmission{
			TenantID: tenantID,
			FormKey:  formKey,
			Data:     rawJSON,
		}

		if err := repo.SaveFormSubmission(r.Context(), sub); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save form submission"})
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]string{"status": "received", "id": sub.ID})
	}
}
