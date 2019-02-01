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
func (a *Adapter) ExportVulnsStatus(uuid string) (VulnExportStatus, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)

	var status VulnExportStatus
	raw, err := u.VulnsExportStatus(uuid)
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
		uuid, err = a.CachedExportUUID()
		if err != nil {
			return err
		}
	}

	chunks, err := a.CachedChunkList(uuid, chunks)
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

	chunks, err := a.CachedChunkList(uuid, chunks)
	if err != nil {
		return err
	}

	for _, chunk := range strings.Split(chunks, ",") {
		filename, err := a.ToCacheFilename(tenable.EndPoints.VulnsExportGet, map[string]string{"ExportUUID": uuid, "ChunkID": chunk})
		if err != nil {
			a.Config.Log.Errorf("error: reading chunk file '%s' ", filename)
			return err
		}

		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.New(fmt.Sprintf("error: cannot read cached file: '%s: %v", filename, err))
		}

		filter := a.JSONQuery(bb, jqex)

		// Output the JQEX filtered JSON response.
		fmt.Println(string(filter))
	}

	return nil
}

// CachedChunkList
func (a *Adapter) CachedChunkList(uuid string, chunks string) (string, error) {
	if chunks != "ALL" || chunks != "" {
		return chunks, nil
	}
	ep := tenable.EndPoints.VulnsExportStatus
	p := map[string]string{"ExportUUID": uuid}
	filename, err := a.ToCacheFilename(ep, p)
	if err != nil {
		a.Config.Log.Errorf("error: reading cached 'status' file '%s' ", filename)
		return "", err
	}

	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error: cannot read cached file: '%s: %v", filename, err))
	}

	var status VulnExportStatus
	err = json.Unmarshal(bb, status)
	if err != nil {
		a.Config.Log.Errorf("error: can get unmarshal '%s': %v", filename, err)
		return "", err
	}

	chunkList := strings.Join(status.Chunks, ",")

	return chunkList, nil
}

// JSONPretty will look for 'jq' to pretty the json input
func (a *Adapter) JSONPretty(json []byte) []byte {
	return a.JSONQuery(json, ".")
}

// UnpackJQExec extracts the jq executable packed in templates.go
func (a *Adapter) UnpackJQExec() (string, error) {
	tempFile, err := ioutil.TempFile("", "jq.")
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
	cmd := exec.Command(jq, "-c", jqex)
	cmd.Stdin = strings.NewReader(string(json))
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		log.Warnf("couldn't parse 'jq' command: %s: %v", jqex, err)
		return []byte("")
	}

	return []byte(stdout.String())
}

// CachedExportUUID reads the service export vulns cache for entries, and returns the first one as exportUUID
func (a *Adapter) CachedExportUUID() (string, error) {
	folder := filepath.Join(a.Config.VM.CacheFolder, "service", "export", "vulns")
	entries, err := a.DirEntries(folder)
	if err != nil || len(entries) == 0 {
		return "", err
	}
	a.Config.Log.Debugf("Returning first entry as uuid: %s", entries[0])
	return entries[0], nil
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
func (a *Adapter) ToCacheFilename(endpoint tenable.EndPointType, p map[string]string) (string, error) {

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
