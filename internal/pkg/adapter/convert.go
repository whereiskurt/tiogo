package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Convert struct {
	Config           *app.Config
	Adapter          *Adapter
	IgnoreScanId     map[string]bool
	IncludeScanId    map[string]bool
	IgnorePluginId   map[string]bool
	IncludePluginId  map[string]bool
	IgnoreHistoryId  map[string]bool
	IncludeHistoryId map[string]bool
	IgnoreAssetId    map[string]bool
	IncludeAssetId   map[string]bool
	IgnoreHostId     map[string]bool
	IncludeHostId    map[string]bool

	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})
}

func NewConvert(a *Adapter) (c *Convert) {
	c = new(Convert)
	c.Config = a.Config
	c.Adapter = a

	c.Errorf = c.Config.Logger.Errorf
	c.Debugf = c.Config.Logger.Debugf
	c.Warnf = c.Config.Logger.Warnf
	c.Infof = c.Config.Logger.Infof

	c.InitSkips()
	return
}

func (c *Convert) InitSkips() {
	config := c.Config.VM
	c.IncludeScanId = make(map[string]bool)
	c.IgnoreScanId = make(map[string]bool)
	c.IncludeHistoryId = make(map[string]bool)
	c.IgnoreHistoryId = make(map[string]bool)
	c.IncludeAssetId = make(map[string]bool)
	c.IgnoreAssetId = make(map[string]bool)
	c.IncludeHostId = make(map[string]bool)
	c.IgnoreHostId = make(map[string]bool)
	c.IncludePluginId = make(map[string]bool)
	c.IgnorePluginId = make(map[string]bool)
	for _, id := range strings.Split(config.ScanID, ",") {
		if id != "" {
			c.IncludeScanId[id] = true
		}
	}
	for _, id := range strings.Split(config.IgnoreScanID, ",") {
		if id != "" {
			c.IgnoreScanId[id] = true
		}
	}
	for _, id := range strings.Split(config.HistoryID, ",") {
		if id != "" {
			c.IncludeHistoryId[id] = true
		}
	}
	for _, id := range strings.Split(config.IgnoreHistoryID, ",") {
		if id != "" {
			c.IgnoreHistoryId[id] = true
		}
	}
	for _, id := range strings.Split(config.AssetUUID, ",") {
		if id != "" {
			c.IncludeAssetId[id] = true
		}
	}
	for _, id := range strings.Split(config.IgnoreAssetUUID, ",") {
		if id != "" {
			c.IgnoreAssetId[id] = true
		}
	}
	for _, id := range strings.Split(config.HostID, ",") {
		if id != "" {
			c.IncludeHostId[id] = true
		}
	}
	for _, id := range strings.Split(config.IgnoreHostID, ",") {
		if id != "" {
			c.IgnoreHostId[id] = true
		}
	}
	for _, id := range strings.Split(config.IgnorePluginID, ",") {
		if id != "" {
			c.IgnorePluginId[id] = true
		}
	}
	for _, id := range strings.Split(config.PluginID, ",") {
		if id != "" {
			c.IncludePluginId[id] = true
		}
	}
}

func (c *Convert) ScanList(t tenable.ScanList) (scans []dao.Scan, err error) {
	for _, s := range t.Scans {

		if c.SkipScan(string(s.Id)) {
			continue
		}

		scanId := string(s.Id)
		scan := new(dao.Scan)
		scan.ScanID = scanId
		scan.UUID = s.UUID
		scan.Name = s.Name
		scan.Status = s.Status
		scan.Owner = s.Owner
		scan.UserPermissions = string(s.UserPermissions)
		scan.Enabled = fmt.Sprintf("%v", s.Enabled)
		scan.RRules = s.RRules
		scan.Timezone = s.Timezone
		scan.StartTime = s.StartTime
		scan.CreationDate = fmt.Sprintf("%v", c.UnixDate(s.CreationDate))
		scan.LastModifiedDate = fmt.Sprintf("%v", c.UnixDate(s.LastModifiedDate))

		scan.Timestamp = string(t.Timestamp)

		scans = append(scans, *scan)
	}

	return
}
func (c *Convert) ScanHistory(s dao.Scan, hist []tenable.ScanDetailHistory, depth int) (sh dao.ScanHistory, err error) {
	sh.Scan = s

	if len(hist) == 0 && depth > 0 {
		c.Adapter.Log.Warnf(fmt.Sprintf("no history for scan \"%s\" (id:%s) - has never run", s.Name, s.ScanID))
		return
	}

	// Get the ScanHistoryDetail for ecah past scan, with-in our depth.
	for i := range hist {
		h := hist[i]
		if i >= depth {
			break
		}
		var lkp dao.ScanHistory
		lkp.Scan = s
		lkp.History = []dao.ScanHistoryDetail{{HistoryID: string(h.HistoryID)}}

		var shd dao.ScanHistoryDetail
		shd, err = c.Adapter.ScanHistoryDetail(lkp, i)
		if err != nil {
			return
		}
		// Add the detail for the history details list at an offset less than depth
		sh.History = append(sh.History, shd)
	}

	return
}

var reStripNewLines = regexp.MustCompile(`\r?\n`)

func (c *Convert) ScanHistoryDetail(h dao.ScanHistory, sd tenable.ScanDetail, offset int) (shd dao.ScanHistoryDetail, err error) {

	start, end := c.UnixStartEnd(sd.Info.Start, sd.Info.End)

	shd.Scan = h.Scan

	shd.HostCount = string(sd.Info.HostCount)

	shd.PolicyName = sd.Info.PolicyName
	shd.Owner = sd.Info.Owner
	shd.ScannerName = sd.Info.ScannerName

	shd.Targets = reStripNewLines.ReplaceAllString(sd.Info.Targets, ",")

	shd.HistoryID = string(sd.History[offset].HistoryID)

	shd.CreationDate = string(sd.History[offset].CreationDate)

	// shd.LastModifiedDate = string(sd.History[offset].LastModifiedDate)
	shd.LastModifiedDate = fmt.Sprintf("%v", c.UnixDate(sd.History[offset].LastModifiedDate))

	shd.Status = sd.History[offset].Status
	shd.ScanStart = fmt.Sprintf("%v", start)
	shd.ScanStartUnix = fmt.Sprintf("%s", string(sd.Info.Start))
	shd.ScanEnd = fmt.Sprintf("%v", end)
	shd.ScanEndUnix = fmt.Sprintf("%s", string(sd.Info.End))
	shd.ScanDuration = fmt.Sprintf("%v", end.Sub(start))

	for k := range sd.Info.AgentTarget {
		t := sd.Info.AgentTarget[k]
		at := dao.ScannerAgentGroup{Name: t.Name, ID: string(t.ID), UUID: t.UUID}
		shd.AgentGroup = append(shd.AgentGroup, at)
	}

	shd.AgentCount = string(sd.Info.AgentCount)
	shd.ScanType = sd.Info.ScanType
	shd.Host = c.HostScanSummary(&shd, &sd)
	shd.HostPlugin = c.PluginSummaryDetail(&shd, &sd)
	shd.HostAssetMap, err = c.HostAsset(&shd, &sd)

	return
}

func (c *Convert) HostScanSummary(shd *dao.ScanHistoryDetail, sd *tenable.ScanDetail) (ss map[string]dao.HostScanSummary) {
	ss = make(map[string]dao.HostScanSummary)

	for _, h := range sd.Hosts {
		if c.SkipHost(string(h.ID)) {
			continue
		}

		crit, _ := strconv.Atoi(string(h.SeverityCritical))
		high, _ := strconv.Atoi(string(h.SeverityHigh))
		med, _ := strconv.Atoi(string(h.SeverityMedium))
		low, _ := strconv.Atoi(string(h.SeverityLow))
		hcrit, _ := strconv.Atoi(shd.PluginCriticalCount)
		hhigh, _ := strconv.Atoi(shd.PluginHighCount)
		hmed, _ := strconv.Atoi(shd.PluginMediumCount)
		hlow, _ := strconv.Atoi(shd.PluginLowCount)

		var s dao.HostScanSummary
		s.HostID = string(h.ID)
		s.AssetID = string(h.AssetID)
		s.HostnameOrIP = h.HostnameOrIP
		s.Progress = h.Progress
		s.ChecksConsidered = string(h.ChecksConsidered)
		s.ChecksTotal = string(h.ChecksTotal)
		s.ScanProgressCurrent = string(h.ProgressCurrent)
		s.ScanProgressTotal = string(h.ProgressTotal)
		s.PluginCriticalCount = fmt.Sprintf("%v", crit)
		s.PluginHighCount = fmt.Sprintf("%v", high)
		s.PluginMediumCount = fmt.Sprintf("%v", med)
		s.PluginLowCount = fmt.Sprintf("%v", low)
		s.PluginTotalCount = fmt.Sprintf("%v", low+med+high+crit)
		s.Score = string(h.Score)

		shd.PluginCriticalCount = fmt.Sprintf("%v", hcrit+crit)
		shd.PluginHighCount = fmt.Sprintf("%v", hhigh+high)
		shd.PluginMediumCount = fmt.Sprintf("%v", hmed+med)
		shd.PluginLowCount = fmt.Sprintf("%v", hlow+low)
		shd.PluginTotalCount = fmt.Sprintf("%v", hlow+low+hmed+med+hhigh+high+hcrit+crit)
		s.ScanHistoryDetail = *shd

		ss[s.HostID] = s
	}

	return
}
func (c *Convert) PluginSummaryDetail(shd *dao.ScanHistoryDetail, sd *tenable.ScanDetail) (pp map[string]dao.Plugin) {
	pp = make(map[string]dao.Plugin)

	for _, v := range sd.Vulnerabilities {
		var p dao.Plugin

		if c.SkipPlugin(string(v.PluginID)) {
			continue
		}

		p.PluginID = string(v.PluginID)
		p.Name = v.Name
		p.FamilyName = v.Family
		p.Count = string(v.Count)
		p.Severity = string(v.Severity)

		pp[p.PluginID] = p
	}

	return
}
func (c *Convert) HostAsset(shd *dao.ScanHistoryDetail, sd *tenable.ScanDetail) (ss map[string]string, err error) {
	a := c.Adapter
	t := a.Tenable
	memcache := a.Memory.Cache
	cache := a.Disk
	config := c.Config

	ss = make(map[string]string)
	var p = map[string]string{
		"baseURL":   a.Config.BaseURL,
		"scanID":    shd.Scan.ScanID,
		"historyID": shd.HistoryID,
	}
	url, terr := a.ToURL("AssetMap", p)
	if terr != nil {
		err = terr
		return
	}

	memhit := memcache.Get(url)
	if memhit != nil {
		ss = memhit.Value().(map[string]string)
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

	var ah tenable.AssetHost
	if len(raw) == 0 {
		ah, raw, err = t.AssetHostMap(url)
		if err != nil {
			return
		}
		a.Disk.Store(filekey, raw, a.Config.CachePretty)

	} else {
		err = json.Unmarshal(raw, &ah)
		if err != nil {
			if config.ClobberCacheMode == false {
				err = errors.New(fmt.Sprintf("noclobber: couldn't unmarshal HostAsset: invalid cache entry: %s: %s", filekey, err)) // Don't overwrite bad cache file.
				return
			}
			cache.Clear(filekey) // Delete bad cache entry - cacheKey likely changed and we're clobbering.

			ah, raw, err = t.AssetHostMap(url)
			if err != nil {
				return
			}
		}
	}

	for _, value := range ah.Assets {
		ss[string(value.HostID)] = value.UUID
	}
	a.Memory.Cache.Set(url, ss, time.Minute*60)

	return
}

func (c *Convert) HostDetail(hss dao.HostScanSummary, hd *tenable.HostDetail) (hsd dao.HostScanDetail, err error) {
	tz := c.GetScannerTZ(hss.ScanHistoryDetail.ScannerName)

	histID := hss.ScanHistoryDetail.HistoryID
	sID := hss.ScanHistoryDetail.Scan.ScanID
	hID := hss.HostID

	hsd.IP = hd.Info.HostIP
	hsd.FQDN = hd.Info.FQDN
	hsd.NetBIOS = hd.Info.NetBIOS
	hsd.MACAddresses = strings.Replace(hd.Info.MACAddress, "\n", ",", -1)
	hsd.OperatingSystems = strings.Join(hd.Info.OperatingSystem, ",")

	start, tmStart, err := c.FromNoTZ(string(hd.Info.HostStart), tz)
	hsd.ScanStartUnix = fmt.Sprintf("%v", tmStart.In(time.Local).Unix())
	hsd.ScanStart = start
	if err != nil {
		c.Warnf("(scanID:%s:histID:%s:hostID:%s) %s", sID, histID, hID, err)
		err = nil
	}

	end, tmEnd, err := c.FromNoTZ(string(hd.Info.HostEnd), tz)
	hsd.ScanEndUnix = fmt.Sprintf("%v", tmEnd.In(time.Local).Unix())
	hsd.ScanEnd = end
	if err != nil {
		c.Warnf("(scanID:%s:histID:%s:hostID:%s) %s ", sID, histID, hID, err)
		err = nil
	}

	hsd.ScanDuration = fmt.Sprintf("%v", tmEnd.Sub(tmStart))

	hsd.PluginMap = make(map[string]dao.Plugin)

	for _, v := range hd.Vulnerabilities {
		var p dao.Plugin
		p.PluginID = string(v.PluginId)
		p.Name = v.PluginName
		p.FamilyName = v.PluginFamily
		p.Count = string(v.Count)
		p.Severity = string(v.Severity)
		// p.Detail GO GET IT!!!
		hsd.PluginMap[p.PluginID] = p
	}
	return
}
func (c *Convert) PluginDetail(pds dao.Plugin, p tenable.Plugin) (pd dao.PluginDetail, err error) {

	pd.Attribute = make(map[string]dao.PluginDetailAttribute)

	for _, v := range p.Attributes {
		var a dao.PluginDetailAttribute
		a.Name = v.Name
		a.Value = v.Value

		pd.Attribute[v.Name] = a

		if a.Name == "fname" {
			pd.FunctionName = a.Value
		} else if a.Name == "patch_publication_date" {
			pd.PatchPublicationDate = a.Value
		} else if a.Name == "risk_factor" {
			pds.Severity = a.Value
		} else if a.Name == "plugin_name" {
			pds.Name = a.Value
		} else if a.Name == "plugin_publication_date" {
			pd.PluginPublicationDate = a.Value
		}
	}
	pds.Count = "1"

	return
}
func (c *Convert) AssetDetail(ai tenable.AssetInfo) (ad dao.Asset, err error) {

	ad.UUID = ai.UUID
	if ad.UUID == "" {
		ad.UUID = ai.ID
	}

	ad.ID = ai.ID
	ad.OperatingSystem = ai.OperatingSystem
	ad.HasAgent = ai.HasAgent
	ad.CreatedAt = ai.CreatedAt
	ad.UpdatedAt = ai.UpdatedAt
	ad.FirstSeenAt = ai.FirstSeenAt
	ad.LastSeenAt = ai.LastSeenAt
	ad.LastAuthenticatedScanAt = ai.LastAuthenticatedScanAt
	ad.LastLicensedScanAt = ai.LastLicensedScanAt
	ad.IPV4 = ai.IPV4
	ad.IPV6 = ai.IPV6
	ad.FQDN = ai.FQDN
	ad.MACAddress = ai.MACAddress
	ad.NetBIOS = ai.NetBIOS
	ad.SystemType = ai.SystemType
	ad.TenableUUID = ai.TenableUUID
	ad.HostName = ai.HostName
	ad.AgentName = ai.AgentName
	ad.BIOSUUID = ai.BIOSUUID
	ad.AWSEC2InstanceID = ai.AWSEC2InstanceId
	ad.AWSEC2InstanceAMIID = ai.AWSEC2InstanceAMIId
	ad.AWSOwnerID = ai.AWSOwnerId
	ad.AWSAvailabilityZone = ai.AWSAvailabilityZone
	ad.AWSRegion = ai.AWSRegion
	ad.AWSVPCID = ai.AWSVPCID
	ad.AWSEC2InstanceGroupName = ai.AWSEC2InstanceGroupName
	ad.AWSEC2InstanceStateName = ai.AWSEC2InstanceStateName
	ad.AWSEC2InstanceType = ai.AWSEC2InstanceType
	ad.AWSSubnetID = ai.AWSSubnetId
	ad.AWSEC2ProductCode = ai.AWSEC2ProductCode
	ad.AWSEC2Name = ai.AWSEC2Name
	ad.AzureVMID = ai.AzureVMId
	ad.AzureResourceID = ai.AzureResourceId
	ad.SSHFingerPrint = ai.SSHFingerPrint
	ad.McafeeEPOGUID = ai.McafeeEPOGUID
	ad.McafeeEPOAgentGUID = ai.McafeeEPOAgentGUID
	ad.QualysHostID = ai.QualysHostId
	ad.QualysAssetID = ai.QualysAssetId
	ad.ServiceNowSystemID = ai.ServiceNowSystemId

	for _, t := range ai.Tags {
		var tag dao.AssetTagDetail
		tag.UUID = t.UUID
		tag.CategoryName = t.CategoryName
		tag.Value = t.Value
		tag.AddedBy = t.AddedBy
		tag.AddedAt = t.AddedAt
		tag.Source = t.Source

		ad.Tags = append(ad.Tags, tag)
	}

	for _, v := range ai.Counts.Vulnerabilities.Severities {
		sev := dao.AssetSev{Name: v.Name, Count: string(v.Count), Level: string(v.Level)}
		ad.VulnSevCount.Severity = append(ad.VulnSevCount.Severity, sev)
	}
	ad.VulnSevCount.Total = string(ai.Counts.Vulnerabilities.Total)

	for _, a := range ai.Counts.Audits.Severities {
		sev := dao.AssetSev{Name: a.Name, Count: string(a.Count), Level: string(a.Level)}
		ad.AuditSevCount.Severity = append(ad.AuditSevCount.Severity, sev)
	}
	ad.AuditSevCount.Total = string(ai.Counts.Audits.Total)

	for _, i := range ai.Interfaces {
		inf := dao.AssetInterface{Name: i.Name, FQDN: i.FQDN, IPV4: i.IPV4, IPV6: i.IPV6, MACAddress: i.MACAddress}
		ad.Interface = append(ad.Interface, inf)
	}

	for _, s := range ai.Sources {
		src := dao.AssetSource{Name: s.Name, FirstSeenAt: s.FirstSeenAt, LastSeenAt: s.LastSeenAt}
		ad.Source = append(ad.Source, src)
	}

	return
}

func (c *Convert) PluginFamily(plugins tenable.PluginFamilies) (pf []dao.PluginFamily, err error) {
	for _, v := range plugins.Families {
		pf = append(pf, dao.PluginFamily{ID: string(v.Id), Name: v.Name, Count: string(v.Count)})
	}
	return
}
func (c *Convert) PluginFamilyPlugin(p tenable.FamilyPlugins) (pfp dao.PluginFamilyDetail, err error) {

	pfp.ID = string(p.ID)
	pfp.Name = p.Name

	for _, v := range p.Plugins {
		id := string(v.ID)
		if c.SkipPlugin(id) {
			continue
		}

		pfp.PluginID = append(pfp.PluginID, id)
	}

	return
}

func (c *Convert) GetScannerTZ(name string) (tz string) {
	tz = c.Config.DefaultTZ
	return
}

func (c *Convert) SkipScan(scanId string) (skip bool) {
	if _, ignore := c.IgnoreScanId[scanId]; ignore {
		skip = true
	}
	if len(c.IncludeScanId) > 0 {
		if include := c.IncludeScanId[scanId]; !include {
			skip = true
		}
	}
	return
}
func (c *Convert) SkipPlugin(pluginId string) (skip bool) {
	if _, ignore := c.IgnorePluginId[pluginId]; ignore {
		skip = true
	}
	if len(c.IncludePluginId) > 0 {
		if include := c.IncludePluginId[pluginId]; !include {
			skip = true
		}
	}
	return
}
func (c *Convert) SkipHistory(hID string) (skip bool) {
	if _, ignore := c.IgnoreHistoryId[hID]; ignore {
		skip = true
	}
	if len(c.IncludeHistoryId) > 0 {
		if include := c.IncludeHistoryId[hID]; !include {
			skip = true
		}
	}
	return
}
func (c *Convert) SkipHost(hID string) (skip bool) {
	if _, ignore := c.IgnoreHostId[hID]; ignore {
		skip = true
	}
	if len(c.IncludeHostId) > 0 {
		if include := c.IncludeHostId[hID]; !include {
			skip = true
		}
	}
	return
}

func (c *Convert) SkipRegex(pattern string, value string) (skip bool) {

	if pattern == "" {
		return
	}

	match, _ := regexp.MatchString(pattern, value)

	return !match
}

var TimeFormatNoTZ = "Mon Jan _2 15:04:05 2006"
var TimeFormatTZ = "2006-01-_2 15:04:05 -0700 MST"

func (c *Convert) UnixStartEnd(start json.Number, end json.Number) (time.Time, time.Time) {
	return c.UnixDate(start), c.UnixDate(end)
}
func (c *Convert) UnixDate(unix json.Number) (t time.Time) {
	rawScanStart, errParseStart := strconv.ParseInt(string(unix), 10, 64)
	if errParseStart != nil {
		rawScanStart = int64(0)
	}
	t = time.Unix(rawScanStart, 0)
	return
}

func (c *Convert) FromNoTZ(dts string, setTZ string) (withTZ string, unix time.Time, err error) {
	dtsInt, err := strconv.ParseInt(dts, 10, 64)
	if err == nil {
		unix = time.Unix(dtsInt, 0)
	} else {
		unix, err = time.Parse(TimeFormatNoTZ, dts)
		if err != nil {
			err = errors.New(fmt.Sprintf("parse error: dts '%v' with TZ '%s': %s", dts, setTZ, err))
			return
		}
	}

	// Render UNIX time (which is UTC 0) and add timezone from scanner.
	// The scanner captured the UNIX time, without a TZ.
	withTZ = strings.Replace(fmt.Sprintf("%v", unix), "+0000 UTC", setTZ, -1)
	tmTZ, err := time.Parse(TimeFormatTZ, withTZ)

	if err != nil {
		err = errors.New(fmt.Sprintf("parse error: date '%v' with TZ '%s': %s", unix, setTZ, err))
		return
	}
	withTZ = fmt.Sprintf("%v", tmTZ)
	return
}

func (c *Convert) ScannerList(sl tenable.ScannerList) (ss []dao.Scanner, err error) {
	for _, x := range sl.Scanners {
		s := dao.Scanner{}
		s.ID = string(x.ID)
		s.Name = x.Name
		s.UUID = x.UUID
		s.Type = x.Type
		s.Owner = x.Owner
		s.EngineVersion = x.EngineVersion
		s.Key = x.Key
		s.LoadedPluginSet = x.LoadedPluginSet
		s.Platform = x.Platform
		s.RegistrationCode = x.RegistrationCode
		s.ScanCount = string(x.ScanCount)
		s.Status = x.Status

		lic := dao.ScannerLicense{}
		lic.Scanners = s.License.Scanners
		lic.Type = s.License.Type
		lic.Agents = s.License.Agents
		lic.IPS = s.License.IPS

		s.License = lic

		ss = append(ss, s)
	}
	return
}

func (c *Convert) ScannerAgents(sa tenable.ScannerAgent) (ss []dao.ScannerAgent, err error) {
	regex := c.Config.VM.Regex

	for _, a := range sa.Agents {
		if c.SkipRegex(regex, a.Name+","+a.IP) {
			continue
		}

		agent := dao.ScannerAgent{
			ID:          string(a.ID),
			Name:        a.Name,
			UUID:        a.UUID,
			Status:      a.Status,
			Platform:    a.Platform,
			CoreBuild:   a.CoreBuild,
			CoreVersion: a.CoreVersion,
			Distro:      a.Distro,
			Feed:        a.Feed,
			IP:          a.IP,
			LastConnect: c.UnixDate(a.LastConnect),
			LastScanned: c.UnixDate(a.LastScanned),
			LinkedOn:    c.UnixDate(a.LinkedOn),
		}

		agent.Groups = make(map[string]dao.ScannerAgentGroup)
		for _, group := range a.Groups {
			if len(group.Name) == 0 {
				continue
			}

			g := dao.ScannerAgentGroup{ID: string(group.ID), Name: group.Name}
			agent.Groups[group.Name] = g
		}

		ss = append(ss, agent)
	}

	return
}

func (c *Convert) AssetVuln(av tenable.AssetVuln) (vv []dao.AssetVuln, err error) {
	for _, vuln := range av.Vulnerabilities {
		if c.SkipPlugin(string(vuln.PluginID)) {
			continue
		}

		v := dao.AssetVuln{
			PluginID:     string(vuln.PluginID),
			PluginName:   vuln.PluginName,
			PluginFamily: vuln.PluginFamily,
			Count:        string(vuln.Count),
			Severity:     string(vuln.Severity),
			State:        string(vuln.State),
		}
		vv = append(vv, v)
	}
	return
}
func (c *Convert) AssetVulnDetail(avi tenable.AssetVulnInfo) (v dao.AssetVulnDetail, err error) {
	v.Severity = string(avi.Info.Severity)
	v.Count = string(avi.Info.Count)
	v.Description = avi.Info.Description
	v.PluginDetails = avi.Info.PluginDetails
	v.ReferenceInfo = avi.Info.ReferenceInfo
	v.Discovery = avi.Info.Discovery
	v.RiskInfo = avi.Info.RiskInfo
	v.SeeAlso = avi.Info.SeeAlso
	v.Solution = avi.Info.Solution
	v.VulnInfo = avi.Info.VulnInfo
	v.Synopsis = avi.Info.Synopsis
	return
}
func (c *Convert) AssetVulnOutput(avo tenable.AssetVulnOutput) (v []dao.AssetVulnOutput, err error) {

	return
}

func (c *Convert) ScannerAgentGroups(groups tenable.ScannerAgentGroups) (sag []dao.ScannerAgentGroup, err error) {
	for _, g := range groups.Groups {
		s := dao.ScannerAgentGroup{
			ID:   string(g.ID),
			Name: g.Name,
			UUID: g.UUID,
		}
		sag = append(sag, s)
	}
	return
}
func (c *Convert) ScannerAgentGroup(group tenable.ScannerAgentGroup) (sag dao.ScannerAgentGroup, err error) {
	sag.ID = string(group.ID)
	sag.Name = group.Name
	sag.UUID = group.UUID
	return
}

func (c *Convert) AssetExportChunk(assets []tenable.AssetExportChunk) (aa []dao.Asset, err error) {
	for _, asset := range assets {
		a := dao.Asset{
			HostName:                asset.HostName,
			FQDN:                    asset.FQDN,
			IPV4:                    asset.IPV4,
			IPV6:                    asset.IPV6,
			OperatingSystem:         asset.OperatingSystem,
			UUID:                    asset.UUID,
			LastAuthenticatedScanAt: asset.LastAuthenticatedScanAt,
			LastLicensedScanAt:      asset.LastLicensedScanAt,
			FirstSeenAt:             asset.FirstSeenAt,
			LastSeenAt:              asset.LastSeenAt,
			AgentName:               asset.AgentNames,
			CreatedAt:               asset.CreatedAt,
			UpdatedAt:               asset.UpdatedAt,
			HasAgent:                asset.HasAgent,
			MACAddress:              asset.MACAddress,
			BIOSUUID:                []string{asset.BIOSUUID},
			NetBIOS:                 asset.NetBIOS,
			SSHFingerPrint:          asset.SSHFingerPrint,
			SystemType:              asset.SystemType,
			TenableUUID:             []string{asset.AgentUUID},
		}
		aa = append(aa, a)
	}
	return
}
func (c *Convert) VulnExportChunk(vulns []tenable.VulnExportChunk) (vv []dao.VulnExportChunk, err error) {
	for _, vuln := range vulns {

		asset := dao.VulnExportChunkAsset{
			UUID:                     vuln.Asset.UUID,
			OperatingSystem:          vuln.Asset.OperatingSystem,
			IPV4:                     vuln.Asset.IPV4,
			FQDN:                     vuln.Asset.FQDN,
			DeviceType:               vuln.Asset.DeviceType,
			HostName:                 vuln.Asset.HostName,
			LastAuthenticatedResults: vuln.Asset.LastAuthenticatedResults,
			NetBIOSWorkgroup:         vuln.Asset.NETBIOSWorkgroup,
			Tracked:                  vuln.Asset.Tracked,
		}
		plugin := dao.VulnExportChunkPlugin{
			Name:             vuln.Plugin.Name,
			ID:               string(vuln.Plugin.PluginID),
			Synopsis:         vuln.Plugin.Synopsis,
			Solution:         vuln.Plugin.Solution,
			Description:      vuln.Plugin.Description,
			Type:             vuln.Plugin.Type,
			Family:           vuln.Plugin.Family,
			FamilyID:         string(vuln.Plugin.FamilyID),
			HasPatch:         vuln.Plugin.HasPatch,
			ModificationDate: vuln.Plugin.ModificationDate,
			PublicationDate:  vuln.Plugin.PublicationDate,
			RiskFactor:       vuln.Plugin.RiskFactor,
			Version:          vuln.Plugin.Version,
		}
		scan := dao.VulnExportChunkScan{
			UUID:         vuln.Scan.UUID,
			CompletedAt:  vuln.Scan.CompletedAt,
			ScheduleUUID: vuln.Scan.ScheduleUUID,
			StartedAt:    vuln.Scan.StartedAt,
		}
		port := dao.VulnExportChunkPort{
			Port:     string(vuln.Port.Port),
			Service:  vuln.Port.Service,
			Protocol: vuln.Port.Protocol,
		}

		v := dao.VulnExportChunk{
			Asset:    asset,
			Output:   vuln.Output,
			Severity: vuln.Severity,
			Plugin:   plugin,
			Scan:     scan,
			Network:  port,
		}

		vv = append(vv, v)
	}

	return
}
