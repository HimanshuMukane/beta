package handlers

import (
	"net/http"
	"time"

	"core/utils"
)

type tenantRequestRecord struct {
	ID          string    `db:"id" json:"id"`
	CompanyName string    `db:"company_name" json:"company_name"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := utils.DBFromContext(ctx)

	var reqs []tenantRequestRecord
	if db != nil {
		query := `SELECT id, company_name, status, created_at FROM tenant_registration_requests ORDER BY created_at DESC LIMIT 50`
		_ = db.SelectContext(ctx, &reqs, query)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"tenants": reqs,
	})
}
