package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/prunus/pkg/middleware"
	transport "github.com/prunus/pkg/transport/http"
)

func AuthRouter(h *transport.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth())
		r.Get("/me", h.GetMe)
		r.Post("/logout", h.Logout)
		r.Post("/refresh-token", h.RefreshToken)
	})

	return r
}
