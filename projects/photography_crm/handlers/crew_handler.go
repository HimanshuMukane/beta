package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func CrewMembersListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Database unavailable",
			},
		})
		return
	}
	tid := getTenantID(r, db)
	if tid == "" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Valid workspace membership required",
			},
		})
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT u.id, u.first_name, u.last_name, u.email, COALESCE(u.phone, ''),
		       COALESCE(u.role_id, ''), COALESCE(r.name, '')
		FROM users u
		INNER JOIN memberships m ON m.user_id = u.id AND m.tenant_id = ? AND m.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL
		WHERE u.deleted_at IS NULL AND (r.name IS NULL OR LOWER(r.name) != 'admin')
		ORDER BY u.created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to list crew members",
			},
		})
		return
	}
	defer rows.Close()

	members := make([]CrewMemberDTO, 0)
	for rows.Next() {
		var id, fn, ln, email, phone, roleID, roleName string
		if err := rows.Scan(&id, &fn, &ln, &email, &phone, &roleID, &roleName); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read crew member"})
			return
		}
		fullName := strings.TrimSpace(fn + " " + ln)
		if fullName == "" {
			fullName = fn
		}
		members = append(members, CrewMemberDTO{
			ID:          id,
			TenantID:    tid,
			FullName:    fullName,
			FirstName:   fn,
			LastName:    ln,
			Email:       email,
			MobileNo:    phone,
			CountryCode: "+91",
			RoleID:      roleID,
			Designation: roleName,
		})
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error iterating crew members"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, members)
}

func CreateCrewMemberHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Database unavailable",
			},
		})
		return
	}
	tid := getTenantID(r, db)
	if tid == "" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Valid workspace membership required",
			},
		})
		return
	}

	var req struct {
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request body",
			},
		})
		return
	}

	fullName := strings.TrimSpace(req.FullName)
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)

	if fullName != "" && firstName == "" {
		parts := strings.Fields(fullName)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = strings.Join(parts[1:], " ")
		}
	}
	if firstName == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "First name or full name is required",
			},
		})
		return
	}

	designation := strings.TrimSpace(req.Designation)
	roleID := strings.TrimSpace(req.RoleID)

	if roleID != "" {
		var roleName string
		if err := db.QueryRowContext(r.Context(), "SELECT name FROM roles WHERE id = ? AND deleted_at IS NULL", roleID).Scan(&roleName); err == nil {
			if strings.ToLower(roleName) == "admin" {
				utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
					"error": map[string]any{
						"code":    "VALIDATION_ERROR",
						"message": "Cannot assign admin role to crew member",
					},
				})
				return
			}
			designation = roleName
		}
	} else if designation != "" {
		if strings.ToLower(designation) == "admin" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": map[string]any{
					"code":    "VALIDATION_ERROR",
					"message": "Cannot assign admin role to crew member",
				},
			})
			return
		}
		var matchedID string
		if err := db.QueryRowContext(r.Context(), "SELECT id FROM roles WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL LIMIT 1", designation).Scan(&matchedID); err == nil {
			roleID = matchedID
		} else {
			roleID = uuid.New().String()
			_, _ = db.ExecContext(r.Context(), "INSERT INTO roles (id, name) VALUES (?, ?)", roleID, designation)
		}
	}

	userID := uuid.New().String()
	email := strings.TrimSpace(req.Email)
	if email == "" {
		email = strings.ToLower(strings.ReplaceAll(firstName, " ", ".")) + "." + userID[:6] + "@especacrew.local"
	}
	phone := strings.TrimSpace(req.MobileNo)

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to start transaction",
			},
		})
		return
	}
	defer tx.Rollback()

	var rID *string
	if roleID != "" {
		rID = &roleID
	}

	_, err = tx.ExecContext(r.Context(), `
		INSERT INTO users (id, first_name, last_name, email, phone, role_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, firstName, lastName, email, phone, rID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}

	membershipID := uuid.New().String()
	roleNameInMembership := designation
	if roleNameInMembership == "" {
		roleNameInMembership = "photographer"
	}

	_, err = tx.ExecContext(r.Context(), `
		INSERT INTO memberships (id, tenant_id, user_id, role, status)
		VALUES (?, ?, ?, ?, 'active')
	`, membershipID, tid, userID, roleNameInMembership)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to assign workspace membership",
			},
		})
		return
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to commit crew member",
			},
		})
		return
	}

	dispFullName := strings.TrimSpace(firstName + " " + lastName)
	utils.WriteJSON(w, http.StatusCreated, CrewMemberDTO{
		ID:          userID,
		TenantID:    tid,
		FullName:    dispFullName,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		MobileNo:    phone,
		CountryCode: "+91",
		RoleID:      roleID,
		Designation: designation,
	})
}

func DeleteCrewMemberHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Database unavailable",
			},
		})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "User ID is required",
			},
		})
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to start transaction",
			},
		})
		return
	}
	defer tx.Rollback()

	_, _ = tx.ExecContext(r.Context(), `
		UPDATE memberships SET deleted_at = NOW(6) WHERE user_id = ? AND tenant_id = ?
	`, id, tid)

	_, err = tx.ExecContext(r.Context(), `
		UPDATE users SET deleted_at = NOW(6) WHERE id = ?
	`, id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete crew member",
			},
		})
		return
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to commit delete",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
