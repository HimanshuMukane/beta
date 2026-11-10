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

type LeadAssigneeDTO struct {
	Name     string `json:"name"`
	Initials string `json:"initials"`
}

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

type UpdateLeadStageRequest struct {
	PipelineStage string `json:"pipelineStage"`
	LeadStage     string `json:"leadStage"`
}

type LeadStageUpdateResponse struct {
	Status        string `json:"status"`
	ID            string `json:"id"`
	PipelineStage string `json:"pipelineStage"`
	LeadStage     string `json:"leadStage"`
}

type LeadDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

type LeadConversionResponse struct {
	Status      string `json:"status"`
	ProjectID   string `json:"projectId"`
	ProjectName string `json:"projectName"`
}

type ProjectFinanceSummary struct {
	Profit             int64 `json:"profit"`
	CurrentProfit      int64 `json:"currentProfit"`
	NetPayableAmount   int64 `json:"netPayableAmount"`
	PaidAmount         int64 `json:"paidAmount"`
	BalanceAmount      int64 `json:"balanceAmount"`
	PaymentDue         int64 `json:"paymentDue"`
	Expense            int64 `json:"expense"`
	EventDiscountAmount int64 `json:"eventDiscountAmount"`
	PaymentHistory     []any `json:"paymentHistory"`
}

func LeadsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, name, mobile_no, country_code, email, address, lead_source, generate,
		       assignee_name, assignee_initials, lead_pipeline, pipeline_stage, lead_stage,
		       quote_amount, followup_date, project_type, quote_type, project_start_date,
		       project_end_date, project_name, project_address, package_name,
		       project_segments, selected_services, selected_deliverables,
		       quotation_template, discount_amount, tax_amount, created_at
		FROM leads
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list leads", queryErr)
		return
	}
	defer rows.Close()

	var leadsList = make([]LeadDTO, 0)
	for rows.Next() {
		var (
			leadID               string
			leadName             string
			mobileNumber         string
			countryCode          string
			email                string
			address              string
			leadSource           string
			generate             string
			assigneeName         string
			assigneeInitials     string
			leadPipeline         string
			pipelineStage        string
			leadStage            string
			followupDate         string
			projectType          string
			quoteType            string
			projectStartDate     string
			projectEndDate       string
			projectName          string
			projectAddress       string
			packageName          string
			quotationTemplate    string
			quoteAmount          int64
			discountAmount       int64
			taxAmount            int64
			createdAt            time.Time
			rawSegmentsBytes     []byte
			rawServicesBytes     []byte
			rawDeliverablesBytes []byte
		)
		var (
			nullableAddress           sql.NullString
			nullableFollowup          sql.NullString
			nullableStartDate         sql.NullString
			nullableEndDate           sql.NullString
			nullableProjectName       sql.NullString
			nullableProjectAddress    sql.NullString
			nullablePackageName       sql.NullString
			nullableQuotationTemplate sql.NullString
		)

		var scanErr error = rows.Scan(
			&leadID, &leadName, &mobileNumber, &countryCode, &email, &nullableAddress, &leadSource, &generate,
			&assigneeName, &assigneeInitials, &leadPipeline, &pipelineStage, &leadStage,
			&quoteAmount, &nullableFollowup, &projectType, &quoteType, &nullableStartDate,
			&nullableEndDate, &nullableProjectName, &nullableProjectAddress, &nullablePackageName,
			&rawSegmentsBytes, &rawServicesBytes, &rawDeliverablesBytes,
			&nullableQuotationTemplate, &discountAmount, &taxAmount, &createdAt,
		)
		if scanErr != nil {
			continue
		}

		if nullableAddress.Valid {
			address = nullableAddress.String
		}
		if nullableFollowup.Valid {
			followupDate = nullableFollowup.String
		}
		if nullableStartDate.Valid {
			projectStartDate = nullableStartDate.String
		}
		if nullableEndDate.Valid {
			projectEndDate = nullableEndDate.String
		}
		if nullableProjectName.Valid {
			projectName = nullableProjectName.String
		}
		if nullableProjectAddress.Valid {
			projectAddress = nullableProjectAddress.String
		}
		if nullablePackageName.Valid {
			packageName = nullablePackageName.String
		}
		if nullableQuotationTemplate.Valid {
			quotationTemplate = nullableQuotationTemplate.String
		}

		var segmentsData any
		var servicesData any
		var deliverablesData any
		if len(rawSegmentsBytes) > 0 {
			_ = json.Unmarshal(rawSegmentsBytes, &segmentsData)
		}
		if len(rawServicesBytes) > 0 {
			_ = json.Unmarshal(rawServicesBytes, &servicesData)
		}
		if len(rawDeliverablesBytes) > 0 {
			_ = json.Unmarshal(rawDeliverablesBytes, &deliverablesData)
		}

		var leadDTO = LeadDTO{
			ID:                   leadID,
			Name:                 leadName,
			MobileNo:             mobileNumber,
			CountryCode:          countryCode,
			Email:                email,
			Address:              address,
			CreatedDate:          createdAt.Format("02/01/2006"),
			LeadSource:           leadSource,
			Generate:             generate,
			Assignee:             LeadAssigneeDTO{Name: assigneeName, Initials: assigneeInitials},
			LeadPipeline:         leadPipeline,
			PipelineStage:        pipelineStage,
			LeadStage:            leadStage,
			QuoteAmount:          quoteAmount,
			FollowupDate:         followupDate,
			ProjectType:          projectType,
			QuoteType:            quoteType,
			ProjectStartDate:     projectStartDate,
			ProjectEndDate:       projectEndDate,
			ProjectName:          projectName,
			ProjectAddress:       projectAddress,
			PackageName:          packageName,
			ProjectSegments:      segmentsData,
			SelectedServices:     servicesData,
			SelectedDeliverables: deliverablesData,
			QuotationTemplate:    quotationTemplate,
			DiscountAmount:       discountAmount,
			TaxAmount:            taxAmount,
		}
		leadsList = append(leadsList, leadDTO)
	}

	utils.WriteJSON(responseWriter, http.StatusOK, leadsList)
}

func GetLeadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var leadID string = chi.URLParam(request, "id")

	var (
		leadName             string
		mobileNumber         string
		countryCode          string
		email                string
		address              string
		leadSource           string
		generate             string
		assigneeName         string
		assigneeInitials     string
		leadPipeline         string
		pipelineStage        string
		leadStage            string
		followupDate         string
		projectType          string
		quoteType            string
		projectStartDate     string
		projectEndDate       string
		projectName          string
		projectAddress       string
		packageName          string
		quotationTemplate    string
		quoteAmount          int64
		discountAmount       int64
		taxAmount            int64
		createdAt            time.Time
		rawSegmentsBytes     []byte
		rawServicesBytes     []byte
		rawDeliverablesBytes []byte
	)
	var (
		nullableAddress           sql.NullString
		nullableFollowup          sql.NullString
		nullableStartDate         sql.NullString
		nullableEndDate           sql.NullString
		nullableProjectName       sql.NullString
		nullableProjectAddress    sql.NullString
		nullablePackageName       sql.NullString
		nullableQuotationTemplate sql.NullString
	)

	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT id, name, mobile_no, country_code, email, address, lead_source, generate,
		       assignee_name, assignee_initials, lead_pipeline, pipeline_stage, lead_stage,
		       quote_amount, followup_date, project_type, quote_type, project_start_date,
		       project_end_date, project_name, project_address, package_name,
		       project_segments, selected_services, selected_deliverables,
		       quotation_template, discount_amount, tax_amount, created_at
		FROM leads
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, leadID, tenantID).Scan(
		&leadID, &leadName, &mobileNumber, &countryCode, &email, &nullableAddress, &leadSource, &generate,
		&assigneeName, &assigneeInitials, &leadPipeline, &pipelineStage, &leadStage,
		&quoteAmount, &nullableFollowup, &projectType, &quoteType, &nullableStartDate,
		&nullableEndDate, &nullableProjectName, &nullableProjectAddress, &nullablePackageName,
		&rawSegmentsBytes, &rawServicesBytes, &rawDeliverablesBytes,
		&nullableQuotationTemplate, &discountAmount, &taxAmount, &createdAt,
	)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Lead not found")
		return
	}

	address = ""
	if nullableAddress.Valid {
		address = nullableAddress.String
	}
	followupDate = ""
	if nullableFollowup.Valid {
		followupDate = nullableFollowup.String
	}
	projectStartDate = ""
	if nullableStartDate.Valid {
		projectStartDate = nullableStartDate.String
	}
	projectEndDate = ""
	if nullableEndDate.Valid {
		projectEndDate = nullableEndDate.String
	}
	projectName = ""
	if nullableProjectName.Valid {
		projectName = nullableProjectName.String
	}
	projectAddress = ""
	if nullableProjectAddress.Valid {
		projectAddress = nullableProjectAddress.String
	}
	packageName = ""
	if nullablePackageName.Valid {
		packageName = nullablePackageName.String
	}
	quotationTemplate = ""
	if nullableQuotationTemplate.Valid {
		quotationTemplate = nullableQuotationTemplate.String
	}

	var segmentsData any
	var servicesData any
	var deliverablesData any
	if len(rawSegmentsBytes) > 0 {
		_ = json.Unmarshal(rawSegmentsBytes, &segmentsData)
	}
	if len(rawServicesBytes) > 0 {
		_ = json.Unmarshal(rawServicesBytes, &servicesData)
	}
	if len(rawDeliverablesBytes) > 0 {
		_ = json.Unmarshal(rawDeliverablesBytes, &deliverablesData)
	}

	var leadDTO = LeadDTO{
		ID:                   leadID,
		Name:                 leadName,
		MobileNo:             mobileNumber,
		CountryCode:          countryCode,
		Email:                email,
		Address:              address,
		CreatedDate:          createdAt.Format("02/01/2006"),
		LeadSource:           leadSource,
		Generate:             generate,
		Assignee:             LeadAssigneeDTO{Name: assigneeName, Initials: assigneeInitials},
		LeadPipeline:         leadPipeline,
		PipelineStage:        pipelineStage,
		LeadStage:            leadStage,
		QuoteAmount:          quoteAmount,
		FollowupDate:         followupDate,
		ProjectType:          projectType,
		QuoteType:            quoteType,
		ProjectStartDate:     projectStartDate,
		ProjectEndDate:       projectEndDate,
		ProjectName:          projectName,
		ProjectAddress:       projectAddress,
		PackageName:          packageName,
		ProjectSegments:      segmentsData,
		SelectedServices:     servicesData,
		SelectedDeliverables: deliverablesData,
		QuotationTemplate:    quotationTemplate,
		DiscountAmount:       discountAmount,
		TaxAmount:            taxAmount,
	}

	utils.WriteJSON(responseWriter, http.StatusOK, leadDTO)
}

func CreateLeadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var leadPayload LeadDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&leadPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	if leadPayload.ID == "" {
		leadPayload.ID = uuid.New().String()
	}

	var assigneeName string = "Rohan Mehta"
	var assigneeInitials string = "RM"
	if assigneeMap, ok := leadPayload.Assignee.(map[string]any); ok {
		if nameVal, ok := assigneeMap["name"].(string); ok && nameVal != "" {
			assigneeName = nameVal
		}
		if initVal, ok := assigneeMap["initials"].(string); ok && initVal != "" {
			assigneeInitials = initVal
		}
	}

	segmentsBytes, _ := json.Marshal(leadPayload.ProjectSegments)
	servicesBytes, _ := json.Marshal(leadPayload.SelectedServices)
	deliverablesBytes, _ := json.Marshal(leadPayload.SelectedDeliverables)

	_, insertErr := database.ExecContext(request.Context(), `
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
		leadPayload.ID, tenantID, leadPayload.Name, leadPayload.MobileNo, leadPayload.CountryCode, leadPayload.Email, leadPayload.Address, leadPayload.LeadSource,
		leadPayload.Generate, assigneeName, assigneeInitials, leadPayload.LeadPipeline, leadPayload.PipelineStage,
		leadPayload.LeadStage, leadPayload.QuoteAmount, leadPayload.FollowupDate, leadPayload.ProjectType, leadPayload.QuoteType,
		leadPayload.ProjectStartDate, leadPayload.ProjectEndDate, leadPayload.ProjectName, leadPayload.ProjectAddress,
		leadPayload.PackageName, segmentsBytes, servicesBytes, deliverablesBytes,
		leadPayload.QuotationTemplate, leadPayload.DiscountAmount, leadPayload.TaxAmount,
	)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create lead", insertErr)
		return
	}

	leadPayload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(responseWriter, http.StatusCreated, leadPayload)
}

func UpdateLeadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var leadID string = chi.URLParam(request, "id")

	var leadPayload LeadDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&leadPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var assigneeName string = "Rohan Mehta"
	var assigneeInitials string = "RM"
	if assigneeMap, ok := leadPayload.Assignee.(map[string]any); ok {
		if nameVal, ok := assigneeMap["name"].(string); ok && nameVal != "" {
			assigneeName = nameVal
		}
		if initVal, ok := assigneeMap["initials"].(string); ok && initVal != "" {
			assigneeInitials = initVal
		}
	}

	segmentsBytes, _ := json.Marshal(leadPayload.ProjectSegments)
	servicesBytes, _ := json.Marshal(leadPayload.SelectedServices)
	deliverablesBytes, _ := json.Marshal(leadPayload.SelectedDeliverables)

	_, updateErr := database.ExecContext(request.Context(), `
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
		leadPayload.Name, leadPayload.MobileNo, leadPayload.CountryCode, leadPayload.Email, leadPayload.Address, leadPayload.LeadSource,
		leadPayload.Generate, assigneeName, assigneeInitials, leadPayload.LeadPipeline,
		leadPayload.PipelineStage, leadPayload.LeadStage, leadPayload.QuoteAmount, leadPayload.FollowupDate,
		leadPayload.ProjectType, leadPayload.QuoteType, leadPayload.ProjectStartDate, leadPayload.ProjectEndDate,
		leadPayload.ProjectName, leadPayload.ProjectAddress, leadPayload.PackageName, segmentsBytes,
		servicesBytes, deliverablesBytes, leadPayload.QuotationTemplate,
		leadPayload.DiscountAmount, leadPayload.TaxAmount,
		leadID, tenantID,
	)

	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update lead", updateErr)
		return
	}

	leadPayload.ID = leadID
	utils.WriteJSON(responseWriter, http.StatusOK, leadPayload)
}

func UpdateLeadStageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var leadID string = chi.URLParam(request, "id")

	var requestPayload UpdateLeadStageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var query string = "UPDATE leads SET pipeline_stage = ?"
	var queryArguments = []any{requestPayload.PipelineStage}
	if requestPayload.LeadStage != "" {
		query += ", lead_stage = ?"
		queryArguments = append(queryArguments, requestPayload.LeadStage)
	}
	query += " WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL"
	queryArguments = append(queryArguments, leadID, tenantID)

	_, updateErr := database.ExecContext(request.Context(), query, queryArguments...)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update lead stage", updateErr)
		return
	}

	var responsePayload = LeadStageUpdateResponse{
		Status:        "ok",
		ID:            leadID,
		PipelineStage: requestPayload.PipelineStage,
		LeadStage:     requestPayload.LeadStage,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

func DeleteLeadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var leadID string = chi.URLParam(request, "id")

	_, deleteErr := database.ExecContext(request.Context(), `
		UPDATE leads SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, leadID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete lead", deleteErr)
		return
	}

	var responsePayload = LeadDeleteResponse{
		Status: "deleted",
		ID:     leadID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

func ConvertLeadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var leadID string = chi.URLParam(request, "id")

	var (
		clientName        string
		mobileNumber      string
		countryCode       string
		email             string
		address           string
		leadSource        string
		projectName       string
		projectType       string
		nullableAddress   sql.NullString
		nullableStartDate sql.NullString
		nullableEndDate   sql.NullString
		rawSegmentsBytes  []byte
	)

	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT name, mobile_no, country_code, email, address, lead_source, project_name,
		       project_type, project_address, project_start_date, project_end_date, project_segments
		FROM leads
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, leadID, tenantID).Scan(
		&clientName, &mobileNumber, &countryCode, &email, &address, &leadSource, &projectName,
		&projectType, &nullableAddress, &nullableStartDate, &nullableEndDate, &rawSegmentsBytes,
	)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Lead not found")
		return
	}

	var projectID string = "PRJ-" + uuid.New().String()[:8]
	if strings.TrimSpace(projectName) == "" {
		projectName = clientName + " Shoot"
	}

	var clientInitials string = "PR"
	var nameParts = strings.Fields(clientName)
	if len(nameParts) > 1 {
		clientInitials = strings.ToUpper(string([]rune(nameParts[0])[:1]) + string([]rune(nameParts[1])[:1]))
	} else if len(nameParts) == 1 && len(nameParts[0]) > 0 {
		clientInitials = strings.ToUpper(string([]rune(nameParts[0])[:1]))
	}

	transaction, txErr := database.BeginTx(request.Context(), nil)
	if txErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to begin transaction", txErr)
		return
	}
	defer transaction.Rollback()

	var clientID string
	var clientQueryErr error = transaction.QueryRowContext(request.Context(), `
		SELECT id FROM clients WHERE tenant_id = ? AND email = ? LIMIT 1
	`, tenantID, email).Scan(&clientID)
	if clientQueryErr != nil || clientID == "" {
		clientID = "CLI-" + uuid.New().String()[:8]
		var firstName string = clientName
		var lastName string = ""
		if len(nameParts) > 1 {
			firstName = nameParts[0]
			lastName = strings.Join(nameParts[1:], " ")
		}
		_, insertClientErr := transaction.ExecContext(request.Context(), `
			INSERT INTO clients (id, tenant_id, first_name, last_name, email, phone, source)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, clientID, tenantID, firstName, lastName, email, mobileNumber, leadSource)
		if insertClientErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create client record", insertClientErr)
			return
		}
	}

	var financeModel = ProjectFinanceSummary{
		Profit:              0,
		CurrentProfit:       0,
		NetPayableAmount:    0,
		PaidAmount:          0,
		BalanceAmount:       0,
		PaymentDue:          0,
		Expense:             0,
		EventDiscountAmount: 0,
		PaymentHistory:      []any{},
	}
	financeBytes, _ := json.Marshal(financeModel)

	var projectAddress string = ""
	if nullableAddress.Valid {
		projectAddress = nullableAddress.String
	}
	var projectStartDate string = ""
	if nullableStartDate.Valid {
		projectStartDate = nullableStartDate.String
	}
	var projectEndDate string = ""
	if nullableEndDate.Valid {
		projectEndDate = nullableEndDate.String
	}

	_, insertProjectErr := transaction.ExecContext(request.Context(), `
		INSERT INTO projects (
			id, tenant_id, client_id, project_name, client_name, client_initials, client_mobile,
			client_country_code, client_email, client_address, lead_source, project_manager,
			project_type, project_address, status, stage, start_date, end_date,
			shoot_hours, shoot_members_count, shoot_expense, galleries_count, notes,
			segments, documents, sessions, tasks, printing, deliverables, expenses, finance
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, NULL,
			?, ?, 'New', 'Booked', ?, ?,
			8.00, 2, 0, 0, '',
			?, '[]', '[]', '[]', '[]', '[]', '[]', ?
		)
	`,
		projectID, tenantID, clientID, projectName, clientName, clientInitials, mobileNumber,
		countryCode, email, address, leadSource,
		projectType, projectAddress, projectStartDate, projectEndDate,
		rawSegmentsBytes, financeBytes,
	)
	if insertProjectErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create project", insertProjectErr)
		return
	}

	_, updateLeadErr := transaction.ExecContext(request.Context(), `
		UPDATE leads SET pipeline_stage = 'converted', lead_stage = 'Booked', status = 'converted',
		converted_client_id = ?, converted_project_id = ?
		WHERE id = ? AND tenant_id = ?
	`, clientID, projectID, leadID, tenantID)
	if updateLeadErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update lead", updateLeadErr)
		return
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit conversion", commitErr)
		return
	}

	var responsePayload = LeadConversionResponse{
		Status:      "ok",
		ProjectID:   projectID,
		ProjectName: projectName,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}
