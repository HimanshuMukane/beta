package handlers

import (
	"net/http"

	"core/utils"
)

func (h *CRMHandlers) SignInPageHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "SignIn endpoint"})
}

func (h *CRMHandlers) SignUpPageHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "SignUp endpoint"})
}

func (h *CRMHandlers) ForgotPasswordPageHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "ForgotPassword endpoint"})
}

func (h *CRMHandlers) NotFoundPageHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func (h *CRMHandlers) LogoutHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     cfg.SessionCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
}
