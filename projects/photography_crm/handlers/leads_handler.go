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

type LeadDTO struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	MobileNo             string `json:"mobileNo"`
	CountryCode          string `json:"countryCode"`
	Email                string `json:"email"`
	Address              string `json:"address"`
	CreatedDate          string `json:"createdDate"`
	LeadSource           string `json:"leadSource"`
	Generate             string `json:"generate"`
	Assignee             any    `json:"assignee"`
	LeadPipeline         string `json:"leadPipeline"`
	PipelineStage        string `json:"pipelineStage"`
	LeadStage            string `json:"leadStage"`
	QuoteAmount          int64  `json:"quoteAmount"`
	FollowupDate         string `json:"followupDate"`
	ProjectType          string `json:"projectType"`
	QuoteType            string `json:"quoteType"`
	ProjectStartDate     string `json:"projectStartDate"`
	ProjectEndDate       string `json:"projectEndDate"`
	ProjectName          string `json:"projectName"`
	ProjectAddress       string `json:"projectAddress"`
	PackageName          string `json:"packageName"`
	ProjectSegments      any    `json:"projectSegments"`
	SelectedServices     any    `json:"selectedServices"`
	SelectedDeliverables any    `json:"selectedDeliverables"`
	QuotationTemplate    string `json:"quotationTemplate"`
	DiscountAmount       int64  `json:"discountAmount"`
	TaxAmount            int64  `json:"taxAmount"`
}

func LeadsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, name, mobile_no, country_code, email, address, lead_source, generate,
		       assignee_name, assignee_initials, lead_pipeline, pipeline_stage, lead_stage,
		       quote_amount, followup_date, project_type, quote_type, project_start_date,
		       project_end_date, project_name, project_address, package_name,
		       project_segments, selected_services, selected_deliverables,
		       quotation_template, discount_amount, tax_amount, created_at
		FROM leads
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	result := make([]LeadDTO, 0)
	for rows.Next() {
		var (
			id, name, mobileNo, countryCode, email, address, leadSource, generate string
			assigneeName, assigneeInit, leadPipeline, pipelineStage, leadStage    string
			followupDate, projectType, quoteType, projStartDate, projEndDate      string
			projName, projAddress, pkgName, quotTemplate                          string
			quoteAmount, discountAmount, taxAmount                                int64
			createdAt                                                             time.Time
			segmentsRaw, servicesRaw, deliverablesRaw                             []byte
		)
		var (
			sAddr, sFollowup, sStartDate, sEndDate, sProjName, sProjAddr, sPkgName, sTemplate sql.NullString
		)

		err := rows.Scan(
			&id, &name, &mobileNo, &countryCode, &email, &sAddr, &leadSource, &generate,
			&assigneeName, &assigneeInit, &leadPipeline, &pipelineStage, &leadStage,
			&quoteAmount, &sFollowup, &projectType, &quoteType, &sStartDate,
			&sEndDate, &sProjName, &sProjAddr, &sPkgName,
			&segmentsRaw, &servicesRaw, &deliverablesRaw,
			&sTemplate, &discountAmount, &taxAmount, &createdAt,
		)
		if err != nil {
			continue
		}

		if sAddr.Valid {
			address = sAddr.String
		}
		if sFollowup.Valid {
			followupDate = sFollowup.String
		}
		if sStartDate.Valid {
			projStartDate = sStartDate.String
		}
		if sEndDate.Valid {
			projEndDate = sEndDate.String
		}
		if sProjName.Valid {
			projName = sProjName.String
		}
		if sProjAddr.Valid {
			projAddress = sProjAddr.String
		}
		if sPkgName.Valid {
			pkgName = sPkgName.String
		}
		if sTemplate.Valid {
			quotTemplate = sTemplate.String
		}

		var segments, services, deliverables any
		if len(segmentsRaw) > 0 {
			_ = json.Unmarshal(segmentsRaw, &segments)
		}
		if len(servicesRaw) > 0 {
			_ = json.Unmarshal(servicesRaw, &services)
		}
		if len(deliverablesRaw) > 0 {
			_ = json.Unmarshal(deliverablesRaw, &deliverables)
		}

		lead := LeadDTO{
			ID:                   id,
			Name:                 name,
			MobileNo:             mobileNo,
			CountryCode:          countryCode,
			Email:                email,
			Address:              address,
			CreatedDate:          createdAt.Format("02/01/2006"),
			LeadSource:           leadSource,
			Generate:             generate,
			Assignee:             map[string]string{"name": assigneeName, "initials": assigneeInit},
			LeadPipeline:         leadPipeline,
			PipelineStage:        pipelineStage,
			LeadStage:            leadStage,
			QuoteAmount:          quoteAmount,
			FollowupDate:         followupDate,
			ProjectType:          projectType,
			QuoteType:            quoteType,
			ProjectStartDate:     projStartDate,
			ProjectEndDate:       projEndDate,
			ProjectName:          projName,
			ProjectAddress:       projAddress,
			PackageName:          pkgName,
			ProjectSegments:      segments,
			SelectedServices:     services,
			SelectedDeliverables: deliverables,
			QuotationTemplate:    quotTemplate,
			DiscountAmount:       discountAmount,
			TaxAmount:            taxAmount,
		}
		result = append(result, lead)
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
		payload.ID = uuid.New().String()
	}

	assigneeName := "Rohan Mehta"
	assigneeInit := "RM"
	if m, ok := payload.Assignee.(map[string]any); ok {
		if n, ok := m["name"].(string); ok && n != "" {
			assigneeName = n
		}
		if in, ok := m["initials"].(string); ok && in != "" {
			assigneeInit = in
		}
	}

	segmentsJSON, _ := json.Marshal(payload.ProjectSegments)
	servicesJSON, _ := json.Marshal(payload.SelectedServices)
	deliverablesJSON, _ := json.Marshal(payload.SelectedDeliverables)

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO leads (
			id, tenant_id, name, mobile_no, country_code, email, address, lead_source,
			generate, assignee_name, assignee_initials, lead_pipeline, pipeline_stage,
			lead_stage, quote_amount, followup_date, project_type, quote_type,
			project_start_date, project_end_date, project_name, project_address,
			package_name, project_segments, selected_services, selected_deliverables,
			quotation_template, discount_amount, tax_amount
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
	`,
		payload.ID, tid, payload.Name, payload.MobileNo, payload.CountryCode, payload.Email, payload.Address, payload.LeadSource,
		payload.Generate, assigneeName, assigneeInit, payload.LeadPipeline, payload.PipelineStage,
		payload.LeadStage, payload.QuoteAmount, payload.FollowupDate, payload.ProjectType, payload.QuoteType,
		payload.ProjectStartDate, payload.ProjectEndDate, payload.ProjectName, payload.ProjectAddress,
		payload.PackageName, segmentsJSON, servicesJSON, deliverablesJSON,
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
	id := chi.URLParam(r, "id")

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
		if in, ok := m["initials"].(string); ok && in != "" {
			assigneeInit = in
		}
	}

	segmentsJSON, _ := json.Marshal(payload.ProjectSegments)
	servicesJSON, _ := json.Marshal(payload.SelectedServices)
	deliverablesJSON, _ := json.Marshal(payload.SelectedDeliverables)

	_, err := db.ExecContext(r.Context(), `
		UPDATE leads SET
			name = ?, mobile_no = ?, country_code = ?, email = ?, address = ?, lead_source = ?,
			generate = ?, assignee_name = ?, assignee_initials = ?, lead_pipeline = ?,
			pipeline_stage = ?, lead_stage = ?, quote_amount = ?, followup_date = ?,
			project_type = ?, quote_type = ?, project_start_date = ?, project_end_date = ?,
			project_name = ?, project_address = ?, package_name = ?, project_segments = ?,
			selected_services = ?, selected_deliverables = ?, quotation_template = ?,
			discount_amount = ?, tax_amount = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		payload.Name, payload.MobileNo, payload.CountryCode, payload.Email, payload.Address, payload.LeadSource,
		payload.Generate, assigneeName, assigneeInit, payload.LeadPipeline,
		payload.PipelineStage, payload.LeadStage, payload.QuoteAmount, payload.FollowupDate,
		payload.ProjectType, payload.QuoteType, payload.ProjectStartDate, payload.ProjectEndDate,
		payload.ProjectName, payload.ProjectAddress, payload.PackageName, segmentsJSON,
		servicesJSON, deliverablesJSON, payload.QuotationTemplate,
		payload.DiscountAmount, payload.TaxAmount,
		id, tid,
	)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.ID = id
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateLeadStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var req struct {
		PipelineStage string `json:"pipelineStage"`
		LeadStage     string `json:"leadStage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	query := "UPDATE leads SET pipeline_stage = ?"
	args := []any{req.PipelineStage}
	if req.LeadStage != "" {
		query += ", lead_stage = ?"
		args = append(args, req.LeadStage)
	}
	query += " WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL"
	args = append(args, id, tid)

	_, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"status":        "ok",
		"id":            id,
		"pipelineStage": req.PipelineStage,
		"leadStage":     req.LeadStage,
	})
}

func DeleteLeadHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(), `
		UPDATE leads SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

func ConvertLeadHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var (
		name, mobileNo, countryCode, email, address, leadSource, projName, projType string
		projAddress, startDate, endDate                                             sql.NullString
		segmentsRaw                                                                 []byte
	)

	err := db.QueryRowContext(r.Context(), `
		SELECT name, mobile_no, country_code, email, address, lead_source, project_name,
		       project_type, project_address, project_start_date, project_end_date, project_segments
		FROM leads
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, id, tid).Scan(
		&name, &mobileNo, &countryCode, &email, &address, &leadSource, &projName,
		&projType, &projAddress, &startDate, &endDate, &segmentsRaw,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Lead not found"})
		return
	}

	projID := "PRJ-" + uuid.New().String()[:8]
	if strings.TrimSpace(projName) == "" {
		projName = name + " Shoot"
	}

	initials := "PR"
	parts := strings.Fields(name)
	if len(parts) > 1 {
		initials = strings.ToUpper(string([]rune(parts[0])[:1]) + string([]rune(parts[1])[:1]))
	} else if len(parts) == 1 && len(parts[0]) > 0 {
		initials = strings.ToUpper(string([]rune(parts[0])[:1]))
	}

	mgrJSON, _ := json.Marshal(map[string]string{"name": "Rohan Mehta", "initials": "RM"})
	financeJSON, _ := json.Marshal(map[string]any{
		"profit": 0, "currentProfit": 0, "netPayableAmount": 0,
		"paidAmount": 0, "balanceAmount": 0, "paymentDue": 0,
		"expense": 0, "eventDiscountAmount": 0, "paymentHistory": []any{},
	})

	var sAddress, sStartDate, sEndDate string
	if projAddress.Valid {
		sAddress = projAddress.String
	}
	if startDate.Valid {
		sStartDate = startDate.String
	}
	if endDate.Valid {
		sEndDate = endDate.String
	}

	_, err = db.ExecContext(r.Context(), `
		INSERT INTO projects (
			id, tenant_id, project_name, client_name, client_initials, client_mobile,
			client_country_code, client_email, client_address, lead_source, project_manager,
			project_type, project_address, status, stage, start_date, end_date,
			shoot_hours, shoot_members_count, shoot_expense, galleries_count, notes,
			segments, documents, sessions, tasks, printing, deliverables, expenses, finance
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, 'New', 'Booked', ?, ?,
			8.00, 2, 0, 0, '',
			?, '[]', '[]', '[]', '[]', '[]', '[]', ?
		)
	`,
		projID, tid, projName, name, initials, mobileNo,
		countryCode, email, address, leadSource, mgrJSON,
		projType, sAddress, sStartDate, sEndDate,
		segmentsRaw, financeJSON,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_, _ = db.ExecContext(r.Context(), `
		UPDATE leads SET pipeline_stage = 'converted', lead_stage = 'Booked'
		WHERE id = ? AND tenant_id = ?
	`, id, tid)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"status":      "ok",
		"projectId":   projID,
		"projectName": projName,
	})
}

