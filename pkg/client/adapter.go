package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/cache"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"os/exec"
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
	jq, err := exec.LookPath("jq")
	if err == nil {
		var pretty bytes.Buffer
		cmd := exec.Command(jq, ".")
		cmd.Stdin = strings.NewReader(string(json))
		cmd.Stdout = &pretty
		err := cmd.Run()
		if err == nil {
			json = []byte(pretty.String())
		}
	}
	return json
}

func (a *Adapter) VulnsExportStatus(exportUUID string) (string, error) {

	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStatus, metrics.Methods.Service.Get)

	u := NewUnmarshal(a.Config, a.Metrics)

	status, err := u.VulnsExportStatus(exportUUID)
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return "", err
	}

	return status, nil
}

func (a *Adapter) VulnsExportStart() (string, error) {
	a.Metrics.ClientInc(metrics.EndPoints.VulnsExportStart, metrics.Methods.Service.Update)

	u := NewUnmarshal(a.Config, a.Metrics)

	json, err := u.VulnsExportStart()
	if err != nil {
		a.Config.Log.Errorf("error: failed to get the export-vulns: %v", err)
		return "", err
	}

	return json, nil
}
