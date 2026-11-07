package handlers

import (
	"net/http"

	"core/utils"
	"superadmin/ui/pages"
)

func (h *FinanceHandlers) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.SignIn().Render(r.Context(), w)
}

func (h *FinanceHandlers) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.SignUp().Render(r.Context(), w)
}

func (h *FinanceHandlers) ForgotPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.ForgotPassword().Render(r.Context(), w)
}

func (h *FinanceHandlers) NotFoundPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = pages.NotFound().Render(r.Context(), w)
}

func (h *FinanceHandlers) LogoutHandler(cfg utils.Config) http.HandlerFunc {
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
