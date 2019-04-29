package server

import (
	"context"
	"encoding/json"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/server/middleware"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"io/ioutil"
	"net/http"
	"time"
)

type CachedTenableCallParams struct {
	w            http.ResponseWriter
	r            *http.Request
	f            func(t tenable.Service) (json []byte, err error)
	endPoint     tenable.EndPointType
	metricType   metrics.EndPointType
	metricMethod metrics.ServiceMethodType
}

func (s *Server) CachedTenableCall(pp CachedTenableCallParams) {
	s.CachedTenableCallSkipSave(pp, true, true)
}

// Check for the server cache file and serve it, otherwise create a Tenable.io service and make the call.
// Cache the results so another immediately call would use the server cache. aka caching proxy server
func (s *Server) CachedTenableCallSkipSave(pp CachedTenableCallParams, skipOnHit bool, writeOnReturn bool) {
	s.Metrics.ServerInc(pp.metricType, metrics.Methods.Service.Get)

	if skipOnHit {
		// Check for a cache hit! :- )
		bb, err := s.cacheFetch(pp.r, pp.endPoint, pp.metricType)
		if err == nil && len(bb) > 0 {
			_, _ = pp.w.Write(bb)
			return
		}
	}

	// Take the AccessKeys and SecretKeys from context
	ak := middleware.AccessKey(pp.r)
	sk := middleware.SecretKey(pp.r)

	t := tenable.NewService(s.ServiceBaseURL, sk, ak)

	json, err := pp.f(t)
	if err != nil {
		http.Error(pp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if writeOnReturn {
		s.cacheStore(pp.w, pp.r, json, pp.endPoint, pp.metricType)
		_, _ = pp.w.Write(json)
	}

	return
}

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

func (s *Server) VulnsExportStart(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.VulnsExportStart
	pp.metricType = metrics.EndPoints.VulnsExportStart
	pp.metricMethod = metrics.Methods.Service.Update
	pp.f = func(t tenable.Service) ([]byte, error) {
		bb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Log.Errorf("couldn't read vuln export body: %v", err)
			return nil, err
		}

		var body tenable.ExportFilter
		err = json.Unmarshal(bb, &body)

		return t.VulnsExportStart(string(body.Filters.Since))
	}
	s.CachedTenableCall(pp)
}
func (s *Server) VulnsExportStatus(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.VulnsExportStatus
	pp.metricType = metrics.EndPoints.VulnsExportStatus
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		exportUUID := middleware.ExportUUID(r)
		return t.VulnsExportStatus(exportUUID, true , true)
	}

	s.CachedTenableCallSkipSave(pp, false, true)

}
func (s *Server) VulnsExportGet(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.f = func(t tenable.Service) ([]byte, error) {
		exportUUID := middleware.ExportUUID(r)
		chunkID := middleware.ChunkID(r)
		return t.VulnsExportGet(exportUUID, chunkID)
	}
	pp.endPoint = tenable.EndPoints.VulnsExportGet
	pp.metricType = metrics.EndPoints.VulnsExportGet
	pp.metricMethod = metrics.Methods.Service.Get
	s.CachedTenableCall(pp)
}

func (s *Server) ScannersList(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.ScannersList
	pp.metricType = metrics.EndPoints.ScannersList
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		return t.ScannersList()
	}
	s.CachedTenableCall(pp)
}
func (s *Server) AgentGroups(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.ScannerAgentGroups
	pp.metricType = metrics.EndPoints.AgentGroups
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		scanner := middleware.ScannerUUID(r)
		return t.ScannerAgentGroups(scanner)
	}
	s.CachedTenableCall(pp)
}
func (s *Server) AgentsList(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.AgentsList
	pp.metricType = metrics.EndPoints.AgentsList
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		scanner := middleware.ScannerUUID(r)
		offset := middleware.Offset(r)
		limit := middleware.Limit(r)
		return t.AgentList(scanner, offset, limit)
	}
	s.CachedTenableCall(pp)
}
