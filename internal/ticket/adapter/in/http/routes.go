package http

import (
	"github.com/go-chi/chi/v5"

	appmw "github.com/jhonsferg/corelog/internal/platform/middleware"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/tickets", func(r chi.Router) {
		r.Use(appmw.Authenticate(h.jwtSecret))

		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}/status", h.UpdateStatus)
		r.Patch("/{id}/assign", h.Assign)
	})
}
