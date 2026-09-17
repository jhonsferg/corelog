package http

import (
	"github.com/go-chi/chi/v5"

	appmw "github.com/jhonsferg/corelog/internal/platform/middleware"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(appmw.Authenticate(h.jwtSecret))

		r.Get("/me", h.Me)
		r.Patch("/me", h.UpdateProfile)
		r.Post("/me/password", h.ChangePassword)
		r.Get("/search", h.Search)

		r.Group(func(r chi.Router) {
			r.Use(appmw.RequireRole("admin"))

			r.Get("/", h.List)
			r.Patch("/{id}/team", h.SetTeam)
		})
	})
}
