package handlers

import (
	"net/http"
	"time"

	"core/utils"
	"photography_crm/ui/pages"
)

type followupRecord struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	DueDate   time.Time `db:"due_date"`
	Status    string    `db:"status"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := utils.DBFromContext(ctx)

	var rows []pages.TaskRow
	if db != nil {
		var records []followupRecord
		query := `SELECT id, title, due_date, status FROM followups ORDER BY due_date ASC LIMIT 50`
		if err := db.SelectContext(ctx, &records, query); err == nil {
			for i, rec := range records {
				statusType := "todo"
				if rec.Status == "completed" {
					statusType = "done"
				} else if rec.Status == "in_progress" {
					statusType = "inprogress"
				}

				rows = append(rows, pages.TaskRow{
					Index:      i + 1,
					Name:       rec.Title,
					Status:     rec.Status,
					StatusType: statusType,
					DueDate:    rec.DueDate.Format("Jan 02, 2006"),
				})
			}
		}
	}

	component := pages.Dashboard("CRM Tasks & Pipeline", rows)
	component.Render(ctx, w)
}
