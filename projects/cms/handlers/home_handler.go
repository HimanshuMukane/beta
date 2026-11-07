package handlers

import (
	"cms/ui/pages"
	"net/http"
	"time"

	"core/utils"
)

type articleRecord struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := utils.DBFromContext(ctx)

	var rows []pages.TaskRow
	if db != nil {
		var articles []articleRecord
		query := `SELECT id, title, status, created_at FROM cms_articles ORDER BY created_at DESC LIMIT 50`
		if err := db.SelectContext(ctx, &articles, query); err == nil {
			for i, art := range articles {
				statusType := "todo"
				if art.Status == "published" {
					statusType = "done"
				} else if art.Status == "review" {
					statusType = "inprogress"
				}

				rows = append(rows, pages.TaskRow{
					Index:      i + 1,
					Name:       art.Title,
					Status:     art.Status,
					StatusType: statusType,
					DueDate:    art.CreatedAt.Format("Jan 02, 2006"),
				})
			}
		}
	}

	component := pages.Dashboard("CMS Articles & Content", rows)
	component.Render(ctx, w)
}
