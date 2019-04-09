package tenable

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/whereiskurt/tiogo/pkg/cache"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"gopkg.in/matryer/try.v1"
	"strings"
	"sync"
)

var EndPoints = endPointTypes{
	Scanners:          EndPointType("Scanners"),
	VulnsExportStart:  EndPointType("VulnsExportStart"),
	VulnsExportStatus: EndPointType("VulnsExportStatus"),
	VulnsExportGet:    EndPointType("VulnsExportGet"),
}

type endPointTypes struct {
	Scanners          EndPointType
	VulnsExportStart  EndPointType
	VulnsExportStatus EndPointType
	VulnsExportGet    EndPointType
}

// ServiceMap defines all the endpoints provided by the ACME service
var ServiceMap = map[EndPointType]ServiceTransport{
	EndPoints.VulnsExportStart: {
		URL:           "/vulns/export",
		CacheFilename: "/export/vulns/request.json",
		MethodTemplate: map[httpMethodType]MethodTemplate{
			HTTP.Post: {`
			"export-request": "export-request",
			"filters": {
				"since": "2019-01-01T00:00:00Z"
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
	EndPoints.Scanners: {
		URL:           "/scanners",
		CacheFilename: "/scanners.json",
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
	SecretKey      string // ACME Secret Keys for API Access (provided by ACME)
	AccessKey      string //             Access Key for API access (provided by ACME)
	RetryIntervals []int  // When a call to a transport fails, this will control the retrying.
	DiskCache      *cache.Disk
	Worker         *sync.WaitGroup // Used by Go routines to control workers (TODO)
	Log            *log.Logger
	Metrics        *metrics.Metrics
}

func NewService(base string, secret string, access string) (s Service) {
	s.BaseURL = strings.TrimSuffix(base, "/")
	s.SecretKey = secret
	s.AccessKey = access
	s.RetryIntervals = DefaultRetryIntervals
	s.Worker = new(sync.WaitGroup)

	s.Log = log.New()

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
func (s *Service) VulnsExportStart() ([]byte, error) {
	var raw []byte

	err := try.Do(func(attempt int) (bool, error) {
		body, status, err := s.update(EndPoints.VulnsExportStart, nil)
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
