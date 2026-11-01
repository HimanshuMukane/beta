package handlers

import (
	"net/http"

	billingPages "chi-mysql-server/ui/pages/billing"
	calendarPages "chi-mysql-server/ui/pages/calendar"
	clientsPages "chi-mysql-server/ui/pages/clients"
	contractsPages "chi-mysql-server/ui/pages/contracts"
	deliverablesPages "chi-mysql-server/ui/pages/deliverables"
	followupsPages "chi-mysql-server/ui/pages/followups"
	inboxPages "chi-mysql-server/ui/pages/inbox"
	leadsPages "chi-mysql-server/ui/pages/leads"
	ordersPages "chi-mysql-server/ui/pages/orders"
	projectsPages "chi-mysql-server/ui/pages/projects"
	reportsPages "chi-mysql-server/ui/pages/reports"
	sessionsPages "chi-mysql-server/ui/pages/sessions"
	settingsPages "chi-mysql-server/ui/pages/settings"
)

func LeadsHandler(w http.ResponseWriter, r *http.Request) {
	stage := r.URL.Query().Get("stage")
	source := r.URL.Query().Get("source")
	_ = leadsPages.Leads(stage, source).Render(r.Context(), w)
}

func LeadDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = leadsPages.LeadDetails(id).Render(r.Context(), w)
}

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	_ = inboxPages.Inbox(filter).Render(r.Context(), w)
}

func FollowupsHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	_ = followupsPages.Followups(filter).Render(r.Context(), w)
}

func ClientsHandler(w http.ResponseWriter, r *http.Request) {
	_ = clientsPages.Clients().Render(r.Context(), w)
}

func ClientDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = clientsPages.ClientDetails(id).Render(r.Context(), w)
}

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	_ = ordersPages.Orders(status).Render(r.Context(), w)
}

func OrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = ordersPages.OrderDetails(id).Render(r.Context(), w)
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	_ = projectsPages.Projects(filter).Render(r.Context(), w)
}

func ProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = projectsPages.ProjectDetails(id).Render(r.Context(), w)
}

func SessionsHandler(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	_ = sessionsPages.Sessions(tab).Render(r.Context(), w)
}

func CalendarHandler(w http.ResponseWriter, r *http.Request) {
	view := r.URL.Query().Get("view")
	if view == "" {
		view = "month"
	}

	filters := []calendarPages.CalendarFilterOption{
		{ID: "meet", Label: "Google Meet & Calls", Color: "#7c5dfa", Count: 3, Checked: true},
		{ID: "session", Label: "Photo Sessions", Color: "#00cec9", Count: 2, Checked: true},
		{ID: "busy", Label: "Busy / Unavailable", Color: "#ff7675", Count: 2, Checked: true},
	}

	upcoming := []calendarPages.UpcomingMeetingItem{
		{Title: "Discovery Consultation", Time: "Today, 10:00 AM", MeetType: "meet", Client: "Emma Vance", IsBusy: false},
		{Title: "Thornton Nuptials Shoot", Time: "In 2 Days, 2:00 PM", MeetType: "session", Client: "Thornton Family", IsBusy: false},
		{Title: "Blocked / Out of Studio", Time: "Tomorrow, 1:00 PM", MeetType: "busy", Client: "Internal", IsBusy: true},
	}

	_ = calendarPages.CalendarLayout(view, filters, upcoming).Render(r.Context(), w)
}

func ContractsHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("status")
	_ = contractsPages.Contracts(filter).Render(r.Context(), w)
}

func DeliverablesHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("status")
	_ = deliverablesPages.Deliverables(filter).Render(r.Context(), w)
}

func BillingOverviewHandler(w http.ResponseWriter, r *http.Request) {
	_ = billingPages.BillingOverview().Render(r.Context(), w)
}

func PaymentsHandler(w http.ResponseWriter, r *http.Request) {
	_ = billingPages.Payments().Render(r.Context(), w)
}

func InvoicesHandler(w http.ResponseWriter, r *http.Request) {
	_ = billingPages.Invoices().Render(r.Context(), w)
}

func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	_ = billingPages.Expenses().Render(r.Context(), w)
}

func RefundsHandler(w http.ResponseWriter, r *http.Request) {
	_ = billingPages.Refunds().Render(r.Context(), w)
}

func CancelledHandler(w http.ResponseWriter, r *http.Request) {
	_ = billingPages.CancelledTransactions().Render(r.Context(), w)
}

func ReportsHandler(w http.ResponseWriter, r *http.Request) {
	_ = reportsPages.Reports().Render(r.Context(), w)
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	_ = settingsPages.Settings(tab).Render(r.Context(), w)
}

func CrewSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.CrewSettings().Render(r.Context(), w)
}

func SequenceIDsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.SequenceIDsSettings().Render(r.Context(), w)
}

func ServicesSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.ServicesSettings().Render(r.Context(), w)
}

func PackagesSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.PackagesSettings().Render(r.Context(), w)
}

func CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.PackageEditor(false, "").Render(r.Context(), w)
}

func EditPackageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = settingsPages.PackageEditor(true, id).Render(r.Context(), w)
}

func OrderPresetsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.OrderPresetsSettings().Render(r.Context(), w)
}

func DeliverablePresetsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.DeliverablePresetsSettings().Render(r.Context(), w)
}

func FormsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.FormsSettings().Render(r.Context(), w)
}

func QuotationTemplatesSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.QuotationTemplatesSettings().Render(r.Context(), w)
}

func InvoiceTemplatesSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.InvoiceTemplatesSettings().Render(r.Context(), w)
}

func ContractTemplatesSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.ContractTemplatesSettings().Render(r.Context(), w)
}

func CreateContractTemplateHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.ContractTemplateEditor(false, "").Render(r.Context(), w)
}

func EditContractTemplateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = settingsPages.ContractTemplateEditor(true, id).Render(r.Context(), w)
}

func FeedbackTemplatesSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.FeedbackTemplatesSettings().Render(r.Context(), w)
}

func CreateFeedbackTemplateHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.FeedbackTemplateEditor(false, "").Render(r.Context(), w)
}

func EditFeedbackTemplateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_ = settingsPages.FeedbackTemplateEditor(true, id).Render(r.Context(), w)
}

func WorkflowSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.WorkflowSettings().Render(r.Context(), w)
}

func IntegrationsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.IntegrationsSettings().Render(r.Context(), w)
}

func BusinessSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_ = settingsPages.BusinessSettings().Render(r.Context(), w)
}
