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

		r.Route("/vulns", func(r chi.Router) {
			r.Route("/export", func(r chi.Router) {
				r.Post("/", s.VulnsExportStart)

				r.Route("/{ExportUUID}", func(r chi.Router) {
					r.Use(middleware.ExportCtx)
					r.Get("/status", s.VulnsExportStatus)

					r.Route("/chunks/{ChunkID}", func(r chi.Router) {
						r.Use(middleware.ExportChunkCtx)
						r.Get("/", s.VulnsExportGet)
					})
				})
			})
		})

		r.Route("/scanners", func(r chi.Router) {
			r.Get("/", s.ScannersList)
			r.Route("/{ScannerUUID}", func(r chi.Router) {
				r.Use(middleware.ScannersCtx)
				r.Get("/agents", s.AgentsList)
				r.Get("/agent-groups", s.AgentGroups)

			})
		})

	})
}
