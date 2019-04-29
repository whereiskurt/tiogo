package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/whereiskurt/tiogo/pkg/cache"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
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

func (a *Adapter) Scanners() ([]Scanner, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)
	var scanners []Scanner
	raw, err := u.Scanners()
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the scanners list: %v", err)
		return scanners, err
	}

	convert := NewConvert()
	scanners, err = convert.ToScanners(raw)

	return scanners, err
}

var MagicAgentScanner = "00000000-0000-0000-0000-00000000000000000000000000001"

func (a *Adapter) Agents() ([]ScannerAgent, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AgentsList, metrics.Methods.Service.Get)

	scanners, err := a.Scanners()
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the scanners list for agents list: %v", err)
		return nil, err
	}

	u := NewUnmarshal(a.Config, a.Metrics)
	convert := NewConvert()

	var agents []ScannerAgent

	limit := 5000
	for i := range scanners {
		if scanners[i].UUID != MagicAgentScanner {
			continue
		}

		// Using the MagicalScanner ;-)
		totalAgents, err := strconv.Atoi(scanners[i].License.AgentsUsed)
		if err != nil {
			log.Fatalf("error: invalid agents_used:%s:%s", scanners[i].License.AgentsUsed,err)
		}

		offset, loops := 0, 0
		for {
			// The API doc says to use ID but in practice WebGUI uses UUID...
			// NOTE: ANY VALUE WILL WORK!!! LITERALLY!
			uuid := scanners[i].ID

			raw, err := u.Agents(uuid, fmt.Sprintf("%d", offset), fmt.Sprintf("%d", limit))
			if err != nil {
				a.Config.Log.Errorf("error: failed to get the agents: uuid: %s: %v", uuid, err)
			}

			agents, err := convert.ToAgents(raw)
			if err != nil {
				a.Config.Log.Errorf("error: failed to convert agents: uuid: %s: %v", uuid, err)
			}

			scanners[i].Agents = append(scanners[i].Agents, agents...)

			if limit * (loops + 1) >= totalAgents { break }

			loops = loops + 1
			offset = loops * limit
		}
		agents = scanners[i].Agents
		break
	}

	return agents, err
}

func (a *Adapter) AgentGroups() ([]AgentGroup, error) {
	a.Metrics.ClientInc(metrics.EndPoints.AgentGroups, metrics.Methods.Service.Get)
	u := NewUnmarshal(a.Config, a.Metrics)

	scanners, err := a.Scanners()
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the agent scanners list : %v", err)
		return nil, err
	}

	var agentGroups []AgentGroup
	for i := range scanners {
		if scanners[i].UUID != MagicAgentScanner {
			continue
		}

		id := scanners[i].ID
		raw, err := u.AgentGroups(id)

		if err != nil {
			a.Config.Log.Errorf("error: failed to get the scanners agent groups: %v", err)
			return agentGroups, err
		}

		convert := NewConvert()
		agentGroups, err = convert.ToAgentGroups(raw)

		break
	}

	return agentGroups, err
}

func (a *Adapter) ExportVulnsStart() (string, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStart, metrics.Methods.Service.Update)

	u := NewUnmarshal(a.Config, a.Metrics)

	raw, err := u.VulnsExportStart()
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return "", err
	}

	var export tenable.VulnExportStart
	err = json.Unmarshal(raw, &export)
	if err != nil {
		a.Config.Log.Errorf("error: failed to unmarshal response from start export-vulns: %v", err)
		return "", err
	}

	return export.UUID, nil
}
func (a *Adapter) ExportVulnsStatus(uuid string, skipOnHit bool, writeOnReturn bool) (VulnExportStatus, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)

	var status VulnExportStatus
	raw, err := u.VulnsExportStatus(uuid, skipOnHit,writeOnReturn)
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return status, err
	}

	convert := NewConvert()
	status, err = convert.ToVulnExportStatus(raw)

	return status, err
}
func (a *Adapter) ExportVulnsGet(uuid string, chunks string) error {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportGet, metrics.Methods.Service.Get)

	if uuid == "" {
		var err error
		uuid, err = a.ExportCachedUUID()
		if err != nil {
			return err
		}
	}

	chunks, err := a.ExportCachedChunks(uuid, chunks)
	if err != nil {
		return err
	}

	u := NewUnmarshal(a.Config, a.Metrics)

	cc := strings.Split(chunks, ",")
	a.Config.Log.Infof("ChunkSize='%d' for uuid='%s", len(cc), uuid)
	for _, chunk := range cc {
		raw, err := u.VulnsExportGet(uuid, chunk)
		if err != nil {
			a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
			return err
		}

		a.Config.Log.Infof("Downloaded chunk '%s', file size '%d' bytes", chunk, len(raw))
	}

	return nil
}
func (a *Adapter) ExportVulnsQuery(uuid string, chunks string, jqex string) error {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportQuery, metrics.Methods.Service.Get)

	if uuid == "" {
		return errors.New("error: uuid not specified")
	}

	chunks, err := a.ExportCachedChunks(uuid, chunks)
	if err != nil {
		return err
	}

	for _, chunk := range strings.Split(chunks, ",") {
		ep := tenable.EndPoints.VulnsExportGet
		p := map[string]string{"ExportUUID": uuid, "ChunkID": chunk}
		filename, err := a.CachedFilename(ep, p)

		if err != nil {
			a.Config.Log.Errorf("error: reading chunk file '%s' ", filename)
			return err
		}

		a.Config.Log.Debugf("read chunk file: %s", filename)

		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.New(fmt.Sprintf("error: cannot read cached file: '%s: %v", filename, err))
		}

		filter := a.JSONQuery(bb, jqex)
		if len(filter) > 1 {
			fmt.Println(string(filter))
		}
	}

	return nil
}

// ExportCachedChunks reads hte chunks for uuid from the cached file
func (a *Adapter) ExportCachedChunks(uuid string, chunks string) (string, error) {
	if !(chunks == "ALL" || chunks == "") {
		return chunks, nil
	}
	ep := tenable.EndPoints.VulnsExportStatus
	p := map[string]string{"ExportUUID": uuid}
	filename, err := a.CachedFilename(ep, p)
	if err != nil {
		a.Config.Log.Errorf("error: reading cached 'status' file '%s' ", filename)
		return "", err
	}

	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error: cannot read cached file: '%s: %v", filename, err))
	}

	var status tenable.VulnExportStatus
	err = json.Unmarshal(bb, &status)
	if err != nil {
		a.Config.Log.Errorf("error: can get unmarshal '%s': %v", filename, err)
		return "", err
	}

	var cc []string
	for i := range status.Chunks {
		cc = append(cc, fmt.Sprintf("%s", string(status.Chunks[i])))
	}
	c := strings.Join(cc, ",")
	return c, nil
}

// ExportCachedUUID reads the service export vulns cache for entries, and returns the first one as exportUUID
func (a *Adapter) ExportCachedUUID() (string, error) {
	folder := filepath.Join(a.Config.VM.CacheFolder, "service", "export", "vulns")
	entries, err := a.DirEntries(folder)
	if err != nil || len(entries) == 0 {
		return "", err
	}
	a.Config.Log.Debugf("Returning first entry as uuid: %s", entries[0])
	return entries[0], nil
}

// CachedFilename will output the filename on disk for that end-point requested with a parameter map p
func (a *Adapter) CachedFilename(endpoint tenable.EndPointType, p map[string]string) (string, error) {

	filename, err := tenable.ToCacheFilename(endpoint, p)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error: can't get chunk filename for '%s'", filename))
	}
	filename = filepath.Join(a.Config.VM.CacheFolder, "service", filename)
	if _, stat := os.Stat(filename); os.IsNotExist(stat) {
		return "", errors.New(fmt.Sprintf("Cannot read cached file: '%s", filename))
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
			jqint = "binary/jq/linux/jq-linux32"
		case "amd64":
			jqint = "binary/jq/linux/jq-linux64"
		}
	case "osx":
		switch runtime.GOARCH {
		case "amd64":
			jqint = "binary/jq/osx/jq-osx-amd64"
		}
	case "windows":
		// NOTE: We cannot execute a non-exec
		defer os.Remove(tempFile.Name())
		jqexe = jqexe + ".exe"
		switch runtime.GOARCH {
		case "386":
			jqint = "binary/jq/windows/jq-win32.exe"
		case "amd64":
			jqint = "binary/jq/windows/jq-win64.exe"
		}
	}

	if jqint == "" {
		err = errors.New("error: jq not found in path, and cannot self-extract")
		return "", err
	}

	a.Config.Log.Debugf("Creating temporary file for jq executable: %s from %s", jqexe, jqint)
	file, err := config.TemplateFolder.Open(jqint)
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

// JSON Query will pipe bytes through jq and return results.
func (a *Adapter) JSONQuery(json []byte, jqex string) []byte {
	jq, err := exec.LookPath("jq")
	if err != nil {
		a.Config.Log.Infof("'jq' exec not found in path: extracting from self")
		jq, err = a.UnpackJQExec()
		if err != nil {
			log.Errorf("error: cannot unpack jq from self and not in path.")
			return []byte("")
		}
		defer os.Remove(jq)
	}

	var stdout bytes.Buffer
	cmd := exec.Command(jq, "-c", "-r", jqex)
	cmd.Stdin = strings.NewReader(string(json))
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		log.Warnf("couldn't parse 'jq' command: %s: %v", jqex, err)
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
