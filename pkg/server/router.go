package server

import (
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/whereiskurt/tiogo/pkg/server/middleware"
)

// EnableDefaultRouter defines routes with middleware for request tracking, logging, param contexts
func (s *Server) EnableDefaultRouter() {
	s.Router.Use(chimiddleware.RequestID)
	s.Router.Use(middleware.NewStructuredLogger(s.Log))
	s.Router.Use(chimiddleware.Recoverer)

	s.Router.Route("/", func(r chi.Router) {

		r.Use(middleware.InitialCtx)
		r.Use(middleware.PrettyResponseCtx)

		r.Get("/shutdown", s.Shutdown)

	})
}
