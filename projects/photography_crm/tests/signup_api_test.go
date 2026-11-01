package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignupAPI_Success(t *testing.T) {
	handler := setupLoginAPITest(t)

	payload := map[string]string{
		"first_name": "Alice",
		"last_name":  "Wonderland",
		"email":      "alice@example.com",
		"password":   "StrongPass123!",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Message string `json:"message"`
		User    struct {
			ID        string `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.User.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", resp.User.Email)
	}

	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "photography_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("expected session cookie to be set on signup")
	}

	loginPayload := map[string]string{
		"email":    "alice@example.com",
		"password": "StrongPass123!",
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()

	handler.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Fatalf("expected status 200 on login after signup, got %d: %s", loginW.Code, loginW.Body.String())
	}
}

func TestSignupAPI_DuplicateEmail(t *testing.T) {
	handler := setupLoginAPITest(t)

	payload := map[string]string{
		"first_name": "Duplicate",
		"last_name":  "User",
		"email":      testLoginEmail,
		"password":   "SomePass123!",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSignupAPI_ValidationErrors(t *testing.T) {
	handler := setupLoginAPITest(t)

	tests := []struct {
		name        string
		payload     map[string]string
		contentType string
		wantCode    int
	}{
		{
			name: "short password",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"email":      "short@example.com",
				"password":   "123",
			},
			contentType: "application/json",
			wantCode:    http.StatusBadRequest,
		},
		{
			name: "missing email",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"password":   "ValidPass123!",
			},
			contentType: "application/json",
			wantCode:    http.StatusBadRequest,
		},
		{
			name: "missing first name",
			payload: map[string]string{
				"last_name": "User",
				"email":     "nofirst@example.com",
				"password":  "ValidPass123!",
			},
			contentType: "application/json",
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(b))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)
			if w.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d: %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}
