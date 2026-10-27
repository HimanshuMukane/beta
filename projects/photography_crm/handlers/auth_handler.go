package handlers

import (
	"net/http"

	"chi-mysql-server/ui/pages"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	pages.SignIn().Render(r.Context(), w)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	pages.SignUp().Render(r.Context(), w)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	pages.ForgotPassword().Render(r.Context(), w)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	pages.NotFound().Render(r.Context(), w)
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	pages.InternalError().Render(r.Context(), w)
}
