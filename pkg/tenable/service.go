package tenable

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/whereiskurt/tiogo/pkg/cache"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"gopkg.in/matryer/try.v1"
)

var EndPoints = endPointTypes{
	VulnsExportStart:   EndPointType("VulnsExportStart"),
	VulnsExportStatus:  EndPointType("VulnsExportStatus"),
	VulnsExportGet:     EndPointType("VulnsExportGet"),
	AssetsExportStart:  EndPointType("AssetsExportStart"),
	AssetsExportStatus: EndPointType("AssetsExportStatus"),
	AssetsExportGet:    EndPointType("AssetsExportGet"),
	ScannersList:       EndPointType("ScannersList"),
	AgentsList:         EndPointType("AgentsList"),
	ScannerAgentGroups: EndPointType("ScannerAgentGroups"),
	AgentsGroup:        EndPointType("AgentsGroup"),
}

type endPointTypes struct {
	VulnsExportStart  EndPointType
	VulnsExportStatus EndPointType
	VulnsExportGet    EndPointType

	AssetsExportStart  EndPointType
	AssetsExportStatus EndPointType
	AssetsExportGet    EndPointType

	ScannersList       EndPointType
	AgentsList         EndPointType
	ScannerAgentGroups EndPointType
	AgentsGroup        EndPointType
	AgentsUngroup      EndPointType
}

// ServiceMap defines all the endpoints provided by the ACME service
var ServiceMap = map[EndPointType]ServiceTransport{

	EndPoints.ScannerAgentGroups: {
		URL:           "/scanners/{{.ScannerID}}/agent-groups",
		CacheFilename: "/scanners/{{.ScannerID}}/agent.groups.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},

	EndPoints.AgentsList: {
		URL:           "/scanners/{{.ScannerID}}/agents?offset={{.Offset}}&limit={{.Limit}}",
		CacheFilename: "/scanners/{{.ScannerID}}/agents.{{.Offset}}.{{.Limit}}.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},

	EndPoints.AgentsGroup: {
		URL: "/scanners/{{.ScannerID}}/agent-groups/{{.GroupID}}/agents/{{.AgentID}}",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Put: {},
		},
	},
	EndPoints.AgentsUngroup: {
		URL: "/scanners/{{.ScannerID}}/agent-groups/{{.GroupID}}/agents/{{.AgentID}}",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Delete: {},
		},
	},

	EndPoints.ScannersList: {
		URL:           "/scanners",
		CacheFilename: "/scanners/scanners.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},

	EndPoints.VulnsExportStart: {
		URL:           "/vulns/export",
		CacheFilename: "/export/vulns/request.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Post: {`
{
	"export-request": "export-request",
	"num_assets": {{.Limit}},
	"filters": {
		"since": {{.Since}}
	}
}`},
		},
	},
	EndPoints.VulnsExportStatus: {
		URL:           "/vulns/export/{{.ExportUUID}}/status",
		CacheFilename: "/export/vulns/{{.ExportUUID}}/status.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},
	EndPoints.VulnsExportGet: {
		URL:           "/vulns/export/{{.ExportUUID}}/chunks/{{.ChunkID}}",
		CacheFilename: "/export/vulns/{{.ExportUUID}}/chunk.{{.ChunkID}}.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},
	EndPoints.AssetsExportStart: {
		URL:           "/assets/export",
		CacheFilename: "/export/assets/request.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Post: {`
{
	"export-request": "export-request",
	"chunk_size": {{.Limit}},
	"filters": {
		"last_assessed": {{.LastAssessed}} 
	}
}`},
		},
	},
	EndPoints.AssetsExportStatus: {
		URL:           "/assets/export/{{.ExportUUID}}/status",
		CacheFilename: "/export/assets/{{.ExportUUID}}/status.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},
	EndPoints.AssetsExportGet: {
		URL:           "/assets/export/{{.ExportUUID}}/chunks/{{.ChunkID}}",
		CacheFilename: "/export/assets/{{.ExportUUID}}/chunk.{{.ChunkID}}.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Get: {},
		},
	},
}

// ServiceTransport describes a URL endpoint that can be called ACME. Depending on the HTTP method (GET/POST/DELETE)
// we will render the appropriate MethodTemplate
type ServiceTransport struct {
	URL            string
	CacheFilename  string
	MethodTemplate map[httpMethodType]MethodTemplate
}

// MethodTemplate for each GET/PUT/POST/DELETE that is called this template is rendered
// For POST it is Put in the BODY, for GET it is added after "?" on the URL, f
type MethodTemplate struct {
	Template string
}

// Service exposes ACME services by converting the JSON results to to Go []structures
type Service struct {
	BaseURL        string // Put in front of every transport call
	SecretKey      string //
	AccessKey      string //
	RetryIntervals []int  // When a call to a transport fails, this will control the retrying.
	DiskCache      *cache.Disk
	Worker         *sync.WaitGroup // Used by Go routines to control workers (TODO)
	Log            *log.Logger
	Metrics        *metrics.Metrics
	SkipOnHit      bool
	WriteOnReturn  bool
}

func NewService(base string, secret string, access string, log *log.Logger) (s Service) {
	s.BaseURL = strings.TrimSuffix(base, "/")
	s.SecretKey = secret
	s.AccessKey = access
	s.RetryIntervals = DefaultRetryIntervals
	s.Worker = new(sync.WaitGroup)

	s.Log = log

	return
}

// EnableMetrics will produce prometheus metrics calls
func (s *Service) EnableMetrics(metrics *metrics.Metrics) {
	s.Metrics = metrics
}

// EnableCache will create a new Disk Cache for all request.
func (s *Service) EnableCache(cacheFolder string, cryptoKey string) {
	var useCrypto = false
	if cryptoKey != "" {
		useCrypto = true
	}
	s.DiskCache = cache.NewDisk(cacheFolder, cryptoKey, useCrypto)
	return
}

func ToCacheFilename(name EndPointType, p map[string]string) (string, error) {
	sMap, ok := ServiceMap[name]
	if !ok {
		return "", fmt.Errorf("invalid name '%s' for cache filename lookup", name)
	}
	return toTemplate(name, p, sMap.CacheFilename)
}

func (s *Service) ScannersList() ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.ScannersList, nil)
		if err != nil {
			s.Log.Infof("failed to scanners list: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		if status != 200 {
			msg := fmt.Sprintf("error not implemented! status: %d, %v", status, err)
			s.Log.Error(msg)
			return false, errors.New(msg)
		}

		raw = body
		return false, nil
	})

	return raw, err
}
func (s *Service) ScannerAgentGroups(uuid string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.ScannerAgentGroups, map[string]string{"ScannerID": uuid})
		if err != nil {
			s.Log.Infof("failed to agent groups list: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		if status != 200 {
			msg := fmt.Sprintf("error not implemented! status: %d, %v", status, err)
			s.Log.Error(msg)
			return false, errors.New(msg)
		}

		raw = body
		return false, nil
	})

	return raw, err
}

func (s *Service) AgentList(scannerId string, offset string, limit string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.AgentsList, map[string]string{"ScannerID": scannerId, "Offset": offset, "Limit": limit})
		if err != nil {
			s.Log.Infof("failed to agent list: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		if status != 200 {
			msg := fmt.Sprintf("error not implemented! status: %d, %v", status, err)
			s.Log.Error(msg)
			return false, errors.New(msg)
		}

		raw = body
		return false, nil
	})

	return raw, err
}

// AgentGroup calls Service (proxy or Tenable.io) and will assign an agentID to a groupID, givent the scannerID
func (s *Service) AgentGroup(agentID string, groupID string, scannerID string) ([]byte, error) {
	//TODO: Consider reording method sig. agent,scanner,group
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.put(EndPoints.AgentsGroup, map[string]string{"ScannerID": scannerID, "GroupID": groupID, "AgentID": agentID})
		if err != nil {
			s.Log.Infof("failed to group agent: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		if status != 200 {
			msg := fmt.Sprintf("error not implemented: status: %d", status)
			s.Log.Error(msg)
			return false, errors.New(msg)
		}

		raw = body
		return false, nil
	})

	return raw, err
}
func (s *Service) AgentUngroup(agentId string, groupId string, scannerId string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.delete(EndPoints.AgentsGroup, map[string]string{"ScannerID": scannerId, "GroupID": groupId, "AgentID": agentId})
		if err != nil {
			s.Log.Infof("failed to ungroup agent: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		if status != 200 {
			msg := fmt.Sprintf("error not implemented: status: %d", status)
			s.Log.Error(msg)
			return false, errors.New(msg)
		}

		raw = body
		return false, nil
	})

	return raw, err
}

func (s *Service) VulnsExportStatus(exportUUID string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.VulnsExportStatus, map[string]string{"ExportUUID": exportUUID})
		if err != nil {
			s.Log.Infof("failed to get export-vulns status: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		raw = body
		s.Log.Debugf("JSON body from start export-vulns status: %s", raw)

		return false, nil
	})

	return raw, err
}
func (s *Service) VulnsExportStart(limit string, sinceUnix string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.update(EndPoints.VulnsExportStart, map[string]string{"Since": sinceUnix, "Limit": limit})
		if err != nil {
			s.Log.Infof("failed to export-vulns start: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		raw = body
		s.Log.Debugf("JSON body from export-vulns start: %s", string(raw))

		return false, nil
	})

	return raw, err
}
func (s *Service) VulnsExportGet(exportUUID string, chunk string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.VulnsExportGet, map[string]string{"ExportUUID": exportUUID, "ChunkID": chunk})
		if err != nil {
			s.Log.Infof("failed to get export-vulns status: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		raw = body
		s.Log.Debugf("JSON body from get export-vulns length: %d", len(raw))

		return false, nil
	})

	return raw, err
}

func (s *Service) AssetsExportStatus(exportUUID string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.AssetsExportStatus, map[string]string{"ExportUUID": exportUUID})
		if err != nil {
			s.Log.Infof("failed to get export-assets status: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		raw = body
		s.Log.Debugf("JSON body from start export-assets status: %s", raw)

		return false, nil
	})

	return raw, err
}
func (s *Service) AssetsExportStart(limit string, lastAssessedUnix string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.update(EndPoints.AssetsExportStart, map[string]string{"Limit": limit, "LastAssessed": lastAssessedUnix})
		if err != nil {
			s.Log.Infof("failed to export-assets start: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		raw = body
		s.Log.Debugf("JSON body from export-assets start: %s", string(raw))

		return false, nil
	})

	return raw, err
}
func (s *Service) AssetsExportGet(exportUUID string, chunk string) ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.get(EndPoints.AssetsExportGet, map[string]string{"ExportUUID": exportUUID, "ChunkID": chunk})
		if err != nil {
			s.Log.Infof("failed to get export-assets status: http status: %d: %s", status, err)
			retry := s.sleepBeforeRetry(attempt)
			return retry, err
		}

		raw = body
		s.Log.Debugf("JSON body from get export-assets length: %d", len(raw))

		return false, nil
	})

	return raw, err
}
