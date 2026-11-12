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
	"github.com/jmoiron/sqlx"
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

type UpdateProjectStageRequest struct {
	Stage  string `json:"stage"`
	Status string `json:"status"`
}

type ProjectStageUpdateResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	Stage  string `json:"stage"`
}

type ProjectDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func parseProjectManager(rawBytes []byte) ProjectManagerDTO {
	if len(rawBytes) == 0 {
		return ProjectManagerDTO{Name: "Unassigned", Initials: "UA"}
	}
	var projectManager ProjectManagerDTO
	var unmarshalErr error = json.Unmarshal(rawBytes, &projectManager)
	if unmarshalErr == nil && (projectManager.Name != "" || projectManager.Initials != "") {
		if projectManager.Initials == "" && len(projectManager.Name) > 0 {
			projectManager.Initials = strings.ToUpper(projectManager.Name[:1])
		}
		return projectManager
	}
	var stringName string
	var stringUnmarshalErr error = json.Unmarshal(rawBytes, &stringName)
	if stringUnmarshalErr == nil && stringName != "" {
		var initials string = "PM"
		if len(stringName) > 0 {
			initials = strings.ToUpper(stringName[:1])
		}
		return ProjectManagerDTO{Name: stringName, Initials: initials}
	}
	return ProjectManagerDTO{Name: "Unassigned", Initials: "UA"}
}

func calculateTasksCount(tasks any) ProjectTasksCountDTO {
	var completedCount int = 0
	var totalCount int = 0
	var taskList []any
	var isTaskList bool
	taskList, isTaskList = tasks.([]any)
	if isTaskList {
		totalCount = len(taskList)
		var taskIndex int
		for taskIndex = 0; taskIndex < len(taskList); taskIndex++ {
			var taskItem any = taskList[taskIndex]
			var taskMap map[string]any
			var isTaskMap bool
			taskMap, isTaskMap = taskItem.(map[string]any)
			if isTaskMap {
				var rawStatusAny any = taskMap["status"]
				var statusString string
				var isStatusString bool
				statusString, isStatusString = rawStatusAny.(string)
				if isStatusString && strings.EqualFold(statusString, "completed") {
					completedCount++
				} else {
					var rawCompletedAny any = taskMap["completed"]
					var isCompleted bool
					var isBool bool
					isCompleted, isBool = rawCompletedAny.(bool)
					if isBool && isCompleted {
						completedCount++
					}
				}
			}
		}
	}
	return ProjectTasksCountDTO{Completed: completedCount, Total: totalCount}
}

func calculateDeliverablesCount(deliverables any) ProjectDeliverablesCountDTO {
	var readyCount int = 0
	var totalCount int = 0
	var deliverableList []any
	var isDeliverableList bool
	deliverableList, isDeliverableList = deliverables.([]any)
	if isDeliverableList {
		totalCount = len(deliverableList)
		var deliverableIndex int
		for deliverableIndex = 0; deliverableIndex < len(deliverableList); deliverableIndex++ {
			var deliverableItem any = deliverableList[deliverableIndex]
			var deliverableMap map[string]any
			var isDeliverableMap bool
			deliverableMap, isDeliverableMap = deliverableItem.(map[string]any)
			if isDeliverableMap {
				var rawStatusAny any = deliverableMap["status"]
				var statusString string
				var isStatusString bool
				statusString, isStatusString = rawStatusAny.(string)
				if isStatusString && (strings.EqualFold(statusString, "ready") || strings.EqualFold(statusString, "completed")) {
					readyCount++
				}
			}
		}
	}
	return ProjectDeliverablesCountDTO{Ready: readyCount, Total: totalCount}
}

func ProjectsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var statusFilter string = strings.TrimSpace(request.URL.Query().Get("status"))
	var stageFilter string = strings.TrimSpace(request.URL.Query().Get("stage"))
	var projectTypeFilter string = strings.TrimSpace(request.URL.Query().Get("projectType"))
	if projectTypeFilter == "" {
		projectTypeFilter = strings.TrimSpace(request.URL.Query().Get("type"))
	}
	var searchFilter string = strings.TrimSpace(request.URL.Query().Get("search"))

	var querySql string = `
		SELECT p.id, p.project_name, p.client_name, p.client_initials, p.client_mobile,
		       p.client_country_code, p.client_email, p.client_address, p.lead_source,
		       p.project_manager, p.project_type, p.project_address, p.status, p.stage,
		       p.start_date, p.end_date, p.shoot_hours, p.shoot_members_count, p.shoot_expense,
		       p.galleries_count, p.notes, p.tasks, p.deliverables, p.created_at,
		       (SELECT COUNT(*) FROM project_sessions ps WHERE ps.project_id = p.id AND ps.deleted_at IS NULL) as live_sessions_count
		FROM projects p
		WHERE p.tenant_id = ? AND p.deleted_at IS NULL
	`
	var queryParams []interface{}
	queryParams = append(queryParams, tenantID)

	if statusFilter != "" && !strings.EqualFold(statusFilter, "all") {
		querySql += " AND p.status = ?"
		queryParams = append(queryParams, statusFilter)
	}
	if stageFilter != "" && !strings.EqualFold(stageFilter, "all") {
		querySql += " AND p.stage = ?"
		queryParams = append(queryParams, stageFilter)
	}
	if projectTypeFilter != "" && !strings.EqualFold(projectTypeFilter, "all") {
		querySql += " AND p.project_type = ?"
		queryParams = append(queryParams, projectTypeFilter)
	}
	if searchFilter != "" {
		querySql += " AND (p.project_name LIKE ? OR p.client_name LIKE ? OR p.client_email LIKE ? OR p.client_mobile LIKE ?)"
		var wildcard string = "%" + searchFilter + "%"
		queryParams = append(queryParams, wildcard, wildcard, wildcard, wildcard)
	}
	querySql += " ORDER BY p.created_at DESC"

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), querySql, queryParams...)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to load projects", queryErr)
		return
	}
	defer rows.Close()

	var projectsList = make([]ProjectDTO, 0)
	for rows.Next() {
		var (
			projectID         string
			projectName       string
			clientName        string
			clientInitials    string
			clientMobile      string
			clientCountryCode string
			clientEmail       string
			projectStatus     string
			projectStage      string
			projectType       string
			shootHours        float64
			shootMembersCount int
			galleriesCount    int
			liveSessionsCount int
			shootExpense      int64
			createdAt         time.Time
			rawManagerBytes   []byte
			rawTasksBytes     []byte
			rawDelivBytes     []byte
		)
		var (
			nullableClientAddress  sql.NullString
			nullableLeadSource     sql.NullString
			nullableProjectAddress sql.NullString
			nullableStartDate      sql.NullString
			nullableEndDate        sql.NullString
			nullableNotes          sql.NullString
		)

		var scanErr error = rows.Scan(
			&projectID, &projectName, &clientName, &clientInitials, &clientMobile,
			&clientCountryCode, &clientEmail, &nullableClientAddress, &nullableLeadSource,
			&rawManagerBytes, &projectType, &nullableProjectAddress, &projectStatus, &projectStage,
			&nullableStartDate, &nullableEndDate, &shootHours, &shootMembersCount, &shootExpense,
			&galleriesCount, &nullableNotes, &rawTasksBytes, &rawDelivBytes, &createdAt,
			&liveSessionsCount,
		)
		if scanErr != nil {
			continue
		}

		var clientAddress string = ""
		if nullableClientAddress.Valid {
			clientAddress = nullableClientAddress.String
		}
		var leadSource string = ""
		if nullableLeadSource.Valid {
			leadSource = nullableLeadSource.String
		}
		var projectAddress string = ""
		if nullableProjectAddress.Valid {
			projectAddress = nullableProjectAddress.String
		}
		var startDate string = ""
		if nullableStartDate.Valid {
			startDate = nullableStartDate.String
		}
		var endDate string = ""
		if nullableEndDate.Valid {
			endDate = nullableEndDate.String
		}
		var notes string = ""
		if nullableNotes.Valid {
			notes = nullableNotes.String
		}

		var tasksData any
		var deliverablesData any

		if len(rawTasksBytes) > 0 {
			_ = json.Unmarshal(rawTasksBytes, &tasksData)
		}
		if len(rawDelivBytes) > 0 {
			_ = json.Unmarshal(rawDelivBytes, &deliverablesData)
		}

		var projectDTO = ProjectDTO{
			ID:                projectID,
			ProjectName:       projectName,
			ClientName:        clientName,
			ClientInitials:    clientInitials,
			ClientMobile:      clientMobile,
			ClientCountryCode: clientCountryCode,
			ClientEmail:       clientEmail,
			ClientAddress:     clientAddress,
			LeadSource:        leadSource,
			ProjectManager:    parseProjectManager(rawManagerBytes),
			ProjectType:       projectType,
			ProjectAddress:    projectAddress,
			SessionsCount:     liveSessionsCount,
			TasksCount:        calculateTasksCount(tasksData),
			DeliverablesCount: calculateDeliverablesCount(deliverablesData),
			Status:            projectStatus,
			Stage:             projectStage,
			StartDate:         startDate,
			EndDate:           endDate,
			CreatedDate:       createdAt.Format("02/01/2006"),
			ShootHours:        shootHours,
			ShootMembersCount: shootMembersCount,
			ShootExpense:      shootExpense,
			GalleriesCount:    galleriesCount,
			Notes:             notes,
		}
		projectsList = append(projectsList, projectDTO)
	}

	utils.WriteJSON(responseWriter, http.StatusOK, projectsList)
}

func GetProjectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var projectID string = chi.URLParam(request, "id")

	var (
		projectName       string
		clientName        string
		clientInitials    string
		clientMobile      string
		clientCountryCode string
		clientEmail       string
		projectStatus     string
		projectStage      string
		projectType       string
		shootHours        float64
		shootMembersCount int
		galleriesCount    int
		liveSessionsCount int
		shootExpense      int64
		createdAt         time.Time
		rawManagerBytes   []byte
		rawSegmentsBytes  []byte
		rawDocumentsBytes []byte
		rawSessionsBytes  []byte
		rawTasksBytes     []byte
		rawPrintingBytes  []byte
		rawDelivBytes     []byte
		rawExpensesBytes  []byte
		rawFinanceBytes   []byte
	)
	var (
		nullableClientAddress  sql.NullString
		nullableLeadSource     sql.NullString
		nullableProjectAddress sql.NullString
		nullableStartDate      sql.NullString
		nullableEndDate        sql.NullString
		nullableNotes          sql.NullString
	)

	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT p.id, p.project_name, p.client_name, p.client_initials, p.client_mobile,
		       p.client_country_code, p.client_email, p.client_address, p.lead_source,
		       p.project_manager, p.project_type, p.project_address, p.status, p.stage,
		       p.start_date, p.end_date, p.shoot_hours, p.shoot_members_count, p.shoot_expense,
		       p.galleries_count, p.notes, p.segments, p.documents, p.sessions, p.tasks,
		       p.printing, p.deliverables, p.expenses, p.finance, p.created_at,
		       (SELECT COUNT(*) FROM project_sessions ps WHERE ps.project_id = p.id AND ps.deleted_at IS NULL) as live_sessions_count
		FROM projects p
		WHERE p.id = ? AND p.tenant_id = ? AND p.deleted_at IS NULL
	`, projectID, tenantID).Scan(
		&projectID, &projectName, &clientName, &clientInitials, &clientMobile,
		&clientCountryCode, &clientEmail, &nullableClientAddress, &nullableLeadSource,
		&rawManagerBytes, &projectType, &nullableProjectAddress, &projectStatus, &projectStage,
		&nullableStartDate, &nullableEndDate, &shootHours, &shootMembersCount, &shootExpense,
		&galleriesCount, &nullableNotes, &rawSegmentsBytes, &rawDocumentsBytes, &rawSessionsBytes, &rawTasksBytes,
		&rawPrintingBytes, &rawDelivBytes, &rawExpensesBytes, &rawFinanceBytes, &createdAt,
		&liveSessionsCount,
	)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Project not found")
		return
	}

	var clientAddress string = ""
	if nullableClientAddress.Valid {
		clientAddress = nullableClientAddress.String
	}
	var leadSource string = ""
	if nullableLeadSource.Valid {
		leadSource = nullableLeadSource.String
	}
	var projectAddress string = ""
	if nullableProjectAddress.Valid {
		projectAddress = nullableProjectAddress.String
	}
	var startDate string = ""
	if nullableStartDate.Valid {
		startDate = nullableStartDate.String
	}
	var endDate string = ""
	if nullableEndDate.Valid {
		endDate = nullableEndDate.String
	}
	var notes string = ""
	if nullableNotes.Valid {
		notes = nullableNotes.String
	}

	var managerData any
	var segmentsData any
	var documentsData any
	var sessionsData any
	var tasksData any
	var printingData any
	var deliverablesData any
	var expensesData any
	var financeData any

	if len(rawManagerBytes) > 0 {
		_ = json.Unmarshal(rawManagerBytes, &managerData)
	}
	if len(rawSegmentsBytes) > 0 {
		_ = json.Unmarshal(rawSegmentsBytes, &segmentsData)
	}
	if len(rawDocumentsBytes) > 0 {
		_ = json.Unmarshal(rawDocumentsBytes, &documentsData)
	}
	if len(rawSessionsBytes) > 0 {
		_ = json.Unmarshal(rawSessionsBytes, &sessionsData)
	}
	if len(rawTasksBytes) > 0 {
		_ = json.Unmarshal(rawTasksBytes, &tasksData)
	}
	if len(rawPrintingBytes) > 0 {
		_ = json.Unmarshal(rawPrintingBytes, &printingData)
	}
	if len(rawDelivBytes) > 0 {
		_ = json.Unmarshal(rawDelivBytes, &deliverablesData)
	}
	if len(rawExpensesBytes) > 0 {
		_ = json.Unmarshal(rawExpensesBytes, &expensesData)
	}
	if len(rawFinanceBytes) > 0 {
		_ = json.Unmarshal(rawFinanceBytes, &financeData)
	}

	var projectDTO = ProjectDTO{
		ID:                projectID,
		ProjectName:       projectName,
		ClientName:        clientName,
		ClientInitials:    clientInitials,
		ClientMobile:      clientMobile,
		ClientCountryCode: clientCountryCode,
		ClientEmail:       clientEmail,
		ClientAddress:     clientAddress,
		LeadSource:        leadSource,
		ProjectManager:    parseProjectManager(rawManagerBytes),
		ProjectType:       projectType,
		ProjectAddress:    projectAddress,
		SessionsCount:     liveSessionsCount,
		TasksCount:        calculateTasksCount(tasksData),
		DeliverablesCount: calculateDeliverablesCount(deliverablesData),
		Status:            projectStatus,
		Stage:             projectStage,
		StartDate:         startDate,
		EndDate:           endDate,
		CreatedDate:       createdAt.Format("02/01/2006"),
		ShootHours:        shootHours,
		ShootMembersCount: shootMembersCount,
		ShootExpense:      shootExpense,
		GalleriesCount:    galleriesCount,
		Notes:             notes,
		Segments:          segmentsData,
		Documents:         documentsData,
		Sessions:          sessionsData,
		Tasks:             tasksData,
		Printing:          printingData,
		Deliverables:      deliverablesData,
		Expenses:          expensesData,
		Finance:           financeData,
	}

	utils.WriteJSON(responseWriter, http.StatusOK, projectDTO)
}

func CreateProjectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var projectPayload ProjectDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&projectPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	if projectPayload.ID == "" {
		projectPayload.ID = "PRJ-" + uuid.New().String()[:8]
	}
	if projectPayload.ClientInitials == "" && len(projectPayload.ClientName) > 0 {
		projectPayload.ClientInitials = strings.ToUpper(projectPayload.ClientName[0:1])
	}

	var managerBytes []byte
	managerBytes, _ = json.Marshal(projectPayload.ProjectManager)
	var segmentsBytes []byte
	segmentsBytes, _ = json.Marshal(projectPayload.Segments)
	var documentsBytes []byte
	documentsBytes, _ = json.Marshal(projectPayload.Documents)
	var sessionsBytes []byte
	sessionsBytes, _ = json.Marshal(projectPayload.Sessions)
	var tasksBytes []byte
	tasksBytes, _ = json.Marshal(projectPayload.Tasks)
	var printingBytes []byte
	printingBytes, _ = json.Marshal(projectPayload.Printing)
	var deliverablesBytes []byte
	deliverablesBytes, _ = json.Marshal(projectPayload.Deliverables)
	var expensesBytes []byte
	expensesBytes, _ = json.Marshal(projectPayload.Expenses)
	var financeBytes []byte
	financeBytes, _ = json.Marshal(projectPayload.Finance)

	var insertErr error
	_, insertErr = database.ExecContext(request.Context(), `
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
		projectPayload.ID, tenantID, projectPayload.ProjectName, projectPayload.ClientName, projectPayload.ClientInitials, projectPayload.ClientMobile,
		projectPayload.ClientCountryCode, projectPayload.ClientEmail, projectPayload.ClientAddress, projectPayload.LeadSource, managerBytes,
		projectPayload.ProjectType, projectPayload.ProjectAddress, projectPayload.Status, projectPayload.Stage, projectPayload.StartDate, projectPayload.EndDate,
		projectPayload.ShootHours, projectPayload.ShootMembersCount, projectPayload.ShootExpense, projectPayload.GalleriesCount, projectPayload.Notes,
		segmentsBytes, documentsBytes, sessionsBytes, tasksBytes, printingBytes, deliverablesBytes, expensesBytes, financeBytes,
	)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create project", insertErr)
		return
	}

	var sessionList []any
	var isSessionList bool
	sessionList, isSessionList = projectPayload.Sessions.([]any)
	if isSessionList && len(sessionList) > 0 {
		var sessionIndex int
		for sessionIndex = 0; sessionIndex < len(sessionList); sessionIndex++ {
			var sessionObject any = sessionList[sessionIndex]
			var sessionMap map[string]any
			var isSessionMap bool
			sessionMap, isSessionMap = sessionObject.(map[string]any)
			if isSessionMap {
				var sessionID string
				var rawIDAny any = sessionMap["id"]
				var idVal string
				var isIDString bool
				idVal, isIDString = rawIDAny.(string)
				if isIDString {
					sessionID = idVal
				}
				if sessionID == "" {
					sessionID = "SES-" + uuid.New().String()[:8]
				}

				var sessionTitle string
				var rawTitleAny any = sessionMap["title"]
				var titleVal string
				var isTitleString bool
				titleVal, isTitleString = rawTitleAny.(string)
				if isTitleString {
					sessionTitle = titleVal
				}

				var sessionDate string
				var rawDateAny any = sessionMap["date"]
				var dateVal string
				var isDateString bool
				dateVal, isDateString = rawDateAny.(string)
				if isDateString {
					sessionDate = dateVal
				}

				var sessionTimeRange string
				var rawTimeRangeAny any = sessionMap["timeRange"]
				var timeVal string
				var isTimeRangeString bool
				timeVal, isTimeRangeString = rawTimeRangeAny.(string)
				if isTimeRangeString {
					sessionTimeRange = timeVal
				}

				var venue string
				var rawVenueAny any = sessionMap["venue"]
				var venueVal string
				var isVenueString bool
				venueVal, isVenueString = rawVenueAny.(string)
				if isVenueString {
					venue = venueVal
				}

				var city string
				var rawCityAny any = sessionMap["city"]
				var cityVal string
				var isCityString bool
				cityVal, isCityString = rawCityAny.(string)
				if isCityString {
					city = cityVal
				}

				var sessionStatus string
				var rawStatusAny any = sessionMap["status"]
				var statusVal string
				var isStatusString bool
				statusVal, isStatusString = rawStatusAny.(string)
				if isStatusString {
					sessionStatus = statusVal
				}
				if sessionStatus == "" {
					sessionStatus = "New"
				}

				var sessionMembersBytes []byte
				sessionMembersBytes, _ = json.Marshal(sessionMap["members"])

				_, _ = database.ExecContext(request.Context(), `
					INSERT INTO project_sessions (
						id, tenant_id, project_id, title, project_name, session_manager,
						project_type, shoot_timing, status, date, time_range, venue, city,
						members, shoots_activity_count
					) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
					ON DUPLICATE KEY UPDATE title = VALUES(title), date = VALUES(date), venue = VALUES(venue)
				`, sessionID, tenantID, projectPayload.ID, sessionTitle, projectPayload.ProjectName, "especa",
					projectPayload.ProjectType, "Evening", sessionStatus, sessionDate, sessionTimeRange, venue, city, sessionMembersBytes)
			}
		}
	}

	projectPayload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(responseWriter, http.StatusCreated, projectPayload)
}

func UpdateProjectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var projectID string = chi.URLParam(request, "id")

	var projectPayload ProjectDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&projectPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var managerBytes []byte
	managerBytes, _ = json.Marshal(projectPayload.ProjectManager)
	var segmentsBytes []byte
	segmentsBytes, _ = json.Marshal(projectPayload.Segments)
	var documentsBytes []byte
	documentsBytes, _ = json.Marshal(projectPayload.Documents)
	var sessionsBytes []byte
	sessionsBytes, _ = json.Marshal(projectPayload.Sessions)
	var tasksBytes []byte
	tasksBytes, _ = json.Marshal(projectPayload.Tasks)
	var printingBytes []byte
	printingBytes, _ = json.Marshal(projectPayload.Printing)
	var deliverablesBytes []byte
	deliverablesBytes, _ = json.Marshal(projectPayload.Deliverables)
	var expensesBytes []byte
	expensesBytes, _ = json.Marshal(projectPayload.Expenses)
	var financeBytes []byte
	financeBytes, _ = json.Marshal(projectPayload.Finance)

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE projects SET
			project_name = ?, client_name = ?, client_initials = ?, client_mobile = ?,
			client_country_code = ?, client_email = ?, client_address = ?, lead_source = ?,
			project_manager = ?, project_type = ?, project_address = ?, status = ?, stage = ?,
			start_date = ?, end_date = ?, shoot_hours = ?, shoot_members_count = ?,
			shoot_expense = ?, galleries_count = ?, notes = ?, segments = ?, documents = ?,
			sessions = ?, tasks = ?, printing = ?, deliverables = ?, expenses = ?, finance = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		projectPayload.ProjectName, projectPayload.ClientName, projectPayload.ClientInitials, projectPayload.ClientMobile,
		projectPayload.ClientCountryCode, projectPayload.ClientEmail, projectPayload.ClientAddress, projectPayload.LeadSource,
		managerBytes, projectPayload.ProjectType, projectPayload.ProjectAddress, projectPayload.Status, projectPayload.Stage,
		projectPayload.StartDate, projectPayload.EndDate, projectPayload.ShootHours, projectPayload.ShootMembersCount,
		projectPayload.ShootExpense, projectPayload.GalleriesCount, projectPayload.Notes, segmentsBytes, documentsBytes,
		sessionsBytes, tasksBytes, printingBytes, deliverablesBytes, expensesBytes, financeBytes,
		projectID, tenantID,
	)

	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update project", updateErr)
		return
	}

	projectPayload.ID = projectID
	utils.WriteJSON(responseWriter, http.StatusOK, projectPayload)
}

func UpdateProjectStageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var projectID string = chi.URLParam(request, "id")

	var requestPayload UpdateProjectStageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var query string = "UPDATE projects SET stage = ?"
	var queryArguments = []any{requestPayload.Stage}
	if requestPayload.Status != "" {
		query += ", status = ?"
		queryArguments = append(queryArguments, requestPayload.Status)
	}
	query += " WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL"
	queryArguments = append(queryArguments, projectID, tenantID)

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), query, queryArguments...)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update project stage", updateErr)
		return
	}

	var responsePayload = ProjectStageUpdateResponse{
		Status: "ok",
		ID:     projectID,
		Stage:  requestPayload.Stage,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

func DeleteProjectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var projectID string = chi.URLParam(request, "id")

	var deleteErr error
	_, deleteErr = database.ExecContext(request.Context(), `
		UPDATE projects SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, projectID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete project", deleteErr)
		return
	}

	_, _ = database.ExecContext(request.Context(), `
		UPDATE project_sessions SET deleted_at = NOW() WHERE project_id = ? AND tenant_id = ?
	`, projectID, tenantID)

	var responsePayload = ProjectDeleteResponse{
		Status: "deleted",
		ID:     projectID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}
