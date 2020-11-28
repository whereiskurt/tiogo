package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/whereiskurt/tiogo/pkg/cache"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Adapter is used to call ACME services and convert them to Gopher/Things in Go structures we like.
type Adapter struct {
	Config    *config.Config
	Metrics   *metrics.Metrics
	Unmarshal Unmarshal
	Filter    *Filter
	Convert   Converter
	Worker    *sync.WaitGroup
	DiskCache *cache.Disk
}

// NewAdapter manages calls the remote services, converts the results and manages a memory/disk cache.
func NewAdapter(config *config.Config, metrics *metrics.Metrics) (a *Adapter) {
	a = new(Adapter)
	a.Config = config
	a.Metrics = metrics
	a.Worker = new(sync.WaitGroup)
	a.Unmarshal = NewUnmarshal(config, metrics)
	a.Filter = NewFilter(config)
	a.Convert = NewConvert()
	if a.Config.VM.CacheResponse {
		a.DiskCache = cache.NewDisk(a.Config.VM.CacheFolder, a.Config.VM.CacheKey, a.Config.VM.CacheKey != "")
	}

	return
}

// Scanners returns all the Tenable.io scanners
func (a *Adapter) Scanners(skipOnHit bool, writeOnReturn bool) ([]Scanner, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)
	var scanners []Scanner
	raw, err := u.Scanners(skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scanners list: %v", err)
		return scanners, err
	}

	convert := NewConvert()
	scanners, err = convert.ToScanners(raw)

	return scanners, err
}

// MagicAgentScanner is the Tenable.io number for the scanner that has agents attached
var MagicAgentScanner = "00000000-0000-0000-0000-00000000000000000000000000001"

// Agents uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) Agents(skipOnHit bool, writeOnReturn bool) ([]ScannerAgent, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AgentsList, metrics.Methods.Service.Get)

	AgentsPerRequest, _ := strconv.Atoi(a.Config.VM.ExportLimit)

	scanners, err := a.Scanners(skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scanners list for agents list: %v", err)
		return nil, err
	}

	u := NewUnmarshal(a.Config, a.Metrics)
	convert := NewConvert()

	var agents []ScannerAgent

	for i := range scanners {
		if scanners[i].UUID != MagicAgentScanner {
			continue
		}

		// Using the MagicalScanner ;-)
		totalAgents, err := strconv.Atoi(scanners[i].License.AgentsUsed)
		if err != nil {
			log.Fatalf("error: invalid agents_used:%s:%s", scanners[i].License.AgentsUsed, err)
		}

		offset, loops := 0, 0
		for {
			// The API doc says to use ID but in practice WebGUI uses UUID...
			// NOTE: ANY VALUE WILL WORK!!! LITERALLY!
			id := scanners[i].ID

			raw, err := u.Agents(id, fmt.Sprintf("%d", offset), fmt.Sprintf("%d", AgentsPerRequest), skipOnHit, writeOnReturn)
			if err != nil {
				a.Config.VM.Log.Errorf("error: failed to get the agents: uuid: %s: %v", id, err)
			}

			agents, err := convert.ToAgents(scanners[i], raw)
			if err != nil {
				a.Config.VM.Log.Errorf("error: failed to convert agents: uuid: %s: %v", id, err)
			}

			scanners[i].Agents = append(scanners[i].Agents, agents...)

			if AgentsPerRequest*(loops+1) >= totalAgents {
				break
			}

			loops = loops + 1
			offset = loops * AgentsPerRequest
		}
		agents = scanners[i].Agents
		break
	}

	return agents, err
}

// AgentGroups uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) AgentGroups(skipOnHit bool, writeOnReturn bool) ([]AgentGroup, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AgentGroups, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	scanners, err := a.Scanners(skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the agent scanners list : %v", err)
		return nil, err
	}

	var agentGroups []AgentGroup
	for i := range scanners {
		// We only want the MagicAgentScanner ... can probably replay w/ 00000-000 hardcoded
		// If this scanner in our collection isnt the Agent Scanner
		if scanners[i].UUID != MagicAgentScanner {
			continue // Try the next one!
		}

		// Use thie ID to get all our Agent Groups
		id := scanners[i].ID
		raw, err := u.ScannerAgentGroups(id, skipOnHit, writeOnReturn)

		if err != nil {
			a.Config.VM.Log.Errorf("error: failed to get the scanners agent groups: %v", err)
			return agentGroups, err
		}

		convert := NewConvert()
		agentGroups, err = convert.ToAgentGroups(raw)

		break // DONE!
	}

	return agentGroups, err
}

// AgentAssignGroup uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) AgentAssignGroup(agentID string, groupID string, scannerID string) error {
	a.Metrics.ClientInc(metrics.EndPoints.AgentsGroup, metrics.Methods.Service.Update)
	u := NewUnmarshal(a.Config, a.Metrics)

	_, err := u.AgentGroup(agentID, groupID, scannerID)

	return err
}

// AgentUnassignGroup uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) AgentUnassignGroup(agentID string, groupID string, scannerID string) error {
	a.Metrics.ClientInc(metrics.EndPoints.AgentsUngroup, metrics.Methods.Service.Update)
	u := NewUnmarshal(a.Config, a.Metrics)

	_, err := u.AgentUngroup(agentID, groupID, scannerID)

	return err
}

// ExportVulnsStart uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportVulnsStart() (string, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStart, metrics.Methods.Service.Update)

	u := NewUnmarshal(a.Config, a.Metrics)

	raw, err := u.VulnsExportStart()
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return "", err
	}

	var export tenable.VulnExportStart
	err = json.Unmarshal(raw, &export)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to unmarshal response from start export-vulns: %v", err)
		return "", err
	}

	return export.UUID, nil
}

// ExportVulnsStatus uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportVulnsStatus(uuid string, skipOnHit bool, writeOnReturn bool) (VulnExportStatus, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)

	var status VulnExportStatus
	raw, err := u.VulnsExportStatus(uuid, skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return status, err
	}

	convert := NewConvert()
	status, err = convert.ToVulnExportStatus(raw)

	return status, err
}

// ExportVulnsGet uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportVulnsGet(uuid string, chunks string) error {
	ep := tenable.EndPoints.VulnsExportStatus

	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportGet, metrics.Methods.Service.Get)

	if uuid == "" {
		var err error
		uuid, err = a.LastCachedExportVulnUUID()
		if err != nil {
			return err
		}
	}

	chunks, err := a.ExportCachedChunks(uuid, chunks, ep)
	if err != nil {
		return err
	}

	u := NewUnmarshal(a.Config, a.Metrics)

	cc := strings.Split(chunks, ",")
	a.Config.VM.Log.Infof("ChunkSize='%d' for uuid='%s", len(cc), uuid)
	for _, chunk := range cc {
		raw, err := u.VulnsExportGet(uuid, chunk)
		if err != nil {
			a.Config.VM.Log.Errorf("error: failed to get the export-vulns: %v", err)
			return err
		}

		a.Config.VM.Log.Infof("Downloaded chunk '%s', file size '%d' bytes", chunk, len(raw))
	}

	return nil
}

// ExportVulnsQuery uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportVulnsQuery(uuid string, chunks string, jqex string) error {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportQuery, metrics.Methods.Service.Get)
	ep := tenable.EndPoints.VulnsExportStatus

	if uuid == "" {
		return errors.New("error: uuid not specified")
	}

	chunks, err := a.ExportCachedChunks(uuid, chunks, ep)
	if err != nil {
		return err
	}

	for _, chunk := range strings.Split(chunks, ",") {
		ep := tenable.EndPoints.VulnsExportGet
		p := map[string]string{"ExportUUID": uuid, "ChunkID": chunk}
		filename, err := a.CachedFilename(ep, p)

		if err != nil {
			a.Config.VM.Log.Errorf("error: reading chunk file '%s' ", filename)
			return err
		}

		a.Config.VM.Log.Debugf("read chunk file: %s", filename)

		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error: cannot read cached file: '%s: %v", filename, err)
		}

		filter := a.JSONQuery(bb, jqex)
		if len(filter) > 1 {
			fmt.Println(string(filter))
		}
	}

	return nil
}

// ExportAssetsStart uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportAssetsStart() (string, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AssetsExportStart, metrics.Methods.Service.Update)

	u := NewUnmarshal(a.Config, a.Metrics)

	raw, err := u.AssetsExportStart()
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the export-assets: %v", err)
		return "", err
	}

	var export tenable.AssetExportStart
	err = json.Unmarshal(raw, &export)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to unmarshal response from start export-assets: %v", err)
		return "", err
	}

	return export.UUID, nil
}

// ExportAssetsStatus uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportAssetsStatus(uuid string, skipOnHit bool, writeOnReturn bool) (AssetExportStatus, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AssetsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)

	var status AssetExportStatus
	raw, err := u.AssetsExportStatus(uuid, skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the export-assets: %v", err)
		return status, err
	}

	convert := NewConvert()
	status, err = convert.ToAssetExportStatus(raw)

	return status, err
}

// ExportAssetsGet uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportAssetsGet(uuid string, chunks string) error {
	a.Metrics.ClientInc(metrics.EndPoints.AssetsExportGet, metrics.Methods.Service.Get)
	ep := tenable.EndPoints.AssetsExportStatus

	if uuid == "" {
		var err error
		uuid, err = a.LastCachedExportAssetUUID()
		if err != nil {
			return err
		}
	}

	chunks, err := a.ExportCachedChunks(uuid, chunks, ep)
	if err != nil {
		return err
	}

	u := NewUnmarshal(a.Config, a.Metrics)

	cc := strings.Split(chunks, ",")
	a.Config.VM.Log.Infof("ChunkSize='%d' for uuid='%s", len(cc), uuid)
	for _, chunk := range cc {
		raw, err := u.AssetsExportGet(uuid, chunk)
		if err != nil {
			a.Config.VM.Log.Errorf("error: failed to get the export-assets: %v", err)
			return err
		}

		a.Config.VM.Log.Infof("Downloaded chunk '%s', file size '%d' bytes", chunk, len(raw))
	}

	return nil
}

// ExportAssetsQuery uses an Unmarshaler and Converter to return DTO or error
func (a *Adapter) ExportAssetsQuery(uuid string, chunks string, jqex string) error {
	a.Metrics.ClientInc(metrics.EndPoints.AssetsExportQuery, metrics.Methods.Service.Get)
	ep := tenable.EndPoints.AssetsExportStatus

	if uuid == "" {
		return errors.New("error: uuid not specified")
	}

	chunks, err := a.ExportCachedChunks(uuid, chunks, ep)
	if err != nil {
		return err
	}

	for _, chunk := range strings.Split(chunks, ",") {
		ep := tenable.EndPoints.AssetsExportGet
		p := map[string]string{"ExportUUID": uuid, "ChunkID": chunk}
		filename, err := a.CachedFilename(ep, p)

		if err != nil {
			a.Config.VM.Log.Errorf("error: reading chunk file '%s' ", filename)
			return err
		}

		a.Config.VM.Log.Debugf("read chunk file: %s", filename)

		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error: cannot read cached file: '%s: %v", filename, err)
		}

		filter := a.JSONQuery(bb, jqex)
		if len(filter) > 1 {
			fmt.Println(string(filter))
		}
	}

	return nil
}

// ExportCachedChunks reads the chunks for uuid from the cached file
func (a *Adapter) ExportCachedChunks(uuid string, chunks string, ep tenable.EndPointType) (string, error) {
	if !(chunks == "ALL" || chunks == "") {
		return chunks, nil
	}
	p := map[string]string{"ExportUUID": uuid}
	filename, err := a.CachedFilename(ep, p)
	if err != nil {
		a.Config.VM.Log.Errorf("error: reading cached 'status' file '%s' ", filename)
		return "", err
	}

	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error: cannot read cached file: '%s: %v", filename, err)
	}

	var status tenable.VulnExportStatus
	err = json.Unmarshal(bb, &status)
	if err != nil {
		a.Config.VM.Log.Errorf("error: can get unmarshal '%s': %v", filename, err)
		return "", err
	}

	var cc []string
	for i := range status.Chunks {
		cc = append(cc, fmt.Sprintf("%s", string(status.Chunks[i])))
	}
	c := strings.Join(cc, ",")
	return c, nil
}

// LastCachedExportVulnUUID reads the service export vulns cache for entries, and returns the first one as exportUUID
func (a *Adapter) LastCachedExportVulnUUID() (string, error) {
	folder := filepath.Join(a.Config.VM.CacheFolder, "service", "export", "vulns")
	entries, err := a.DirEntries(folder)
	if err != nil || len(entries) == 0 {
		return "", err
	}
	a.Config.VM.Log.Debugf("Returning first entry as uuid: %s", entries[0])
	return entries[0], nil
}

// LastCachedExportAssetUUID reads the service export assets cache for entries, and returns the first one as exportUUID
func (a *Adapter) LastCachedExportAssetUUID() (string, error) {
	folder := filepath.Join(a.Config.VM.CacheFolder, "service", "export", "assets")
	entries, err := a.DirEntries(folder)
	if err != nil || len(entries) == 0 {
		return "", err
	}
	a.Config.VM.Log.Debugf("Returning first entry as uuid: %s", entries[0])
	return entries[0], nil
}

// CachedFilename will output the filename on disk for that end-point requested with a parameter map p
func (a *Adapter) CachedFilename(endpoint tenable.EndPointType, p map[string]string) (string, error) {

	filename, err := tenable.ToCacheFilename(endpoint, p)
	if err != nil {
		return "", fmt.Errorf("error: can't get chunk filename for '%s'", filename)
	}
	filename = filepath.Join(a.Config.VM.CacheFolder, "service", filename)
	if _, stat := os.Stat(filename); os.IsNotExist(stat) {
		return "", fmt.Errorf(fmt.Sprintf("Cannot read cached file: '%s", filename))
	}

	return filename, nil
}

// JSONPretty will look for 'jq' to pretty the json input
func (a *Adapter) JSONPretty(json []byte) []byte {
	return a.JSONQuery(json, ".")
}

// UnpackJQExec extracts the jq executable packed in templates.go
func (a *Adapter) UnpackJQExec() (string, error) {
	tempFile, err := ioutil.TempFile("", "jq.")
	tempFile.Close()

	if err != nil {
		log.Fatal(err)
	}
	jqint := ""
	jqexe := tempFile.Name()
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "386":
			jqint = "embed/jq/linux/jq-linux32"
		case "amd64":
			jqint = "embed/jq/linux/jq-linux64"
		}
	case "osx":
		switch runtime.GOARCH {
		case "amd64":
			jqint = "embed/jq/osx/jq-osx-amd64"
		}
	case "windows":
		defer os.Remove(tempFile.Name())
		// NOTE: We cannot execute a non-exec
		jqexe = jqexe + ".exe"
		switch runtime.GOARCH {
		case "386":
			jqint = "embed/jq/windows/jq-win32.exe"
		case "amd64":
			jqint = "embed/jq/windows/jq-win64.exe"
		}
	}

	if jqint == "" {
		err = errors.New("error: jq not found in path, and cannot self-extract")
		return "", err
	}

	a.Config.VM.Log.Debugf("Creating temporary file for jq executable: %s from %s", jqexe, jqint)
	file, err := config.BinaryEmbedFolder.Open(jqint)
	if err != nil {
		log.Error(err)
		return "", err
	}

	bb, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return "", err
	}

	err = ioutil.WriteFile(jqexe, bb, 0777)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return jqexe, nil
}

// JSONQuery will pipe bytes through jq and return results.
func (a *Adapter) JSONQuery(json []byte, jqex string) []byte {
	log := a.Config.VM.Log
	if a.Config.VM.JQExec == "" {
		jq, err := exec.LookPath("jq")
		if err != nil {
			log.Infof("'jq' exec not found in path: extracting from self")
			jq, err = a.UnpackJQExec()
			if err != nil {
				log.Errorf("error: cannot unpack jq from self and not in path.")
				return []byte("")
			}
		}
		a.Config.VM.JQExec = jq
		log.Infof("using jq from path '%s'", jq)
	}

	jqexec := a.Config.VM.JQExec

	var stdout bytes.Buffer
	// var cmds = []string{"-c", "-r", jqex}
	var cmds = []string{jqex}

	cmd := exec.Command(jqexec, cmds...)
	cmd.Stdin = strings.NewReader(string(json))
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Warnf("couldn't parse 'jq' command: %s: %+v: %v", jqex, cmd, err)
		return []byte("")
	}

	return []byte(stdout.String())
}

// DirEntries returns an array of files in a folder or error
func (a *Adapter) DirEntries(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ModTime().After(list[j].ModTime()) })

	var files []string
	for _, file := range list {
		files = append(files, file.Name())
	}

	return files, nil
}

// Scans will list all the scans matching --name or --regex
func (a *Adapter) Scans(skipOnHit bool, writeOnReturn bool) ([]Scan, error) {
	a.Metrics.ClientInc(metrics.EndPoints.ScansList, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)
	var scans []Scan
	raw, err := u.ScansList(skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scan list: %v", err)
		return scans, err
	}

	convert := NewConvert()
	scans, err = convert.ToScans(raw)

	return scans, err
}

// AuditLogV1 calls Tenable.io to get Audit Log
func (a *Adapter) AuditLogV1(skipOnHit bool, writeOnReturn bool) ([]AuditLogEvent, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AuditLogV1, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)
	var logs []AuditLogEvent
	raw, err := u.AuditLogV1(skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scan list: %v", err)
		return logs, err
	}

	convert := NewConvert()
	logs, err = convert.ToAuditLogV1Event(raw)

	return logs, err
}

// ScanDetails will list all the scans matching --name or --regex
func (a *Adapter) ScanDetails(s *Scan, skipOnHit bool, writeOnReturn bool) (details ScanHistoryDetail, err error) {
	a.Metrics.ClientInc(metrics.EndPoints.ScansList, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)
	raw, err := u.ScanDetails(s.ScanID, skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scan list: %v", err)
		return details, err
	}

	groups, err := a.AgentGroups(true, true)
	convert := NewConvert()
	details, err = convert.ToScanDetails(raw, groups)

	details.Scan = s
	return details, err
}

// ScansExportStart extract scan for histuuid
func (a *Adapter) ScansExportStart(s *Scan, histid string, format string, chapters string, skipOnHit bool, writeOnReturn bool) (ScansExportStart, error) {
	var export ScansExportStart

	a.Metrics.ClientInc(metrics.EndPoints.ScansExportStart, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	raw, err := u.ScansExportStart(s.ScanID, histid, format, chapters, skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scan export: %v", err)
		return export, err
	}

	convert := NewConvert()
	export, err = convert.ToScansExportStart(format, raw)

	return export, err
}

// ScansExportStatus gets the status of the export-scan that was started
func (a *Adapter) ScansExportStatus(s *Scan, histid string, format string, chapters string, skipOnHit bool, writeOnReturn bool) (ScansExportStatus, error) {
	a.Metrics.ClientInc(metrics.EndPoints.ScansExportStatus, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	started, err := a.ScansExportStart(s, histid, format, chapters, skipOnHit, writeOnReturn)
	if err != nil {
		return ScansExportStatus{}, err
	}

	var export ScansExportStatus
	var fileuuid string
	fileuuid = started.FileUUID
	raw, err := u.ScansExportStatus(s.ScanID, fileuuid, skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scan export: %v", err)
		return export, err
	}

	convert := NewConvert()
	export, err = convert.ToScansExportStatus(fileuuid, raw)

	return export, err
}

// ScansExportGet downloads the prepared
func (a *Adapter) ScansExportGet(s *Scan, histid string, format string, chapters string, skipOnHit bool, writeOnReturn bool) (ScansExportGet, error) {
	a.Metrics.ClientInc(metrics.EndPoints.ScansExportGet, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	status, err := a.ScansExportStatus(s, histid, format, chapters, skipOnHit, writeOnReturn)
	if err != nil {
		return ScansExportGet{}, err
	}
	if strings.ToUpper(status.Status) != "READY" {
		return ScansExportGet{}, fmt.Errorf("error: can't download until the status is 'ready'")
	}
	fileuuid := status.FileUUID
	raw, err := u.ScansExportGet(s.ScanID, fileuuid, skipOnHit, writeOnReturn)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to get the scan export: %v", err)
		return ScansExportGet{}, err
	}

	var export ScansExportGet

	convert := NewConvert()
	// DTO from Tenable.io
	if format == "nessus" {
		export, err = convert.ToScansExportGet(&raw)
	}

	// Store the scan indentifiers with export
	export.ScanID = s.ScanID
	export.ScheduleUUID = s.ScheduleUUID
	export.HistoryID = histid
	export.SourceFile.FileUUID = fileuuid
	filename, _ := tenable.ToCacheFilename(tenable.EndPoints.ScansExportGet, map[string]string{"ScanID": s.ScanID, "FileUUID": fileuuid})
	filename = filepath.Join(a.Config.VM.CacheFolder, "service", filename)
	export.SourceFile.CachedFileName = filename

	return export, err
}

// ScansExportLargeGet downloads the prepared
func (a *Adapter) ScansExportLargeGet(s *Scan, histid string, format string, chapters string) (string, error) {
	a.Metrics.ClientInc(metrics.EndPoints.ScansExportGet, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	status, err := a.ScansExportStatus(s, histid, format, chapters, true, true)
	if err != nil {
		return "", err
	}
	if strings.ToUpper(status.Status) != "READY" {
		return "", fmt.Errorf("error: can't download until the status is 'ready'")
	}

	filename, err := u.ScansExportStream(s.ScanID, status.FileUUID, true, true)
	if err != nil {
		a.Config.VM.Log.Errorf("error: failed to stream the scan export: %v", err)
		return "", err
	}

	return filename, nil
}

// TagValueCreate will make a new category:value, and return a TagValue UUID
func (a *Adapter) TagValueCreate(category string, value string, skipOnHit bool, writeOnReturn bool) (TagValue, error) {
	a.Metrics.ClientInc(metrics.EndPoints.TagValueCreate, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	var tag TagValue
	raw, err := u.TagValueCreate(category, value, skipOnHit, writeOnReturn)
	if err != nil {
		// This happens when we've already created this category:value pair
		a.Config.VM.Log.Debugf("Lookup the TAG value because it already exists: %v...", string(raw))
		if strings.Contains(string(raw), "Duplicate tag value") {

			// TODO: Make call to get TagValues, and return the TagValueUUID
			//raw, err = u.GetTagValue(category, value, skipOnHit, writeOnReturn)
			//1) Get category:value uuid
			//2) Populate tag with values
		}

		return tag, err
	}
	convert := NewConvert()

	tag, err = convert.ToTagValue(&raw)

	return tag, err
}

// TagBulkApply will chunk [] into 1000 units and submit.
func (a *Adapter) TagBulkApply(assetUUID []string, tagUUID []string) (jobID string, err error) {
	a.Metrics.ClientInc(metrics.EndPoints.TagValueCreate, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	raw, err := u.TagBulkApply(assetUUID, tagUUID, false, false)
	if err != nil {
		log.Fatalf("error: bbbbfailed to TagBulkApply: %v:%v", raw, err)
		return "", nil
	}

	convert := NewConvert()
	jobID, err = convert.ToTagTagBulkJobID(&raw)

	return jobID, err
}
