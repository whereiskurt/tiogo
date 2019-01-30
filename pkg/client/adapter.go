package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/cache"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// CacheLabel is the type for where to store the response
type CachePathLabel string

func (c CachePathLabel) String() string {
	return "adapter/" + string(c)
}

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

func (a *Adapter) diskStore(label CachePathLabel, obj interface{}) {
	j, err := json.Marshal(obj)
	if err == nil {
		_ = a.DiskCache.Store(fmt.Sprintf("%s.json", label), PrettyJSON(j))
	}
}

// PrettyJSON will look for 'jq' to pretty the json input
func PrettyJSON(json []byte) []byte {
	return JSONQuery(json, ".")
}

// JSON Query will pipe bytes through jq and return results.
func JSONQuery(json []byte, jqex string) []byte {
	jq, err := exec.LookPath("jq")
	if err == nil {
		var pretty bytes.Buffer
		cmd := exec.Command(jq, "-c", jqex)
		cmd.Stdin = strings.NewReader(string(json))
		cmd.Stdout = &pretty
		err := cmd.Run()
		if err == nil {
			json = []byte(pretty.String())
		}
	}
	return json
}

func (a *Adapter) VulnsExportStart() (string, error) {
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
func (a *Adapter) VulnsExportStatus(exportUUID string) (VulnExportStatus, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)

	var status VulnExportStatus
	raw, err := u.VulnsExportStatus(exportUUID)
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return status, err
	}

	convert := NewConvert()
	status, err = convert.ToVulnExportStatus(raw)

	return status, err
}
func (a *Adapter) VulnsExportGet(exportUUID string, chunks string) error {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportGet, metrics.Methods.Service.Get)

	if exportUUID == "" {
		return errors.New("error: empty uuid: must provide uuid using '--uid=12344-1231-23323'")
	}

	chunks, err := a.ChunkList(exportUUID, chunks)
	if err != nil {
		return err
	}

	u := NewUnmarshal(a.Config, a.Metrics)

	chunkSize := len(strings.Split(chunks, ","))
	a.Config.Log.Infof("ChunkSize='%d' for uuid='%s", chunkSize, exportUUID)
	for _, chunk := range strings.Split(chunks, ",") {
		raw, err := u.VulnsExportGet(exportUUID, chunk)
		if err != nil {
			a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
			return err
		}

		a.Config.Log.Infof("Downloaded chunk '%s', file size '%d' bytes", chunk, len(raw))
	}

	return nil
}
func (a *Adapter) VulnsExportQuery(exportUUID string, chunks string, jqex string) error {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportQuery, metrics.Methods.Service.Get)

	if exportUUID == "" {
		return errors.New("error: empty uuid: must provide uuid using '--uid=12344-1231-23323'")
	}

	chunks, err := a.ChunkList(exportUUID, chunks)
	if err != nil {
		return err
	}

	for _, chunk := range strings.Split(chunks, ",") {
		filename, err := tenable.ToCacheFilename(tenable.EndPoints.VulnsExportGet, map[string]string{"ExportUUID": exportUUID, "ChunkID": chunk})
		if err != nil {
			return errors.New(fmt.Sprintf("error: can't get chunk filename for uuid='%s', chunk='%s'", exportUUID, chunk))
		}
		filename = filepath.Join(a.Config.VM.CacheFolder, "service", filename)

		a.Config.Log.Infof("Reading chunk file '%s' ", filename)

		if _, stat := os.Stat(filename); os.IsNotExist(stat) {
			// File doesn't exist return no error
			return errors.New(fmt.Sprintf("Cannot read cached file: '%s", filename))
		}

		bb, err := ioutil.ReadFile(filename)
		filt := JSONQuery(bb, jqex)
		fmt.Printf("%s\n", string(filt))

	}

	return nil
}

func (a *Adapter) ChunkList(exportUUID string, chunks string) (string, error) {
	chunkList := chunks
	if chunkList == "ALL" || chunkList == "" {
		status, err := a.VulnsExportStatus(exportUUID)
		if err != nil {
			a.Config.Log.Errorf("error: can get status for uuid='%s': %v", exportUUID, err)
			return "", err
		}
		chunkList = strings.Join(status.Chunks, ",")
	}
	return chunkList, nil
}
