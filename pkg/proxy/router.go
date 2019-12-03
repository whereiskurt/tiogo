package proxy

import (
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/whereiskurt/tiogo/pkg/proxy/middleware"
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

		r.Route("/assets", func(r chi.Router) {
			r.Route("/export", func(r chi.Router) {
				r.Post("/", s.AssetsExportStart)

				r.Route("/{ExportUUID}", func(r chi.Router) {
					r.Use(middleware.ExportCtx)
					r.Get("/status", s.AssetsExportStatus)

					r.Route("/chunks/{ChunkID}", func(r chi.Router) {
						r.Use(middleware.ExportChunkCtx)
						r.Get("/", s.AssetsExportGet)
					})
				})
			})
		})

		r.Route("/scanners", func(r chi.Router) {
			r.Get("/", s.ScannersList)
			r.Route("/{ScannerID}", func(r chi.Router) {
				r.Use(middleware.ScannersCtx)
				r.Get("/agents", s.AgentsList)

				r.Route("/agent-groups", func(r chi.Router) {
					r.Get("/", s.AgentGroups)

					r.Route("/{GroupID}", func(r chi.Router) {
						r.Use(middleware.AgentGroupCtx)
						r.Route("/agents/{AgentID}", func(r chi.Router) {
							r.Use(middleware.AgentCtx)
							r.Put("/", s.AgentsGroup)
							r.Delete("/", s.AgentsUngroup)
						})
					})
				})
			})
		})

		// //Adding /scans endpoints
		// r.Route("/scans", func(r chi.Router) {
		// 	r.Get("/", s.Scans)
		// 	r.Route("/{ScanID}", func(r chi.Router) {
		// 		r.Use(middleware.ScansCtx)
		// 		r.Get("/", s.ScansList)
		// 		r.Get("/export", s.ScanExport)
		// 	})
		// })

	})
}
