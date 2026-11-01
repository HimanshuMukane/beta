package handlers

import (
	"net/http"

	"superadmin/ui/pages"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.SignIn().Render(r.Context(), w)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.SignUp().Render(r.Context(), w)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.ForgotPassword().Render(r.Context(), w)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = pages.NotFound().Render(r.Context(), w)
}
