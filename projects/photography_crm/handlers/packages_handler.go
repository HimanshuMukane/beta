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

type PackageServiceItemDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Cost     int64  `json:"cost"`
	Quantity int    `json:"quantity"`
	IsCustom bool   `json:"isCustom,omitempty"`
}

type PackageServiceGroupDTO struct {
	ID        string                  `json:"id"`
	GroupName string                  `json:"groupName"`
	Quantity  int                     `json:"quantity"`
	Services  []PackageServiceItemDTO `json:"services"`
}

type PackageDTO struct {
	ID            string                   `json:"id"`
	TenantID      string                   `json:"tenantId"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	TotalCost     int64                    `json:"totalCost"`
	ServiceGroups []PackageServiceGroupDTO `json:"serviceGroups"`
	Deliverables  []string                 `json:"deliverables"`
}

type CreatePackageRequest struct {
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	TotalCost     int64                    `json:"totalCost"`
	ServiceGroups []PackageServiceGroupDTO `json:"serviceGroups"`
	Deliverables  []string                 `json:"deliverables"`
}

type UpdatePackageRequest struct {
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	TotalCost     int64                    `json:"totalCost"`
	ServiceGroups []PackageServiceGroupDTO `json:"serviceGroups"`
	Deliverables  []string                 `json:"deliverables"`
}

func getStringVal(value any) string {
	if value == nil {
		return ""
	}
	if stringVal, ok := value.(string); ok {
		return stringVal
	}
	return ""
}

func getInt64Val(value any) int64 {
	if value == nil {
		return 0
	}
	switch numberVal := value.(type) {
	case float64:
		return int64(numberVal)
	case int64:
		return numberVal
	case int:
		return int64(numberVal)
	}
	return 0
}

func parsePackageServiceGroups(rawBytes []byte) ([]PackageServiceGroupDTO, error) {
	if len(rawBytes) == 0 {
		return []PackageServiceGroupDTO{}, nil
	}
	var directGroups []PackageServiceGroupDTO
	var unmarshalErr error = json.Unmarshal(rawBytes, &directGroups)
	if unmarshalErr != nil {
		var rawGroups []map[string]any
		var legacyErr error = json.Unmarshal(rawBytes, &rawGroups)
		if legacyErr == nil {
			var result = make([]PackageServiceGroupDTO, 0, len(rawGroups))
			for _, groupMap := range rawGroups {
				var groupName string = getStringVal(groupMap["group_name"])
				if groupName == "" {
					groupName = getStringVal(groupMap["groupName"])
				}
				var groupDTO = PackageServiceGroupDTO{
					GroupName: groupName,
					Services:  make([]PackageServiceItemDTO, 0),
				}
				if servicesList, ok := groupMap["services"].([]any); ok {
					for _, serviceItem := range servicesList {
						if serviceMap, ok := serviceItem.(map[string]any); ok {
							var cost int64 = getInt64Val(serviceMap["cost"])
							if cost == 0 {
								cost = getInt64Val(serviceMap["amount"])
							}
							var quantity int = int(getInt64Val(serviceMap["quantity"]))
							if quantity == 0 {
								quantity = int(getInt64Val(serviceMap["qty"]))
							}
							if quantity == 0 {
								quantity = 1
							}
							groupDTO.Services = append(groupDTO.Services, PackageServiceItemDTO{
								Name:     getStringVal(serviceMap["name"]),
								Cost:     cost,
								Quantity: quantity,
							})
						}
					}
				}
				result = append(result, groupDTO)
			}
			return result, nil
		}
		return []PackageServiceGroupDTO{}, unmarshalErr
	}
	if directGroups == nil {
		return []PackageServiceGroupDTO{}, nil
	}
	return directGroups, nil
}

func parsePackageDeliverables(rawBytes []byte) ([]string, error) {
	if len(rawBytes) == 0 {
		return []string{}, nil
	}
	var directDeliverables []string
	var unmarshalErr error = json.Unmarshal(rawBytes, &directDeliverables)
	if unmarshalErr == nil && directDeliverables != nil {
		return directDeliverables, nil
	}
	var rawItems []any
	var fallbackErr error = json.Unmarshal(rawBytes, &rawItems)
	if fallbackErr == nil {
		var result = make([]string, 0, len(rawItems))
		for _, item := range rawItems {
			if stringVal, ok := item.(string); ok {
				result = append(result, stringVal)
			} else if mapVal, ok := item.(map[string]any); ok {
				if nameVal, ok := mapVal["name"].(string); ok && nameVal != "" {
					result = append(result, nameVal)
				}
			}
		}
		return result, nil
	}
	return []string{}, fallbackErr
}

func PackagesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
		SELECT id, tenant_id, name, COALESCE(description, ''), total_cost, service_groups, deliverables
		FROM packages
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list packages", queryErr)
		return
	}
	defer rows.Close()

	var packagesList = make([]PackageDTO, 0)
	for rows.Next() {
		var packageItem PackageDTO
		var rawServiceGroupsBytes []byte
		var rawDeliverablesBytes []byte
		var scanErr error = rows.Scan(
			&packageItem.ID,
			&packageItem.TenantID,
			&packageItem.Name,
			&packageItem.Description,
			&packageItem.TotalCost,
			&rawServiceGroupsBytes,
			&rawDeliverablesBytes,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read package", scanErr)
			return
		}
		parsedServiceGroups, sgErr := parsePackageServiceGroups(rawServiceGroupsBytes)
		if sgErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Package contains corrupted service groups data", sgErr)
			return
		}
		parsedDeliverables, delivErr := parsePackageDeliverables(rawDeliverablesBytes)
		if delivErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Package contains corrupted deliverables data", delivErr)
			return
		}
		packageItem.ServiceGroups = parsedServiceGroups
		packageItem.Deliverables = parsedDeliverables
		packagesList = append(packagesList, packageItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating packages", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, packagesList)
}

func CreatePackageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload CreatePackageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var packageName string = strings.TrimSpace(requestPayload.Name)
	if packageName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Package name is required")
		return
	}

	if requestPayload.ServiceGroups == nil {
		requestPayload.ServiceGroups = []PackageServiceGroupDTO{}
	}
	if requestPayload.Deliverables == nil {
		requestPayload.Deliverables = []string{}
	}

	serviceGroupsBytes, _ := json.Marshal(requestPayload.ServiceGroups)
	deliverablesBytes, _ := json.Marshal(requestPayload.Deliverables)

	var packageID string = uuid.New().String()
	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO packages (id, tenant_id, name, description, total_cost, service_groups, deliverables)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, packageID, tenantID, packageName, requestPayload.Description, requestPayload.TotalCost, serviceGroupsBytes, deliverablesBytes)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create package", insertErr)
		return
	}

	var createdPackage = PackageDTO{
		ID:            packageID,
		TenantID:      tenantID,
		Name:          packageName,
		Description:   requestPayload.Description,
		TotalCost:     requestPayload.TotalCost,
		ServiceGroups: requestPayload.ServiceGroups,
		Deliverables:  requestPayload.Deliverables,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdPackage)
}

func UpdatePackageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var packageID string = chi.URLParam(request, "id")
	if packageID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Package ID is required")
		return
	}

	var requestPayload UpdatePackageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var packageName string = strings.TrimSpace(requestPayload.Name)
	if packageName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Package name is required")
		return
	}

	if requestPayload.ServiceGroups == nil {
		requestPayload.ServiceGroups = []PackageServiceGroupDTO{}
	}
	if requestPayload.Deliverables == nil {
		requestPayload.Deliverables = []string{}
	}

	serviceGroupsBytes, _ := json.Marshal(requestPayload.ServiceGroups)
	deliverablesBytes, _ := json.Marshal(requestPayload.Deliverables)

	_, updateErr := database.ExecContext(request.Context(), `
		UPDATE packages
		SET name = ?, description = ?, total_cost = ?, service_groups = ?, deliverables = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, packageName, requestPayload.Description, requestPayload.TotalCost, serviceGroupsBytes, deliverablesBytes, packageID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update package", updateErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func DeletePackageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var packageID string = chi.URLParam(request, "id")
	if packageID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Package ID is required")
		return
	}

	_, deleteErr := database.ExecContext(request.Context(), `
		UPDATE packages SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, packageID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete package", deleteErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}
