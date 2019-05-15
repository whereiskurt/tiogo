package tenable

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// DefaultRetryIntervals values in here we control the re-try of the Service
var DefaultRetryIntervals = []int{0, 500, 500, 500, 500, 1000, 1000, 1000, 1000, 1000, 3000}

type EndPointType string

func (c EndPointType) String() string {
	return "pkg.tenable.endpoints." + string(c)
}

func (s *Service) put(endPoint EndPointType, p map[string]string) ([]byte, int, error) {

	url, err := toURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	data, _ := toJSON(endPoint, HTTP.Put, p)

	body, status, err := t.Put(url, data, "application/json")

	if err != nil {
		return nil, status, err
	}

	return body, status, err
}

func (s *Service) get(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	if s.SkipOnHit == true {
		// Check for a cache hit
		if s.DiskCache != nil {
			// We have initialized a cache then write to it.
			filename, err := ToCacheFilename(endPoint, p)
			if err != nil {
				return nil, 0, err
			}

			body, err := s.DiskCache.Fetch(filename)
			if err == nil && len(body) > 0 {
				return body, 200, nil
			}
		}
	}

	url, err := toURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	body, status, err := t.Get(url, s.SkipOnHit, s.WriteOnReturn)

	if err != nil {
		return nil, status, err
	}

	if s.WriteOnReturn == true {
		// If we have a DiskCache it means we will write out responses to disk.
		if s.DiskCache != nil {
			// We have initialized a cache then write to it.
			filename, err := ToCacheFilename(endPoint, p)
			if err != nil {
				return nil, status, err
			}

			err = s.DiskCache.Store(filename, body)
			if err != nil {
				return nil, status, err
			}
		}
	}

	return body, status, err
}
func (s *Service) delete(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	url, err := toURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}
	t := NewTransport(s)
	body, status, err := t.Delete(url)
	if err != nil {
		return nil, status, err
	}

	return body, status, err
}
func (s *Service) update(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	url, err := toURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}

	j, err := toJSON(endPoint, HTTP.Post, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	body, status, err := t.Post(url, j, "application/json")
	if err != nil {
		return nil, status, err
	}

	return body, status, err
}
func (s *Service) add(endPoint EndPointType, p map[string]string) ([]byte, int, error) {
	url, err := toURL(s.BaseURL, endPoint, p)
	if err != nil {
		return nil, 0, err
	}
	j, err := toJSON(endPoint, HTTP.Put, p)
	if err != nil {
		return nil, 0, err
	}

	t := NewTransport(s)
	body, status, err := t.Put(url, j, "application/json")
	if err != nil {
		return nil, status, err
	}

	return body, status, err
}

func toURL(baseURL string, name EndPointType, p map[string]string) (string, error) {
	sMap, hasMethod := ServiceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid name '%s' for URL lookup", name)
	}

	if p == nil {
		p = make(map[string]string)
	}
	p["ServiceBaseURL"] = baseURL

	// Append the ServiceBaseURL to the URL
	url := fmt.Sprintf("%s%s", baseURL, sMap.URL)

	return toTemplate(name, p, url)
}

func toJSON(name EndPointType, method httpMethodType, p map[string]string) (string, error) {
	sMap, hasMethod := ServiceMap[name]
	if !hasMethod {
		return "", fmt.Errorf("invalid method '%s' for name '%s'", method, name)
	}

	mMap, hasTemplate := sMap.MethodTemplate[method]
	if !hasTemplate {
		// return "", fmt.Errorf("invalid template for method '%s' for name '%s'", method, name)
		return "", nil
	}

	tmpl := mMap.Template
	return toTemplate(name, p, tmpl)
}
func toTemplate(name EndPointType, data map[string]string, tmpl string) (string, error) {
	var rawURL bytes.Buffer
	t, terr := template.New(fmt.Sprintf("%s", name)).Parse(tmpl)
	if terr != nil {
		err := fmt.Errorf("error: failed to parse template for %s: %v", name, terr)
		return "", err
	}
	err := t.Execute(&rawURL, data)
	if err != nil {
		return "", err
	}

	url := rawURL.String()

	return url, nil
}

func (s *Service) sleepBeforeRetry(attempt int) (shouldReRun bool) {
	if attempt < len(s.RetryIntervals) {
		s.Log.Infof("Failure leading to sleep='%dms'", s.RetryIntervals[attempt])
		time.Sleep(time.Duration(s.RetryIntervals[attempt]) * time.Millisecond)
		shouldReRun = true
	}
	return
}
