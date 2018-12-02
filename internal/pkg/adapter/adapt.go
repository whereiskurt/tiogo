package adapter

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/cache"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
)

var AllAgentScannerUUID = "00000000-0000-0000-0000-00000000000000000000000000001"

type Adapter struct {
	Config    *app.Config
	Disk      *cache.Disk
	Memory    *cache.Memory
	Worker    *sync.WaitGroup
	Log       *app.Logger
	Convert   *Convert
	Unmarshal *Unmarshal
	Tenable   *tenable.Service
}

var MaxCipherByteLength = 4

type RegexFilename struct {
	Name     string
	Regexp   *regexp.Regexp
	Filename string
}

var urlServiceTmpl = map[string]string{
	"ScanList":                "{{.baseURL}}/scans",
	"Scan":                    "{{.baseURL}}/scans/{{.scanID}}",
	"ScanHistoryDetail":       "{{.baseURL}}/scans/{{.scanID}}?history_id={{.historyID}}",
	"AssetMap":                "{{.baseURL}}/private/scans/{{.scanID}}/assets/vulnerabilities?history_id={{.historyID}}",
	"HostScanDetail":          "{{.baseURL}}/scans/{{.scanID}}/hosts/{{.hostID}}?history_id={{.historyID}}",
	"PluginFamilies":          "{{.baseURL}}/plugins/families",
	"PluginFamily":            "{{.baseURL}}/plugins/families/{{.familyID}}",
	"Plugin":                  "{{.baseURL}}/plugins/plugin/{{.pluginID}}",
	"Asset":                   "{{.baseURL}}/workbenches/assets/{{.assetUUID}}/info",
	"Scanners":                "{{.baseURL}}/scanners",
	"ScannerAgents":           "{{.baseURL}}/scanners/{{.scannerID}}/agents?offset={{.offset}}&limit={{.limit}}",
	"VulnSummary":             "{{.baseURL}}/workbenches/assets/{{.assetUUID}}/vulnerabilities",
	"VulnDetail":              "{{.baseURL}}/workbenches/assets/{{.assetUUID}}/vulnerabilities/{{.pluginID}}/info",
	"VulnOutput":              "{{.baseURL}}/workbenches/assets/{{.assetUUID}}/vulnerabilities/{{.pluginID}}/outputs",
	"ScannerAgentGroup":       "{{.baseURL}}/scanners/{{.scannerID}}/agent-groups",
	"AssignScannerAgentGroup": "{{.baseURL}}/scanners/{{.scannerID}}/agent-groups/{{.groupID}}/agents/{{.agentID}}",
	"AssetExport":             "{{.baseURL}}/assets/export",
	"AssetExportStatus":       "{{.baseURL}}/assets/export/{{.exportUUID}}/status",
	"AssetExportChunk":        "{{.baseURL}}/assets/export/{{.exportUUID}}/chunks/{{.chunkID}}",
	"VulnExport":              "{{.baseURL}}/vulns/export",
	"VulnExportStatus":        "{{.baseURL}}/vulns/export/{{.exportUUID}}/status",
	"VulnExportChunk":         "{{.baseURL}}/vulns/export/{{.exportUUID}}/chunks/{{.chunkID}}",
}
var jsonBodyTmpl = map[string]string{
	"ScannerAgentGroup": `{ "name": "{{.name}}" }`,
	"AssetExport": `{ "chunk_size": {{.chunkSize}},	"filters": { } }`,
	"VulnExport": `{ "num_assets": {{.chunkSize}},	"filters": { } }`,
}
var urlCacheFilename = []RegexFilename{
	{"ScanList", regexp.MustCompile("^.*?/scans$"), "tenable/scan/list.json"},
	{"Scan", regexp.MustCompile("^.*?/scans/(\\d+)$"), "tenable/scan/$1/history.json"},
	{"ScanHistoryDetail", regexp.MustCompile("^.*?/scans/(\\d+)\\?history_id=(\\d+)$"), "tenable/scan/$1/$2/scan.json"},
	{"AssetMap", regexp.MustCompile("^.*?/scans/(\\d+)/assets/vulnerabilities\\?history_id=(\\d+)$"), "tenable/scan/$1/$2/hostassetmap.json"},
	{"HostScanDetail", regexp.MustCompile("^.*?/scans/(\\d+)/hosts/(\\d+)\\?history_id=(\\d+)$"), "tenable/scan/$1/$3/host/$2/host.json"},
	{"PluginFamilies", regexp.MustCompile("^.*?/plugins/families$"), "tenable/plugin/family/families.json"},
	{"PluginFamily", regexp.MustCompile("^.*?/plugins/families/(\\d+)$"), "tenable/plugin/family/$1/family.json"},
	{"Plugin", regexp.MustCompile("^.*?/plugins/plugin/(\\d+)$"), "tenable/plugin/$1/plugin.json"},
	{"Scanners", regexp.MustCompile("^.*?/scanners$"), "tenable/scanners/list.json"},
	{"ScannerAgents", regexp.MustCompile("^.*?/scanners/(\\d+)/agents\\?offset=(\\d+)&limit=(\\d+)$"), "tenable/scanners/$1/agents.$2.$3.json"},
	{"VulnSummary", regexp.MustCompile("^.*?/workbenches/assets/(.+)/vulnerabilities$"), "tenable/asset/$1/vuln/summary.json"},
	{"VulnDetail", regexp.MustCompile("^.*?/workbenches/assets/(.+)/vulnerabilities/(.+)/info$"), "tenable/asset/$1/vuln/$2/detail.json"},
	{"VulnOutput", regexp.MustCompile("^.*?/workbenches/assets/(.+)/vulnerabilities/(.+)/outputs$"), "tenable/asset/$1/vuln/$2/output.json"},
	{"Asset", regexp.MustCompile("^.*?/workbenches/assets/(.+?)/info$"), "tenable/asset/$1/asset.json"},
	{"AssetExportChunk", regexp.MustCompile("^.*?/assets/export/(.+?)/chunks/(.+?)$"), "tenable/export/asset/$1/chunk.$2.json"},
	{"VulnExportChunk", regexp.MustCompile("^.*?/vulns/export/(.+?)/chunks/(.+?)$"), "tenable/export/vuln/$1/chunk.$2.json"},

	// Needed for linking an asset under the scan AND the asset
	// TODO: Figure out if we reallllly need this. I think we don't....
	{"AssetDetailLink", regexp.MustCompile("^.*?/workbenches/assets/(.+)/info/(\\d+)/(\\d+)/(\\d+)$"), "tenable/scan/$2/$3/host/$4/asset/$1/asset.json"},
}

// ToFilename converts the Tenable.IO a.ToURL to a local filename to store cache
func (a *Adapter) ToFilename(url string) (filename string, err error) {
	crypto := a.Config.CryptoCacheMode
	key := a.Config.CacheKey
	folder := a.Config.CacheFolder

	for i := range urlCacheFilename {
		if m := urlCacheFilename[i].Regexp.FindStringSubmatch(url); m != nil {
			if crypto {
				// Encrypt parameters
				for j := range m {
					cipher := sha256.Sum256([]byte(fmt.Sprintf("%s%s", key, m[j])))
					m[j] = fmt.Sprintf("%x", cipher[:MaxCipherByteLength])
				}
			}
			filename = urlCacheFilename[i].Filename
			for k := range m {
				filename = strings.Replace(filename, fmt.Sprintf("$%d", k), m[k], -1)
			}
			filename = fmt.Sprintf("%s/%s", folder, filename)
			return
		}
	}
	err = errors.New("can't create filename from url: " + url)
	return
}
func (a *Adapter) ToURL(name string, p map[string]string) (url string, err error) {
	return a.ToTemplated(name, p, urlServiceTmpl)
}
func (a *Adapter) ToJSON(name string, p map[string]string) (url string, err error) {
	return a.ToTemplated(name, p, jsonBodyTmpl)
}

func (a *Adapter) ToTemplated(name string, p map[string]string, tmap map[string]string) (url string, err error) {
	var rawURL bytes.Buffer
	t, terr := template.New(name).Parse(tmap[name])
	if terr != nil {
		err = errors.New(fmt.Sprintf("error: failed to parse template for %s: %v", name, err))
		return
	}
	err = t.Execute(&rawURL, p)
	if err != nil {
		return
	}

	url = rawURL.String()

	return
}

// NewAdapater manages calls the remote services, converts the results and manages a memory/disk cache.
func NewAdapter(config *app.Config) (a *Adapter) {
	a = new(Adapter)
	a.Config = config
	a.Worker = new(sync.WaitGroup)
	a.Memory = cache.NewMemoryCache()
	a.Disk = cache.NewDisk(a.Config.CryptoCacheMode, a.Config.CacheKey, a.Config.CacheFolder)
	a.Log = config.Logger
	a.Convert = NewConvert(a)
	a.Unmarshal = NewUnmarshal(a)
	a.Tenable = tenable.NewService(a.Config.BaseURL, a.Config.SecretKey, a.Config.AccessKey, a.Log)
	return
}

func (a *Adapter) Scans() (ss []dao.Scan, err error) {
	memcache := a.Memory.Cache

	var p = map[string]string{
		"baseURL": a.Config.BaseURL,
	}
	url, terr := a.ToURL("ScanList", p)
	if terr != nil {
		err = terr
		return
	}

	memhit := memcache.Get(url)
	if memhit != nil {
		ss = memhit.Value().([]dao.Scan)
		return
	}

	ss, err = a.Unmarshal.ScanList(url)
	if err == nil {
		memcache.Set(url, ss, time.Minute*60)
	}

	return
}
func (a *Adapter) ScanHistory(s dao.Scan, depth int) (sh dao.ScanHistory, err error) {
	if depth < 1 {
		err = errors.New("error: scan history depth must be larger than zero")
		return
	}

	memcache := a.Memory.Cache

	var p = map[string]string{
		"baseURL": a.Config.BaseURL,
		"scanID":  s.ScanID,
	}
	url, terr := a.ToURL("Scan", p)
	if terr != nil {
		err = terr
		return
	}

	hit := memcache.Get(url)
	if hit != nil {
		sh = hit.Value().(dao.ScanHistory)
		return
	}

	sh, err = a.Unmarshal.ScanHistory(s, depth, url)

	if err != nil {
		return
	}

	memcache.Set(url, sh, time.Minute*60)
	return
}
func (a *Adapter) ScanHistoryDetail(sh dao.ScanHistory, offset int) (shd dao.ScanHistoryDetail, err error) {
	if len(sh.History) == 0 {
		return
	}

	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"scanID":    sh.Scan.ScanID,
		"historyID": sh.History[0].HistoryID,
	}
	url, terr := a.ToURL("ScanHistoryDetail", p)
	if terr != nil {
		err = terr
		return
	}

	memkey := url
	memhit := a.Memory.Cache.Get(memkey)
	if memhit != nil {
		shd = memhit.Value().(dao.ScanHistoryDetail)
		return
	}

	filekey, diskerr := a.ToFilename(url)
	if diskerr != nil {
		return
	}

	raw, diskerr := a.Disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && a.Config.ClobberCacheMode == true {
		a.Disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var sd tenable.ScanDetail
	err = json.Unmarshal(raw, &sd)
	if err != nil {
		if len(raw) > 0 {
			if a.Config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScanHistoryDetail: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			a.Disk.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.
		}

		sd, raw, err = a.Tenable.ScanDetail(url)
		if err != nil {
			return
		}
	}

	if offset >= len(sd.History) {
		err = errors.New("can't fetch ScanHistoryDetail for offset past end of  history")
		return
	}

	shd, err = a.Convert.ScanHistoryDetail(sh, sd, offset)
	if err != nil {
		return
	}

	a.Memory.Cache.Set(memkey, shd, time.Minute*60)
	a.Disk.Store(filekey, raw, a.Config.CachePretty)

	return
}
func (a *Adapter) HostScanDetail(hss dao.HostScanSummary) (hsd dao.HostScanDetail, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"scanID":    hss.ScanHistoryDetail.Scan.ScanID,
		"hostID":    hss.HostID,
		"historyID": hss.ScanHistoryDetail.HistoryID,
	}
	url, terr := a.ToURL("HostScanDetail", p)
	if terr != nil {
		err = terr
		return
	}

	memkey := url
	memhit := a.Memory.Cache.Get(memkey)
	if memhit != nil {
		hsd = memhit.Value().(dao.HostScanDetail)
		return
	}

	hsd, err = a.Unmarshal.HostScanDetail(hss, url)
	if err != nil {
		a.Log.Warnf("failed to download HostScanDetail: %s: %v", url, err)
		err = nil
	}

	hss.HostDetail = hsd
	a.Memory.Cache.Set(memkey, hsd, time.Minute*60)
	return
}
func (a *Adapter) AssetDetail(hss dao.HostScanSummary, uuid string) (ad dao.Asset, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"assetUUID": uuid,
	}

	url, err := a.ToURL("Asset", p)
	if err != nil {
		return
	}

	ad, err = a.Unmarshal.AssetDetail(hss, url)
	if err != nil {
		a.Log.Warnf("failed to download Asset details: %s: %v", url, err)
		err = nil
	}
	return
}

func (a *Adapter) Plugins() (pp map[string]dao.Plugin, err error) {
	pp = make(map[string]dao.Plugin)

	var p = map[string]string{
		"baseURL": a.Config.BaseURL,
	}
	url, err := a.ToURL("PluginFamilies", p)
	if err != nil {
		return
	}

	// NOTE: This family doesn' exist!
	// 		"family_name": "Port scanners",
	//  	"name": "Nessus SYN scanner",
	//  	"id": 11219

	// Get all the plugin families
	ff, famerr := a.Unmarshal.PluginFamilies(url)
	if famerr != nil {
		err = famerr
		a.Log.Errorf("failed to UNMARSHAL PluginFamiles: %s:%+v", err)
		return
	}

	// For each family get all the plugins IDs
	for _, f := range ff {

		var p = map[string]string{
			"baseURL":  a.Config.BaseURL,
			"familyID": string(f.ID),
		}

		url, perr := a.ToURL("PluginFamily", p)
		if perr != nil {
			err = perr
			return
		}

		pff, perr := a.Unmarshal.PluginFamily(url)
		if perr != nil {
			err = perr
			a.Log.Errorf("failed to UNMARSHAL PluginFamily: %s:%+v", string(f.ID), err)
			return
		}
		// Store plugin ID to map, for PluginDetail to complete
		for _, id := range pff.PluginID {
			pp[id] = dao.Plugin{
				PluginID:   id,
				FamilyId:   f.ID,
				FamilyName: f.Name,
			}
		}
	}

	// Get all details for the plugins and write them into pp
	err = a.PluginDetail(pp)

	return
}

// PluginDetail takes a map of dao.Plugin structs and completes the .Detail
func (a *Adapter) PluginDetail(summary map[string]dao.Plugin) (err error) {
	for i := range summary {

		var p = map[string]string{
			"baseURL":  a.Config.BaseURL,
			"pluginID": summary[i].PluginID,
		}
		url, terr := a.ToURL("Plugin", p)
		if terr != nil {
			err = terr
			return
		}

		var pd dao.PluginDetail
		memhit := a.Memory.Cache.Get(url)
		if memhit != nil {
			pd = memhit.Value().(dao.PluginDetail)
		} else {
			pd, err = a.Unmarshal.PluginDetail(summary[i], url)
			if err != nil {
				a.Log.Warnf("failed to download Plugin details: %s: %v: leaving blank.", url, err)
				err = nil
			}
		}

		a.Memory.Cache.Set(url, pd, time.Minute*60)

		// NOTE: Cannot just do 'summary[i].Detail = pd' because summary isn't a map of pointers.
		// This is because of 'summary[i].Detail = ' is a map index expression, and different than a pointer.
		t := summary[i]
		t.Detail = pd
		summary[i] = t

	}
	return
}

func (a *Adapter) Scanners() (ss []dao.Scanner, err error) {
	var p = map[string]string{
		"baseURL": a.Config.BaseURL,
	}
	url, err := a.ToURL("Scanners", p)
	if err != nil {
		return
	}

	ss, err = a.Unmarshal.Scanners(url)

	return
}
func (a *Adapter) AgentScanner(ss []dao.Scanner) (s dao.Scanner, err error) {
	for i := range ss {
		// FACT: There is only one scanner all of the agents are attached to :-)
		if ss[i].UUID == AllAgentScannerUUID {
			s = ss[i]
			return
		}
	}
	err = errors.New(fmt.Sprintf("Cannot find scanner with UUID '%s'.", AllAgentScannerUUID))
	return
}
func (a *Adapter) Agents(ss dao.Scanner) (sa []dao.ScannerAgent, err error) {
	PageSize := 5000
	page := 0

	// Do Pagination for Agents with a PageSize
	for {
		offset := int64(page * PageSize)
		var p = map[string]string{
			"baseURL":   a.Config.BaseURL,
			"scannerID": ss.ID,
			"offset":    fmt.Sprintf("%d", offset),
			"limit":     fmt.Sprintf("%d", PageSize),
		}
		url, uerr := a.ToURL("ScannerAgents", p)
		if uerr != nil {
			err = uerr
			return
		}

		agents, pag, terr := a.Unmarshal.ScannerAgents(ss, url)
		if terr != nil {
			err = terr
			return
		}
		sa = append(sa, agents...)

		limit, _ := pag.Limit.Int64()
		// offset, _ := pag.Offset.Int64()
		total, _ := pag.Total.Int64()

		if int64(len(sa)) >= total || limit == 0 || ((limit * (offset + 1)) >= total) {
			break
		}
		page = page + 1
	}

	return
}
func (a *Adapter) AgentGroups(scanner dao.Scanner) (aa []dao.ScannerAgentGroup, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"scannerID": scanner.ID,
	}
	url, uerr := a.ToURL("ScannerAgentGroup", p)

	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to create url for GET ScannerAgentGroup: %v", err)
		return
	}

	aa, err = a.Unmarshal.ScannerAgentGroups(url)

	return
}
func (a *Adapter) AgentGroupNames(agents []dao.ScannerAgent) (groups []string) {
	p := make(map[string]bool)

	for i := range agents {
		for j := range agents[i].Groups {
			p[agents[i].Groups[j].Name] = true
		}
	}

	for key := range p {
		if p[key] == true {
			groups = append(groups, key)
		}
	}

	return
}
func (a *Adapter) AssignAgentGroup(scanner dao.Scanner, agent dao.ScannerAgent, group dao.ScannerAgentGroup) (err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"scannerID": scanner.ID,
		"agentID":   agent.ID,
		"groupID":   group.ID,
	}
	url, uerr := a.ToURL("AssignScannerAgentGroup", p)

	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to create url for GET AssignScannerAgentGroup: %v", err)
		return
	}
	_, err = a.Unmarshal.AssignScannerAgentGroup(url)

	return
}
func (a *Adapter) CreateAgentGroup(scanner dao.Scanner, name string) (sag dao.ScannerAgentGroup, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"scannerID": scanner.ID,
		"name":      name,
	}
	url, uerr := a.ToURL("ScannerAgentGroup", p)
	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to create url for POST ScannerAgentGroup: %v", err)
		return
	}

	j, jerr := a.ToJSON("ScannerAgentGroup", p)
	if jerr != nil {
		err = jerr
		a.Log.Errorf("failed to JSON body for POST ScannerAgentGroup: %v", err)
		return
	}

	sag, err = a.Unmarshal.CreateScannerAgentGroup(url, j)

	return
}
func (a *Adapter) MatchAgentGroup(scanner dao.Scanner, name string) (group dao.ScannerAgentGroup, ok bool, err error) {
	aa, aerr := a.AgentGroups(scanner)
	if aerr != nil {
		err = aerr
		return
	}

	for i := range aa {
		if aa[i].Name == name {
			group = aa[i]
			ok = true
			break
		}
	}
	return
}

func (a *Adapter) AssetVulnSummary(ad dao.Asset) (vd []dao.AssetVuln, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"assetUUID": ad.UUID,
	}
	url, uerr := a.ToURL("VulnSummary", p)
	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to create url for VulnSummary: %v", err)
		return
	}

	vd, err = a.Unmarshal.VulnSummary(ad, url)
	if err != nil {
		a.Log.Errorf("failed to unmarshal VulnSummary: %s: %v", url, err)
		return
	}

	return
}
func (a *Adapter) AssetVulnDetail(vs dao.AssetVuln, uuid string) (vd dao.AssetVulnDetail, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"assetUUID": uuid,
		"pluginID":  vs.PluginID,
	}
	url, uerr := a.ToURL("VulnDetail", p)
	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to create url for VulnDetail: %v", err)
		return
	}

	vd, verr := a.Unmarshal.VulnDetail(vs, url)
	if verr != nil {
		a.Log.Errorf("failed to unmarshal VulnDetail: %s: %v", url, err)
		err = verr
		return
	}

	return
}
func (a *Adapter) AssetVulnOutput(vs dao.AssetVuln, uuid string) (vo []dao.AssetVulnOutput, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"assetUUID": uuid,
		"pluginID":  vs.PluginID,
	}
	url, uerr := a.ToURL("VulnOutput", p)
	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to create url for VulnOutput: %v", err)
		return
	}

	vo, verr := a.Unmarshal.VulnOutput(vs, url)
	if verr != nil {
		a.Log.Errorf("failed to unmarshal VulnOutput: %s: %v", url, err)
		err = verr
		return
	}

	return
}
func (a *Adapter) AssetExport() (assets []dao.Asset, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"chunkSize": "5000",
	}
	urlExport, uerr := a.ToURL("AssetExport", p)

	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to build url for GET AssetExport: %v", err)
		return
	}

	j, jerr := a.ToJSON("AssetExport", p)
	if jerr != nil {
		err = jerr
		a.Log.Errorf("failed to JSON body for POST ScannerAgentGroup: %v", err)
		return
	}

	var exportUUID string
	exportUUID, err = a.Unmarshal.AssetExport(urlExport, j)
	if err != nil {
		a.Log.Errorf("failed to unmarshal for AssetExport: %v", err)
		return
	}

	assets, err = a.AssetExportDownload(exportUUID)

	return
}
func (a *Adapter) AssetExportDownload(exportUUID string) (assets []dao.Asset, err error) {
	var p = map[string]string{
		"baseURL":    a.Config.BaseURL,
		"chunkSize":  "5000",
		"exportUUID": exportUUID,
	}

	urlStatus, uerr := a.ToURL("AssetExportStatus", p)
	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to build url for GET AssetExportStatus: %v", err)
		return
	}

	var maxretry = 5
	for {
		var export dao.AssetExportStatus
		export, err = a.Unmarshal.AssetExportStatus(urlStatus)

		if export.Status == "FINISHED" {
			assets, err = a.AssetExportChunks(exportUUID, export.Chunks)
			fmt.Println(fmt.Sprintf("Assets Export Status: %d", len(assets)))
			break
		}

		time.Sleep(time.Duration(10000 * time.Millisecond))
		maxretry = maxretry - 1
		if maxretry < 1 {
			a.Log.Errorf("Export not ready for download. Exceeded timeout and retrys.")
			return
		}
	}

	return
}
func (a *Adapter) AssetExportChunks(exportUUID string, chunks []string) (aa []dao.Asset, err error) {
	for c := range chunks {
		var p = map[string]string{
			"baseURL":    a.Config.BaseURL,
			"exportUUID": exportUUID,
			"chunkID":    chunks[c],
		}

		urlStatus, uerr := a.ToURL("AssetExportChunk", p)
		if uerr != nil {
			err = uerr
			a.Log.Errorf("failed to build url for GET AssetExportStatus: %v", err)
			return
		}

		var assets []dao.Asset
		assets, err = a.Unmarshal.AssetExportChunk(urlStatus)
		if err != nil {
			a.Log.Errorf("failed to unmarshall AssetExportChunks: %v", err)
			return
		}
		aa = append(aa, assets...)
	}

	return
}

func (a *Adapter) VulnExport() (aa []dao.VulnExportChunk, err error) {
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"chunkSize": "5000",
	}
	urlExport, uerr := a.ToURL("VulnExport", p)

	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to build url for GET VulnExport: %v", err)
		return
	}

	j, jerr := a.ToJSON("VulnExport", p)
	if jerr != nil {
		err = jerr
		a.Log.Errorf("failed to JSON body for POST ScannerAgentGroup: %v", err)
		return
	}

	var exportUUID string
	exportUUID, err = a.Unmarshal.VulnExport(urlExport, j)
	if err != nil {
		a.Log.Errorf("failed to unmarshal for VulnExport: %v", err)
		return
	}

	aa, err = a.VulnExportDownload(exportUUID)

	return
}
func (a *Adapter) VulnExportDownload(exportUUID string) (assets []dao.VulnExportChunk, err error) {
	var p = map[string]string{
		"baseURL":    a.Config.BaseURL,
		"chunkSize":  "5000",
		"exportUUID": exportUUID,
	}

	urlStatus, uerr := a.ToURL("VulnExportStatus", p)
	if uerr != nil {
		err = uerr
		a.Log.Errorf("failed to build url for GET VulnExportStatus: %v", err)
		return
	}

	var maxretry = 6
	for {
		var export dao.VulnExportStatus
		export, err = a.Unmarshal.VulnExportStatus(urlStatus)

		if export.Status == "FINISHED" {
			assets, err = a.VulnExportChunks(exportUUID, export.Chunks)
			break
		}

		time.Sleep(time.Duration(30000 * time.Millisecond))
		maxretry = maxretry - 1
		if maxretry < 1 {
			a.Log.Errorf("Export not ready for download. Exceeded timeout and retrys.")
			return
		}
	}

	return
}
func (a *Adapter) VulnExportChunks(exportUUID string, chunks []string) (aa []dao.VulnExportChunk, err error) {
	for c := range chunks {
		var p = map[string]string{
			"baseURL":    a.Config.BaseURL,
			"exportUUID": exportUUID,
			"chunkID":    chunks[c],
		}

		urlStatus, uerr := a.ToURL("VulnExportChunk", p)
		if uerr != nil {
			err = uerr
			a.Log.Errorf("failed to build url for GET VulnExportStatus: %v", err)
			return
		}

		var assets []dao.VulnExportChunk
		assets, err = a.Unmarshal.VulnExportChunk(urlStatus)
		if err != nil {
			a.Log.Errorf("failed to unmarshall VulnExportChunks: %v", err)
			return
		}
		aa = append(aa, assets...)
	}

	return
}

func (a *Adapter) HostHandleFunc(detail dao.ScanHistoryDetail, host *dao.HostScanSummary) (err error) {
	host.HostDetail, err = a.HostScanDetail(*host)
	if err != nil {

		a.Log.Errorf("failed to get host scan details: %+v", err)
		return
	}
	err = a.PluginDetail(host.HostDetail.PluginMap)
	if err != nil {
		a.Log.Errorf("failed to get plugin scan details: %+v", err)
		return
	}
	uuid := detail.HostAssetMap[host.HostID]
	host.Asset, err = a.AssetDetail(*host, uuid)
	if err != nil {
		a.Log.Errorf("failed to get asset details: %+v", err)
		return
	}
	//
	// host.Asset.VM, err = a.AssetVulnSummary(host.Asset)
	// if err != nil {
	// 	a.Log.Errorf("failed to get asset vulnerability summary: %+v", err)
	// 	return
	// }
	//
	// for v, vuln := range host.Asset.VM {
	// 	var vd dao.AssetVulnDetail
	// 	var vo []dao.AssetVulnOutput
	//
	// 	vd, err = a.AssetVulnDetail(vuln, host.Asset.UUID)
	// 	if err != nil {
	// 		a.Log.Errorf("failed to get asset vulnerability detail: %+v", err)
	// 		return
	// 	}
	//
	// 	vo, err = a.AssetVulnOutput(vuln, host.Asset.UUID)
	// 	if err != nil {
	// 		a.Log.Errorf("failed to get asset vulnerability output: %+v", err)
	// 		return
	// 	}
	//
	// 	host.Asset.VM[v].Detail = vd
	// 	host.Asset.VM[v].Output = vo
	// }

	return
}
