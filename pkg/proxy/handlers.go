package proxy

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/proxy/middleware"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// CachedTenableCallParams strucuture to hold handler that will/won't be cached on call
type CachedTenableCallParams struct {
	w            http.ResponseWriter
	r            *http.Request
	f            func(t tenable.Service) (json []byte, err error)
	endPoint     tenable.EndPointType
	metricType   metrics.EndPointType
	metricMethod metrics.ServiceMethodType
}

// CachedTenableCall takes a full-hanlder and checks X-Cache-SkipOnHit and X-Cache-WriteOnReturn
func (s *Server) CachedTenableCall(pp CachedTenableCallParams) {
	skip, _ := strconv.ParseBool(middleware.SkipOnHit(pp.r))
	write, _ := strconv.ParseBool(middleware.WriteOnReturn(pp.r))
	s.CallSkipSave(pp, skip, write)
}

// CallSkipSave checks for the server cache file and serve it,
// otherwise create a Tenable.io Service and make the remote call.
// Cache the
// results so another immediately call would use the server cache.
func (s *Server) CallSkipSave(pp CachedTenableCallParams, skipOnHit bool, writeOnReturn bool) {
	s.Metrics.ServerInc(pp.metricType, metrics.Methods.Service.Get)

	if skipOnHit {
		// Check for a cache hit! :- )
		bb, err := s.cacheFetch(pp.r, pp.endPoint, pp.metricType)

		if err == nil && len(bb) > 0 {
			// Cache hit, write to wesponsewriter.
			_, _ = pp.w.Write(bb) // TODO: Check the return values
			return
		}
	}

	// Unpack the AccessKeys and SecretKeys from middleware context
	ak := middleware.AccessKey(pp.r)
	sk := middleware.SecretKey(pp.r)

	t := tenable.NewService(s.ServiceBaseURL, sk, ak, s.Log)

	// Invoke func f to deal with Tenable service call return
	json, err := pp.f(t)
	if err != nil {
		http.Error(pp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	// Save the response to the lookup
	if writeOnReturn {
		s.cacheStore(pp.w, pp.r, json, pp.endPoint, pp.metricType)
		_, _ = pp.w.Write(json) // TODO: Check the return values
	}

	return
}

// Shutdown is called to terminate proxy server
func (s *Server) Shutdown(w http.ResponseWriter, r *http.Request) {
	s.Log.Infof("/Shutdown called - beginning shutdown")

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

// VulnsExportStart handler for starting a export-vulns
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

		return t.VulnsExportStart(string(body.Limit), string(body.Filters.Since))
	}
	s.CachedTenableCall(pp)
}

// VulnsExportStatus handler for getting the status of an export-vulns
func (s *Server) VulnsExportStatus(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.VulnsExportStatus
	pp.metricType = metrics.EndPoints.VulnsExportStatus
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		exportUUID := middleware.ExportUUID(r)
		return t.VulnsExportStatus(exportUUID)
	}

	s.CallSkipSave(pp, false, true)
}

// VulnsExportGet is a handler that brokers a cacheable proxied call
// for a specific chunkID (1..n) and the export UUID (universially unique identifier)
func (s *Server) VulnsExportGet(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{
		w:            w,
		r:            r,
		endPoint:     tenable.EndPoints.VulnsExportGet,
		metricType:   metrics.EndPoints.VulnsExportGet,
		metricMethod: metrics.Methods.Service.Get,

		f: func(t tenable.Service) ([]byte, error) {
			exportUUID := middleware.ExportUUID(r)
			chunkID := middleware.ChunkID(r)
			return t.VulnsExportGet(exportUUID, chunkID)
		},
	}

	s.CachedTenableCall(pp)
}

// AssetsExportStart handler for starting exports-assets
func (s *Server) AssetsExportStart(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.AssetsExportStart
	pp.metricType = metrics.EndPoints.AssetsExportStart
	pp.metricMethod = metrics.Methods.Service.Update
	pp.f = func(t tenable.Service) ([]byte, error) {
		bb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Log.Errorf("couldn't read asset export body: %v", err)
			return nil, err
		}

		var body tenable.ExportFilter
		err = json.Unmarshal(bb, &body)

		return t.AssetsExportStart(string(body.Limit), string(body.Filters.LastAssessed))
	}
	s.CachedTenableCall(pp)
}

// AssetsExportStatus handler for get the status of an exports-assets
func (s *Server) AssetsExportStatus(w http.ResponseWriter, r *http.Request) {
	var skipOnHit = false
	var writeOnReturn = true

	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.AssetsExportStatus
	pp.metricType = metrics.EndPoints.AssetsExportStatus
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		// Override the cache lookup, and make the call!
		//t.SkipOnHit, _ = strconv.ParseBool(middleware.SkipOnHit(r))
		//t.WriteOnReturn, _ = strconv.ParseBool(middleware.WriteOnReturn(r))
		exportUUID := middleware.ExportUUID(r)
		return t.AssetsExportStatus(exportUUID)
	}

	s.CallSkipSave(pp, skipOnHit, writeOnReturn)
}

// AssetsExportGet handler for get chunks to complete a started exports-assets
func (s *Server) AssetsExportGet(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.f = func(t tenable.Service) ([]byte, error) {
		exportUUID := middleware.ExportUUID(r)
		chunkID := middleware.ChunkID(r)
		return t.AssetsExportGet(exportUUID, chunkID)
	}
	pp.endPoint = tenable.EndPoints.AssetsExportGet
	pp.metricType = metrics.EndPoints.AssetsExportGet
	pp.metricMethod = metrics.Methods.Service.Get
	s.CachedTenableCall(pp)
}

// ScannersList handler for outputting all the known scanners
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

// AgentGroups handler for outputting all the known agent groups
func (s *Server) AgentGroups(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.ScannerAgentGroups
	pp.metricType = metrics.EndPoints.AgentGroups
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		scanner := middleware.ScannerID(r)
		return t.ScannerAgentGroups(scanner)
	}
	s.CachedTenableCall(pp)
}

// AgentsGroup handler to group agents
func (s *Server) AgentsGroup(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.AgentsGroup
	pp.metricType = metrics.EndPoints.AgentsGroup
	pp.metricMethod = metrics.Methods.Service.Update
	pp.f = func(t tenable.Service) ([]byte, error) {
		scannerID := middleware.ScannerID(r)
		agentID := middleware.AgentID(r)
		groupID := middleware.GroupID(r)

		return t.AgentGroup(agentID, groupID, scannerID)
	}
	s.CachedTenableCall(pp)
}

// AgentsUngroup handler for ungroup agents
func (s *Server) AgentsUngroup(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.AgentsUngroup
	pp.metricType = metrics.EndPoints.AgentsUngroup
	pp.metricMethod = metrics.Methods.Service.Update
	pp.f = func(t tenable.Service) ([]byte, error) {
		scannerID := middleware.ScannerID(r)
		agentID := middleware.AgentID(r)
		groupID := middleware.GroupID(r)
		return t.AgentUngroup(agentID, groupID, scannerID)
	}
	s.CachedTenableCall(pp)
}

// AgentsList handler for outputting all the known agents
func (s *Server) AgentsList(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.AgentsList
	pp.metricType = metrics.EndPoints.AgentsList
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		scanner := middleware.ScannerID(r)
		offset := middleware.Offset(r)
		limit := middleware.Limit(r)
		return t.AgentList(scanner, offset, limit)
	}
	s.CachedTenableCall(pp)
}

// Scans will call /scans
func (s *Server) Scans(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.ScansList
	pp.metricType = metrics.EndPoints.ScansList
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		return t.ScansList()
	}
	s.CachedTenableCall(pp)
}

// ScanDetail handler for outputting scan details for ScanID
func (s *Server) ScanDetail(w http.ResponseWriter, r *http.Request) {
	// We need ScanID
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.ScanDetails
	pp.metricType = metrics.EndPoints.ScanDetails
	pp.metricMethod = metrics.Methods.Service.Get
	pp.f = func(t tenable.Service) ([]byte, error) {
		uuid := middleware.ScanUUID(r)
		return t.ScanDetails(uuid)
	}
	s.CachedTenableCall(pp)
	return
}

// ScanHistory handler for outputting scan details for ScanID
func (s *Server) ScanHistory(w http.ResponseWriter, r *http.Request) {
	return
}

// ScansExportStart handler for starting exports-scans
func (s *Server) ScansExportStart(w http.ResponseWriter, r *http.Request) {
	var pp = CachedTenableCallParams{w: w, r: r}
	pp.endPoint = tenable.EndPoints.ScansExportStart
	pp.metricType = metrics.EndPoints.ScansExportStart
	pp.metricMethod = metrics.Methods.Service.Update
	pp.f = func(t tenable.Service) ([]byte, error) {
		scanid := middleware.ScanID(r) // We use chi apped UUI
		histid := middleware.HistoryID(r)
		return t.ScansExportStart(scanid, histid)
	}
	s.CachedTenableCall(pp)
}
