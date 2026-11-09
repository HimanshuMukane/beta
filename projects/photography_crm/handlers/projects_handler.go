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

type ProjectDTO struct {
	ID                 string   `json:"id"`
	ProjectName        string   `json:"projectName"`
	ClientName         string   `json:"clientName"`
	ClientInitials     string   `json:"clientInitials"`
	ClientMobile       string   `json:"clientMobile"`
	ClientCountryCode  string   `json:"clientCountryCode"`
	ClientEmail        string   `json:"clientEmail"`
	ClientAddress      string   `json:"clientAddress,omitempty"`
	LeadSource         string   `json:"leadSource,omitempty"`
	ProjectManager     any      `json:"projectManager"`
	ProjectType        string   `json:"projectType"`
	ProjectAddress     string   `json:"projectAddress,omitempty"`
	SessionsCount      int      `json:"sessionsCount"`
	TasksCount         any      `json:"tasksCount"`
	DeliverablesCount  any      `json:"deliverablesCount"`
	Status             string   `json:"status"`
	Stage              string   `json:"stage"`
	StartDate          string   `json:"startDate"`
	EndDate            string   `json:"endDate"`
	CreatedDate        string   `json:"createdDate"`
	Segments           any      `json:"segments"`
	Documents          any      `json:"documents"`
	Sessions           any      `json:"sessions"`
	Tasks              any      `json:"tasks"`
	Printing           any      `json:"printing"`
	Deliverables       any      `json:"deliverables"`
	Expenses           any      `json:"expenses"`
	Notes              string   `json:"notes"`
	ShootHours         float64  `json:"shootHours"`
	ShootMembersCount  int      `json:"shootMembersCount"`
	ShootExpense       int64    `json:"shootExpense"`
	GalleriesCount     int      `json:"galleriesCount"`
	Finance            any      `json:"finance"`
}

func ProjectsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, client_id, client_name, client_mobile, client_email, client_address,
		       lead_source, project_manager, project_type, project_address, title,
		       status, stage, start_date, end_date, shoot_hours, shoot_members_count,
		       shoot_expense, galleries_count, notes, segments, finance_snapshot,
		       total_amount_minor, created_at
		FROM projects
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []ProjectDTO
	for rows.Next() {
		var (
			id, clientID, clientName, clientMobile, clientEmail, clientAddr string
			leadSource, projManager, projType, projAddr, title             string
			status, stage, startDate, endDate, notes                        string
			shootHours                                                     float64
			shootMembersCount, galleriesCount                              int
			shootExpense, totalAmountMinor                                 int64
			createdAt                                                      time.Time
			segmentsRaw, financeRaw                                        []byte
		)
		var (
			sClientName, sClientMobile, sClientEmail, sClientAddr sql.NullString
			sLeadSource, sProjManager, sProjAddr, sStartDate, sEndDate, sNotes sql.NullString
		)

		err := rows.Scan(
			&id, &clientID, &sClientName, &sClientMobile, &sClientEmail, &sClientAddr,
			&sLeadSource, &sProjManager, &projType, &sProjAddr, &title,
			&status, &stage, &sStartDate, &sEndDate, &shootHours, &shootMembersCount,
			&shootExpense, &galleriesCount, &sNotes, &segmentsRaw, &financeRaw,
			&totalAmountMinor, &createdAt,
		)
		if err != nil {
			continue
		}

		if sClientName.Valid {
			clientName = sClientName.String
		}
		if sClientMobile.Valid {
			clientMobile = sClientMobile.String
		}
		if sClientEmail.Valid {
			clientEmail = sClientEmail.String
		}
		if sClientAddr.Valid {
			clientAddr = sClientAddr.String
		}
		if sLeadSource.Valid {
			leadSource = sLeadSource.String
		}
		if sProjManager.Valid {
			projManager = sProjManager.String
		} else {
			projManager = "especa"
		}
		if sProjAddr.Valid {
			projAddr = sProjAddr.String
		}
		if sStartDate.Valid {
			startDate = sStartDate.String
		}
		if sEndDate.Valid {
			endDate = sEndDate.String
		}
		if sNotes.Valid {
			notes = sNotes.String
		}

		initials := "CL"
		if len(clientName) > 0 {
			parts := strings.Split(clientName, " ")
			if len(parts) >= 2 {
				initials = strings.ToUpper(parts[0][:1] + parts[1][:1])
			} else {
				initials = strings.ToUpper(parts[0][:1])
			}
		}

		var segments any = []any{}
		parseJSONColumn(segmentsRaw, &segments)

		var finance any = map[string]any{
			"profit":              0,
			"currentProfit":       0,
			"netPayableAmount":    totalAmountMinor / 100,
			"paidAmount":          0,
			"balanceAmount":       totalAmountMinor / 100,
			"paymentDue":          0,
			"expense":             0,
			"eventDiscountAmount": 0,
			"paymentHistory":      []any{},
		}
		parseJSONColumn(financeRaw, &finance)

		// Count sessions, tasks, deliverables for this project
		var sessionsCount int
		_ = db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM project_sessions WHERE project_id = ? OR project_name = ?", id, title).Scan(&sessionsCount)

		var tasksTotal, tasksCompleted int
		_ = db.QueryRowContext(r.Context(), "SELECT COUNT(*), COALESCE(SUM(completed), 0) FROM project_tasks WHERE project_id = ?", id).Scan(&tasksTotal, &tasksCompleted)

		var delivTotal, delivCompleted int
		_ = db.QueryRowContext(r.Context(), "SELECT COUNT(*), COALESCE(SUM(status = 'Completed' OR status = 'Delivered'), 0) FROM customer_deliverables WHERE project_id = ?", id).Scan(&delivTotal, &delivCompleted)

		result = append(result, ProjectDTO{
			ID:                id,
			ProjectName:       title,
			ClientName:        clientName,
			ClientInitials:    initials,
			ClientMobile:      clientMobile,
			ClientCountryCode: "+91",
			ClientEmail:       clientEmail,
			ClientAddress:     clientAddr,
			LeadSource:        leadSource,
			ProjectManager: map[string]string{
				"name":     projManager,
				"initials": "ES",
			},
			ProjectType:       projType,
			ProjectAddress:    projAddr,
			SessionsCount:     sessionsCount,
			TasksCount:        map[string]int{"completed": tasksCompleted, "total": tasksTotal},
			DeliverablesCount: map[string]int{"ready": delivCompleted, "total": delivTotal},
			Status:            status,
			Stage:             stage,
			StartDate:         startDate,
			EndDate:           endDate,
			CreatedDate:       createdAt.Format("02/01/2006"),
			Segments:          segments,
			Documents:         []any{},
			Sessions:          []any{},
			Tasks:             []any{},
			Printing:          []any{},
			Deliverables:      []any{},
			Expenses:          []any{},
			Notes:             notes,
			ShootHours:        shootHours,
			ShootMembersCount: shootMembersCount,
			ShootExpense:      shootExpense,
			GalleriesCount:    galleriesCount,
			Finance:           finance,
		})
	}

	if result == nil {
		result = []ProjectDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func GetProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		id = r.URL.Query().Get("id")
	}

	var (
		clientID, clientName, clientMobile, clientEmail, clientAddr string
		leadSource, projManager, projType, projAddr, title          string
		status, stage, startDate, endDate, notes                     string
		shootHours                                                  float64
		shootMembersCount, galleriesCount                           int
		shootExpense, totalAmountMinor                              int64
		createdAt                                                   time.Time
		segmentsRaw, financeRaw                                     []byte
	)
	var (
		sClientName, sClientMobile, sClientEmail, sClientAddr sql.NullString
		sLeadSource, sProjManager, sProjAddr, sStartDate, sEndDate, sNotes sql.NullString
	)

	err := db.QueryRowContext(r.Context(), `
		SELECT id, client_id, client_name, client_mobile, client_email, client_address,
		       lead_source, project_manager, project_type, project_address, title,
		       status, stage, start_date, end_date, shoot_hours, shoot_members_count,
		       shoot_expense, galleries_count, notes, segments, finance_snapshot,
		       total_amount_minor, created_at
		FROM projects
		WHERE id = ? AND tenant_id = ?
	`, id, tid).Scan(
		&id, &clientID, &sClientName, &sClientMobile, &sClientEmail, &sClientAddr,
		&sLeadSource, &sProjManager, &projType, &sProjAddr, &title,
		&status, &stage, &sStartDate, &sEndDate, &shootHours, &shootMembersCount,
		&shootExpense, &galleriesCount, &sNotes, &segmentsRaw, &financeRaw,
		&totalAmountMinor, &createdAt,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
		return
	}

	if sClientName.Valid {
		clientName = sClientName.String
	}
	if sClientMobile.Valid {
		clientMobile = sClientMobile.String
	}
	if sClientEmail.Valid {
		clientEmail = sClientEmail.String
	}
	if sClientAddr.Valid {
		clientAddr = sClientAddr.String
	}
	if sLeadSource.Valid {
		leadSource = sLeadSource.String
	}
	if sProjManager.Valid {
		projManager = sProjManager.String
	} else {
		projManager = "especa"
	}
	if sProjAddr.Valid {
		projAddr = sProjAddr.String
	}
	if sStartDate.Valid {
		startDate = sStartDate.String
	}
	if sEndDate.Valid {
		endDate = sEndDate.String
	}
	if sNotes.Valid {
		notes = sNotes.String
	}

	var segments any = []any{}
	parseJSONColumn(segmentsRaw, &segments)

	var finance any = map[string]any{
		"profit":              0,
		"currentProfit":       0,
		"netPayableAmount":    totalAmountMinor / 100,
		"paidAmount":          0,
		"balanceAmount":       totalAmountMinor / 100,
		"paymentDue":          0,
		"expense":             0,
		"eventDiscountAmount": 0,
		"paymentHistory":      []any{},
	}
	parseJSONColumn(financeRaw, &finance)

	// Fetch child sessions
	var sessions []any
	sesRows, err := db.QueryContext(r.Context(), `
		SELECT id, title, project_name, session_manager, project_type, shoot_timing,
		       status, session_date, start_time, duration, time_range, venue, city,
		       location_link, remarks, members, shoots_activity_count
		FROM project_sessions
		WHERE project_id = ? OR project_name = ?
	`, id, title)
	if err == nil {
		defer sesRows.Close()
		for sesRows.Next() {
			var sesID, sesTitle, sesProjName, sesMgr, sesType, sesTiming, sesStatus string
			var sesDate, sesStart, sesDur, sesRange, sesVenue, sesCity, sesLink, sesRem string
			var sesActivity int
			var membersRaw []byte
			var (
				sSesStart, sSesDur, sSesRange, sSesVenue, sSesCity, sSesLink, sSesRem sql.NullString
			)
			if err := sesRows.Scan(
				&sesID, &sesTitle, &sesProjName, &sesMgr, &sesType, &sesTiming,
				&sesStatus, &sesDate, &sSesStart, &sSesDur, &sSesRange, &sSesVenue, &sSesCity,
				&sSesLink, &sSesRem, &membersRaw, &sesActivity,
			); err == nil {
				if sSesStart.Valid {
					sesStart = sSesStart.String
				}
				if sSesDur.Valid {
					sesDur = sSesDur.String
				}
				if sSesRange.Valid {
					sesRange = sSesRange.String
				}
				if sSesVenue.Valid {
					sesVenue = sSesVenue.String
				}
				if sSesCity.Valid {
					sesCity = sSesCity.String
				}
				if sSesLink.Valid {
					sesLink = sSesLink.String
				}
				if sSesRem.Valid {
					sesRem = sSesRem.String
				}
				var members []string
				parseJSONColumn(membersRaw, &members)
				sessions = append(sessions, map[string]any{
					"id":                  sesID,
					"title":               sesTitle,
					"projectName":         sesProjName,
					"sessionManager":      sesMgr,
					"projectType":         sesType,
					"shootTiming":         sesTiming,
					"status":              sesStatus,
					"date":                sesDate,
					"startTime":           sesStart,
					"duration":            sesDur,
					"timeRange":           sesRange,
					"venue":               sesVenue,
					"city":                sesCity,
					"locationLink":        sesLink,
					"remarks":             sesRem,
					"members":             members,
					"shootsActivityCount": sesActivity,
				})
			}
		}
	}

	// Fetch child tasks
	var tasks []any
	taskRows, err := db.QueryContext(r.Context(), "SELECT id, title, assignee, due_date, completed, priority FROM project_tasks WHERE project_id = ?", id)
	if err == nil {
		defer taskRows.Close()
		for taskRows.Next() {
			var tID, tTitle, tAssignee, tDueDate, tPriority string
			var tCompleted bool
			if err := taskRows.Scan(&tID, &tTitle, &tAssignee, &tDueDate, &tCompleted, &tPriority); err == nil {
				tasks = append(tasks, map[string]any{
					"id":        tID,
					"title":     tTitle,
					"assignee":  tAssignee,
					"dueDate":   tDueDate,
					"completed": tCompleted,
					"priority":  tPriority,
				})
			}
		}
	}

	// Fetch child deliverables
	var deliverables []any
	delivRows, err := db.QueryContext(r.Context(), "SELECT id, deliverable, due_date, assignee, status, note FROM customer_deliverables WHERE project_id = ?", id)
	if err == nil {
		defer delivRows.Close()
		for delivRows.Next() {
			var dID, dName, dDueDate, dAssignee, dStatus string
			var sNote sql.NullString
			if err := delivRows.Scan(&dID, &dName, &dDueDate, &dAssignee, &dStatus, &sNote); err == nil {
				deliverables = append(deliverables, map[string]any{
					"id":          dID,
					"deliverable": dName,
					"dueDate":     dDueDate,
					"assignee":    dAssignee,
					"status":      dStatus,
					"note":        sNote.String,
				})
			}
		}
	}

	// Fetch child expenses
	var expenses []any
	expRows, err := db.QueryContext(r.Context(), "SELECT id, category, description, amount, expense_date, paid_by FROM project_expenses WHERE project_id = ?", id)
	if err == nil {
		defer expRows.Close()
		for expRows.Next() {
			var eID, eCat, eDesc, eDate, ePaidBy string
			var eAmount int64
			if err := expRows.Scan(&eID, &eCat, &eDesc, &eAmount, &eDate, &ePaidBy); err == nil {
				expenses = append(expenses, map[string]any{
					"id":          eID,
					"category":    eCat,
					"description": eDesc,
					"amount":      eAmount,
					"date":        eDate,
					"paidBy":      ePaidBy,
				})
			}
		}
	}

	// Fetch child documents
	var documents []any
	docRows, err := db.QueryContext(r.Context(), "SELECT id, type, title, reference_no, amount, status, doc_date FROM project_documents WHERE project_id = ?", id)
	if err == nil {
		defer docRows.Close()
		for docRows.Next() {
			var docID, docType, docTitle, docRef, docStatus, docDate string
			var docAmount int64
			var sRef sql.NullString
			if err := docRows.Scan(&docID, &docType, &docTitle, &sRef, &docAmount, &docStatus, &docDate); err == nil {
				if sRef.Valid {
					docRef = sRef.String
				}
				documents = append(documents, map[string]any{
					"id":          docID,
					"type":        docType,
					"title":       docTitle,
					"referenceNo": docRef,
					"amount":      docAmount,
					"status":      docStatus,
					"date":        docDate,
				})
			}
		}
	}

	initials := "CL"
	if len(clientName) > 0 {
		parts := strings.Split(clientName, " ")
		if len(parts) >= 2 {
			initials = strings.ToUpper(parts[0][:1] + parts[1][:1])
		} else {
			initials = strings.ToUpper(parts[0][:1])
		}
	}

	utils.WriteJSON(w, http.StatusOK, ProjectDTO{
		ID:                id,
		ProjectName:       title,
		ClientName:        clientName,
		ClientInitials:    initials,
		ClientMobile:      clientMobile,
		ClientCountryCode: "+91",
		ClientEmail:       clientEmail,
		ClientAddress:     clientAddr,
		LeadSource:        leadSource,
		ProjectManager: map[string]string{
			"name":     projManager,
			"initials": "ES",
		},
		ProjectType:       projType,
		ProjectAddress:    projAddr,
		SessionsCount:     len(sessions),
		TasksCount:        map[string]int{"completed": 0, "total": len(tasks)},
		DeliverablesCount: map[string]int{"ready": 0, "total": len(deliverables)},
		Status:            status,
		Stage:             stage,
		StartDate:         startDate,
		EndDate:           endDate,
		CreatedDate:       createdAt.Format("02/01/2006"),
		Segments:          segments,
		Documents:         documents,
		Sessions:          sessions,
		Tasks:             tasks,
		Printing:          []any{},
		Deliverables:      deliverables,
		Expenses:          expenses,
		Notes:             notes,
		ShootHours:        shootHours,
		ShootMembersCount: shootMembersCount,
		ShootExpense:      shootExpense,
		GalleriesCount:    galleriesCount,
		Finance:           finance,
	})
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var payload ProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if payload.ID == "" {
		payload.ID = uuid.New().String()
	}
	if payload.Status == "" {
		payload.Status = "Active"
	}
	if payload.Stage == "" {
		payload.Stage = "Booked"
	}
	if payload.ProjectType == "" {
		payload.ProjectType = "Wedding"
	}

	clientID := uuid.New().String()
	names := strings.SplitN(strings.TrimSpace(payload.ClientName), " ", 2)
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = names[1]
	}

	_, _ = db.ExecContext(r.Context(), `
		INSERT INTO clients (id, tenant_id, first_name, last_name, email, phone, source)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, clientID, tid, firstName, lastName, payload.ClientEmail, payload.ClientMobile, payload.LeadSource)

	segmentsJSON := toJSONBytes(payload.Segments)
	financeJSON := toJSONBytes(payload.Finance)

	mgrName := "especa"
	if m, ok := payload.ProjectManager.(map[string]any); ok {
		if n, ok := m["name"].(string); ok && n != "" {
			mgrName = n
		}
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO projects (
			id, tenant_id, client_id, client_name, client_mobile, client_email, client_address,
			lead_source, project_manager, project_type, project_address, title, status, stage,
			start_date, end_date, shoot_hours, shoot_members_count, shoot_expense, galleries_count,
			notes, segments, finance_snapshot, total_amount_minor
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?
		)
	`, payload.ID, tid, clientID, payload.ClientName, payload.ClientMobile, payload.ClientEmail, payload.ClientAddress,
		payload.LeadSource, mgrName, payload.ProjectType, payload.ProjectAddress, payload.ProjectName, payload.Status, payload.Stage,
		payload.StartDate, payload.EndDate, payload.ShootHours, payload.ShootMembersCount, payload.ShootExpense, payload.GalleriesCount,
		payload.Notes, segmentsJSON, financeJSON, 0)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(w, http.StatusCreated, payload)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		id = r.URL.Query().Get("id")
	}

	var payload ProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	segmentsJSON := toJSONBytes(payload.Segments)
	financeJSON := toJSONBytes(payload.Finance)

	mgrName := "especa"
	if m, ok := payload.ProjectManager.(map[string]any); ok {
		if n, ok := m["name"].(string); ok && n != "" {
			mgrName = n
		}
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE projects SET
			client_name = ?, client_mobile = ?, client_email = ?, client_address = ?,
			lead_source = ?, project_manager = ?, project_type = ?, project_address = ?,
			title = ?, status = ?, stage = ?, start_date = ?, end_date = ?,
			shoot_hours = ?, shoot_members_count = ?, shoot_expense = ?, galleries_count = ?,
			notes = ?, segments = ?, finance_snapshot = ?
		WHERE id = ? AND tenant_id = ?
	`, payload.ClientName, payload.ClientMobile, payload.ClientEmail, payload.ClientAddress,
		payload.LeadSource, mgrName, payload.ProjectType, payload.ProjectAddress,
		payload.ProjectName, payload.Status, payload.Stage, payload.StartDate, payload.EndDate,
		payload.ShootHours, payload.ShootMembersCount, payload.ShootExpense, payload.GalleriesCount,
		payload.Notes, segmentsJSON, financeJSON,
		id, tid)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.ID = id
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateProjectStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var req struct {
		Stage  string `json:"stage"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	query := "UPDATE projects SET stage = ?"
	args := []any{req.Stage}
	if req.Status != "" {
		query += ", status = ?"
		args = append(args, req.Status)
	}
	query += " WHERE id = ? AND tenant_id = ?"
	args = append(args, id, tid)

	_, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": id, "stage": req.Stage})
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(), "DELETE FROM projects WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

