package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Converter translates Tenable.io raw JSONs responses into DTO objects
type Converter struct{}

// NewConvert exposes methods to take raw JSON from Unmarshal and transform it
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

// ToAssetExportStatus converst raw unmarshaled Tenable.io ExportStatus and converts to a DTO.
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

	return
}

// ToAgents converts raw Tenable Agents from Scanner to DTO ScannerAgents
func (c *Converter) ToAgents(scanner Scanner, raw []byte) ([]ScannerAgent, error) {
	var src tenable.ScannerAgent
	var scannerID = scanner.ID // We enrich the DTO Agent with the scanner ID. DNE in the Tenable.ScannerAgent object.
	var scannerUUID = scanner.UUID

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

		agent.Scanner.ID = scannerID
		agent.Scanner.UUID = scannerUUID

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

// ToAgentGroups converts raw Tenable AgentsGroups to DTO AgentGroup
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

// ToScanners converts raw Tenable AgentsGroups to DTO AgentGroup
func (c *Converter) ToScanners(raw []byte) (scanners []Scanner, err error) {
	var src tenable.ScannerList

	err = json.Unmarshal(raw, &src)
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
		if len(s.Addresses) > 0 {
			scanner.IP = s.Addresses[0]
		} else {
			scanner.IP = "Unknown"
		}

		scanners = append(scanners, scanner)
	}

	return scanners, err
}

// ToScans convert the /scans to DTO
func (c *Converter) ToScans(raw []byte) (converted []Scan, err error) {
	var src tenable.ScansList

	err = json.Unmarshal(raw, &src)
	if err != nil {
		return converted, err
	}

	for _, s := range src.Scans {
		var scan Scan
		scan.Name = s.Name
		scan.UUID = s.UUID
		scan.ScheduleUUID = s.ScheduleUUID
		scan.ScanID = s.ID.String()
		scan.Type = s.Type
		scan.StartTime = s.StartTime

		scan.RRules = s.RRules
		scan.Enabled = fmt.Sprintf("%v", s.Enabled)
		scan.CreationDate = s.CreationDate.String()
		scan.LastModifiedDate = s.LastModifiedDate.String()
		scan.Status = s.Status
		scan.Owner = s.Owner
		scan.Timezone = s.Timezone
		scan.UserPermissions = s.UserPermissions.String()

		converted = append(converted, scan)
	}

	return
}

// ToScanDetails convert the /scans to DTO
func (c *Converter) ToScanDetails(raw []byte) (converted ScanHistoryDetail, err error) {
	var src tenable.ScanDetail

	err = json.Unmarshal(raw, &src)
	if err != nil {
		return converted, err
	}

	if len(src.History) == 0 {
		return converted, nil
	}

	converted.ScanStartUnix = src.Info.Start.String()
	i, e1 := strconv.ParseInt(converted.ScanStartUnix, 10, 64)
	if e1 == nil {
		tm := time.Unix(i, 0)
		converted.ScanStart = tm.String()
	}

	converted.ScanEndUnix = src.Info.End.String()
	i, e2 := strconv.ParseInt(converted.ScanEndUnix, 10, 64)
	if e2 == nil {
		tm := time.Unix(i, 0)
		converted.ScanEnd = tm.String()
	}

	converted.TimestampUnix = src.Info.Timestamp.String()
	i, e3 := strconv.ParseInt(converted.TimestampUnix, 10, 64)
	if e3 == nil {
		tm := time.Unix(i, 0)
		converted.Timestamp = tm.String()
	}

	converted.ScanType = src.Info.ScanType
	converted.PolicyName = src.Info.PolicyName
	converted.Targets = src.Info.Targets
	converted.ScannerName = src.Info.ScannerName

	for _, h := range src.History {
		var news ScanHistoryItem
		news.HistoryID = h.HistoryID.String()
		news.UUID = h.UUID
		news.Status = h.Status
		news.LastModifiedDate = h.LastModifiedDate.String()
		news.CreationDate = h.CreationDate.String()
		converted.History = append(converted.History, news)
	}
	converted.HistoryCount = fmt.Sprintf("%d", len(converted.History))

	converted.LastModifiedDate = src.History[0].LastModifiedDate.String()
	converted.CreationDate = src.History[0].CreationDate.String()
	converted.Status = src.History[0].Status

	converted.HostCount = fmt.Sprintf("%v", len(src.Hosts))
	converted.Host = make(map[string]HostScanSummary)

	for _, h := range src.Hosts {
		var sd HostScanSummary
		sd.HostID = h.ID.String()
		sd.AssetID = h.AssetID.String()
		sd.HostnameOrIP = h.HostnameOrIP

		critsHist, _ := strconv.Atoi(converted.PluginCriticalCount)
		critsHost, _ := strconv.Atoi(string(h.SeverityCritical))
		converted.PluginCriticalCount = fmt.Sprintf("%v", critsHist+critsHost)
		highHist, _ := strconv.Atoi(converted.PluginHighCount)
		highHost, _ := strconv.Atoi(string(h.SeverityHigh))
		converted.PluginHighCount = fmt.Sprintf("%v", highHist+highHost)
		mediumHist, _ := strconv.Atoi(converted.PluginMediumCount)
		mediumHost, _ := strconv.Atoi(string(h.SeverityMedium))
		converted.PluginMediumCount = fmt.Sprintf("%v", mediumHist+mediumHost)
		lowHist, _ := strconv.Atoi(converted.PluginLowCount)
		lowHost, _ := strconv.Atoi(string(h.SeverityLow))
		converted.PluginLowCount = fmt.Sprintf("%v", lowHist+lowHost)
		infoHist, _ := strconv.Atoi(converted.PluginInfoCount)
		infoHost, _ := strconv.Atoi(string(h.SeverityInfo))
		converted.PluginInfoCount = fmt.Sprintf("%v", infoHist+infoHost)

		converted.PluginTotalCount = fmt.Sprintf("%v", infoHist+infoHost+lowHist+lowHost+mediumHist+mediumHost+highHist+highHost+critsHist+critsHost)

		sd.ScanHistoryDetail = &converted
		converted.Host[h.ID.String()] = sd
	}

	return
}

//ToScansExportStart converts Tenable.io start scan outputs
func (c *Converter) ToScansExportStart(raw []byte) (converted ScansExportStart, err error) {
	var src tenable.ScansExportStart

	err = json.Unmarshal(raw, &src)
	if err != nil {
		return converted, err
	}

	converted.FileUUID = src.FileUUID
	converted.TempToken = src.TempToken

	return converted, err
}

//ToScansExportStatus converts Tenable.io status scan outputs
func (c *Converter) ToScansExportStatus(fileuuid string, raw []byte) (converted ScansExportStatus, err error) {
	var src tenable.ScansExportStatus

	err = json.Unmarshal(raw, &src)
	if err != nil {
		return converted, err
	}

	converted.Status = strings.ToUpper(src.Status)
	converted.FileUUID = fileuuid

	return converted, err
}
