package http

import (
	"github.com/go-chi/chi/v5"

	appmw "github.com/jhonsferg/corelog/internal/platform/middleware"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/teams", func(r chi.Router) {
		r.Use(appmw.Authenticate(h.jwtSecret))

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)

		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireRole("admin"))

			r.Post("/", h.Create)
			r.Patch("/{id}", h.Rename)
			r.Delete("/{id}", h.Delete)
		})
	})
}
