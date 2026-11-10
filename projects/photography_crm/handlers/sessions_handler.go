package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

type UpdateSessionStatusRequest struct {
	Status string `json:"status"`
}

type SessionStatusUpdateResponse struct {
	Status        string `json:"status"`
	ID            string `json:"id"`
	SessionStatus string `json:"sessionStatus"`
}

type SessionDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func SessionsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var projectIDFilter string = request.URL.Query().Get("projectId")

	var query string = `
		SELECT s.id, COALESCE(s.project_id, ''), s.title, COALESCE(p.project_name, s.project_name, ''),
		       s.session_manager, s.project_type, s.shoot_timing, s.status, s.date,
		       s.start_time, s.duration, s.time_range, s.venue, s.city, s.location_link,
		       s.remarks, s.members, s.shoots_activity_count
		FROM project_sessions s
		LEFT JOIN projects p ON p.id = s.project_id AND p.deleted_at IS NULL
		WHERE s.tenant_id = ? AND s.deleted_at IS NULL
	`
	var queryArguments = []any{tenantID}
	if projectIDFilter != "" {
		query += " AND s.project_id = ?"
		queryArguments = append(queryArguments, projectIDFilter)
	}
	query += " ORDER BY s.date ASC"

	rows, queryErr := database.QueryContext(request.Context(), query, queryArguments...)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list project sessions", queryErr)
		return
	}
	defer rows.Close()

	var sessionsList = make([]SessionDTO, 0)
	for rows.Next() {
		var (
			sessionID           string
			projectID           string
			sessionTitle        string
			projectName         string
			sessionManager      string
			projectType         string
			shootTiming         string
			sessionStatus       string
			sessionDate         string
			timeRange           string
			venue               string
			city                string
			shootsActivityCount int
			rawMembersBytes     []byte
		)
		var (
			nullableStartTime    sql.NullString
			nullableDuration     sql.NullString
			nullableLocationLink sql.NullString
			nullableRemarks      sql.NullString
		)

		var scanErr error = rows.Scan(
			&sessionID, &projectID, &sessionTitle, &projectName, &sessionManager, &projectType, &shootTiming, &sessionStatus, &sessionDate,
			&nullableStartTime, &nullableDuration, &timeRange, &venue, &city, &nullableLocationLink, &nullableRemarks,
			&rawMembersBytes, &shootsActivityCount,
		)
		if scanErr != nil {
			continue
		}

		var startTime string = ""
		if nullableStartTime.Valid {
			startTime = nullableStartTime.String
		}
		var duration string = ""
		if nullableDuration.Valid {
			duration = nullableDuration.String
		}
		var locationLink string = ""
		if nullableLocationLink.Valid {
			locationLink = nullableLocationLink.String
		}
		var remarks string = ""
		if nullableRemarks.Valid {
			remarks = nullableRemarks.String
		}

		var membersList = make([]string, 0)
		if len(rawMembersBytes) > 0 {
			_ = json.Unmarshal(rawMembersBytes, &membersList)
		}

		var sessionItem = SessionDTO{
			ID:                  sessionID,
			ProjectID:           projectID,
			Title:               sessionTitle,
			ProjectName:         projectName,
			SessionManager:      sessionManager,
			ProjectType:         projectType,
			ShootTiming:         shootTiming,
			Status:              sessionStatus,
			Date:                sessionDate,
			StartTime:           startTime,
			Duration:            duration,
			TimeRange:           timeRange,
			Venue:               venue,
			City:                city,
			LocationLink:        locationLink,
			Remarks:             remarks,
			Members:             membersList,
			ShootsActivityCount: shootsActivityCount,
		}
		sessionsList = append(sessionsList, sessionItem)
	}

	utils.WriteJSON(responseWriter, http.StatusOK, sessionsList)
}

func CreateSessionHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var sessionPayload SessionDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&sessionPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	if sessionPayload.ID == "" {
		sessionPayload.ID = "SES-" + uuid.New().String()[:8]
	}
	if sessionPayload.ShootsActivityCount == 0 {
		sessionPayload.ShootsActivityCount = 1
	}
	if sessionPayload.SessionManager == "" {
		sessionPayload.SessionManager = "especa"
	}
	if sessionPayload.Status == "" {
		sessionPayload.Status = "New"
	}

	membersBytes, _ := json.Marshal(sessionPayload.Members)

	var projectIDPointer *string = nil
	if sessionPayload.ProjectID != "" {
		projectIDPointer = &sessionPayload.ProjectID
	} else if sessionPayload.ProjectName != "" {
		var matchedProjectID string
		var matchErr error = database.QueryRowContext(request.Context(), `
			SELECT id FROM projects WHERE tenant_id = ? AND project_name = ? AND deleted_at IS NULL LIMIT 1
		`, tenantID, sessionPayload.ProjectName).Scan(&matchedProjectID)
		if matchErr == nil && matchedProjectID != "" {
			projectIDPointer = &matchedProjectID
		}
	}

	_, insertErr := database.ExecContext(request.Context(), `
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
		sessionPayload.ID, tenantID, projectIDPointer, sessionPayload.Title, sessionPayload.ProjectName, sessionPayload.SessionManager,
		sessionPayload.ProjectType, sessionPayload.ShootTiming, sessionPayload.Status, sessionPayload.Date, sessionPayload.StartTime, sessionPayload.Duration,
		sessionPayload.TimeRange, sessionPayload.Venue, sessionPayload.City, sessionPayload.LocationLink, sessionPayload.Remarks, membersBytes, sessionPayload.ShootsActivityCount,
	)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create session", insertErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusCreated, sessionPayload)
}

func UpdateSessionHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var sessionID string = chi.URLParam(request, "id")

	var sessionPayload SessionDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&sessionPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	membersBytes, _ := json.Marshal(sessionPayload.Members)

	var projectIDPointer *string = nil
	if sessionPayload.ProjectID != "" {
		projectIDPointer = &sessionPayload.ProjectID
	} else if sessionPayload.ProjectName != "" {
		var matchedProjectID string
		var matchErr error = database.QueryRowContext(request.Context(), `
			SELECT id FROM projects WHERE tenant_id = ? AND project_name = ? AND deleted_at IS NULL LIMIT 1
		`, tenantID, sessionPayload.ProjectName).Scan(&matchedProjectID)
		if matchErr == nil && matchedProjectID != "" {
			projectIDPointer = &matchedProjectID
		}
	}

	_, updateErr := database.ExecContext(request.Context(), `
		UPDATE project_sessions SET
			project_id = ?, title = ?, project_name = ?, session_manager = ?,
			project_type = ?, shoot_timing = ?, status = ?, date = ?, start_time = ?,
			duration = ?, time_range = ?, venue = ?, city = ?, location_link = ?,
			remarks = ?, members = ?, shoots_activity_count = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		projectIDPointer, sessionPayload.Title, sessionPayload.ProjectName, sessionPayload.SessionManager,
		sessionPayload.ProjectType, sessionPayload.ShootTiming, sessionPayload.Status, sessionPayload.Date, sessionPayload.StartTime,
		sessionPayload.Duration, sessionPayload.TimeRange, sessionPayload.Venue, sessionPayload.City, sessionPayload.LocationLink,
		sessionPayload.Remarks, membersBytes, sessionPayload.ShootsActivityCount,
		sessionID, tenantID,
	)

	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update session", updateErr)
		return
	}

	sessionPayload.ID = sessionID
	utils.WriteJSON(responseWriter, http.StatusOK, sessionPayload)
}

func UpdateSessionStatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var sessionID string = chi.URLParam(request, "id")

	var requestPayload UpdateSessionStatusRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, updateErr := database.ExecContext(request.Context(), `
		UPDATE project_sessions SET status = ? WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, requestPayload.Status, sessionID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update session status", updateErr)
		return
	}

	var statusResponse = SessionStatusUpdateResponse{
		Status:        "ok",
		ID:            sessionID,
		SessionStatus: requestPayload.Status,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, statusResponse)
}

func DeleteSessionHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var sessionID string = chi.URLParam(request, "id")

	_, updateErr := database.ExecContext(request.Context(), `
		UPDATE project_sessions SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, sessionID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete session", updateErr)
		return
	}

	var deleteResponse = SessionDeleteResponse{
		Status: "deleted",
		ID:     sessionID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, deleteResponse)
}
