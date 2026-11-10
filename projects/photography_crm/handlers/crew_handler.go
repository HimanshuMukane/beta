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

type CrewMemberDTO struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenantId"`
	FullName    string  `json:"fullName"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	MobileNo    string  `json:"mobileNo"`
	CountryCode string  `json:"countryCode"`
	RoleID      string  `json:"roleId"`
	Designation string  `json:"designation"`
	Note        *string `json:"note,omitempty"`
}

type CreateCrewMemberRequest struct {
	FullName    string `json:"fullName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	MobileNo    string `json:"mobileNo"`
	CountryCode string `json:"countryCode"`
	Designation string `json:"designation"`
	RoleID      string `json:"roleId"`
	Note        string `json:"note"`
}

func CrewMembersListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), `
		SELECT u.id, u.first_name, u.last_name, u.email, COALESCE(u.phone, ''),
		       COALESCE(u.role_id, ''), COALESCE(r.name, '')
		FROM users u
		INNER JOIN memberships m ON m.user_id = u.id AND m.tenant_id = ? AND m.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL
		WHERE u.deleted_at IS NULL AND (r.name IS NULL OR LOWER(r.name) != 'admin')
		ORDER BY u.created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list crew members", queryErr)
		return
	}
	defer rows.Close()

	var crewMembers = make([]CrewMemberDTO, 0)
	for rows.Next() {
		var memberID string
		var firstName string
		var lastName string
		var email string
		var phone string
		var roleID string
		var roleName string
		var scanErr error = rows.Scan(&memberID, &firstName, &lastName, &email, &phone, &roleID, &roleName)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read crew member", scanErr)
			return
		}
		var fullName string = strings.TrimSpace(firstName + " " + lastName)
		if fullName == "" {
			fullName = firstName
		}
		crewMembers = append(crewMembers, CrewMemberDTO{
			ID:          memberID,
			TenantID:    tenantID,
			FullName:    fullName,
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			MobileNo:    phone,
			CountryCode: "+91",
			RoleID:      roleID,
			Designation: roleName,
		})
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating crew members", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, crewMembers)
}

func CreateCrewMemberHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var requestPayload CreateCrewMemberRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var fullName string = strings.TrimSpace(requestPayload.FullName)
	var firstName string = strings.TrimSpace(requestPayload.FirstName)
	var lastName string = strings.TrimSpace(requestPayload.LastName)

	if fullName != "" && firstName == "" {
		var nameParts = strings.Fields(fullName)
		firstName = nameParts[0]
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}
	}
	if firstName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "First name or full name is required")
		return
	}

	var designation string = strings.TrimSpace(requestPayload.Designation)
	var roleID string = strings.TrimSpace(requestPayload.RoleID)

	if roleID != "" {
		var roleName string
		var roleQueryErr error = database.QueryRowContext(request.Context(), "SELECT name FROM roles WHERE id = ? AND deleted_at IS NULL", roleID).Scan(&roleName)
		if roleQueryErr == nil {
			if strings.ToLower(roleName) == "admin" {
				utils.WriteError(responseWriter, http.StatusBadRequest, "Cannot assign admin role to crew member")
				return
			}
			designation = roleName
		}
	} else if designation != "" {
		if strings.ToLower(designation) == "admin" {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Cannot assign admin role to crew member")
			return
		}
		var matchedRoleID string
		var matchErr error = database.QueryRowContext(request.Context(), "SELECT id FROM roles WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL LIMIT 1", designation).Scan(&matchedRoleID)
		if matchErr == nil {
			roleID = matchedRoleID
		} else {
			roleID = uuid.New().String()
			_, _ = database.ExecContext(request.Context(), "INSERT INTO roles (id, name) VALUES (?, ?)", roleID, designation)
		}
	}

	var userID string = uuid.New().String()
	var email string = strings.TrimSpace(requestPayload.Email)
	if email == "" {
		email = strings.ToLower(strings.ReplaceAll(firstName, " ", ".")) + "." + userID[:6] + "@especacrew.local"
	}
	var phone string = strings.TrimSpace(requestPayload.MobileNo)

	var transaction *sqlx.Tx
	var beginTxErr error
	transaction, beginTxErr = database.BeginTxx(request.Context(), nil)
	if beginTxErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to start transaction", beginTxErr)
		return
	}
	defer transaction.Rollback()

	var roleIDPointer *string = nil
	if roleID != "" {
		roleIDPointer = &roleID
	}

	var insertUserErr error
	_, insertUserErr = transaction.ExecContext(request.Context(), `
		INSERT INTO users (id, first_name, last_name, email, phone, role_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, firstName, lastName, email, phone, roleIDPointer)
	if insertUserErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create user", insertUserErr)
		return
	}

	var membershipID string = uuid.New().String()
	var membershipRole string = designation
	if membershipRole == "" {
		membershipRole = "photographer"
	}

	var insertMembershipErr error
	_, insertMembershipErr = transaction.ExecContext(request.Context(), `
		INSERT INTO memberships (id, tenant_id, user_id, role, status)
		VALUES (?, ?, ?, ?, 'active')
	`, membershipID, tenantID, userID, membershipRole)
	if insertMembershipErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to assign workspace membership", insertMembershipErr)
		return
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit crew member", commitErr)
		return
	}

	var displayFullName string = strings.TrimSpace(firstName + " " + lastName)
	var createdCrewMember = CrewMemberDTO{
		ID:          userID,
		TenantID:    tenantID,
		FullName:    displayFullName,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		MobileNo:    phone,
		CountryCode: "+91",
		RoleID:      roleID,
		Designation: designation,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdCrewMember)
}

func DeleteCrewMemberHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	var memberID string = chi.URLParam(request, "id")
	if memberID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "User ID is required")
		return
	}

	var transaction *sqlx.Tx
	var beginTxErr error
	transaction, beginTxErr = database.BeginTxx(request.Context(), nil)
	if beginTxErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to start transaction", beginTxErr)
		return
	}
	defer transaction.Rollback()

	_, _ = transaction.ExecContext(request.Context(), `
		UPDATE memberships SET deleted_at = NOW(6) WHERE user_id = ? AND tenant_id = ?
	`, memberID, tenantID)

	var updateErr error
	_, updateErr = transaction.ExecContext(request.Context(), `
		UPDATE users SET deleted_at = NOW(6) WHERE id = ?
	`, memberID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete crew member", updateErr)
		return
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit delete", commitErr)
		return
	}

	var successResponse = ActionStatusResponse{
		Success: true,
		Message: "Crew member deleted successfully",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}
