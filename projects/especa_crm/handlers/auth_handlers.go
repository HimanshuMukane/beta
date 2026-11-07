package handlers

import (
	"net/http"

	"core/utils"
	"superadmin/ui/pages"
)

func (h *CRMHandlers) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.SignIn().Render(r.Context(), w)
}

func (h *CRMHandlers) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.SignUp().Render(r.Context(), w)
}

func (h *CRMHandlers) ForgotPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.ForgotPassword().Render(r.Context(), w)
}

func (h *CRMHandlers) NotFoundPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = pages.NotFound().Render(r.Context(), w)
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

