package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectManagerDTO struct {
	Name     string `json:"name"`
	Initials string `json:"initials"`
	Avatar   string `json:"avatar,omitempty"`
}

type ProjectTasksCountDTO struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
}

type ProjectDeliverablesCountDTO struct {
	Ready int `json:"ready"`
	Total int `json:"total"`
}

type ProjectDTO struct {
	ID                string                      `json:"id"`
	ProjectName       string                      `json:"projectName"`
	ClientID          string                      `json:"clientId,omitempty"`
	ClientName        string                      `json:"clientName"`
	ClientInitials    string                      `json:"clientInitials"`
	ClientMobile      string                      `json:"clientMobile"`
	ClientCountryCode string                      `json:"clientCountryCode"`
	ClientEmail       string                      `json:"clientEmail"`
	ClientAddress     string                      `json:"clientAddress,omitempty"`
	LeadSource        string                      `json:"leadSource,omitempty"`
	ProjectManager    ProjectManagerDTO           `json:"projectManager"`
	ProjectType       string                      `json:"projectType"`
	ProjectAddress    string                      `json:"projectAddress,omitempty"`
	SessionsCount     int                         `json:"sessionsCount"`
	TasksCount        ProjectTasksCountDTO        `json:"tasksCount"`
	DeliverablesCount ProjectDeliverablesCountDTO `json:"deliverablesCount"`
	Status            string                      `json:"status"`
	Stage             string                      `json:"stage"`
	StartDate         string                      `json:"startDate"`
	EndDate           string                      `json:"endDate"`
	CreatedDate       string                      `json:"createdDate"`
	ShootHours        float64                     `json:"shootHours"`
	ShootMembersCount int                         `json:"shootMembersCount"`
	ShootExpense      int64                       `json:"shootExpense"`
	GalleriesCount    int                         `json:"galleriesCount"`
	Notes             string                      `json:"notes"`
	Segments          any                         `json:"segments"`
	Documents         any                         `json:"documents"`
	Sessions          any                         `json:"sessions"`
	Tasks             any                         `json:"tasks"`
	Printing          any                         `json:"printing"`
	Deliverables      any                         `json:"deliverables"`
	Expenses          any                         `json:"expenses"`
	Finance           any                         `json:"finance"`
}

func parseProjectManager(raw []byte) ProjectManagerDTO {
	if len(raw) == 0 {
		return ProjectManagerDTO{Name: "Unassigned", Initials: "UA"}
	}
	var pm ProjectManagerDTO
	if err := json.Unmarshal(raw, &pm); err == nil && (pm.Name != "" || pm.Initials != "") {
		if pm.Initials == "" && len(pm.Name) > 0 {
			pm.Initials = strings.ToUpper(pm.Name[:1])
		}
		return pm
	}
	var strName string
	if err := json.Unmarshal(raw, &strName); err == nil && strName != "" {
		initials := "PM"
		if len(strName) > 0 {
			initials = strings.ToUpper(strName[:1])
		}
		return ProjectManagerDTO{Name: strName, Initials: initials}
	}
	return ProjectManagerDTO{Name: "Unassigned", Initials: "UA"}
}

func calculateTasksCount(tasks any) ProjectTasksCountDTO {
	completed := 0
	total := 0
	if list, ok := tasks.([]any); ok {
		total = len(list)
		for _, item := range list {
			if m, ok := item.(map[string]any); ok {
				if s, ok := m["status"].(string); ok && strings.EqualFold(s, "completed") {
					completed++
				} else if b, ok := m["completed"].(bool); ok && b {
					completed++
				}
			}
		}
	}
	return ProjectTasksCountDTO{Completed: completed, Total: total}
}

func calculateDeliverablesCount(deliverables any) ProjectDeliverablesCountDTO {
	ready := 0
	total := 0
	if list, ok := deliverables.([]any); ok {
		total = len(list)
		for _, item := range list {
			if m, ok := item.(map[string]any); ok {
				if s, ok := m["status"].(string); ok && (strings.EqualFold(s, "ready") || strings.EqualFold(s, "completed")) {
					ready++
				}
			}
		}
	}
	return ProjectDeliverablesCountDTO{Ready: ready, Total: total}
}

func ProjectsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT p.id, p.project_name, p.client_name, p.client_initials, p.client_mobile,
		       p.client_country_code, p.client_email, p.client_address, p.lead_source,
		       p.project_manager, p.project_type, p.project_address, p.status, p.stage,
		       p.start_date, p.end_date, p.shoot_hours, p.shoot_members_count, p.shoot_expense,
		       p.galleries_count, p.notes, p.segments, p.documents, p.sessions, p.tasks,
		       p.printing, p.deliverables, p.expenses, p.finance, p.created_at,
		       (SELECT COUNT(*) FROM project_sessions ps WHERE ps.project_id = p.id AND ps.deleted_at IS NULL) as live_sessions_count
		FROM projects p
		WHERE p.tenant_id = ? AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC
	`, tid)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to load projects", err)
		return
	}
	defer rows.Close()

	result := make([]ProjectDTO, 0)
	for rows.Next() {
		var (
			id, projName, clientName, clientInit, clientMobile, clientCode, clientEmail string
			status, stage, projType                                                     string
			shootHours                                                                  float64
			shootMembersCount, galleriesCount, liveSessionsCount                        int
			shootExpense                                                                int64
			createdAt                                                                   time.Time
			mgrRaw, segRaw, docRaw, sesRaw, taskRaw, printRaw, delivRaw, expRaw, finRaw []byte
		)
		var (
			sClientAddr, sLeadSource, sProjAddr, sStartDate, sEndDate, sNotes sql.NullString
		)

		err := rows.Scan(
			&id, &projName, &clientName, &clientInit, &clientMobile,
			&clientCode, &clientEmail, &sClientAddr, &sLeadSource,
			&mgrRaw, &projType, &sProjAddr, &status, &stage,
			&sStartDate, &sEndDate, &shootHours, &shootMembersCount, &shootExpense,
			&galleriesCount, &sNotes, &segRaw, &docRaw, &sesRaw, &taskRaw,
			&printRaw, &delivRaw, &expRaw, &finRaw, &createdAt,
			&liveSessionsCount,
		)
		if err != nil {
			continue
		}

		clientAddress := ""
		if sClientAddr.Valid {
			clientAddress = sClientAddr.String
		}
		leadSource := ""
		if sLeadSource.Valid {
			leadSource = sLeadSource.String
		}
		projectAddress := ""
		if sProjAddr.Valid {
			projectAddress = sProjAddr.String
		}
		startDate := ""
		if sStartDate.Valid {
			startDate = sStartDate.String
		}
		endDate := ""
		if sEndDate.Valid {
			endDate = sEndDate.String
		}
		notes := ""
		if sNotes.Valid {
			notes = sNotes.String
		}

		var mgr, segments, docs, sessions, tasks, printing, deliverables, expenses, finance any
		if len(mgrRaw) > 0 {
			_ = json.Unmarshal(mgrRaw, &mgr)
		}
		if len(segRaw) > 0 {
			_ = json.Unmarshal(segRaw, &segments)
		}
		if len(docRaw) > 0 {
			_ = json.Unmarshal(docRaw, &docs)
		}
		if len(sesRaw) > 0 {
			_ = json.Unmarshal(sesRaw, &sessions)
		}
		if len(taskRaw) > 0 {
			_ = json.Unmarshal(taskRaw, &tasks)
		}
		if len(printRaw) > 0 {
			_ = json.Unmarshal(printRaw, &printing)
		}
		if len(delivRaw) > 0 {
			_ = json.Unmarshal(delivRaw, &deliverables)
		}
		if len(expRaw) > 0 {
			_ = json.Unmarshal(expRaw, &expenses)
		}
		if len(finRaw) > 0 {
			_ = json.Unmarshal(finRaw, &finance)
		}

		finalSessionsCount := liveSessionsCount
		if finalSessionsCount == 0 {
			if sList, ok := sessions.([]any); ok {
				finalSessionsCount = len(sList)
			}
		}

		dto := ProjectDTO{
			ID:                id,
			ProjectName:       projName,
			ClientName:        clientName,
			ClientInitials:    clientInit,
			ClientMobile:      clientMobile,
			ClientCountryCode: clientCode,
			ClientEmail:       clientEmail,
			ClientAddress:     clientAddress,
			LeadSource:        leadSource,
			ProjectManager:    parseProjectManager(mgrRaw),
			ProjectType:       projType,
			ProjectAddress:    projectAddress,
			SessionsCount:     finalSessionsCount,
			TasksCount:        calculateTasksCount(tasks),
			DeliverablesCount: calculateDeliverablesCount(deliverables),
			Status:            status,
			Stage:             stage,
			StartDate:         startDate,
			EndDate:           endDate,
			CreatedDate:       createdAt.Format("02/01/2006"),
			ShootHours:        shootHours,
			ShootMembersCount: shootMembersCount,
			ShootExpense:      shootExpense,
			GalleriesCount:    galleriesCount,
			Notes:             notes,
			Segments:          segments,
			Documents:         docs,
			Sessions:          sessions,
			Tasks:             tasks,
			Printing:          printing,
			Deliverables:      deliverables,
			Expenses:          expenses,
			Finance:           finance,
		}
		result = append(result, dto)
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
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

	var (
		projName, clientName, clientInit, clientMobile, clientCode, clientEmail string
		status, stage, projType                                                 string
		shootHours                                                              float64
		shootMembersCount, galleriesCount, liveSessionsCount                    int
		shootExpense                                                            int64
		createdAt                                                               time.Time
		mgrRaw, segRaw, docRaw, sesRaw, taskRaw, printRaw, delivRaw, expRaw, finRaw []byte
	)
	var (
		sClientAddr, sLeadSource, sProjAddr, sStartDate, sEndDate, sNotes sql.NullString
	)

	err := db.QueryRowContext(r.Context(), `
		SELECT p.id, p.project_name, p.client_name, p.client_initials, p.client_mobile,
		       p.client_country_code, p.client_email, p.client_address, p.lead_source,
		       p.project_manager, p.project_type, p.project_address, p.status, p.stage,
		       p.start_date, p.end_date, p.shoot_hours, p.shoot_members_count, p.shoot_expense,
		       p.galleries_count, p.notes, p.segments, p.documents, p.sessions, p.tasks,
		       p.printing, p.deliverables, p.expenses, p.finance, p.created_at,
		       (SELECT COUNT(*) FROM project_sessions ps WHERE ps.project_id = p.id AND ps.deleted_at IS NULL) as live_sessions_count
		FROM projects p
		WHERE p.id = ? AND p.tenant_id = ? AND p.deleted_at IS NULL
	`, id, tid).Scan(
		&id, &projName, &clientName, &clientInit, &clientMobile,
		&clientCode, &clientEmail, &sClientAddr, &sLeadSource,
		&mgrRaw, &projType, &sProjAddr, &status, &stage,
		&sStartDate, &sEndDate, &shootHours, &shootMembersCount, &shootExpense,
		&galleriesCount, &sNotes, &segRaw, &docRaw, &sesRaw, &taskRaw,
		&printRaw, &delivRaw, &expRaw, &finRaw, &createdAt,
		&liveSessionsCount,
	)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Project not found")
		return
	}

	clientAddress := ""
	if sClientAddr.Valid {
		clientAddress = sClientAddr.String
	}
	leadSource := ""
	if sLeadSource.Valid {
		leadSource = sLeadSource.String
	}
	projectAddress := ""
	if sProjAddr.Valid {
		projectAddress = sProjAddr.String
	}
	startDate := ""
	if sStartDate.Valid {
		startDate = sStartDate.String
	}
	endDate := ""
	if sEndDate.Valid {
		endDate = sEndDate.String
	}
	notes := ""
	if sNotes.Valid {
		notes = sNotes.String
	}

	var mgr, segments, docs, sessions, tasks, printing, deliverables, expenses, finance any
	if len(mgrRaw) > 0 {
		_ = json.Unmarshal(mgrRaw, &mgr)
	}
	if len(segRaw) > 0 {
		_ = json.Unmarshal(segRaw, &segments)
	}
	if len(docRaw) > 0 {
		_ = json.Unmarshal(docRaw, &docs)
	}
	if len(sesRaw) > 0 {
		_ = json.Unmarshal(sesRaw, &sessions)
	}
	if len(taskRaw) > 0 {
		_ = json.Unmarshal(taskRaw, &tasks)
	}
	if len(printRaw) > 0 {
		_ = json.Unmarshal(printRaw, &printing)
	}
	if len(delivRaw) > 0 {
		_ = json.Unmarshal(delivRaw, &deliverables)
	}
	if len(expRaw) > 0 {
		_ = json.Unmarshal(expRaw, &expenses)
	}
	if len(finRaw) > 0 {
		_ = json.Unmarshal(finRaw, &finance)
	}

	dto := ProjectDTO{
		ID:                id,
		ProjectName:       projName,
		ClientName:        clientName,
		ClientInitials:    clientInit,
		ClientMobile:      clientMobile,
		ClientCountryCode: clientCode,
		ClientEmail:       clientEmail,
		ClientAddress:     clientAddress,
		LeadSource:        leadSource,
		ProjectManager:    parseProjectManager(mgrRaw),
		ProjectType:       projType,
		ProjectAddress:    projectAddress,
		SessionsCount:     liveSessionsCount,
		TasksCount:        calculateTasksCount(tasks),
		DeliverablesCount: calculateDeliverablesCount(deliverables),
		Status:            status,
		Stage:             stage,
		StartDate:         startDate,
		EndDate:           endDate,
		CreatedDate:       createdAt.Format("02/01/2006"),
		ShootHours:        shootHours,
		ShootMembersCount: shootMembersCount,
		ShootExpense:      shootExpense,
		GalleriesCount:    galleriesCount,
		Notes:             notes,
		Segments:          segments,
		Documents:         docs,
		Sessions:          sessions,
		Tasks:             tasks,
		Printing:          printing,
		Deliverables:      deliverables,
		Expenses:          expenses,
		Finance:           finance,
	}

	utils.WriteJSON(w, http.StatusOK, dto)
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	var payload ProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if payload.ID == "" {
		payload.ID = "PRJ-" + uuid.New().String()[:8]
	}
	if payload.ClientInitials == "" && len(payload.ClientName) > 0 {
		payload.ClientInitials = string([]rune(payload.ClientName)[:1])
	}

	mgrJSON, _ := json.Marshal(payload.ProjectManager)
	segJSON, _ := json.Marshal(payload.Segments)
	docJSON, _ := json.Marshal(payload.Documents)
	sesJSON, _ := json.Marshal(payload.Sessions)
	taskJSON, _ := json.Marshal(payload.Tasks)
	printJSON, _ := json.Marshal(payload.Printing)
	delivJSON, _ := json.Marshal(payload.Deliverables)
	expJSON, _ := json.Marshal(payload.Expenses)
	finJSON, _ := json.Marshal(payload.Finance)

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO projects (
			id, tenant_id, project_name, client_name, client_initials, client_mobile,
			client_country_code, client_email, client_address, lead_source, project_manager,
			project_type, project_address, status, stage, start_date, end_date,
			shoot_hours, shoot_members_count, shoot_expense, galleries_count, notes,
			segments, documents, sessions, tasks, printing, deliverables, expenses, finance
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`,
		payload.ID, tid, payload.ProjectName, payload.ClientName, payload.ClientInitials, payload.ClientMobile,
		payload.ClientCountryCode, payload.ClientEmail, payload.ClientAddress, payload.LeadSource, mgrJSON,
		payload.ProjectType, payload.ProjectAddress, payload.Status, payload.Stage, payload.StartDate, payload.EndDate,
		payload.ShootHours, payload.ShootMembersCount, payload.ShootExpense, payload.GalleriesCount, payload.Notes,
		segJSON, docJSON, sesJSON, taskJSON, printJSON, delivJSON, expJSON, finJSON,
	)

	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to create project", err)
		return
	}

	if sesList, ok := payload.Sessions.([]any); ok && len(sesList) > 0 {
		for _, sObj := range sesList {
			if sMap, ok := sObj.(map[string]any); ok {
				sesID, _ := sMap["id"].(string)
				if sesID == "" {
					sesID = "SES-" + uuid.New().String()[:8]
				}
				title, _ := sMap["title"].(string)
				sDate, _ := sMap["date"].(string)
				sTimeRange, _ := sMap["timeRange"].(string)
				venue, _ := sMap["venue"].(string)
				city, _ := sMap["city"].(string)
				status, _ := sMap["status"].(string)
				if status == "" {
					status = "New"
				}
				membersJSON, _ := json.Marshal(sMap["members"])

				_, _ = db.ExecContext(r.Context(), `
					INSERT INTO project_sessions (
						id, tenant_id, project_id, title, project_name, session_manager,
						project_type, shoot_timing, status, date, time_range, venue, city,
						members, shoots_activity_count
					) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
					ON DUPLICATE KEY UPDATE title = VALUES(title), date = VALUES(date), venue = VALUES(venue)
				`, sesID, tid, payload.ID, title, payload.ProjectName, "especa",
					payload.ProjectType, "Evening", status, sDate, sTimeRange, venue, city, membersJSON)
			}
		}
	}

	payload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(w, http.StatusCreated, payload)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
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

	var payload ProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	mgrJSON, _ := json.Marshal(payload.ProjectManager)
	segJSON, _ := json.Marshal(payload.Segments)
	docJSON, _ := json.Marshal(payload.Documents)
	sesJSON, _ := json.Marshal(payload.Sessions)
	taskJSON, _ := json.Marshal(payload.Tasks)
	printJSON, _ := json.Marshal(payload.Printing)
	delivJSON, _ := json.Marshal(payload.Deliverables)
	expJSON, _ := json.Marshal(payload.Expenses)
	finJSON, _ := json.Marshal(payload.Finance)

	_, err := db.ExecContext(r.Context(), `
		UPDATE projects SET
			project_name = ?, client_name = ?, client_initials = ?, client_mobile = ?,
			client_country_code = ?, client_email = ?, client_address = ?, lead_source = ?,
			project_manager = ?, project_type = ?, project_address = ?, status = ?, stage = ?,
			start_date = ?, end_date = ?, shoot_hours = ?, shoot_members_count = ?,
			shoot_expense = ?, galleries_count = ?, notes = ?, segments = ?, documents = ?,
			sessions = ?, tasks = ?, printing = ?, deliverables = ?, expenses = ?, finance = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		payload.ProjectName, payload.ClientName, payload.ClientInitials, payload.ClientMobile,
		payload.ClientCountryCode, payload.ClientEmail, payload.ClientAddress, payload.LeadSource,
		mgrJSON, payload.ProjectType, payload.ProjectAddress, payload.Status, payload.Stage,
		payload.StartDate, payload.EndDate, payload.ShootHours, payload.ShootMembersCount,
		payload.ShootExpense, payload.GalleriesCount, payload.Notes, segJSON, docJSON,
		sesJSON, taskJSON, printJSON, delivJSON, expJSON, finJSON,
		id, tid,
	)

	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to update project", err)
		return
	}

	payload.ID = id
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateProjectStageHandler(w http.ResponseWriter, r *http.Request) {
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
		Stage  string `json:"stage"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	query := "UPDATE projects SET stage = ?"
	args := []any{req.Stage}
	if req.Status != "" {
		query += ", status = ?"
		args = append(args, req.Status)
	}
	query += " WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL"
	args = append(args, id, tid)

	_, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to update project stage", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": id, "stage": req.Stage})
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
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
		UPDATE projects SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to delete project", err)
		return
	}

	_, _ = db.ExecContext(r.Context(), `
		UPDATE project_sessions SET deleted_at = NOW() WHERE project_id = ? AND tenant_id = ?
	`, id, tid)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

