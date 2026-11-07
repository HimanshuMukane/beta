package handlers

import (
	"net/http"
	"superadmin/ui/pages"
	"time"

	"core/utils"
)

type tenantRequestRecord struct {
	ID          string    `db:"id"`
	CompanyName string    `db:"company_name"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := utils.DBFromContext(ctx)

	var rows []pages.TaskRow
	if db != nil {
		var reqs []tenantRequestRecord
		query := `SELECT id, company_name, status, created_at FROM tenant_registration_requests ORDER BY created_at DESC LIMIT 50`
		if err := db.SelectContext(ctx, &reqs, query); err == nil {
			for i, tReq := range reqs {
				statusType := "todo"
				if tReq.Status == "approved" {
					statusType = "done"
				} else if tReq.Status == "in_review" {
					statusType = "inprogress"
				}

				rows = append(rows, pages.TaskRow{
					Index:      i + 1,
					Name:       "Tenant Request: " + tReq.CompanyName,
					Status:     tReq.Status,
					StatusType: statusType,
					DueDate:    tReq.CreatedAt.Format("Jan 02, 2006"),
				})
			}
		}
	}

	component := pages.Dashboard("Superadmin Control Center", rows)
	component.Render(ctx, w)
}
