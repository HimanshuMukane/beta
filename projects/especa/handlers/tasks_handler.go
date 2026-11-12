package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChecklistItem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TaskDTO struct {
	ID           string          `json:"id"`
	TaskId       string          `json:"taskId"`
	TenantID     string          `json:"tenantId"`
	ProjectID    string          `json:"projectId"`
	TaskName     string          `json:"taskName"`
	CategoryName string          `json:"categoryName"`
	Phase        string          `json:"phase"`
	Status       string          `json:"status"`
	Priority     string          `json:"priority"`
	DueDate      string          `json:"dueDate"`
	Assignee     string          `json:"assignee"`
	Notes        string          `json:"notes"`
	Checklist    []ChecklistItem `json:"checklist"`
	OrderIndex   int             `json:"orderIndex"`
}

type CreateTaskRequest struct {
	ProjectID    string          `json:"projectId"`
	TaskName     string          `json:"taskName"`
	CategoryName string          `json:"categoryName"`
	Phase        string          `json:"phase"`
	Status       string          `json:"status"`
	Priority     string          `json:"priority"`
	DueDate      string          `json:"dueDate"`
	Assignee     string          `json:"assignee"`
	Notes        string          `json:"notes"`
	Checklist    []ChecklistItem `json:"checklist"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}

type TaskStatusUpdateResponse struct {
	Status     string `json:"status"`
	ID         string `json:"id"`
	TaskStatus string `json:"taskStatus"`
}

type TaskDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

var defaultPreProductionTasks = []struct {
	category string
	tasks    []string
}{
	{category: "Client Information", tasks: []string{"Collect details", "Event details", "Notes"}},
	{category: "Mood Board & Theme", tasks: []string{"Discuss style & theme", "Create mood board", "Prepare shot list"}},
	{category: "Schedule & Crew", tasks: []string{"Create Sessions", "Assign Crew", "Schedule pre-shoot meeting"}},
	{category: "Equipment Checklist", tasks: []string{"Check cameras & lenses", "Format memory cards", "Charge batteries", "Test audio & lights"}},
}

var defaultPostProductionTasks = []struct {
	category string
	tasks    []string
}{
	{category: "Data Management", tasks: []string{"Transfer all File", "Upload to cloud storage", "Organize files"}},
	{category: "Culling & Selection", tasks: []string{"Cull raw photos", "Select highlights", "Create web gallery", "Send gallery to client"}},
	{category: "Editing & Retouching", tasks: []string{"Color grading", "Retouch portraits", "Apply studio presets", "Quality check"}},
	{category: "Album Design", tasks: []string{"Layout draft album", "Client album proofing", "Make revisions", "Send to print lab"}},
	{category: "Delivery", tasks: []string{"Export high-res photos", "Prepare USB drive", "Package physical prints", "Final delivery to client"}},
}

func seedDefaultProjectTasks(database *sqlx.DB, tenantID string, projectID string) {
	var taskCount int = 0
	var countErr error = database.QueryRow("SELECT COUNT(*) FROM project_tasks WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL", tenantID, projectID).Scan(&taskCount)
	if countErr != nil || taskCount > 0 {
		return
	}

	var orderCounter int = 0
	var catIndex int
	for catIndex = 0; catIndex < len(defaultPreProductionTasks); catIndex++ {
		var cat = defaultPreProductionTasks[catIndex]
		var taskIndex int
		for taskIndex = 0; taskIndex < len(cat.tasks); taskIndex++ {
			var taskTitle string = cat.tasks[taskIndex]
			var taskID string = uuid.New().String()
			_, _ = database.Exec(`
				INSERT INTO project_tasks (id, tenant_id, project_id, task_name, category_name, phase, status, priority, order_index)
				VALUES (?, ?, ?, ?, ?, 'Pre-production', 'Pending', 'Medium', ?)
			`, taskID, tenantID, projectID, taskTitle, cat.category, orderCounter)
			orderCounter++
		}
	}

	for catIndex = 0; catIndex < len(defaultPostProductionTasks); catIndex++ {
		var cat = defaultPostProductionTasks[catIndex]
		var taskIndex int
		for taskIndex = 0; taskIndex < len(cat.tasks); taskIndex++ {
			var taskTitle string = cat.tasks[taskIndex]
			var taskID string = uuid.New().String()
			_, _ = database.Exec(`
				INSERT INTO project_tasks (id, tenant_id, project_id, task_name, category_name, phase, status, priority, order_index)
				VALUES (?, ?, ?, ?, ?, 'Post-production', 'Pending', 'Medium', ?)
			`, taskID, tenantID, projectID, taskTitle, cat.category, orderCounter)
			orderCounter++
		}
	}
}

func ProjectTasksListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var projectID string = chi.URLParam(request, "projectId")
	if projectID == "" {
		projectID = strings.TrimSpace(request.URL.Query().Get("projectId"))
	}

	if projectID != "" {
		seedDefaultProjectTasks(database, tenantID, projectID)
	}

	var query string = `
		SELECT id, project_id, task_name, category_name, phase, status, priority,
		       COALESCE(due_date, ''), COALESCE(assignee, ''), COALESCE(notes, ''),
		       checklist, order_index
		FROM project_tasks
		WHERE tenant_id = ? AND deleted_at IS NULL
	`
	var queryArguments = []any{tenantID}
	if projectID != "" {
		query += " AND project_id = ?"
		queryArguments = append(queryArguments, projectID)
	}
	var phaseFilter string = strings.TrimSpace(request.URL.Query().Get("phase"))
	if phaseFilter != "" && !strings.EqualFold(phaseFilter, "all") {
		query += " AND phase = ?"
		queryArguments = append(queryArguments, phaseFilter)
	}
	var statusFilter string = strings.TrimSpace(request.URL.Query().Get("status"))
	if statusFilter != "" && !strings.EqualFold(statusFilter, "all") {
		query += " AND status = ?"
		queryArguments = append(queryArguments, statusFilter)
	}
	var priorityFilter string = strings.TrimSpace(request.URL.Query().Get("priority"))
	if priorityFilter != "" && !strings.EqualFold(priorityFilter, "all") {
		query += " AND priority = ?"
		queryArguments = append(queryArguments, priorityFilter)
	}
	var assigneeFilter string = strings.TrimSpace(request.URL.Query().Get("assignee"))
	if assigneeFilter != "" && !strings.EqualFold(assigneeFilter, "all") {
		query += " AND assignee = ?"
		queryArguments = append(queryArguments, assigneeFilter)
	}
	var categoryFilter string = strings.TrimSpace(request.URL.Query().Get("category"))
	if categoryFilter != "" && !strings.EqualFold(categoryFilter, "all") {
		query += " AND category_name = ?"
		queryArguments = append(queryArguments, categoryFilter)
	}

	query += " ORDER BY order_index ASC, created_at ASC"

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), query, queryArguments...)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list project tasks", queryErr)
		return
	}
	defer rows.Close()

	var tasks = make([]TaskDTO, 0)
	for rows.Next() {
		var taskID string
		var projID string
		var taskName string
		var categoryName string
		var phase string
		var status string
		var priority string
		var dueDate string
		var assignee string
		var notes string
		var checklistRaw []byte
		var orderIndex int

		var scanErr error = rows.Scan(
			&taskID, &projID, &taskName, &categoryName, &phase, &status, &priority,
			&dueDate, &assignee, &notes, &checklistRaw, &orderIndex,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read task record", scanErr)
			return
		}

		var checklistItems = make([]ChecklistItem, 0)
		if len(checklistRaw) > 0 {
			_ = json.Unmarshal(checklistRaw, &checklistItems)
		}

		tasks = append(tasks, TaskDTO{
			ID:           taskID,
			TaskId:       taskID,
			TenantID:     tenantID,
			ProjectID:    projID,
			TaskName:     taskName,
			CategoryName: categoryName,
			Phase:        phase,
			Status:       status,
			Priority:     priority,
			DueDate:      dueDate,
			Assignee:     assignee,
			Notes:        notes,
			Checklist:    checklistItems,
			OrderIndex:   orderIndex,
		})
	}

	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating tasks", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, tasks)
}

func CreateProjectTaskHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var requestPayload CreateTaskRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var projectID string = chi.URLParam(request, "projectId")
	if projectID == "" {
		projectID = chi.URLParam(request, "id")
	}
	if projectID == "" {
		projectID = strings.TrimSpace(requestPayload.ProjectID)
	}
	if projectID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Project ID is required")
		return
	}

	var taskName string = strings.TrimSpace(requestPayload.TaskName)
	if taskName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Task name is required")
		return
	}

	var categoryName string = strings.TrimSpace(requestPayload.CategoryName)
	if categoryName == "" {
		categoryName = "General"
	}
	var phase string = strings.TrimSpace(requestPayload.Phase)
	if phase == "" {
		phase = "Pre-production"
	}
	var status string = strings.TrimSpace(requestPayload.Status)
	if status == "" {
		status = "Pending"
	}
	var priority string = strings.TrimSpace(requestPayload.Priority)
	if priority == "" {
		priority = "Medium"
	}

	var checklistBytes []byte
	var marshalErr error
	if len(requestPayload.Checklist) > 0 {
		checklistBytes, marshalErr = json.Marshal(requestPayload.Checklist)
		if marshalErr != nil {
			checklistBytes = []byte("[]")
		}
	} else {
		checklistBytes = []byte("[]")
	}

	var taskID string = uuid.New().String()
	var insertErr error
	_, insertErr = database.ExecContext(request.Context(), `
		INSERT INTO project_tasks (id, tenant_id, project_id, task_name, category_name, phase, status, priority, due_date, assignee, notes, checklist, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 999)
	`, taskID, tenantID, projectID, taskName, categoryName, phase, status, priority, requestPayload.DueDate, requestPayload.Assignee, requestPayload.Notes, checklistBytes)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create project task", insertErr)
		return
	}

	var createdTask = TaskDTO{
		ID:           taskID,
		TaskId:       taskID,
		TenantID:     tenantID,
		ProjectID:    projectID,
		TaskName:     taskName,
		CategoryName: categoryName,
		Phase:        phase,
		Status:       status,
		Priority:     priority,
		DueDate:      requestPayload.DueDate,
		Assignee:     requestPayload.Assignee,
		Notes:        requestPayload.Notes,
		Checklist:    requestPayload.Checklist,
		OrderIndex:   999,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdTask)
}

func GetTaskHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var taskID string = chi.URLParam(request, "id")
	if taskID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Task ID is required")
		return
	}

	var projID, taskName, categoryName, phase, status, priority, dueDate, assignee, notes string
	var checklistRaw []byte
	var orderIndex int

	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT project_id, task_name, category_name, phase, status, priority,
		       COALESCE(due_date, ''), COALESCE(assignee, ''), COALESCE(notes, ''),
		       checklist, order_index
		FROM project_tasks
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, taskID, tenantID).Scan(
		&projID, &taskName, &categoryName, &phase, &status, &priority,
		&dueDate, &assignee, &notes, &checklistRaw, &orderIndex,
	)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Task not found")
		return
	}

	var checklistItems = make([]ChecklistItem, 0)
	if len(checklistRaw) > 0 {
		_ = json.Unmarshal(checklistRaw, &checklistItems)
	}

	var task = TaskDTO{
		ID:           taskID,
		TaskId:       taskID,
		TenantID:     tenantID,
		ProjectID:    projID,
		TaskName:     taskName,
		CategoryName: categoryName,
		Phase:        phase,
		Status:       status,
		Priority:     priority,
		DueDate:      dueDate,
		Assignee:     assignee,
		Notes:        notes,
		Checklist:    checklistItems,
		OrderIndex:   orderIndex,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, task)
}

func UpdateTaskHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var taskID string = chi.URLParam(request, "id")
	if taskID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Task ID is required")
		return
	}

	var requestPayload CreateTaskRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var checklistBytes []byte
	var marshalErr error
	if len(requestPayload.Checklist) > 0 {
		checklistBytes, marshalErr = json.Marshal(requestPayload.Checklist)
		if marshalErr != nil {
			checklistBytes = []byte("[]")
		}
	} else {
		checklistBytes = []byte("[]")
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE project_tasks
		SET task_name = COALESCE(NULLIF(?, ''), task_name),
		    category_name = COALESCE(NULLIF(?, ''), category_name),
		    phase = COALESCE(NULLIF(?, ''), phase),
		    status = COALESCE(NULLIF(?, ''), status),
		    priority = COALESCE(NULLIF(?, ''), priority),
		    due_date = ?,
		    assignee = ?,
		    notes = ?,
		    checklist = ?,
		    updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, requestPayload.TaskName, requestPayload.CategoryName, requestPayload.Phase,
		requestPayload.Status, requestPayload.Priority, requestPayload.DueDate,
		requestPayload.Assignee, requestPayload.Notes, checklistBytes, taskID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update task", updateErr)
		return
	}

	var updatedTask TaskDTO
	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT id, project_id, task_name, category_name, phase, status, priority,
		       COALESCE(due_date, ''), COALESCE(assignee, ''), COALESCE(notes, ''),
		       checklist, order_index
		FROM project_tasks
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, taskID, tenantID).Scan(
		&updatedTask.ID, &updatedTask.ProjectID, &updatedTask.TaskName, &updatedTask.CategoryName,
		&updatedTask.Phase, &updatedTask.Status, &updatedTask.Priority, &updatedTask.DueDate,
		&updatedTask.Assignee, &updatedTask.Notes, &checklistBytes, &updatedTask.OrderIndex,
	)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to fetch updated task", queryErr)
		return
	}
	updatedTask.TaskId = updatedTask.ID
	updatedTask.TenantID = tenantID
	_ = json.Unmarshal(checklistBytes, &updatedTask.Checklist)

	utils.WriteJSON(responseWriter, http.StatusOK, updatedTask)
}

func UpdateTaskStatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var taskID string = chi.URLParam(request, "id")
	if taskID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Task ID is required")
		return
	}

	var requestPayload UpdateTaskStatusRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var targetStatus string = strings.TrimSpace(requestPayload.Status)
	if targetStatus == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Status is required")
		return
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE project_tasks
		SET status = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, targetStatus, taskID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update task status", updateErr)
		return
	}

	var successResponse = TaskStatusUpdateResponse{
		Status:     "ok",
		ID:         taskID,
		TaskStatus: targetStatus,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func DeleteTaskHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var taskID string = chi.URLParam(request, "id")
	if taskID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Task ID is required")
		return
	}

	var deleteErr error
	_, deleteErr = database.ExecContext(request.Context(), `
		UPDATE project_tasks
		SET deleted_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, taskID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete task", deleteErr)
		return
	}

	var successResponse = TaskDeleteResponse{
		Status: "ok",
		ID:     taskID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}
