package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/server/middleware"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"net/http"
	"time"
)

func (s *Server) Shutdown(w http.ResponseWriter, r *http.Request) {
	s.Log.Debugf("/Shutdown called - beginning s Shutdown")

	_, _ = w.Write([]byte(ui.Gopher()))
	_, _ = w.Write([]byte("\n...bye felicia\n"))

	timeout, cancel := context.WithTimeout(s.Context, 5*time.Second)
	err := s.HTTP.Shutdown(timeout)
	if err != nil && err != context.Canceled {
		s.Log.Errorf("server error during Shutdown: %+v", err)
	}

	s.Finished()
	cancel()
}

func (s *Server) VulnsExport(w http.ResponseWriter, r *http.Request) {
	endPoint := tenable.EndPoints.VulnsExport
	metricType := metrics.EndPoints.VulnsExport

	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Update)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, metricType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

}
func (s *Server) VulnsExportStatus(w http.ResponseWriter, r *http.Request) {
	s.Metrics.ServerInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Update)
	exportUUID := chi.URLParam(r, "ExportUUID")
	s.Log.Infof("Check status of ExportUUDI: %s", exportUUID)

}
func (s *Server) VulnsExportChunk(w http.ResponseWriter, r *http.Request) {
	s.Metrics.ServerInc(metrics.EndPoints.VulnsExportChunk, metrics.Methods.Service.Update)
	exportUUID := chi.URLParam(r, "ExportUUID")
	chunkID := chi.URLParam(r, "ChunkID")
	s.Log.Infof("Fetching ExportUUID: %s, ChunkID: %s", exportUUID, chunkID)

	ak := middleware.ContextMap(r)["SecretKey"]
	sk := middleware.ContextMap(r)["AccessKey"]

	bb := fmt.Sprintf("AccessKey: %s, SecretKey: %s", ak, sk)

	_, _ = w.Write([]byte(bb))

}
