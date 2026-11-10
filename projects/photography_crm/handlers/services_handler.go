package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ServiceItemDTO struct {
	ID          string `json:"id"`
	ServiceName string `json:"serviceName"`
	Quantity    int    `json:"quantity"`
	Amount      int64  `json:"amount"`
	IsCustom    bool   `json:"isCustom,omitempty"`
}

type ServiceGroupDTO struct {
	ID        string           `json:"id"`
	TenantID  string           `json:"tenantId"`
	GroupName string           `json:"groupName"`
	Services  []ServiceItemDTO `json:"services"`
}

type CreateServiceGroupRequest struct {
	GroupName string           `json:"groupName"`
	Services  []ServiceItemDTO `json:"services"`
}

type UpdateServiceGroupRequest struct {
	GroupName string           `json:"groupName"`
	Services  []ServiceItemDTO `json:"services"`
}

type BatchSaveServiceGroupItem struct {
	ID        string           `json:"id"`
	GroupName string           `json:"groupName"`
	Services  []ServiceItemDTO `json:"services"`
}

type BatchSaveServicesRequest struct {
	Groups []BatchSaveServiceGroupItem `json:"groups"`
}

func parseServicesList(rawBytes []byte) []ServiceItemDTO {
	if len(rawBytes) == 0 {
		return []ServiceItemDTO{}
	}
	var directServices []ServiceItemDTO
	var unmarshalErr error = json.Unmarshal(rawBytes, &directServices)
	if unmarshalErr == nil && directServices != nil {
		return directServices
	}
	var rawItems []map[string]any
	var fallbackErr error = json.Unmarshal(rawBytes, &rawItems)
	if fallbackErr != nil {
		return []ServiceItemDTO{}
	}
	var servicesList = make([]ServiceItemDTO, 0, len(rawItems))
	for _, rawItem := range rawItems {
		var serviceID string
		if idVal, ok := rawItem["id"].(string); ok {
			serviceID = idVal
		}
		var serviceName string
		if nameVal, ok := rawItem["serviceName"].(string); ok {
			serviceName = nameVal
		} else if nameVal, ok := rawItem["name"].(string); ok {
			serviceName = nameVal
		}
		var quantity int = 1
		if qtyVal, ok := rawItem["quantity"].(float64); ok {
			quantity = int(qtyVal)
		} else if qtyVal, ok := rawItem["qty"].(float64); ok {
			quantity = int(qtyVal)
		}
		if quantity <= 0 {
			quantity = 1
		}
		var amount int64 = 0
		if amtVal, ok := rawItem["amount"].(float64); ok {
			amount = int64(amtVal)
		} else if amtVal, ok := rawItem["cost"].(float64); ok {
			amount = int64(amtVal)
		}
		var isCustom bool = false
		if customVal, ok := rawItem["isCustom"].(bool); ok {
			isCustom = customVal
		}
		servicesList = append(servicesList, ServiceItemDTO{
			ID:          serviceID,
			ServiceName: serviceName,
			Quantity:    quantity,
			Amount:      amount,
			IsCustom:    isCustom,
		})
	}
	return servicesList
}

func ServicesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
		SELECT id, tenant_id, group_name, services
		FROM services
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list services", queryErr)
		return
	}
	defer rows.Close()

	var serviceGroups = make([]ServiceGroupDTO, 0)
	for rows.Next() {
		var serviceGroup ServiceGroupDTO
		var rawServicesBytes []byte
		var scanErr error = rows.Scan(&serviceGroup.ID, &serviceGroup.TenantID, &serviceGroup.GroupName, &rawServicesBytes)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read service group", scanErr)
			return
		}
		serviceGroup.Services = parseServicesList(rawServicesBytes)
		serviceGroups = append(serviceGroups, serviceGroup)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating service groups", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, serviceGroups)
}

func CreateServiceGroupHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload CreateServiceGroupRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var groupName string = strings.TrimSpace(requestPayload.GroupName)
	if groupName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Group name is required")
		return
	}

	if requestPayload.Services == nil {
		requestPayload.Services = []ServiceItemDTO{}
	}
	servicesBytes, _ := json.Marshal(requestPayload.Services)

	var serviceGroupID string = uuid.New().String()
	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO services (id, tenant_id, group_name, services)
		VALUES (?, ?, ?, ?)
	`, serviceGroupID, tenantID, groupName, servicesBytes)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create service group", insertErr)
		return
	}

	var createdGroup = ServiceGroupDTO{
		ID:        serviceGroupID,
		TenantID:  tenantID,
		GroupName: groupName,
		Services:  requestPayload.Services,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdGroup)
}

func UpdateServiceGroupHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var serviceGroupID string = chi.URLParam(request, "id")
	if serviceGroupID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Service group ID is required")
		return
	}

	var requestPayload UpdateServiceGroupRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var groupName string = strings.TrimSpace(requestPayload.GroupName)
	if groupName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Group name is required")
		return
	}

	if requestPayload.Services == nil {
		requestPayload.Services = []ServiceItemDTO{}
	}
	servicesBytes, _ := json.Marshal(requestPayload.Services)

	_, updateErr := database.ExecContext(request.Context(), `
		UPDATE services
		SET group_name = ?, services = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, groupName, servicesBytes, serviceGroupID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update service group", updateErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func DeleteServiceGroupHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var serviceGroupID string = chi.URLParam(request, "id")
	if serviceGroupID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Service group ID is required")
		return
	}

	_, deleteErr := database.ExecContext(request.Context(), `
		UPDATE services SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, serviceGroupID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete service group", deleteErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func BatchSaveServicesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload BatchSaveServicesRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, beginTxErr := database.BeginTxx(request.Context(), nil)
	if beginTxErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to start database transaction", beginTxErr)
		return
	}
	defer transaction.Rollback()

	_, resetErr := transaction.ExecContext(request.Context(), "UPDATE services SET deleted_at = NOW(6) WHERE tenant_id = ? AND deleted_at IS NULL", tenantID)
	if resetErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update existing services", resetErr)
		return
	}

	for _, groupItem := range requestPayload.Groups {
		var groupName string = strings.TrimSpace(groupItem.GroupName)
		if groupName == "" {
			continue
		}
		var groupID string = groupItem.ID
		if groupID == "" || strings.HasPrefix(groupID, "temp-") || strings.HasPrefix(groupID, "srv-") {
			groupID = uuid.New().String()
		}
		var services = groupItem.Services
		if services == nil {
			services = []ServiceItemDTO{}
		}
		servicesBytes, _ := json.Marshal(services)

		_, insertErr := transaction.ExecContext(request.Context(), `
			INSERT INTO services (id, tenant_id, group_name, services)
			VALUES (?, ?, ?, ?)
		`, groupID, tenantID, groupName, servicesBytes)
		if insertErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to save service group", insertErr)
			return
		}
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit services", commitErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}
