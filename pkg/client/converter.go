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
	for _, a := range src.Agents {
		var agent ScannerAgent
		agent.ID = string(a.ID)
		agent.Name = a.Name
		agent.UUID = a.UUID
		agent.Platform = a.Platform
		agent.Status = a.Status
		agent.CoreBuild = a.CoreBuild
		agent.CoreVersion = a.CoreVersion
		agent.Feed = a.Feed
		agent.Distro = a.Distro
		agent.IP = a.IP

		agent.Scanner.ID = scanner.ID
		agent.Scanner.UUID = scanner.UUID

		agent.Groups = make(map[string]AgentGroup)
		for _, g := range a.Groups {
			group := AgentGroup{ID: string(g.ID), Name: g.Name}
			agent.Groups[group.Name] = group
		}

		lastConnect, err := strconv.ParseInt(string(a.LastConnect), 10, 64)
		if err == nil {
			agent.LastConnect = time.Unix(lastConnect, 0)
		}

		linkedOn, err := strconv.ParseInt(string(a.LinkedOn), 10, 64)
		if err == nil {
			agent.LinkedOn = time.Unix(linkedOn, 0)
		}

		lastScanned, err := strconv.ParseInt(string(a.LastScanned), 10, 64)
		if err == nil {
			agent.LastScanned = time.Unix(lastScanned, 0)
		}

		agents = append(agents, agent)
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
		var scanner Scanner
		scanner.License.Agents = string(s.License.Agents)
		scanner.License.Scanners = string(s.License.Scanners)
		scanner.License.IPS = string(s.License.IPS)
		scanner.License.Type = string(s.License.Type)
		scanner.License.ScannersUsed = string(s.License.ScannersUsed)
		scanner.License.AgentsUsed = string(s.License.AgentsUsed)
		scanner.Type = s.Type
		scanner.Name = s.Name
		scanner.UUID = s.UUID
		scanner.Status = s.Status
		scanner.RegistrationCode = s.RegistrationCode
		scanner.Key = s.Key
		scanner.LoadedPluginSet = s.LoadedPluginSet
		scanner.ID = string(s.ID)
		scanner.EngineVersion = s.EngineVersion
		scanner.Owner = s.Owner
		scanner.ScanCount = string(s.ScanCount)
		scanner.Platform = s.Platform
		scanners = append(scanners, scanner)
	}

	return scanners, err
}
