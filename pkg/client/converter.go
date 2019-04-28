package client

import (
	"encoding/json"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Converter does need any other objects or references
type Converter struct{}

// NewConvert returns a converter, used by the adapter
func NewConvert() (convert Converter) { return }

func (c *Converter) ToVulnExportStatus(raw []byte) (VulnExportStatus, error) {
	var src tenable.VulnExportStatus
	var dst VulnExportStatus

	err := json.Unmarshal(raw, &src)
	if err != nil {
		return dst, err
	}

	dst.Status = src.Status

	for i := range src.Chunks {
		dst.Chunks = append(dst.Chunks, string(src.Chunks[i]))
	}
	for i := range src.ChunksCancelled {
		dst.ChunksCancelled = append(dst.ChunksCancelled, string(src.ChunksCancelled[i]))
	}
	for i := range src.ChunksFailed {
		dst.ChunksFailed = append(dst.ChunksFailed, string(src.ChunksFailed[i]))
	}

	return dst, nil
}

func (c *Converter) ToScanners(raw []byte) ([]Scanner, error) {
	var src tenable.ScannerList
	var scanners []Scanner

	err := json.Unmarshal(raw, &src)
	if err != nil {
		return scanners, err
	}

	for _, s := range src.Scanners {
		var dst Scanner
		dst.License.Agents = string(s.License.Agents)
		dst.License.Scanners = string(s.License.Scanners)
		dst.License.IPS = string(s.License.IPS)
		dst.License.Type = string(s.License.Type)
		dst.Type = s.Type
		dst.Name = s.Name
		dst.UUID = s.UUID
		dst.Status = s.Status
		dst.RegistrationCode = s.RegistrationCode
		dst.Key = s.Key
		dst.LoadedPluginSet = s.LoadedPluginSet
		dst.ID = string(s.ID)
		dst.EngineVersion = s.EngineVersion
		dst.Owner = s.Owner
		dst.ScanCount = string(s.ScanCount)
		dst.Platform = s.Platform
		scanners = append(scanners, dst)
	}

	return scanners, err
}
