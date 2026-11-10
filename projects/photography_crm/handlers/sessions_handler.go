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
	ProjectID           string   `json:"projectId,omitempty"`
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
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}
	projectIDFilter := r.URL.Query().Get("projectId")

	query := `
		SELECT s.id, COALESCE(s.project_id, ''), s.title, COALESCE(p.project_name, s.project_name, ''),
		       s.session_manager, s.project_type, s.shoot_timing, s.status, s.date,
		       s.start_time, s.duration, s.time_range, s.venue, s.city, s.location_link,
		       s.remarks, s.members, s.shoots_activity_count
		FROM project_sessions s
		LEFT JOIN projects p ON p.id = s.project_id AND p.deleted_at IS NULL
		WHERE s.tenant_id = ? AND s.deleted_at IS NULL
	`
	args := []any{tid}
	if projectIDFilter != "" {
		query += " AND s.project_id = ?"
		args = append(args, projectIDFilter)
	}
	query += " ORDER BY s.date ASC"

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to list project sessions", err)
		return
	}
	defer rows.Close()

	result := make([]SessionDTO, 0)
	for rows.Next() {
		var (
			id, projID, title, projName, mgr, pType, timing, status, date, timeRange, venue, city string
			shootsCount                                                                           int
			membersRaw                                                                            []byte
		)
		var (
			sStartTime, sDuration, sLoc, sRemarks sql.NullString
		)

		err := rows.Scan(
			&id, &projID, &title, &projName, &mgr, &pType, &timing, &status, &date,
			&sStartTime, &sDuration, &timeRange, &venue, &city, &sLoc, &sRemarks,
			&membersRaw, &shootsCount,
		)
		if err != nil {
			continue
		}

		startTime := ""
		if sStartTime.Valid {
			startTime = sStartTime.String
		}
		duration := ""
		if sDuration.Valid {
			duration = sDuration.String
		}
		locationLink := ""
		if sLoc.Valid {
			locationLink = sLoc.String
		}
		remarks := ""
		if sRemarks.Valid {
			remarks = sRemarks.String
		}

		members := make([]string, 0)
		if len(membersRaw) > 0 {
			_ = json.Unmarshal(membersRaw, &members)
		}

		dto := SessionDTO{
			ID:                  id,
			ProjectID:           projID,
			Title:               title,
			ProjectName:         projName,
			SessionManager:      mgr,
			ProjectType:         pType,
			ShootTiming:         timing,
			Status:              status,
			Date:                date,
			StartTime:           startTime,
			Duration:            duration,
			TimeRange:           timeRange,
			Venue:               venue,
			City:                city,
			LocationLink:        locationLink,
			Remarks:             remarks,
			Members:             members,
			ShootsActivityCount: shootsCount,
		}
		result = append(result, dto)
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	var payload SessionDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if payload.ID == "" {
		payload.ID = "SES-" + uuid.New().String()[:8]
	}
	if payload.ShootsActivityCount == 0 {
		payload.ShootsActivityCount = 1
	}
	if payload.SessionManager == "" {
		payload.SessionManager = "especa"
	}
	if payload.Status == "" {
		payload.Status = "New"
	}

	membersJSON, _ := json.Marshal(payload.Members)

	var projIDVal *string
	if payload.ProjectID != "" {
		projIDVal = &payload.ProjectID
	} else if payload.ProjectName != "" {
		var matchedID string
		if err := db.QueryRowContext(r.Context(), `
			SELECT id FROM projects WHERE tenant_id = ? AND project_name = ? AND deleted_at IS NULL LIMIT 1
		`, tid, payload.ProjectName).Scan(&matchedID); err == nil && matchedID != "" {
			projIDVal = &matchedID
		}
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO project_sessions (
			id, tenant_id, project_id, title, project_name, session_manager,
			project_type, shoot_timing, status, date, start_time, duration,
			time_range, venue, city, location_link, remarks, members, shoots_activity_count
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?
		)
	`,
		payload.ID, tid, projIDVal, payload.Title, payload.ProjectName, payload.SessionManager,
		payload.ProjectType, payload.ShootTiming, payload.Status, payload.Date, payload.StartTime, payload.Duration,
		payload.TimeRange, payload.Venue, payload.City, payload.LocationLink, payload.Remarks, membersJSON, payload.ShootsActivityCount,
	)

	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to create session", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, payload)
}

func UpdateSessionHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")

	var payload SessionDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	membersJSON, _ := json.Marshal(payload.Members)

	var projIDVal *string
	if payload.ProjectID != "" {
		projIDVal = &payload.ProjectID
	} else if payload.ProjectName != "" {
		var matchedID string
		if err := db.QueryRowContext(r.Context(), `
			SELECT id FROM projects WHERE tenant_id = ? AND project_name = ? AND deleted_at IS NULL LIMIT 1
		`, tid, payload.ProjectName).Scan(&matchedID); err == nil && matchedID != "" {
			projIDVal = &matchedID
		}
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE project_sessions SET
			project_id = ?, title = ?, project_name = ?, session_manager = ?,
			project_type = ?, shoot_timing = ?, status = ?, date = ?, start_time = ?,
			duration = ?, time_range = ?, venue = ?, city = ?, location_link = ?,
			remarks = ?, members = ?, shoots_activity_count = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		projIDVal, payload.Title, payload.ProjectName, payload.SessionManager,
		payload.ProjectType, payload.ShootTiming, payload.Status, payload.Date, payload.StartTime,
		payload.Duration, payload.TimeRange, payload.Venue, payload.City, payload.LocationLink,
		payload.Remarks, membersJSON, payload.ShootsActivityCount,
		id, tid,
	)

	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to update session", err)
		return
	}

	payload.ID = id
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE project_sessions SET status = ? WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, req.Status, id, tid)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to update session status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": id, "sessionStatus": req.Status})
}

func DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(), `
		UPDATE project_sessions SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to delete session", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

