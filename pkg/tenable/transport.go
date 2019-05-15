package tenable

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var HTTP = httpMethodTypes{
	Get:    httpMethodType("Get"),
	Delete: httpMethodType("Delete"),
	Head:   httpMethodType("Head"),
	Post:   httpMethodType("Post"),
	Put:    httpMethodType("Put"),
}

type httpMethodType string
type httpMethodTypes struct {
	Get    httpMethodType
	Put    httpMethodType
	Post   httpMethodType
	Delete httpMethodType
	Head   httpMethodType
}

func (c httpMethodType) String() string {
	return "pkg.tenable.transport." + string(c)
}

var tr = &http.Transport{
	MaxIdleConns:        20,
	IdleConnTimeout:     20 * time.Second,
	TLSHandshakeTimeout: 10 * time.Second,
}

// Transport defines the HTTP details for the API call.
type Transport struct {
	BaseURL     string
	AccessKey   string
	SecretKey   string
	Log         *log.Logger
	WorkerCount int
	ThreadSafe  *sync.Mutex
}

// NewTransport handles the HTTP methods GET/POST/PUT/DELETE
func NewTransport(s *Service) (p Transport) {
	p.BaseURL = s.BaseURL
	p.AccessKey = s.AccessKey
	p.SecretKey = s.SecretKey
	p.ThreadSafe = new(sync.Mutex)
	p.Log = s.Log
	return
}

// Inserts the AccessKey and SecretKey into the request authHeaderValue.
// AccessKey/SecretKey may be equally lengthed comma separated values that are rotated through each call.
// headerCallCount is thread-safely incremented allowing multiple-requests from multiple-credentials (access/secret keys.)
var headerCallCount int

func (t *Transport) authHeaderKey() string {
	return "X-ApiKeys"
}
func (t *Transport) authHeaderValue() string {
	ak := strings.Split(t.AccessKey, ",")
	sk := strings.Split(t.SecretKey, ",")

	if len(ak) != len(sk) {
		t.Log.Fatalf("error: equal amount of accesskeys and secretkeys must be specified.")
		return ""
	}

	// Ensure incremental non-overlapping count
	t.ThreadSafe.Lock()
	headerCallCount = headerCallCount + 1
	mod := headerCallCount % len(ak)
	t.ThreadSafe.Unlock()

	return fmt.Sprintf("accessKey=%s;secretKey=%s", ak[mod], sk[mod])
}

// Get will HTTP GET for the url provided, returning the body, status, and error associated with the call.
func (t *Transport) Get(url string, skipOnHit bool, writeOnReturn bool) ([]byte, int, error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Minute, // Really big downloads, we will wait 5mins max.
	}

	var err error
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add(t.authHeaderKey(), t.authHeaderValue())
	req.Header.Add("X-Cache-SkipOnHit", fmt.Sprintf("%v", skipOnHit))
	req.Header.Add("X-Cache-WriteOnReturn", fmt.Sprintf("%v", writeOnReturn))

	log.Debugf("tenable.transport.Get - URL='%s'", url)

	resp, err = client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode

	log.Debugf("tenable.transport.Get - URL='%s', stuat='%d'", url, http.StatusAccepted)
	if status == http.StatusTooManyRequests {
		err = errors.New("error: we need to slow down")
		return nil, status, err
	}
	if status == http.StatusForbidden {
		err = errors.New("error: credentials no longer authorized")
		return nil, status, err
	}
	if status != http.StatusOK {
		err = fmt.Errorf("error: status code '%d' does not appear successful for '%s'", status, url)
		return nil, status, err
	}

	respBody := resp.Body
	//
	// switch strings.ToLower(resp.Header.Get("Content-Encoding")) {
	// case "gzip":
	// 	respBody, err = gzip.NewReader(resp.Body)
	// }

	var body []byte
	body, err = ioutil.ReadAll(respBody)

	return body, status, err
}

// Get will HTTP POST for the url provided, returning the body, status, and error associated with the call.
func (t *Transport) Post(url string, data string, datatype string) ([]byte, int, error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Minute, // Really big downloads, we will wait 5mins max.
	}

	var err error
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add(t.authHeaderKey(), t.authHeaderValue())
	req.Header.Set("Content-Type", datatype)
	// req.Header.Add("Accept-Encoding", "gzip, deflate")

	resp, err = client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	status := resp.StatusCode

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}

	return body, status, err
}

// Get will HTTP PUT for the url provided, returning the body, status, and error associated with the call.
func (t *Transport) Put(url string, data string, datatype string) ([]byte, int, error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Minute, // Really big downloads, we will wait 5mins max.
	}

	var err error
	req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add(t.authHeaderKey(), t.authHeaderValue())
	// req.Header.Add("Accept-Encoding", "gzip, deflate")

	resp, err = client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	status := resp.StatusCode

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}
	return body, status, err
}

// Get will HTTP DELETE for the url provided, returning the body, status, and error associated with the call.
func (t *Transport) Delete(url string) ([]byte, int, error) {
	var req *http.Request
	var resp *http.Response

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Minute, // Really big downloads, we will wait 5mins max.
	}

	var err error
	req, err = http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add(t.authHeaderKey(), t.authHeaderValue())
	// req.Header.Add("Accept-Encoding", "gzip, deflate")

	resp, err = client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	status := resp.StatusCode

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}

	return body, status, err
}
