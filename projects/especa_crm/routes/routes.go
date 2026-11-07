package routes

import (
	"net/http"

	coreHandlers "core/handlers"
	coreMiddleware "core/middleware"
	"core/utils"
	"especa_crm/handlers"
	"especa_crm/ui/static"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(database *sqlx.DB, companyID string) http.Handler {
	return SetupRouterWithAuth(database, companyID, utils.AppConfig)
}

func SetupRouterWithAuth(database *sqlx.DB, companyID string, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(coreMiddleware.RequestID)
	r.Use(coreMiddleware.RequestLogger)
	r.Use(coreMiddleware.Recoverer)
	r.Use(coreMiddleware.SecurityHeaders)
	r.Use(coreMiddleware.CSRFProtection)

	fs := http.StripPrefix("/static/", http.FileServer(static.FS()))
	r.Handle("/static/*", fs)

	r.Get("/api/health", func(w http.ResponseWriter, req *http.Request) {
		if database == nil || database.PingContext(req.Context()) != nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "db": "connected"})
	})

	r.Route("/api/auth", func(api chi.Router) {
		api.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		api.Post("/signup", coreHandlers.SignupAPIHandler(cfg, "workspace_user"))
		api.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))
		api.With(coreMiddleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
	})

	h := handlers.NewCRMHandlers(database, companyID)

	r.Get("/signin", h.SignInPageHandler)
	r.Get("/signup", h.SignUpPageHandler)
	r.Get("/forgot-password", h.ForgotPasswordPageHandler)
	r.Get("/logout", h.LogoutHandler(cfg))

	r.Group(func(authRouter chi.Router) {
		authRouter.Use(coreMiddleware.RequireAuthPage(cfg))

		authRouter.Get("/", h.DashboardHandler)

		authRouter.Get("/leads", h.LeadsListHandler)
		authRouter.Get("/leads/new", h.LeadCreateHandler)
		authRouter.Post("/leads/new", h.LeadCreateHandler)
		authRouter.Get("/leads/{id}", h.LeadDetailHandler)

		authRouter.Get("/deals", h.DealsPipelineHandler)
		authRouter.Get("/deals/new", h.DealCreateHandler)
		authRouter.Post("/deals/new", h.DealCreateHandler)
		authRouter.Get("/deals/{id}", h.DealDetailHandler)
		authRouter.Post("/deals/{id}/convert-project", h.ConvertDealToProjectHandler)

		authRouter.Get("/projects", h.ProjectsListHandler)
		authRouter.Get("/projects/new", h.ProjectCreateHandler)
		authRouter.Post("/projects/new", h.ProjectCreateHandler)
		authRouter.Get("/projects/{id}", h.ProjectDetailHandler)
		authRouter.Post("/projects/{id}/billing-event", h.TriggerProjectBillingEventHandler)
		authRouter.Post("/milestones/{id}/complete", h.CompleteMilestoneHandler)

		authRouter.Get("/schedule", h.ScheduleCalendarHandler)

		authRouter.Get("/quotes", h.QuotesListHandler)
		authRouter.Get("/quotes/new", h.QuoteCreateHandler)
		authRouter.Post("/quotes/new", h.QuoteCreateHandler)

		authRouter.Get("/contacts", h.ContactsListHandler)
		authRouter.Get("/contacts/new", h.ContactCreateHandler)
		authRouter.Post("/contacts/new", h.ContactCreateHandler)
		authRouter.Get("/contacts/{id}", h.ContactDetailHandler)
	})

	r.NotFound(h.NotFoundPageHandler)
	return r
}
