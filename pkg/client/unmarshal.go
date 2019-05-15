package client

import (
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"path/filepath"
	"time"
)

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

func (u *Unmarshal) service() (s tenable.Service) {
	s = tenable.NewService(u.Config.VM.BaseURL, u.Config.VM.SecretKey, u.Config.VM.AccessKey, u.Config.VM.Log)
	s.EnableMetrics(u.Metrics)

	if u.Config.VM.CacheResponse {
		serviceCacheFolder := filepath.Join(u.Config.VM.CacheFolder, "service/")
		s.EnableCache(serviceCacheFolder, u.Config.VM.CacheKey)
	}
	s.Log = u.Config.VM.Log
	s.WriteOnReturn = true
	s.SkipOnHit = true
	return
}

func (u *Unmarshal) ScannerAgentGroups(scannerId string) ([]byte, error) {
	s := u.service()
	raw, err := s.ScannerAgentGroups(scannerId)
	return raw, err
}

func (u *Unmarshal) AgentGroup(agentId string, groupId string, scannerId string) ([]byte, error) {
	s := u.service()
	raw, err := s.AgentGroup(agentId, groupId, scannerId)

	return raw, err
}
func (u *Unmarshal) AgentUngroup(agentId string, groupId string, scannerId string) ([]byte, error) {
	s := u.service()
	raw, err := s.AgentUngroup(agentId, groupId, scannerId)

	return raw, err
}

func (u *Unmarshal) Scanners() ([]byte, error) {
	s := u.service()
	raw, err := s.ScannersList()
	return raw, err
}

func (u *Unmarshal) Agents(scannerId string, offset string, limit string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service()

	s.WriteOnReturn = writeOnReturn
	s.SkipOnHit = skipOnHit
	raw, err := s.AgentList(scannerId, offset, limit)
	return raw, err
}

func (u *Unmarshal) VulnsExportStart() ([]byte, error) {
	s := u.service()

	// Convert Human dates into Unix()
	since := u.Config.VM.AfterDate
	tt, err := time.Parse(config.DateLayout, since)
	if err != nil {
		s.Log.Errorf("failed to export-vulns start: invalid since value: %s: %s", since, err)
		return nil, err
	}
	sinceUnix := fmt.Sprintf("%d", tt.Unix())

	raw, err := s.VulnsExportStart(sinceUnix)

	return raw, err
}
func (u *Unmarshal) VulnsExportStatus(uuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service()
	s.SkipOnHit = skipOnHit
	s.WriteOnReturn = writeOnReturn
	raw, err := s.VulnsExportStatus(uuid, skipOnHit, writeOnReturn)
	return raw, err
}
func (u *Unmarshal) VulnsExportGet(uuid string, chunk string) ([]byte, error) {
	s := u.service()
	raw, err := s.VulnsExportGet(uuid, chunk)
	return raw, err
}

func (u *Unmarshal) AssetsExportStart(limit string) ([]byte, error) {
	s := u.service()

	raw, err := s.AssetsExportStart(limit)

	return raw, err
}
func (u *Unmarshal) AssetsExportStatus(uuid string, skipOnHit bool, writeOnReturn bool) ([]byte, error) {
	s := u.service()
	s.SkipOnHit = skipOnHit
	s.WriteOnReturn = writeOnReturn
	raw, err := s.AssetsExportStatus(uuid, skipOnHit, writeOnReturn)
	return raw, err
}
func (u *Unmarshal) AssetsExportGet(uuid string, chunk string) ([]byte, error) {
	s := u.service()
	raw, err := s.AssetsExportGet(uuid, chunk)
	return raw, err
}
