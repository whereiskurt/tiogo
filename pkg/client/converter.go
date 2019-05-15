package client

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"strconv"
	"time"
)

// Converter does not need any other objects or references
type Converter struct{}

// NewConvert returns a converter, used by the adapter
func NewConvert() (convert Converter) { return }

// ToVulnExportStatus takes a raw byte array of JSON from Tenable.io,
// marshal's into typed structure, and then we convert to our data objects.
func (c *Converter) ToVulnExportStatus(raw []byte) (converted VulnExportStatus, err error) {
	var tenableStatus tenable.VulnExportStatus

	err = json.Unmarshal(raw, &tenableStatus)
	if err != nil {
		return
	}

	// Convert Tenable chunks to tiogo structures
	for i := range tenableStatus.Chunks {
		converted.Chunks = append(converted.Chunks, string(tenableStatus.Chunks[i]))
	}
	for i := range tenableStatus.ChunksCancelled {
		converted.ChunksCancelled = append(converted.ChunksCancelled, string(tenableStatus.ChunksCancelled[i]))
	}
	for i := range tenableStatus.ChunksFailed {
		converted.ChunksFailed = append(converted.ChunksFailed, string(tenableStatus.ChunksFailed[i]))
	}
	converted.Status = tenableStatus.Status

	return converted, nil
}

func (c *Converter) ToAssetExportStatus(raw []byte) (converted AssetExportStatus, err error) {
	var tenableStatus tenable.AssetExportStatus

	err = json.Unmarshal(raw, &tenableStatus)
	if err != nil {
		return
	}

	// Convert Tenable chunks to tiogo structures
	for i := range tenableStatus.Chunks {
		converted.Chunks = append(converted.Chunks, string(tenableStatus.Chunks[i]))
	}
	for i := range tenableStatus.ChunksCancelled {
		converted.ChunksCancelled = append(converted.ChunksCancelled, string(tenableStatus.ChunksCancelled[i]))
	}
	for i := range tenableStatus.ChunksFailed {
		converted.ChunksFailed = append(converted.ChunksFailed, string(tenableStatus.ChunksFailed[i]))
	}
	converted.Status = tenableStatus.Status

	return converted, nil
}

func (c *Converter) ToAgents(scanner Scanner, raw []byte) ([]ScannerAgent, error) {
	var src tenable.ScannerAgent

	log.Debug(fmt.Sprintf("%s", raw))

	err := json.Unmarshal(raw, &src)
	if err != nil {
		return nil, err
	}

	var agents []ScannerAgent
	for _, agent := range src.Agents {
		var tgt ScannerAgent
		tgt.ID = string(agent.ID)
		tgt.Name = agent.Name
		tgt.UUID = agent.UUID
		tgt.Platform = agent.Platform
		tgt.Status = agent.Status
		tgt.CoreBuild = agent.CoreBuild
		tgt.CoreVersion = agent.CoreVersion
		tgt.Feed = agent.Feed
		tgt.Distro = agent.Distro
		tgt.IP = agent.IP

		tgt.Scanner.ID = scanner.ID
		tgt.Scanner.UUID = scanner.UUID

		tgt.Groups = make(map[string]AgentGroup)
		for _, group := range agent.Groups {
			g := AgentGroup{ID: string(group.ID), Name: group.Name}
			tgt.Groups[g.Name] = g
		}

		lastConnect, err := strconv.ParseInt(string(agent.LastConnect), 10, 64)
		if err == nil {
			tgt.LastConnect = time.Unix(lastConnect, 0)
		}

		linkedOn, err := strconv.ParseInt(string(agent.LinkedOn), 10, 64)
		if err == nil {
			tgt.LinkedOn = time.Unix(linkedOn, 0)
		}

		lastScanned, err := strconv.ParseInt(string(agent.LastScanned), 10, 64)
		if err == nil {
			tgt.LastScanned = time.Unix(lastScanned, 0)
		}

		agents = append(agents, tgt)
	}

	return agents, err
}

func (c *Converter) ToAgentGroups(raw []byte) ([]AgentGroup, error) {
	var src tenable.ScannerAgentGroups
	var groups []AgentGroup

	err := json.Unmarshal(raw, &src)
	if err != nil {
		return groups, err
	}

	for _, g := range src.Groups {
		var group AgentGroup

		group.UUID = g.UUID
		group.ID = string(g.ID)
		group.Name = g.Name
		group.AgentsCount = string(g.AgentCount)

		groups = append(groups, group)
	}

	return groups, err
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
		dst.License.ScannersUsed = string(s.License.ScannersUsed)
		dst.License.AgentsUsed = string(s.License.AgentsUsed)
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
