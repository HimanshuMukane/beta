package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SessionDTO struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	ProjectName         string   `json:"projectName"`
	SessionManager      string   `json:"sessionManager"`
	ProjectType         string   `json:"projectType"`
	ShootTiming         string   `json:"shootTiming"`
	Status              string   `json:"status"`
	Date                string   `json:"date"`
	StartTime           string   `json:"startTime,omitempty"`
	Duration            string   `json:"duration,omitempty"`
	TimeRange           string   `json:"timeRange"`
	Venue               string   `json:"venue"`
	City                string   `json:"city"`
	LocationLink        string   `json:"locationLink,omitempty"`
	Remarks             string   `json:"remarks,omitempty"`
	Members             []string `json:"members"`
	ShootsActivityCount int      `json:"shootsActivityCount"`
}

func SessionsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, title, project_name, session_manager, project_type, shoot_timing,
		       status, session_date, start_time, duration, time_range, venue, city,
		       location_link, remarks, members, shoots_activity_count
		FROM project_sessions
		WHERE tenant_id = ?
		ORDER BY session_date DESC, start_time ASC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []SessionDTO
	for rows.Next() {
		var id, title, projName, mgr, pType, timing, status, date, start, dur, tRange, venue, city, link, remarks string
		var activity int
		var membersRaw []byte
		var (
			sStart, sDur, sRange, sVenue, sCity, sLink, sRemarks sql.NullString
		)

		err := rows.Scan(
			&id, &title, &projName, &mgr, &pType, &timing,
			&status, &date, &sStart, &sDur, &sRange, &sVenue, &sCity,
			&sLink, &sRemarks, &membersRaw, &activity,
		)
		if err != nil {
			continue
		}

		if sStart.Valid {
			start = sStart.String
		}
		if sDur.Valid {
			dur = sDur.String
		}
		if sRange.Valid {
			tRange = sRange.String
		}
		if sVenue.Valid {
			venue = sVenue.String
		}
		if sCity.Valid {
			city = sCity.String
		}
		if sLink.Valid {
			link = sLink.String
		}
		if sRemarks.Valid {
			remarks = sRemarks.String
		}

		var members []string = []string{}
		parseJSONColumn(membersRaw, &members)

		result = append(result, SessionDTO{
			ID:                  id,
			Title:               title,
			ProjectName:         projName,
			SessionManager:      mgr,
			ProjectType:         pType,
			ShootTiming:         timing,
			Status:              status,
			Date:                date,
			StartTime:           start,
			Duration:            dur,
			TimeRange:           tRange,
			Venue:               venue,
			City:                city,
			LocationLink:        link,
			Remarks:             remarks,
			Members:             members,
			ShootsActivityCount: activity,
		})
	}

	if result == nil {
		result = []SessionDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var payload SessionDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if payload.ID == "" {
		payload.ID = "ses-" + uuid.New().String()[:8]
	}
	if payload.Status == "" {
		payload.Status = "Scheduled"
	}
	if payload.ShootTiming == "" {
		payload.ShootTiming = "Morning"
	}
	if payload.SessionManager == "" {
		payload.SessionManager = "especa"
	}

	membersJSON := toJSONBytes(payload.Members)

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO project_sessions (
			id, tenant_id, project_name, title, session_manager, project_type,
			shoot_timing, status, session_date, start_time, duration, time_range,
			venue, city, location_link, remarks, members, shoots_activity_count
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?
		)
	`, payload.ID, tid, payload.ProjectName, payload.Title, payload.SessionManager, payload.ProjectType,
		payload.ShootTiming, payload.Status, payload.Date, payload.StartTime, payload.Duration, payload.TimeRange,
		payload.Venue, payload.City, payload.LocationLink, payload.Remarks, membersJSON, payload.ShootsActivityCount)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, payload)
}

func UpdateSessionHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var payload SessionDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	membersJSON := toJSONBytes(payload.Members)

	_, err := db.ExecContext(r.Context(), `
		UPDATE project_sessions SET
			project_name = ?, title = ?, session_manager = ?, project_type = ?,
			shoot_timing = ?, status = ?, session_date = ?, start_time = ?,
			duration = ?, time_range = ?, venue = ?, city = ?, location_link = ?,
			remarks = ?, members = ?, shoots_activity_count = ?
		WHERE id = ? AND tenant_id = ?
	`, payload.ProjectName, payload.Title, payload.SessionManager, payload.ProjectType,
		payload.ShootTiming, payload.Status, payload.Date, payload.StartTime,
		payload.Duration, payload.TimeRange, payload.Venue, payload.City, payload.LocationLink,
		payload.Remarks, membersJSON, payload.ShootsActivityCount,
		id, tid)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.ID = id
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err := db.ExecContext(r.Context(), "UPDATE project_sessions SET status = ? WHERE id = ? AND tenant_id = ?", req.Status, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": id, "sessionStatus": req.Status})
}

func DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(), "DELETE FROM project_sessions WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

