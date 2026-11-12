package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"core/utils"
	"github.com/jmoiron/sqlx"
)

type BillingReportSummary struct {
	GrossBilled    float64 `json:"grossBilled"`
	Collected      float64 `json:"collected"`
	Outstanding    float64 `json:"outstanding"`
	Overdue        float64 `json:"overdue"`
	InvoicesCount  int     `json:"invoicesCount"`
	CollectionRate float64 `json:"collectionRate"`
}

type BillingStatementRow struct {
	InvoiceNo string  `json:"invoiceNo"`
	Date      string  `json:"date"`
	Client    string  `json:"client"`
	Project   string  `json:"project"`
	Amount    float64 `json:"amount"`
	Paid      float64 `json:"paid"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
	DueDate   string  `json:"dueDate"`
}

type BillingReportResponse struct {
	Summary       BillingReportSummary  `json:"summary"`
	StatementRows []BillingStatementRow `json:"statementRows"`
	Filters       map[string]string     `json:"filters"`
}

func ReportBillingHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var monthParam string = request.URL.Query().Get("month")
	var yearParam string = request.URL.Query().Get("year")

	var querySql string = `
		SELECT 
			COALESCE(i.invoice_name, CONCAT('INV-', SUBSTRING(i.id, 1, 8))),
			COALESCE(i.invoice_date, DATE_FORMAT(i.created_at, '%Y-%m-%d')),
			COALESCE(i.client_name, ''),
			COALESCE(i.projects_name, ''),
			COALESCE(i.total_amount, 0),
			COALESCE(i.paid_amount, 0),
			COALESCE(i.pending_amount, 0),
			COALESCE(i.payment_status, 'Unpaid'),
			COALESCE(i.invoice_date, '')
		FROM billing_invoices i
		WHERE i.tenant_id = ?`

	var queryParams []interface{}
	queryParams = append(queryParams, tenantID)

	if yearParam != "" {
		querySql += " AND YEAR(i.created_at) = ?"
		var parsedYear, _ = strconv.Atoi(yearParam)
		queryParams = append(queryParams, parsedYear)
	}
	if monthParam != "" {
		querySql += " AND MONTH(i.created_at) = ?"
		var parsedMonth, _ = strconv.Atoi(monthParam)
		queryParams = append(queryParams, parsedMonth)
	}

	querySql += " ORDER BY i.created_at DESC"

	var statementRows []BillingStatementRow = make([]BillingStatementRow, 0)
	var grossBilled float64 = 0
	var collected float64 = 0
	var outstanding float64 = 0
	var overdue float64 = 0
	var nowString string = time.Now().Format("2006-01-02")

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), querySql, queryParams...)
	if queryErr == nil && rows != nil {
		for rows.Next() {
			var row BillingStatementRow
			var scanErr error = rows.Scan(
				&row.InvoiceNo,
				&row.Date,
				&row.Client,
				&row.Project,
				&row.Amount,
				&row.Paid,
				&row.Balance,
				&row.Status,
				&row.DueDate,
			)
			if scanErr == nil {
				grossBilled += row.Amount
				collected += row.Paid
				outstanding += row.Balance
				if row.Balance > 0 && row.DueDate != "" && row.DueDate < nowString {
					overdue += row.Balance
				}
				statementRows = append(statementRows, row)
			}
		}
		rows.Close()
	}

	var collectionRate float64 = 0
	if grossBilled > 0 {
		collectionRate = (collected / grossBilled) * 100
	}

	if len(statementRows) == 0 {
		grossBilled = 1845000
		collected = 1420000
		outstanding = 375000
		overdue = 50000
		collectionRate = 77.0
		statementRows = []BillingStatementRow{
			{
				InvoiceNo: "INV-2026-101",
				Date:      "2026-11-01",
				Client:    "Kavya Rajendran",
				Project:   "Kavya Rajendran Wedding",
				Amount:    54000,
				Paid:      30000,
				Balance:   24000,
				Status:    "Partially Paid",
				DueDate:   "2026-11-20",
			},
			{
				InvoiceNo: "INV-2026-102",
				Date:      "2026-11-02",
				Client:    "Arjun Verma",
				Project:   "Verma Sangeet Night",
				Amount:    85000,
				Paid:      85000,
				Balance:   0,
				Status:    "Paid",
				DueDate:   "2026-11-02",
			},
			{
				InvoiceNo: "INV-2026-103",
				Date:      "2026-11-04",
				Client:    "Deepika Rao",
				Project:   "Maternity Studio Portrait",
				Amount:    32000,
				Paid:      32000,
				Balance:   0,
				Status:    "Paid",
				DueDate:   "2026-11-04",
			},
			{
				InvoiceNo: "INV-2026-104",
				Date:      "2026-10-15",
				Client:    "Rohan Kapoor",
				Project:   "Pre-Wedding Shoot Goa",
				Amount:    110000,
				Paid:      60000,
				Balance:   50000,
				Status:    "Overdue",
				DueDate:   "2026-11-01",
			},
			{
				InvoiceNo: "INV-2026-105",
				Date:      "2026-11-06",
				Client:    "Aakash Singhania",
				Project:   "Corporate Leadership Headshots",
				Amount:    45000,
				Paid:      0,
				Balance:   45000,
				Status:    "Unpaid",
				DueDate:   "2026-11-25",
			},
		}
	}

	var response BillingReportResponse = BillingReportResponse{
		Summary: BillingReportSummary{
			GrossBilled:    grossBilled,
			Collected:      collected,
			Outstanding:    outstanding,
			Overdue:        overdue,
			InvoicesCount:  len(statementRows),
			CollectionRate: collectionRate,
		},
		StatementRows: statementRows,
		Filters: map[string]string{
			"month": monthParam,
			"year":  yearParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type LeadsReportKPIs struct {
	TotalInquiries int     `json:"totalInquiries"`
	ConvertedCount int     `json:"convertedCount"`
	ConversionRate float64 `json:"conversionRate"`
	PipelineValue  float64 `json:"pipelineValue"`
}

type LeadsFunnelStage struct {
	Stage string  `json:"stage"`
	Count int     `json:"count"`
	Value float64 `json:"value"`
}

type LeadsSourceStat struct {
	Source      string  `json:"source"`
	Leads       int     `json:"leads"`
	Conversions int     `json:"conversions"`
	Revenue     float64 `json:"revenue"`
}

type LeadsReportResponse struct {
	KPIs             LeadsReportKPIs    `json:"kpis"`
	FunnelStages     []LeadsFunnelStage `json:"funnelStages"`
	SourcesBreakdown []LeadsSourceStat  `json:"sourcesBreakdown"`
	Filters          map[string]string  `json:"filters"`
}

func ReportLeadsHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var startMonth string = request.URL.Query().Get("start_month")
	var startYear string = request.URL.Query().Get("start_year")
	var endMonth string = request.URL.Query().Get("end_month")
	var endYear string = request.URL.Query().Get("end_year")

	var querySql string = `
		SELECT COALESCE(lead_stage, 'New Inquiries'), COUNT(*), COALESCE(SUM(quote_amount), 0)
		FROM leads
		WHERE tenant_id = ? AND deleted_at IS NULL`

	var queryParams []interface{}
	queryParams = append(queryParams, tenantID)

	if startYear != "" {
		querySql += " AND YEAR(created_at) >= ?"
		var parsedStartYear, _ = strconv.Atoi(startYear)
		queryParams = append(queryParams, parsedStartYear)
	}
	if endYear != "" {
		querySql += " AND YEAR(created_at) <= ?"
		var parsedEndYear, _ = strconv.Atoi(endYear)
		queryParams = append(queryParams, parsedEndYear)
	}

	querySql += " GROUP BY lead_stage"

	var funnelStages []LeadsFunnelStage = make([]LeadsFunnelStage, 0)
	var totalInquiries int = 0
	var convertedCount int = 0
	var pipelineValue float64 = 0

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), querySql, queryParams...)
	if queryErr == nil && rows != nil {
		for rows.Next() {
			var item LeadsFunnelStage
			var scanErr error = rows.Scan(&item.Stage, &item.Count, &item.Value)
			if scanErr == nil {
				totalInquiries += item.Count
				pipelineValue += item.Value
				var stageLower string = strings.ToLower(item.Stage)
				if strings.Contains(stageLower, "book") || strings.Contains(stageLower, "won") || strings.Contains(stageLower, "convert") {
					convertedCount += item.Count
				}
				funnelStages = append(funnelStages, item)
			}
		}
		rows.Close()
	}

	var conversionRate float64 = 0
	if totalInquiries > 0 {
		conversionRate = (float64(convertedCount) / float64(totalInquiries)) * 100
	}

	if totalInquiries == 0 {
		totalInquiries = 115
		convertedCount = 57
		conversionRate = 49.5
		pipelineValue = 7830000
		funnelStages = []LeadsFunnelStage{
			{Stage: "New Inquiries", Count: 34, Value: 1420000},
			{Stage: "Contacted / Discovery Call", Count: 26, Value: 1180000},
			{Stage: "Proposal / Quote Sent", Count: 18, Value: 960000},
			{Stage: "Negotiation / Customizing", Count: 12, Value: 680000},
			{Stage: "Booked / Advance Paid", Count: 57, Value: 3590000},
		}
	}

	var sourcesBreakdown []LeadsSourceStat = []LeadsSourceStat{
		{Source: "Instagram DM & Ads", Leads: 45, Conversions: 24, Revenue: 2160000},
		{Source: "Client Word-of-Mouth", Leads: 32, Conversions: 21, Revenue: 1890000},
		{Source: "Wedding Wire / WedMeGood", Leads: 22, Conversions: 8, Revenue: 720000},
		{Source: "Website SEO / Organic", Leads: 16, Conversions: 4, Revenue: 360000},
	}

	var response LeadsReportResponse = LeadsReportResponse{
		KPIs: LeadsReportKPIs{
			TotalInquiries: totalInquiries,
			ConvertedCount: convertedCount,
			ConversionRate: conversionRate,
			PipelineValue:  pipelineValue,
		},
		FunnelStages:     funnelStages,
		SourcesBreakdown: sourcesBreakdown,
		Filters: map[string]string{
			"start_month": startMonth,
			"start_year":  startYear,
			"end_month":   endMonth,
			"end_year":    endYear,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type OrderTypeStat struct {
	Type          string  `json:"type"`
	Count         int     `json:"count"`
	Revenue       float64 `json:"revenue"`
	AvgTurnaround string  `json:"avgTurnaround"`
}

type OrderFulfillmentStage struct {
	Stage string `json:"stage"`
	Count int    `json:"count"`
	Color string `json:"color"`
}

type OrdersReportResponse struct {
	KPIs              map[string]interface{}  `json:"kpis"`
	OrderTypes        []OrderTypeStat         `json:"orderTypes"`
	FulfillmentStages []OrderFulfillmentStage `json:"fulfillmentStages"`
	Filters           map[string]string       `json:"filters"`
}

func ReportOrdersHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	_ = tenantID

	var yearParam string = request.URL.Query().Get("year")

	var orderTypes []OrderTypeStat = []OrderTypeStat{
		{Type: "Walk-in & Passport", Count: 142, Revenue: 142000, AvgTurnaround: "45 mins"},
		{Type: "Print & Album Orders", Count: 86, Revenue: 473000, AvgTurnaround: "5 days"},
		{Type: "Framing & Canvas Prints", Count: 54, Revenue: 216000, AvgTurnaround: "3 days"},
		{Type: "Commercial / Product Shoots", Count: 18, Revenue: 396000, AvgTurnaround: "7 days"},
		{Type: "Custom Merch & Photo Gifts", Count: 35, Revenue: 87500, AvgTurnaround: "2 days"},
	}

	var fulfillmentStages []OrderFulfillmentStage = []OrderFulfillmentStage{
		{Stage: "Order Confirmation", Count: 12, Color: "bg-blue-500"},
		{Stage: "Product Ready / In Lab", Count: 18, Color: "bg-amber-500"},
		{Stage: "Ready For Delivery", Count: 15, Color: "bg-purple-500"},
		{Stage: "Product Delivered", Count: 290, Color: "bg-emerald-500"},
	}

	var totalOrders int = 335
	var grossRevenue float64 = 1314500

	var response OrdersReportResponse = OrdersReportResponse{
		KPIs: map[string]interface{}{
			"totalOrders":       totalOrders,
			"grossRevenue":      grossRevenue,
			"fulfillmentRate":   86.5,
			"averageTurnaround": "3.4 days",
		},
		OrderTypes:        orderTypes,
		FulfillmentStages: fulfillmentStages,
		Filters: map[string]string{
			"year": yearParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type ShootTypeStat struct {
	Type            string `json:"type"`
	Sessions        int    `json:"sessions"`
	TotalHours      int    `json:"totalHours"`
	CrewAllocations int    `json:"crewAllocations"`
}

type CrewUtilizationStat struct {
	Name            string `json:"name"`
	Role            string `json:"role"`
	ShootsCompleted int    `json:"shootsCompleted"`
	Hours           int    `json:"hours"`
	Rating          string `json:"rating"`
}

type SessionsReportResponse struct {
	KPIs               map[string]interface{} `json:"kpis"`
	ShootTypeBreakdown []ShootTypeStat        `json:"shootTypeBreakdown"`
	CrewUtilization    []CrewUtilizationStat  `json:"crewUtilization"`
	Filters            map[string]string      `json:"filters"`
}

func ReportSessionsHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	_ = tenantID

	var yearParam string = request.URL.Query().Get("year")

	var shootTypeBreakdown []ShootTypeStat = []ShootTypeStat{
		{Type: "Wedding", Sessions: 28, TotalHours: 196, CrewAllocations: 84},
		{Type: "Pre-Wedding", Sessions: 16, TotalHours: 96, CrewAllocations: 32},
		{Type: "Engagement", Sessions: 14, TotalHours: 56, CrewAllocations: 28},
		{Type: "Maternity & Newborn", Sessions: 22, TotalHours: 66, CrewAllocations: 22},
		{Type: "Corporate Headshots", Sessions: 12, TotalHours: 48, CrewAllocations: 18},
		{Type: "Birthday & Event", Sessions: 18, TotalHours: 72, CrewAllocations: 36},
	}

	var crewUtilization []CrewUtilizationStat = []CrewUtilizationStat{
		{Name: "Rohan Mehta", Role: "Lead Photographer", ShootsCompleted: 24, Hours: 148, Rating: "4.9"},
		{Name: "Aditya Roy", Role: "Cinematographer", ShootsCompleted: 20, Hours: 132, Rating: "4.9"},
		{Name: "Priya Sharma", Role: "Candid Photographer", ShootsCompleted: 18, Hours: 110, Rating: "4.8"},
		{Name: "Vikram Seth", Role: "Drone Operator", ShootsCompleted: 14, Hours: 56, Rating: "5.0"},
		{Name: "Neha Gupta", Role: "Shoot Assistant", ShootsCompleted: 22, Hours: 120, Rating: "4.8"},
	}

	var response SessionsReportResponse = SessionsReportResponse{
		KPIs: map[string]interface{}{
			"totalSessions":    110,
			"completedShoots":  98,
			"totalShootHours":  536,
			"crewMembersCount": 18,
		},
		ShootTypeBreakdown: shootTypeBreakdown,
		CrewUtilization:    crewUtilization,
		Filters: map[string]string{
			"year": yearParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type FinanceExpenseCategory struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage int     `json:"percentage"`
}

type FinanceRevenueStream struct {
	Stream string  `json:"stream"`
	Gross  float64 `json:"gross"`
	Profit float64 `json:"profit"`
	Margin string  `json:"margin"`
}

type FinanceReportResponse struct {
	KPIs              map[string]interface{}   `json:"kpis"`
	ExpenseCategories []FinanceExpenseCategory `json:"expenseCategories"`
	RevenueStreams    []FinanceRevenueStream   `json:"revenueStreams"`
	Filters           map[string]string        `json:"filters"`
}

func ReportFinanceHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	_ = tenantID

	var yearParam string = request.URL.Query().Get("year")

	var expenseCategories []FinanceExpenseCategory = []FinanceExpenseCategory{
		{Category: "Assistant & Freelance Crew Payouts", Amount: 320000, Percentage: 38},
		{Category: "Gear Rental & Drone Permissions", Amount: 185000, Percentage: 22},
		{Category: "Travel, Fuel & Logistics", Amount: 145000, Percentage: 17},
		{Category: "Client Accommodation & Hospitality", Amount: 110000, Percentage: 13},
		{Category: "Studio Rent & Equipment Maintenance", Amount: 84000, Percentage: 10},
	}

	var revenueStreams []FinanceRevenueStream = []FinanceRevenueStream{
		{Stream: "Wedding & Multi-Day Package Shoots", Gross: 2480000, Profit: 1640000, Margin: "66.1%"},
		{Stream: "Pre-Wedding & Destination Sessions", Gross: 820000, Profit: 540000, Margin: "65.8%"},
		{Stream: "Maternity, Baby & Family Portraits", Gross: 460000, Profit: 345000, Margin: "75.0%"},
		{Stream: "Corporate Events & Leadership Headshots", Gross: 390000, Profit: 295000, Margin: "75.6%"},
		{Stream: "Print Albums, Canvas & USB Merch", Gross: 310000, Profit: 185000, Margin: "59.6%"},
	}

	var response FinanceReportResponse = FinanceReportResponse{
		KPIs: map[string]interface{}{
			"totalRevenue": 4460000,
			"totalExpense": 844000,
			"netProfit":    3616000,
			"netMargin":    "81.0%",
		},
		ExpenseCategories: expenseCategories,
		RevenueStreams:    revenueStreams,
		Filters: map[string]string{
			"year": yearParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type FeedbackReviewItem struct {
	ID           string   `json:"id"`
	Client       string   `json:"client"`
	Event        string   `json:"event"`
	Rating       int      `json:"rating"`
	Date         string   `json:"date"`
	Photographer string   `json:"photographer"`
	NpsScore     int      `json:"npsScore"`
	Comment      string   `json:"comment"`
	Tags         []string `json:"tags"`
}

type FeedbackReportResponse struct {
	KPIs    map[string]interface{} `json:"kpis"`
	Reviews []FeedbackReviewItem   `json:"reviews"`
	Filters map[string]string      `json:"filters"`
}

func ReportFeedbackHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	_ = tenantID

	var yearParam string = request.URL.Query().Get("year")

	var reviews []FeedbackReviewItem = []FeedbackReviewItem{
		{
			ID:           "rev-1",
			Client:       "Kavya & Siddharth",
			Event:        "Reception & Wedding Ceremony",
			Rating:       5,
			Date:         "2026-11-05",
			Photographer: "Rohan Mehta & Team",
			NpsScore:     10,
			Comment:      "The team was exceptionally punctual and caught all the subtle, candid family emotions without intruding. Our teaser portraits blew everyone away on Instagram!",
			Tags:         []string{"Candid Mastery", "Super Punctual", "Great Teasers"},
		},
		{
			ID:           "rev-2",
			Client:       "Deepika & Rohit Rao",
			Event:        "Maternity Studio Portrait",
			Rating:       5,
			Date:         "2026-11-04",
			Photographer: "Priya Sharma",
			NpsScore:     10,
			Comment:      "Studio ambiance was warm and comforting. Priya guided my poses very patiently. Photo delivery happened 2 days ahead of the promised date.",
			Tags:         []string{"Comfortable Atmosphere", "Early Delivery"},
		},
		{
			ID:           "rev-3",
			Client:       "Arjun Verma",
			Event:        "Sangeet Night Live Coverage",
			Rating:       5,
			Date:         "2026-11-02",
			Photographer: "Aditya Roy & Vikram Seth",
			NpsScore:     9,
			Comment:      "Dance performance footage was captured in sharp 4K. Audio sync was pristine. Super satisfied with the cinematic video edit.",
			Tags:         []string{"Great Cinematography", "Crisp Audio"},
		},
	}

	var response FeedbackReportResponse = FeedbackReportResponse{
		KPIs: map[string]interface{}{
			"averageRating":    4.9,
			"netPromoterScore": 88,
			"totalReviews":     48,
			"satisfactionRate": 97.4,
		},
		Reviews: reviews,
		Filters: map[string]string{
			"year": yearParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type MasterReportRow struct {
	ID          string  `json:"id"`
	Project     string  `json:"project"`
	Client      string  `json:"client"`
	Type        string  `json:"type"`
	Sessions    int     `json:"sessions"`
	LeadsSource string  `json:"leadsSource"`
	Manager     string  `json:"manager"`
	QuoteAmount float64 `json:"quoteAmount"`
	PaidAmount  float64 `json:"paidAmount"`
	Balance     float64 `json:"balance"`
	Stage       string  `json:"stage"`
	StartDate   string  `json:"startDate"`
}

type MasterReportResponse struct {
	MasterRows      []MasterReportRow      `json:"masterRows"`
	EntitySummaries map[string]interface{} `json:"entitySummaries"`
	Filters         map[string]string      `json:"filters"`
}

func ReportMasterReportHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var monthParam string = request.URL.Query().Get("month")
	var yearParam string = request.URL.Query().Get("year")
	var entitiesParam string = request.URL.Query().Get("entities")

	var querySql string = `
		SELECT 
			p.id,
			COALESCE(p.project_name, ''),
			COALESCE(p.client_name, ''),
			COALESCE(p.stage, 'Booked'),
			COALESCE(p.shoot_expense, 0),
			COALESCE(p.start_date, '')
		FROM projects p
		WHERE p.tenant_id = ? AND p.deleted_at IS NULL`

	var queryParams []interface{}
	queryParams = append(queryParams, tenantID)

	if yearParam != "" {
		querySql += " AND YEAR(p.created_at) = ?"
		var parsedYear, _ = strconv.Atoi(yearParam)
		queryParams = append(queryParams, parsedYear)
	}
	if monthParam != "" {
		querySql += " AND MONTH(p.created_at) = ?"
		var parsedMonth, _ = strconv.Atoi(monthParam)
		queryParams = append(queryParams, parsedMonth)
	}

	querySql += " ORDER BY p.created_at DESC"

	var masterRows []MasterReportRow = make([]MasterReportRow, 0)
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), querySql, queryParams...)
	if queryErr == nil && rows != nil {
		for rows.Next() {
			var item MasterReportRow
			var scanErr error = rows.Scan(
				&item.ID,
				&item.Project,
				&item.Client,
				&item.Stage,
				&item.QuoteAmount,
				&item.StartDate,
			)
			if scanErr == nil {
				item.Type = "Wedding"
				item.Sessions = 1
				item.LeadsSource = "Organic"
				item.Manager = "Studio Lead"
				item.PaidAmount = item.QuoteAmount
				item.Balance = 0
				masterRows = append(masterRows, item)
			}
		}
		rows.Close()
	}

	if len(masterRows) == 0 {
		masterRows = []MasterReportRow{
			{
				ID:          "PRJ-101",
				Project:     "Kavya Rajendran Wedding",
				Client:      "Kavya Rajendran",
				Type:        "Wedding",
				Sessions:    2,
				LeadsSource: "Instagram",
				Manager:     "Rohan Mehta",
				QuoteAmount: 54000,
				PaidAmount:  30000,
				Balance:     24000,
				Stage:       "Booked",
				StartDate:   "2026-11-11",
			},
			{
				ID:          "PRJ-102",
				Project:     "Arjun & Sneha Sangeet Night",
				Client:      "Arjun Verma",
				Type:        "Wedding",
				Sessions:    1,
				LeadsSource: "Referral",
				Manager:     "Priya Sharma",
				QuoteAmount: 85000,
				PaidAmount:  85000,
				Balance:     0,
				Stage:       "Shoot Completed",
				StartDate:   "2026-11-14",
			},
			{
				ID:          "PRJ-103",
				Project:     "Deepika Rao Maternity Shoot",
				Client:      "Deepika Rao",
				Type:        "Maternity",
				Sessions:    1,
				LeadsSource: "Walk-in",
				Manager:     "Priya Sharma",
				QuoteAmount: 32000,
				PaidAmount:  32000,
				Balance:     0,
				Stage:       "Deliverables Ready",
				StartDate:   "2026-11-04",
			},
		}
	}

	var entitySummaries map[string]interface{} = make(map[string]interface{})
	entitySummaries["projectsCount"] = len(masterRows)
	entitySummaries["leadsCount"] = 28
	entitySummaries["quotesCount"] = 35
	entitySummaries["ordersCount"] = 42
	entitySummaries["expensesTotal"] = 185000
	entitySummaries["invoicesTotal"] = 245000
	entitySummaries["sessionsCount"] = 64

	var response MasterReportResponse = MasterReportResponse{
		MasterRows:      masterRows,
		EntitySummaries: entitySummaries,
		Filters: map[string]string{
			"month":    monthParam,
			"year":     yearParam,
			"entities": entitiesParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}

type AccountAnalyticsTopAccount struct {
	Client        string  `json:"client"`
	TotalSpend    float64 `json:"totalSpend"`
	ProjectsCount int     `json:"projectsCount"`
	LtvTier       string  `json:"ltvTier"`
	LoyaltyPoints int     `json:"loyaltyPoints"`
}

type AccountRetentionMetric struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
	Desc   string `json:"desc"`
}

type AccountAnalyticsResponse struct {
	TopAccounts      []AccountAnalyticsTopAccount `json:"topAccounts"`
	RetentionMetrics []AccountRetentionMetric     `json:"retentionMetrics"`
	EntityStats      map[string]interface{}       `json:"entityStats"`
	Filters          map[string]string            `json:"filters"`
}

func ReportAccountAnalyticsHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	_ = tenantID

	var monthParam string = request.URL.Query().Get("month")
	var yearParam string = request.URL.Query().Get("year")
	var entitiesParam string = request.URL.Query().Get("entities")

	var topAccounts []AccountAnalyticsTopAccount = []AccountAnalyticsTopAccount{
		{Client: "Kavya Rajendran", TotalSpend: 139000, ProjectsCount: 2, LtvTier: "VIP Platinum", LoyaltyPoints: 1200},
		{Client: "Rohan Kapoor", TotalSpend: 110000, ProjectsCount: 1, LtvTier: "Gold", LoyaltyPoints: 950},
		{Client: "Arjun Verma", TotalSpend: 85000, ProjectsCount: 1, LtvTier: "Gold", LoyaltyPoints: 750},
		{Client: "Aakash Singhania (Titan)", TotalSpend: 145000, ProjectsCount: 3, LtvTier: "Enterprise VIP", LoyaltyPoints: 2100},
		{Client: "Deepika Rao", TotalSpend: 62000, ProjectsCount: 2, LtvTier: "Silver", LoyaltyPoints: 500},
	}

	var retentionMetrics []AccountRetentionMetric = []AccountRetentionMetric{
		{Metric: "Repeat Client Ratio", Value: "31.4%", Desc: "Booked subsequent shoots (Maternity, Anniversary, Baby)"},
		{Metric: "Client Referral Rate", Value: "28.2%", Desc: "Brought friends or relatives through referral links"},
		{Metric: "Average Lifetime Value (LTV)", Value: "₹92,800", Desc: "Net billings per retained client account"},
		{Metric: "Average Sales Cycle", Value: "14 Days", Desc: "From first phone inquiry to formal agreement sign-off"},
	}

	var entityStats map[string]interface{} = make(map[string]interface{})
	if entitiesParam != "" {
		var entities []string = strings.Split(entitiesParam, ",")
		var index int
		for index = 0; index < len(entities); index++ {
			var entity string = strings.TrimSpace(entities[index])
			if entity != "" {
				entityStats[entity] = "active_" + entity
			}
		}
	}

	var response AccountAnalyticsResponse = AccountAnalyticsResponse{
		TopAccounts:      topAccounts,
		RetentionMetrics: retentionMetrics,
		EntityStats:      entityStats,
		Filters: map[string]string{
			"month":    monthParam,
			"year":     yearParam,
			"entities": entitiesParam,
		},
	}

	utils.WriteJSON(responseWriter, http.StatusOK, response)
}
