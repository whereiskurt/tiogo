package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"time"
)

type Unmarshal struct {
	Adapter *Adapter
	Infof   func(fmt string, args ...interface{})
	Debugf  func(fmt string, args ...interface{})
	Warnf   func(fmt string, args ...interface{})
	Errorf  func(fmt string, args ...interface{})
}

func NewUnmarshal(a *Adapter) (u *Unmarshal) {
	u = new(Unmarshal)
	u.Adapter = a

	u.Errorf = a.Config.Logger.Errorf
	u.Debugf = a.Config.Logger.Debugf
	u.Warnf = a.Config.Logger.Warnf
	u.Infof = a.Config.Logger.Infof

	return
}

func (u *Unmarshal) ScanList(url string) (ss []dao.Scan, err error) {
	var s tenable.ScanList

	a := u.Adapter
	cache := a.Disk
	convert := a.Convert
	config := a.Config

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := cache.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		cache.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	// If no raw fetched from cache or raw fetched wasn't unmarshalable (ie. changed cacheKey?)
	if len(raw) == 0 {
		s, raw, err = a.Tenable.ScanList(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &s)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScanList: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			cache.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.

			s, raw, err = a.Tenable.ScanList(url)
			if err != nil {
				return
			}
		}
	}

	ss, err = convert.ScanList(s)

	if err == nil {
		cache.Store(filekey, raw, config.CachePretty) // Tenable disk cache store
	}
	return
}
func (u *Unmarshal) ScanHistory(s dao.Scan, depth int, url string) (sh dao.ScanHistory, err error) {

	a := u.Adapter
	t := a.Tenable
	cache := a.Disk
	convert := a.Convert
	config := a.Config

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := cache.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		cache.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	// Without having a history_id, we will get the list of all scans
	var sd tenable.ScanDetail

	if len(raw) == 0 {
		sd, raw, err = t.ScanDetail(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &sd)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScanHistoryDetail: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			cache.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.

			sd, raw, err = t.ScanDetail(url)
			if err != nil {
				return
			}
		}
	}

	// We now have our history_id, so cache it!
	if len(sd.History) > 0 {

		var p = map[string]string{
			"baseURL":   a.Config.BaseURL,
			"scanID":    s.ScanID,
			"historyID": string(sd.History[0].HistoryID),
		}
		url, terr := a.ToURL("ScanHistoryDetail", p)
		if terr != nil {
			err = terr
			return
		}

		filename, ferr := a.ToFilename(url)
		if ferr != nil {
			err = ferr
			return
		}
		cache.Store(filename, raw, config.CachePretty)
	}

	sh, err = convert.ScanHistory(s, sd.History, depth)

	if err == nil {
		cache.Store(filekey, raw, config.CachePretty) // Tenable disk cache store
	}

	return
}
func (u *Unmarshal) HostScanDetail(hss dao.HostScanSummary, url string) (hsd dao.HostScanDetail, err error) {
	a := u.Adapter
	memcache := a.Memory.Cache
	cache := a.Disk
	config := a.Config

	memkey := url
	memhit := memcache.Get(memkey)
	if memhit != nil {
		hsd = memhit.Value().(dao.HostScanDetail)
		return
	}
	filekey, diskerr := a.ToFilename(url)
	if diskerr != nil {
		return
	}

	raw, diskerr := cache.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		cache.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var hd tenable.HostDetail

	if len(raw) == 0 {
		hd, raw, err = a.Tenable.HostDetail(url)
		if err != nil {
			return
		}
		a.Disk.Store(filekey, raw, a.Config.CachePretty)
	} else {
		err = json.Unmarshal(raw, &hd)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal HostScanDetail: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			cache.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.

			hd, raw, err = a.Tenable.HostDetail(url)
			if err != nil {
				return
			}
		}
	}

	hsd, err = a.Convert.HostDetail(hss, &hd)
	if err != nil {
		return
	}

	memcache.Set(memkey, hsd, 60*time.Minute)
	cache.Store(filekey, raw, config.CachePretty) // Tenable disk cache store
	return
}
func (u *Unmarshal) PluginDetail(pds dao.Plugin, url string) (pd dao.PluginDetail, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	cache := a.Disk
	config := a.Config

	memkey := url
	memhit := memcache.Get(memkey)
	if memhit != nil {
		pd = memhit.Value().(dao.PluginDetail)
		return
	}
	filekey, diskerr := a.ToFilename(url)
	if diskerr != nil {
		return
	}

	raw, diskerr := cache.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		cache.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var p tenable.Plugin
	if len(raw) == 0 {
		p, raw, err = a.Tenable.Plugin(url)
		if err != nil {
			return
		}
		cache.Store(filekey, raw, a.Config.CachePretty)
	} else {
		err = json.Unmarshal(raw, &p)

		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScanList: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			cache.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.

			p, raw, err = a.Tenable.Plugin(url)
			if err != nil {
				return
			}
			cache.Store(filekey, raw, a.Config.CachePretty)
		}
	}

	pd, err = a.Convert.PluginDetail(pds, p)
	if err != nil {
		return
	}

	memcache.Set(memkey, pd, 60*time.Minute)
	cache.Store(filekey, raw, config.CachePretty) // Tenable disk cache store

	return
}
func (u *Unmarshal) AssetDetail(hss dao.HostScanSummary, url string) (ad dao.Asset, err error) {
	a := u.Adapter
	cache := a.Disk
	config := a.Config

	// Add Scan,Hist,Host IDs to the URL to make a cache key.
	// Assets aren't per scan,hist,host, but we make it that way to capture historical ai snapshot.
	sID := hss.ScanHistoryDetail.Scan.ScanID
	histID := hss.ScanHistoryDetail.HistoryID
	hID := hss.HostID

	memkey := fmt.Sprintf("%s/%s/%s/%s", url, sID, histID, hID)
	memhit := a.Memory.Cache.Get(memkey)
	if memhit != nil {
		ad = memhit.Value().(dao.Asset)
		return
	}

	filekey, diskerr := a.ToFilename(memkey)
	if diskerr != nil {
		return
	}

	raw, diskerr := cache.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		cache.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var asset tenable.Asset

	if len(raw) == 0 {
		// HTTP GET Asset
		asset, raw, err = a.Tenable.Asset(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &asset)

		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScanList: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			// Delete bad cache entry - cacheKey likely changed and we're clobbering.
			cache.Clear(filekey)
			// HTTP GET Asset
			asset, raw, err = a.Tenable.Asset(url)
			if err != nil {
				return
			}
		}
	}

	ai := asset.Info
	// CONVERT from Tenable JSON to tiogo.DTO structs
	ad, err = a.Convert.AssetDetail(ai)
	if err != nil {
		return
	}
	a.Memory.Cache.Set(memkey, ad, time.Minute*60)

	// FILE CACHE based on
	filekey, diskerr = a.ToFilename(url)
	if diskerr != nil {
		return
	}

	a.Disk.Store(filekey, raw, a.Config.CachePretty)

	return
}
func (u *Unmarshal) PluginFamilies(url string) (ff []dao.PluginFamily, err error) {
	a := u.Adapter
	t := a.Tenable
	cache := a.Disk
	config := a.Config
	memcache := a.Memory.Cache

	memkey := url
	memhit := a.Memory.Cache.Get(memkey)

	var pf tenable.PluginFamilies
	if memhit != nil {
		pf = memhit.Value().(tenable.PluginFamilies)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	// convert := u.Adapter.Convert
	raw, diskerr := cache.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		cache.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	if len(raw) == 0 {
		pf, raw, err = t.PluginFamiles(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &pf)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal PluginFamilies: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			cache.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.

			pf, raw, err = t.PluginFamiles(url)
			if err != nil {
				return
			}
		}
	}

	if err != nil {
		return
	}

	ff, err = a.Convert.PluginFamily(pf)
	if err != nil {
		return
	}

	memcache.Set(memkey, pf, 60*time.Minute)
	cache.Store(filekey, raw, config.CachePretty) // Tenable disk cache store

	return
}
func (u *Unmarshal) PluginFamily(url string) (pf dao.PluginFamilyDetail, err error) {

	a := u.Adapter
	t := a.Tenable
	disk := u.Adapter.Disk
	config := u.Adapter.Config
	memcache := a.Memory.Cache

	memhit := memcache.Get(url)

	if memhit != nil {
		pf = memhit.Value().(dao.PluginFamilyDetail)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var fam tenable.FamilyPlugins
	if len(raw) == 0 {
		fam, raw, err = t.PluginFamily(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &fam)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal PluginFamily: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.

			fam, raw, err = t.PluginFamily(url)
			if err != nil {
				return
			}
		}
	}

	pf, err = a.Convert.PluginFamilyPlugin(fam)
	if err != nil {
		return
	}

	memcache.Set(url, pf, 60*time.Minute)
	disk.Store(filekey, raw, config.CachePretty)

	return
}
func (u *Unmarshal) Scanners(url string) (ss []dao.Scanner, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		ss = memhit.Value().([]dao.Scanner)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var tsl tenable.ScannerList

	if len(raw) == 0 {
		tsl, raw, err = a.Tenable.Scanners(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &tsl)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal Scanners: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			tsl, raw, err = a.Tenable.Scanners(url)

			if err != nil {
				return
			}
		}
	}

	ss, err = a.Convert.ScannerList(tsl)
	if err != nil {
		return
	}

	memcache.Set(url, ss, 60*time.Minute)
	disk.Store(filekey, raw, config.CachePretty)

	return
}
func (u *Unmarshal) ScannerAgents(s dao.Scanner, url string) (sa []dao.ScannerAgent, pag tenable.Pagination, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		sa = memhit.Value().([]dao.ScannerAgent)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var agents tenable.ScannerAgent

	if len(raw) == 0 {
		agents, raw, err = a.Tenable.ScannerAgents(url)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(raw, &agents)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScannerAgents: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			agents, raw, err = a.Tenable.ScannerAgents(url)

			if err != nil {
				return
			}
		}
	}

	sa, err = a.Convert.ScannerAgents(agents)
	if err != nil {
		return
	}

	memcache.Set(url, agents, 60*time.Minute)
	disk.Store(filekey, raw, config.CachePretty)

	pag = agents.Pagination

	return
}
func (u *Unmarshal) VulnSummary(ad dao.Asset, url string) (vs []dao.AssetVuln, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		vs = memhit.Value().([]dao.AssetVuln)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var asv tenable.AssetVuln

	if len(raw) == 0 {
		asv, raw, err = a.Tenable.AssetVuln(url)
		if err != nil {
			return
		}

	} else {
		err = json.Unmarshal(raw, &asv)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal AssetVuln: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			asv, raw, err = a.Tenable.AssetVuln(url)

			if err != nil {
				return
			}
		}
	}

	vs, err = a.Convert.AssetVuln(asv)
	if err != nil {
		return
	}

	memcache.Set(url, vs, 60*time.Minute)
	disk.Store(filekey, raw, config.CachePretty)

	return
}
func (u *Unmarshal) VulnDetail(vs dao.AssetVuln, url string) (vd dao.AssetVulnDetail, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		vd = memhit.Value().(dao.AssetVulnDetail)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var asvi tenable.AssetVulnInfo

	if len(raw) == 0 {
		asvi, raw, err = a.Tenable.AssetVulnInfo(url)
		if err != nil {
			return
		}

	} else {
		err = json.Unmarshal(raw, &asvi)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal AssetVulnInfo: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			asvi, raw, err = a.Tenable.AssetVulnInfo(url)

			if err != nil {
				return
			}
		}
	}

	vd, err = a.Convert.AssetVulnDetail(asvi)
	if err != nil {
		return
	}

	memcache.Set(url, vs, 60*time.Minute)
	disk.Store(filekey, raw, config.CachePretty)

	return
}
func (u *Unmarshal) VulnOutput(vs dao.AssetVuln, url string) (vo []dao.AssetVulnOutput, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		vo = memhit.Value().([]dao.AssetVulnOutput)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var avo tenable.AssetVulnOutput

	if len(raw) == 0 {
		avo, raw, err = a.Tenable.AssetVulnOutput(url)
		if err != nil {
			return
		}

	} else {
		err = json.Unmarshal(raw, &avo)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal AssetVulnOutput: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			avo, raw, err = a.Tenable.AssetVulnOutput(url)

			if err != nil {
				return
			}
		}
	}

	vo, err = a.Convert.AssetVulnOutput(avo)
	if err != nil {
		return
	}

	memcache.Set(url, vs, 60*time.Minute)
	disk.Store(filekey, raw, config.CachePretty)

	return
}

func (u *Unmarshal) ScannerAgentGroups(url string) (sag []dao.ScannerAgentGroup, err error) {
	a := u.Adapter

	var src tenable.ScannerAgentGroups
	src, _, err = a.Tenable.AgentGroups(url) // Not caching these calls
	if err != nil {
		return
	}

	sag, err = a.Convert.ScannerAgentGroups(src)

	return
}
func (u *Unmarshal) CreateScannerAgentGroup(url string, json string) (sag dao.ScannerAgentGroup, err error) {
	a := u.Adapter

	var src tenable.ScannerAgentGroup

	// We don't need the RAW because we don't cache POST
	src, _, err = a.Tenable.CreateAgentGroup(url, json)
	if err != nil {
		return
	}

	sag, err = a.Convert.ScannerAgentGroup(src)

	return
}
func (u *Unmarshal) AssignScannerAgentGroup(url string) (body []byte, err error) {
	a := u.Adapter
	body, err = a.Tenable.AssignAgentGroup(url)
	return
}

func (u *Unmarshal) AssetExport(url string, json string) (uuid string, err error) {
	a := u.Adapter

	var src tenable.AssetExport

	// We don't need the RAW because we don't cache POST
	src, _, err = a.Tenable.AssetExport(url, json)
	if err != nil {
		return
	}

	uuid = src.UUID

	return
}
func (u *Unmarshal) AssetExportStatus(url string) (status dao.AssetExportStatus, err error) {
	a := u.Adapter

	var src tenable.AssetExportStatus

	// TODO: Add filecache concepts here
	src, _, err = a.Tenable.AssetExportStatus(url)
	if err != nil {
		return
	}

	status.Status = src.Status
	for c := range src.Chunks {
		status.Chunks = append(status.Chunks, string(src.Chunks[c]))
	}

	if src.Status == "FINISHED" {
		// Write the cache
	}

	return
}
func (u *Unmarshal) AssetExportChunk(url string) (assets []dao.Asset, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		assets = memhit.Value().([]dao.Asset)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var src []tenable.AssetExportChunk
	if len(raw) == 0 {
		src, raw, err = a.Tenable.AssetExportChunk(url)
		if err != nil {
			return
		}
		disk.Store(filekey, raw, config.CachePretty)
	} else {
		// We don't need the RAW because we don't cache POST
		err = json.Unmarshal(raw, &src)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal ScannerAgents: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			src, raw, err = a.Tenable.AssetExportChunk(url)

			if err != nil {
				return
			}
			disk.Store(filekey, raw, config.CachePretty)
		}
	}

	assets, err = a.Convert.AssetExportChunk(src)
	if err != nil {
		return
	}

	memcache.Set(url, assets, 60*time.Minute)

	return
}

func (u *Unmarshal) VulnExport(url string, json string) (uuid string, err error) {
	a := u.Adapter

	var src tenable.VulnExport

	// We don't need the RAW because we don't cache POST
	src, _, err = a.Tenable.VulnExport(url, json)
	if err != nil {
		return
	}

	uuid = src.UUID

	return
}
func (u *Unmarshal) VulnExportStatus(url string) (status dao.VulnExportStatus, err error) {
	a := u.Adapter

	var src tenable.VulnExportStatus

	src, _, err = a.Tenable.VulnExportStatus(url)
	if err != nil {
		return
	}

	status.Status = src.Status
	for c := range src.Chunks {
		status.Chunks = append(status.Chunks, string(src.Chunks[c]))
	}

	if src.Status == "FINISHED" {
		// Write the cache
	}

	return
}
func (u *Unmarshal) VulnExportChunk(url string) (vulns []dao.VulnExportChunk, err error) {
	a := u.Adapter

	memcache := a.Memory.Cache
	disk := u.Adapter.Disk
	config := u.Adapter.Config

	memhit := memcache.Get(url)
	if memhit != nil {
		vulns = memhit.Value().([]dao.VulnExportChunk)
		return
	}

	filekey, keyerr := a.ToFilename(url)
	if keyerr != nil {
		err = keyerr
		return
	}

	raw, diskerr := disk.Fetch(filekey)
	if diskerr != nil && len(raw) > 0 && config.ClobberCacheMode == true {
		// Cache entry doesn't match cacheKey, delete entry.
		disk.Clear(filekey)
		raw = nil
	} else if diskerr != nil {
		err = diskerr
		return
	}

	var src []tenable.VulnExportChunk
	if len(raw) == 0 {
		src, raw, err = a.Tenable.VulnExportChunk(url)
		if err != nil {
			return
		}
		disk.Store(filekey, raw, config.CachePretty)
	} else {

		err = json.Unmarshal(raw, &src)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal VulnExportChunk: invalid disk entry: %s: %s", filekey, err)) // Don't overwrite bad disk file.
				return
			}
			disk.Clear(filekey) // Delete bad disk entry - cacheKey likely changed and we're clobbering.
			src, raw, err = a.Tenable.VulnExportChunk(url)

			if err != nil {
				return
			}
			disk.Store(filekey, raw, config.CachePretty)
		}
	}

	vulns, err = a.Convert.VulnExportChunk(src)
	if err != nil {
		return
	}

	memcache.Set(url, vulns, 60*time.Minute)

	return
}
