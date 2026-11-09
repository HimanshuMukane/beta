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

func getTenantID(r *http.Request, db *sqlx.DB) string {
	if tc := utils.TenantFromContext(r.Context()); tc != nil && tc.TenantID != uuid.Nil {
		return tc.TenantID.String()
	}
	if cookie, err := r.Cookie("active_tenant_id"); err == nil && cookie.Value != "" {
		if _, err := uuid.Parse(cookie.Value); err == nil {
			return cookie.Value
		}
	}
	if db != nil {
		var tid string
		if err := db.QueryRowContext(r.Context(), "SELECT tenant_id FROM leads WHERE tenant_id IS NOT NULL LIMIT 1").Scan(&tid); err == nil && tid != "" {
			return tid
		}
		if err := db.QueryRowContext(r.Context(), "SELECT tenant_id FROM projects WHERE tenant_id IS NOT NULL LIMIT 1").Scan(&tid); err == nil && tid != "" {
			return tid
		}
	}
	return "01a1b654-4c11-7a64-98e2-a13d108ace67"
}

func parseJSONColumn(data any, target any) {
	if data == nil {
		return
	}
	switch v := data.(type) {
	case []byte:
		if len(v) > 0 {
			_ = json.Unmarshal(v, target)
		}
	case string:
		if len(v) > 0 {
			_ = json.Unmarshal([]byte(v), target)
		}
	}
}

func toJSONBytes(data any) []byte {
	if data == nil {
		return []byte("[]")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return []byte("[]")
	}
	return b
}

type LeadDTO struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	MobileNo             string   `json:"mobileNo"`
	CountryCode          string   `json:"countryCode"`
	Email                string   `json:"email"`
	Address              string   `json:"address,omitempty"`
	CreatedDate          string   `json:"createdDate"`
	LeadSource           string   `json:"leadSource"`
	Generate             string   `json:"generate"`
	Assignee             any      `json:"assignee"`
	LeadPipeline         string   `json:"leadPipeline"`
	PipelineStage        string   `json:"pipelineStage"`
	LeadStage            string   `json:"leadStage"`
	QuoteAmount          int64    `json:"quoteAmount"`
	FollowupDate         string   `json:"followupDate"`
	ProjectType          string   `json:"projectType"`
	QuoteType            string   `json:"quoteType"`
	ProjectStartDate     string   `json:"projectStartDate"`
	ProjectEndDate       string   `json:"projectEndDate"`
	ProjectName          string   `json:"projectName,omitempty"`
	ProjectAddress       string   `json:"projectAddress,omitempty"`
	ProjectSegments      any      `json:"projectSegments"`
	PackageName          string   `json:"packageName,omitempty"`
	SelectedServices     any      `json:"selectedServices"`
	SelectedDeliverables []string `json:"selectedDeliverables"`
	QuotationTemplate    string   `json:"quotationTemplate,omitempty"`
	DiscountAmount       int64    `json:"discountAmount,omitempty"`
	TaxAmount            int64    `json:"taxAmount,omitempty"`
}

func LeadsHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, contact_name, contact_phone, country_code, contact_email, address,
		       created_at, source, generate_type, assignee_name, assignee_initials,
		       lead_pipeline, status, lead_stage, quote_amount, followup_date,
		       project_type, quote_type, project_start_date, project_end_date,
		       project_name, location, project_segments, package_name,
		       selected_services, selected_deliverables, quotation_template,
		       discount_amount, tax_amount
		FROM leads
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []LeadDTO
	for rows.Next() {
		var (
			id, contactName, contactPhone, countryCode, contactEmail, address string
			source, generateType, assigneeName, assigneeInitials               string
			leadPipeline, status, leadStage                                   string
			quoteAmount                                                       int64
			followupDate, projectType, quoteType                              string
			projStartDate, projEndDate, projectName, location                 string
			packageName, quotationTemplate                                    string
			discountAmount, taxAmount                                         int64
			createdAt                                                         time.Time
			segmentsRaw, servicesRaw, deliverablesRaw                         []byte
		)
		var (
			sContactPhone, sCountryCode, sContactEmail, sAddress sql.NullString
			sSource, sGenerateType, sAssigneeName, sAssigneeInit sql.NullString
			sLeadPipeline, sStatus, sLeadStage                   sql.NullString
			sFollowupDate, sProjectType, sQuoteType              sql.NullString
			sProjStartDate, sProjEndDate, sProjectName, sLoc     sql.NullString
			sPackageName, sQuotationTemplate                     sql.NullString
		)

		err := rows.Scan(
			&id, &contactName, &sContactPhone, &sCountryCode, &sContactEmail, &sAddress,
			&createdAt, &sSource, &sGenerateType, &sAssigneeName, &sAssigneeInit,
			&sLeadPipeline, &sStatus, &sLeadStage, &quoteAmount, &sFollowupDate,
			&sProjectType, &sQuoteType, &sProjStartDate, &sProjEndDate,
			&sProjectName, &sLoc, &segmentsRaw, &sPackageName,
			&servicesRaw, &deliverablesRaw, &sQuotationTemplate,
			&discountAmount, &taxAmount,
		)
		if err != nil {
			continue
		}

		if sContactPhone.Valid {
			contactPhone = sContactPhone.String
		}
		if sCountryCode.Valid {
			countryCode = sCountryCode.String
		} else {
			countryCode = "+91"
		}
		if sContactEmail.Valid {
			contactEmail = sContactEmail.String
		}
		if sAddress.Valid {
			address = sAddress.String
		}
		if sSource.Valid {
			source = sSource.String
		} else {
			source = "Manual Leads"
		}
		if sGenerateType.Valid {
			generateType = sGenerateType.String
		} else {
			generateType = "Direct Entry"
		}
		if sAssigneeName.Valid {
			assigneeName = sAssigneeName.String
		} else {
			assigneeName = "Rohan Mehta"
		}
		if sAssigneeInit.Valid {
			assigneeInitials = sAssigneeInit.String
		} else {
			assigneeInitials = "RM"
		}
		if sLeadPipeline.Valid {
			leadPipeline = sLeadPipeline.String
		} else {
			leadPipeline = "Wedding Inquiries"
		}
		if sStatus.Valid {
			status = sStatus.String
		} else {
			status = "new_leads"
		}
		if sLeadStage.Valid {
			leadStage = sLeadStage.String
		} else {
			leadStage = "New Lead"
		}
		if sFollowupDate.Valid {
			followupDate = sFollowupDate.String
		}
		if sProjectType.Valid {
			projectType = sProjectType.String
		} else {
			projectType = "Wedding"
		}
		if sQuoteType.Valid {
			quoteType = sQuoteType.String
		} else {
			quoteType = "Standard Package"
		}
		if sProjStartDate.Valid {
			projStartDate = sProjStartDate.String
		}
		if sProjEndDate.Valid {
			projEndDate = sProjEndDate.String
		}
		if sProjectName.Valid {
			projectName = sProjectName.String
		}
		if sLoc.Valid {
			location = sLoc.String
		}
		if sPackageName.Valid {
			packageName = sPackageName.String
		}
		if sQuotationTemplate.Valid {
			quotationTemplate = sQuotationTemplate.String
		}

		var segments any = []any{}
		parseJSONColumn(segmentsRaw, &segments)

		var services any = []any{}
		parseJSONColumn(servicesRaw, &services)

		var deliverables []string = []string{}
		parseJSONColumn(deliverablesRaw, &deliverables)

		result = append(result, LeadDTO{
			ID:          id,
			Name:        contactName,
			MobileNo:    contactPhone,
			CountryCode: countryCode,
			Email:       contactEmail,
			Address:     address,
			CreatedDate: createdAt.Format("02/01/2006"),
			LeadSource:  source,
			Generate:    generateType,
			Assignee: map[string]string{
				"name":     assigneeName,
				"initials": assigneeInitials,
			},
			LeadPipeline:         leadPipeline,
			PipelineStage:        status,
			LeadStage:            leadStage,
			QuoteAmount:          quoteAmount,
			FollowupDate:         followupDate,
			ProjectType:          projectType,
			QuoteType:            quoteType,
			ProjectStartDate:     projStartDate,
			ProjectEndDate:       projEndDate,
			ProjectName:          projectName,
			ProjectAddress:       location,
			ProjectSegments:      segments,
			PackageName:          packageName,
			SelectedServices:     services,
			SelectedDeliverables: deliverables,
			QuotationTemplate:    quotationTemplate,
			DiscountAmount:       discountAmount,
			TaxAmount:            taxAmount,
		})
	}

	if result == nil {
		result = []LeadDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateLeadHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var payload LeadDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if payload.ID == "" {
		payload.ID = "LD-" + uuid.New().String()[:8]
	}
	if payload.PipelineStage == "" {
		payload.PipelineStage = "new_leads"
	}
	if payload.LeadStage == "" {
		payload.LeadStage = "New Lead"
	}
	if payload.LeadPipeline == "" {
		payload.LeadPipeline = "Wedding Inquiries"
	}
	if payload.CountryCode == "" {
		payload.CountryCode = "+91"
	}
	if payload.Generate == "" {
		payload.Generate = "Direct Entry"
	}

	assigneeName := "Rohan Mehta"
	assigneeInit := "RM"
	if m, ok := payload.Assignee.(map[string]any); ok {
		if n, ok := m["name"].(string); ok && n != "" {
			assigneeName = n
		}
		if init, ok := m["initials"].(string); ok && init != "" {
			assigneeInit = init
		}
	}

	segmentsJSON := toJSONBytes(payload.ProjectSegments)
	servicesJSON := toJSONBytes(payload.SelectedServices)
	deliverablesJSON := toJSONBytes(payload.SelectedDeliverables)

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO leads (
			id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email, contact_phone,
			country_code, address, source, generate_type, assignee_name, assignee_initials,
			lead_pipeline, status, lead_stage, quote_amount, followup_date, project_type,
			quote_type, project_start_date, project_end_date, project_name, location,
			project_segments, package_name, selected_services, selected_deliverables,
			quotation_template, discount_amount, tax_amount
		) VALUES (
			?, ?, 'wf-default', 'stage-new', ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
	`,
		payload.ID, tid, payload.Name, payload.Name, payload.Email, payload.MobileNo,
		payload.CountryCode, payload.Address, payload.LeadSource, payload.Generate, assigneeName, assigneeInit,
		payload.LeadPipeline, payload.PipelineStage, payload.LeadStage, payload.QuoteAmount, payload.FollowupDate, payload.ProjectType,
		payload.QuoteType, payload.ProjectStartDate, payload.ProjectEndDate, payload.ProjectName, payload.ProjectAddress,
		segmentsJSON, payload.PackageName, servicesJSON, deliverablesJSON,
		payload.QuotationTemplate, payload.DiscountAmount, payload.TaxAmount,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(w, http.StatusCreated, payload)
}

func UpdateLeadHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	leadID := chi.URLParam(r, "id")
	if leadID == "" {
		leadID = r.URL.Query().Get("id")
	}

	var payload LeadDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	assigneeName := "Rohan Mehta"
	assigneeInit := "RM"
	if m, ok := payload.Assignee.(map[string]any); ok {
		if n, ok := m["name"].(string); ok && n != "" {
			assigneeName = n
		}
		if init, ok := m["initials"].(string); ok && init != "" {
			assigneeInit = init
		}
	}

	segmentsJSON := toJSONBytes(payload.ProjectSegments)
	servicesJSON := toJSONBytes(payload.SelectedServices)
	deliverablesJSON := toJSONBytes(payload.SelectedDeliverables)

	_, err := db.ExecContext(r.Context(), `
		UPDATE leads SET
			contact_name = ?, contact_email = ?, contact_phone = ?, country_code = ?, address = ?,
			source = ?, generate_type = ?, assignee_name = ?, assignee_initials = ?,
			lead_pipeline = ?, status = ?, lead_stage = ?, quote_amount = ?, followup_date = ?,
			project_type = ?, quote_type = ?, project_start_date = ?, project_end_date = ?,
			project_name = ?, location = ?, project_segments = ?, package_name = ?,
			selected_services = ?, selected_deliverables = ?, quotation_template = ?,
			discount_amount = ?, tax_amount = ?
		WHERE id = ? AND tenant_id = ?
	`,
		payload.Name, payload.Email, payload.MobileNo, payload.CountryCode, payload.Address,
		payload.LeadSource, payload.Generate, assigneeName, assigneeInit,
		payload.LeadPipeline, payload.PipelineStage, payload.LeadStage, payload.QuoteAmount, payload.FollowupDate,
		payload.ProjectType, payload.QuoteType, payload.ProjectStartDate, payload.ProjectEndDate,
		payload.ProjectName, payload.ProjectAddress, segmentsJSON, payload.PackageName,
		servicesJSON, deliverablesJSON, payload.QuotationTemplate,
		payload.DiscountAmount, payload.TaxAmount,
		leadID, tid,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.ID = leadID
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateLeadStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	leadID := chi.URLParam(r, "id")

	var req struct {
		PipelineStage string `json:"pipelineStage"`
		LeadStage     string `json:"leadStage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE leads
		SET status = ?, lead_stage = ?
		WHERE id = ? AND tenant_id = ?
	`, req.PipelineStage, req.LeadStage, leadID, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": leadID, "pipelineStage": req.PipelineStage, "leadStage": req.LeadStage})
}

func DeleteLeadHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	leadID := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(), "DELETE FROM leads WHERE id = ? AND tenant_id = ?", leadID, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": leadID})
}

func ConvertLeadHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	leadID := chi.URLParam(r, "id")

	var (
		contactName, contactPhone, contactEmail, address, source, projectType, projectName string
		projStart, projEnd                                                                 string
		quoteAmount                                                                        int64
		segmentsRaw                                                                        []byte
	)
	var (
		sPhone, sEmail, sAddr, sSource, sProjType, sProjName, sStart, sEnd sql.NullString
	)

	err := db.QueryRowContext(r.Context(), `
		SELECT contact_name, contact_phone, contact_email, address, source,
		       project_type, project_name, project_start_date, project_end_date,
		       quote_amount, project_segments
		FROM leads
		WHERE id = ? AND tenant_id = ?
	`, leadID, tid).Scan(
		&contactName, &sPhone, &sEmail, &sAddr, &sSource,
		&sProjType, &sProjName, &sStart, &sEnd,
		&quoteAmount, &segmentsRaw,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Lead not found"})
		return
	}

	if sPhone.Valid {
		contactPhone = sPhone.String
	}
	if sEmail.Valid {
		contactEmail = sEmail.String
	}
	if sAddr.Valid {
		address = sAddr.String
	}
	if sSource.Valid {
		source = sSource.String
	}
	if sProjType.Valid {
		projectType = sProjType.String
	} else {
		projectType = "Wedding"
	}
	if sProjName.Valid {
		projectName = sProjName.String
	} else {
		projectName = contactName + " Shoot"
	}
	if sStart.Valid {
		projStart = sStart.String
	} else {
		projStart = time.Now().Format("2006-01-02")
	}
	if sEnd.Valid {
		projEnd = sEnd.String
	} else {
		projEnd = projStart
	}

	clientID := uuid.New().String()
	projectID := uuid.New().String()

	names := strings.SplitN(strings.TrimSpace(contactName), " ", 2)
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = names[1]
	}

	_, _ = db.ExecContext(r.Context(), `
		INSERT INTO clients (id, tenant_id, first_name, last_name, email, phone, source)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, clientID, tid, firstName, lastName, contactEmail, contactPhone, source)

	_, err = db.ExecContext(r.Context(), `
		INSERT INTO projects (
			id, tenant_id, client_id, client_name, client_mobile, client_email, client_address,
			title, project_type, project_address, status, stage, start_date, end_date,
			total_amount_minor, segments
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, 'Active', 'Booked', ?, ?,
			?, ?
		)
	`, projectID, tid, clientID, contactName, contactPhone, contactEmail, address,
		projectName, projectType, address, projStart, projEnd,
		quoteAmount*100, segmentsRaw)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_, _ = db.ExecContext(r.Context(), `
		UPDATE leads
		SET status = 'converted', lead_stage = 'Booked', converted_client_id = ?, converted_project_id = ?
		WHERE id = ? AND tenant_id = ?
	`, clientID, projectID, leadID, tid)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"status":      "converted",
		"clientId":    clientID,
		"projectId":   projectID,
		"projectName": projectName,
	})
}
