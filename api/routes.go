package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Routes(h AnalyticsHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/all", h.Get)
	r.Post("/api/analytics", h.Post)

	return r
}
