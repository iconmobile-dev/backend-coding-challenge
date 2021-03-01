package user

import (
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Route("/users", func(r chi.Router) {
		r.Post("/v1/UserCreate", s.userCreateRoute)
		r.Post("/v1/UserGet", s.userGetRoute)
		r.Post("/v1/UserDelete", s.userDeleteRoute)
	})

	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/v1/Login", s.loginRoute)
		r.Post("/v1/Logout", s.logoutRoute)
	})
}
