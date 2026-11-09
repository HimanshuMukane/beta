package handlers

import (
	"net/http"

	"core/utils"
)

func LeadsHandler(w http.ResponseWriter, r *http.Request) {
	stage := r.URL.Query().Get("stage")
	source := r.URL.Query().Get("source")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "stage": stage, "source": source})
}

func LeadDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "id": id})
}

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "filter": filter})
}

func FollowupsHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "filter": filter})
}

func ClientsHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "clients": []any{}})
}

func ClientDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "id": id})
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "project_status": status})
}

func ProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "id": id})
}

func SessionsHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "date": date})
}

func SessionDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "id": id})
}

func OrdersHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "orders": []any{}})
}

func DeliverablesHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "deliverables": []any{}})
}

func ContractsHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "contracts": []any{}})
}

func InvoicesHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "invoices": []any{}})
}

func PaymentsHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "payments": []any{}})
}

func CalendarHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "events": []any{}})
}

func ReportsHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "reports": []any{}})
}

func SettingsHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "settings": map[string]any{}})
}
