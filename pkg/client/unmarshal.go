package client

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Unmarshal is responsible for constructing Tenable.io service, with or without cache,
// to make calls against Tenable.io and potential retrieve raw data. Raw data is converted
// by the Coverter
// it is the process
type Unmarshal struct {
	Config  *config.Config
	Metrics *metrics.Metrics
}

// NewUnmarshal calls the ACME EndPoints and returns ACME JSONs to the adapter
func NewUnmarshal(config *config.Config, metrics *metrics.Metrics) (u Unmarshal) {
	u.Config = config
	u.Metrics = metrics
	return
}

// NewService checks for previous cache hits, and if not present, calls and writes-cache
func (u *Unmarshal) NewService() (s tenable.Service) {
	return u.service(true, true)
}

// NewServiceSaveOnly does not check for previous cache hits, and makes a fresh call everytime
func (u *Unmarshal) NewServiceSaveOnly() (s tenable.Service) {
	return u.service(true, false)
}

// DefaultServiceFolder is appended to the cache to allow for many cache types in the cache folder
var DefaultServiceFolder = "service"

// service wraps the Tenable service calls in
func (u *Unmarshal) service(skipOnHit bool, writeOnReturn bool) (s tenable.Service) {
	s = tenable.NewService(u.Config.VM.BaseURL, u.Config.VM.SecretKey, u.Config.VM.AccessKey, u.Config.VM.Log)
	s.EnableMetrics(u.Metrics)

	if u.Config.VM.CacheResponse {
		serviceCacheFolder := filepath.Join(u.Config.VM.CacheFolder, DefaultServiceFolder)
		s.EnableCache(serviceCacheFolder, u.Config.VM.CacheKey)
	}
	s.Log = u.Config.VM.Log
	s.SkipOnHit = skipOnHit
	s.WriteOnReturn = writeOnReturn

	return
}

// ScannerAgentGroups outputs all Agent Groups associated with scanner
func (u *Unmarshal) ScannerAgentGroups(scannerID string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.ScannerAgentGroups(scannerID)
	return raw, err
}

// AgentGroup assigns an AgentID a GroupID given a ScannerID. [The opposite of AgentUngroup.]
func (u *Unmarshal) AgentGroup(agentID string, groupID string, scannerID string) ([]byte, error) {
	s := u.NewService()
	raw, err := s.AgentGroup(agentID, groupID, scannerID)
	return raw, err
}

// AgentUngroup unassigns an AgentID a GroupID given a ScannerID. [The opposite of AgentGroup.]
func (u *Unmarshal) AgentUngroup(agentID string, groupID string, scannerID string) ([]byte, error) {
	s := u.NewService()
	raw, err := s.AgentUngroup(agentID, groupID, scannerID)
	return raw, err
}

// Scanners returns all scanners registered in Tenable.io
func (u *Unmarshal) Scanners(skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.ScannersList()
	return raw, err
}

// Agents returns raw Agents from a givent offset and limit
func (u *Unmarshal) Agents(scannerID string, offset string, limit string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.AgentList(scannerID, offset, limit)
	return raw, err
}

// VulnsExportStart start vuln export, if not already started (i.e. cached)
func (u *Unmarshal) VulnsExportStart() ([]byte, error) {
	s := u.NewService()
	limit := u.Config.VM.ExportLimit

	// Convert Human dates into Unix()
	since := u.Config.VM.AfterDate
	tt, err := time.Parse(config.DateLayout, since)
	if err != nil {
		s.Log.Errorf("failed to export-vulns start: invalid since value: %s: %s", since, err)
		return nil, err
	}
	sinceUnix := fmt.Sprintf("%d", tt.Unix())

	raw, err := s.VulnsExportStart(limit, sinceUnix)

	return raw, err
}

// VulnsExportStatus will return the raw status of the vuln export
func (u *Unmarshal) VulnsExportStatus(uuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.VulnsExportStatus(uuid)
	return raw, err
}

// VulnsExportGet will return the raw chunk file for the vuln export uuid
func (u *Unmarshal) VulnsExportGet(uuid string, chunk string) ([]byte, error) {
	s := u.NewService()
	raw, err := s.VulnsExportGet(uuid, chunk)
	return raw, err
}

// AssetsExportStart creates request with limit and lastAssessed based on Config
func (u *Unmarshal) AssetsExportStart() ([]byte, error) {
	limit := u.Config.VM.ExportLimit
	lastAssessed := u.Config.VM.AfterDate

	s := u.NewService()

	tt, err := time.Parse(config.DateLayout, lastAssessed)
	if err != nil {
		s.Log.Errorf("failed to export-vulns start: invalid since value: %s: %s", lastAssessed, err)
		return nil, err
	}
	lastAssessedUnix := fmt.Sprintf("%d", tt.Unix())
	raw, err := s.AssetsExportStart(limit, lastAssessedUnix)

	return raw, err
}

// AssetsExportStatus will return the raw status of the assets export
func (u *Unmarshal) AssetsExportStatus(uuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.AssetsExportStatus(uuid)
	return raw, err
}

// AssetsExportGet will return the raw chunk file for the asset export uuid
func (u *Unmarshal) AssetsExportGet(uuid string, chunk string) ([]byte, error) {
	s := u.NewService()
	raw, err := s.AssetsExportGet(uuid, chunk)
	return raw, err
}

// ScansList will retrieve all scans
func (u *Unmarshal) ScansList(skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.ScansList()
	return raw, err
}

// ScanDetails will retrieve scan details for current scan
func (u *Unmarshal) ScanDetails(uuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)
	raw, err := s.ScanDetails(uuid)
	return raw, err
}

// ScansExportStart creates request with limit and lastAssessed based on Config
func (u *Unmarshal) ScansExportStart(scanid string, histid string, format string, chapters string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)

	raw, err := s.ScansExportStart(scanid, histid, format, chapters)
	return raw, err
}

// ScansExportStatus gets the status for the export-scan, returns 'ready' on done.
func (u *Unmarshal) ScansExportStatus(scanid string, fileuuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)

	raw, err := s.ScansExportStatus(scanid, fileuuid)
	return raw, err
}

// ScansExportGet downloads the file
func (u *Unmarshal) ScansExportGet(scanid string, fileuuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)

	raw, err := s.ScansExportGet(scanid, fileuuid)
	return raw, err
}

// TagValueCreate creates new tag category (if necessary) and new value for that category
func (u *Unmarshal) TagValueCreate(category string, value string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)

	raw, err := s.TagCategoryValueCreate(category, value)
	return raw, err
}

// TagBulkApply creates new tag category (if necessary) and new value for that category
func (u *Unmarshal) TagBulkApply(assetUUID []string, tagUUID []string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service(skipOnHit, writeOnReturn)

	raw, err := s.TagBulkApply(assetUUID, tagUUID)
	return raw, err
}
