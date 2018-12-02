package tenable

import (
	"encoding/json"
	"errors"
	"github.com/whereiskurt/tiogo/internal/app"
	"gopkg.in/matryer/try.v1"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var DefaultMaxRety = []int{0, 500, 1000, 2000}

type Service struct {
	BaseUrl   string
	SecretKey string
	AccessKey string
	MaxRetry  []int
	Worker    *sync.WaitGroup
	Infof     func(fmt string, args ...interface{})
	Debugf    func(fmt string, args ...interface{})
	Warnf     func(fmt string, args ...interface{})
	Errorf    func(fmt string, args ...interface{})
}

func NewService(base string, secret string, access string, log *app.Logger) (s *Service) {
	s = new(Service)
	s.BaseUrl = strings.TrimSuffix(base, "/")
	s.SecretKey = secret
	s.AccessKey = access
	s.Infof = log.Infof
	s.Debugf = log.Debugf
	s.Warnf = log.Warnf
	s.Errorf = log.Errorf
	s.MaxRetry = DefaultMaxRety
	return
}

func (s *Service) SleepAndRerun(attempt int, name string) (rerun bool) {
	if attempt < len(s.MaxRetry) {
		s.Warnf("will retry %dx  more time(s), sleeping %dms in %s", len(s.MaxRetry)-attempt, s.MaxRetry[attempt], name)
		time.Sleep(time.Duration(s.MaxRetry[attempt]) * time.Millisecond)
		rerun = true
	}
	return
}

func (s *Service) ScanList(url string) (scans ScanList, raw []byte, err error) {
	tenable := NewPortal(s)
	raw, err = tenable.GET(url)
	if err == nil {
		err = json.Unmarshal(raw, &scans)
	}
	return
}
func (s *Service) ScanDetail(url string) (sd ScanDetail, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for ScanDetails: %s", url, err)
			rerun = s.SleepAndRerun(attempt, "ScanDetails")
			return
		}

		err = json.Unmarshal(raw, &sd)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL ScanDetails: %s:%s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "ScanDetails")
			return
		}

		// Sort histories by creation date DESC, to get offset history_id
		sort.Slice(sd.History, func(i, j int) bool {
			iv, _ := strconv.ParseInt(string(sd.History[i].CreationDate), 10, 64)
			jv, _ := strconv.ParseInt(string(sd.History[j].CreationDate), 10, 64)
			return iv > jv
		})

		return
	})

	if err != nil {
		s.Errorf("%s failed ScanDetails: %v", err)
		return
	}

	return
}
func (s *Service) HostDetail(url string) (hd HostDetail, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for HostDetails: %s", url, err)
			rerun = s.SleepAndRerun(attempt, "HostDetails")
			return
		}

		err = json.Unmarshal(raw, &hd)
		// If it failed it might be in the 'old' format.
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL HostDetails: %s: %s [len:%d]", url, err, string(raw), len(raw))

			var legacy HostDetailLegacyV2
			var start time.Time
			var end time.Time
			err = json.Unmarshal([]byte(string(raw)), &legacy)
			if err != nil {
				// Potentially third format?! Consider fatals here.
				s.Warnf("'%s' failed to UNMARSHAL HostDetailsv2: %s :%s [len:%d]", url, err, string(raw), len(raw))
				rerun = s.SleepAndRerun(attempt, "HostDetails")
				return
			}
			start, err = time.Parse(time.ANSIC, legacy.Info.HostStart)
			if err != nil {
				s.Warnf("'%s' failed to UNMARSHAL HostDetailsv2: %s :%s [len:%d]", url, err, string(raw), len(raw))
				rerun = s.SleepAndRerun(attempt, "HostDetails")
				return
			}
			end, err = time.Parse(time.ANSIC, legacy.Info.HostEnd)
			if err != nil {
				s.Warnf("'%s' failed to UNMARSHAL HostDetailsv2: %s :%s [len:%d]", url, err, string(raw), len(raw))
				rerun = s.SleepAndRerun(attempt, "HostDetails")
				return
			}
			// Copy legacy into hd
			hd.Info.OperatingSystem = append(hd.Info.OperatingSystem, legacy.Info.OperatingSystem)
			hd.Info.FQDN = legacy.Info.FQDN
			hd.Info.NetBIOS = legacy.Info.NetBIOS
			hd.Vulnerabilities = legacy.Vulnerabilities
			hd.Info.HostStart = json.Number(start.Format(time.ANSIC))
			hd.Info.HostEnd = json.Number(end.Format(time.ANSIC))
			// Rewrite raw with the new JSON
			raw, err = json.Marshal(hd)
			if err != nil {
				s.Warnf("'%s' failed to convert HostDetailsv2 to HostDetails UNMARSHAL: %s: %s [len:%d]", url, err, string(raw), len(raw))
				rerun = s.SleepAndRerun(attempt, "HostDetails")
				return
			}
		}
		return
	})

	if err != nil {
		s.Errorf("%s failed HostDetails: %v", err)
		return
	}

	return
}
func (s *Service) AssetHostMap(url string) (ah AssetHost, raw []byte, err error) {
	tenable := NewPortal(s)
	raw, err = tenable.GET(url)
	if err == nil {
		err = json.Unmarshal(raw, &ah)
	}
	return
}
func (s *Service) Asset(url string) (ai Asset, raw []byte, err error) {
	tenable := NewPortal(s)
	err = try.Do(func(attempt int) (bool, error) {
		raw, err = tenable.GET(url)
		if err != nil || strings.Contains(string(raw), `"error": "Asset`) {
			s.Warnf("'%s' failed to GET for Asset: %s", url, err)
			rr := s.SleepAndRerun(attempt, "Asset")
			return rr, err
		}

		err = json.Unmarshal(raw, &ai)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL Asset: %s: %s:%d", url, err, string(raw), len(raw))
			rr := s.SleepAndRerun(attempt, "Asset")
			return rr, err
		}

		sort.Slice(ai.Info.Tags, func(i, j int) bool {
			if ai.Info.Tags[i].CategoryName == ai.Info.Tags[j].CategoryName {
				return ai.Info.Tags[i].Value < ai.Info.Tags[j].Value
			}
			return ai.Info.Tags[i].CategoryName < ai.Info.Tags[j].CategoryName
		})

		return attempt < len(s.MaxRetry), err

	})

	if err != nil {
		s.Errorf("%s failed Asset: %v", err)
		return
	}

	return
}
func (s *Service) Plugin(url string) (p Plugin, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for Plugin : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "Plugin")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &p)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for Plugin: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "Plugin")
			return
		}

		return rerun, err
	})

	if err != nil {
		s.Errorf("%s permanently failed Plugin: %v", url, err)
		return
	}

	return
}
func (s *Service) PluginFamiles(url string) (ff PluginFamilies, raw []byte, err error) {

	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for PluginFamilies : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "PluginFamilies")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &ff)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for PluginFamilies: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "PluginFamilies")
			return
		}

		return rerun, err
	})

	return
}
func (s *Service) PluginFamily(url string) (ff FamilyPlugins, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for PluginFamily : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "PluginFamily")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &ff)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for PluginFamily: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "PluginFamily")
			return
		}

		return rerun, err
	})

	return
}
func (s *Service) Scanners(url string) (sl ScannerList, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for Scanners : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "Scanners")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &sl)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for Scanners: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "Scanners")
			return
		}

		return rerun, err
	})

	if err != nil {
		s.Errorf("%s permanently failed Scanners: %v", url, err)
		return
	}

	return
}
func (s *Service) ScannerAgents(url string) (sa ScannerAgent, raw []byte, err error) {

	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for ScannerAgents : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "ScannerAgents")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &sa)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for ScannerAgents: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "ScannerAgents")
			return
		}

		return rerun, err
	})

	return
}

func (s *Service) AssetVuln(url string) (as AssetVuln, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for AssetVuln : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "AssetVuln")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &as)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for AssetVuln: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "AssetVuln")
			return
		}

		return rerun, err
	})

	return
}
func (s *Service) AssetVulnInfo(url string) (asvi AssetVulnInfo, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for AssetVulnerabiltyInfo : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "AssetVulnerabiltyInfo")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &asvi)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for AssetVulnerabiltyInfo: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "AssetVulnerabiltyInfo")
			return
		}

		return rerun, err
	})

	return
}
func (s *Service) AssetVulnOutput(url string) (asvi AssetVulnOutput, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for AssetVulnOutput : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "AssetVulnOutput")
			return
		}
		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &asvi)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for AssetVulnOutput: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "AssetVulnOutput")
			return
		}

		return rerun, err
	})

	return
}

func (s *Service) AgentGroups(url string) (group ScannerAgentGroups, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to GET for AgentGroups : %s", url, err)
			rerun = s.SleepAndRerun(attempt, "AgentGroups")
			return
		}

		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &group)
		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for AgentGroups: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "AgentGroups")
			return
		}

		return rerun, err
	})

	return
}

func (s *Service) CreateAgentGroup(url string, j string) (group ScannerAgentGroup, raw []byte, err error) {
	tenable := NewPortal(s)

	raw, err = tenable.POST(url, j, "application/json")
	if err != nil {
		s.Errorf("'%s' failed to POST AgentGroup from POST: %s: %s: %s [len:%d]", url, j, err, string(raw), len(raw))
		return
	}

	if strings.Contains(string(raw), `{"error":"Agent Group with name`) {
		err = errors.New("error: Agent group already exists")
		return
	}

	err = json.Unmarshal(raw, &group)
	if err != nil {
		s.Errorf("'%s' failed to UNMARSHAL AgentGroup from POST: %s: %s: %s [len:%d]", url, j, err, string(raw), len(raw))
		return
	}

	return
}
func (s *Service) AssignAgentGroup(url string) (raw []byte, err error) {
	tenable := NewPortal(s)

	raw, err = tenable.PUT(url)
	if err != nil {
		s.Errorf("'%s' failed to AssignAgentGroup from PUT: %s: %s: %s [len:%d]", url, err, string(raw), len(raw))
		return
	}

	return
}

func (s *Service) AssetExport(url string, j string) (export AssetExport, raw []byte, err error) {
	tenable := NewPortal(s)

	raw, err = tenable.POST(url, j, "application/json")
	if err != nil {
		s.Errorf("'%s' failed to POST AssetExport from POST: %s: %s: %s: %s [len:%d]", url, j, err, string(raw), len(raw))
		return
	}

	err = json.Unmarshal(raw, &export)

	return
}

func (s *Service) AssetExportStatus(url string) (export AssetExportStatus, raw []byte, err error) {
	tenable := NewPortal(s)

	raw, err = tenable.GET(url)
	if err != nil {
		s.Errorf("'%s' failed to POST AssetExportStatus from POST: %s: %s: %s [len:%d]", url, err, string(raw), len(raw))
		return
	}

	err = json.Unmarshal(raw, &export)

	return
}

func (s *Service) AssetExportChunk(url string) (export []AssetExportChunk, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to POST AssetExportChunk from POST: %s: %s: %s [len:%d]", url, err, string(raw), len(raw))
			rerun = s.SleepAndRerun(attempt, "AssetExportChunk")
			return
		}

		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &export)

		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for AssetExportChunk: %s: [len:%d]", url, err, len(raw))
			rerun = s.SleepAndRerun(attempt, "AssetExportChunk")
			return
		}

		return rerun, err
	})

	return
}

func (s *Service) VulnExport(url string, j string) (export VulnExport, raw []byte, err error) {
	tenable := NewPortal(s)

	raw, err = tenable.POST(url, j, "application/json")
	if err != nil {
		s.Errorf("'%s' failed to POST VulnExport from POST: %s: %s: %s: %s [len:%d]", url, j, err, string(raw), len(raw))
		return
	}

	err = json.Unmarshal(raw, &export)

	return
}

func (s *Service) VulnExportStatus(url string) (export VulnExportStatus, raw []byte, err error) {
	tenable := NewPortal(s)

	raw, err = tenable.GET(url)
	if err != nil {
		s.Errorf("'%s' failed to POST VulnExportStatus from POST: %s: %s: %s [len:%d]", url, err, string(raw), len(raw))
		return
	}

	err = json.Unmarshal(raw, &export)

	return
}

func (s *Service) VulnExportChunk(url string) (export []VulnExportChunk, raw []byte, err error) {
	tenable := NewPortal(s)

	err = try.Do(func(attempt int) (rerun bool, err error) {
		// Fetch from Tenable.IO REST API
		raw, err = tenable.GET(url)
		if err != nil {
			s.Warnf("'%s' failed to POST VulnExportChunk from POST: %s: %s: %s [len:%d]", url, err, len(raw))
			rerun = s.SleepAndRerun(attempt, "VulnExportChunk")
			return
		}

		// Convert the response JSON bytes to Golang struct
		err = json.Unmarshal(raw, &export)

		if err != nil {
			s.Warnf("'%s' failed to UNMARSHAL for VulnExportChunk: %s: [len:%d]", url, err, len(raw))
			rerun = s.SleepAndRerun(attempt, "VulnExportChunk")
			return
		}

		return rerun, err
	})

	return
}
