package tenable

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var tr = &http.Transport{
	MaxIdleConns:    20,
	IdleConnTimeout: 30 * time.Second,
}

type Portal struct {
	BaseUrl     string
	AccessKey   string
	SecretKey   string
	WorkerCount int
	ThreadSafe  *sync.Mutex
	Service     *Service
}

func NewPortal(s *Service) (p *Portal) {
	p = new(Portal)
	p.Service = s
	p.BaseUrl = s.BaseUrl
	p.AccessKey = s.AccessKey
	p.SecretKey = s.SecretKey
	p.ThreadSafe = new(sync.Mutex)

	return
}

// hCalls is incremented every call to construct the Tenable.IO HTTP reques header
var hCalls int

// TenableXHeader inserts the AccessKey and SecretKey into the request header.
// hCalls is thread-safely incremented allowing multiple-requests from multiple-credentials (access/secret keys.)
func (p *Portal) TenableXHeader() (header string) {
	akeys := strings.Split(p.AccessKey, ",")
	skeys := strings.Split(p.SecretKey, ",")

	if len(akeys) != len(skeys) {
		return
	}

	p.ThreadSafe.Lock()
	hCalls = hCalls + 1
	mod := hCalls % len(akeys)
	p.ThreadSafe.Unlock()

	header = fmt.Sprintf("accessKey=%s;secretKey=%s", akeys[mod], skeys[mod])
	return
}

func (p *Portal) GET(endPoint string) ([]byte, error) {

	p.Service.Debugf("pkg.Tenable.GET(%s)", endPoint)
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-ApiKeys", p.TenableXHeader())

	resp, err := client.Do(req) // <-------HTTPS GET Request!
	if err != nil {
		return nil, err
	}

	// TODO:Improve some error detection here for 429 codes etc.
	if resp.StatusCode == 429 {
		err := errors.New("error: we need to slow down")
		return nil, err
	}
	if resp.StatusCode == 403 {
		err := errors.New("error: creds no longer authorized")
		return nil, err
	}

	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("error: status code does not appear successful: %d", resp.StatusCode))
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// VulnOutput contains some of the failure words, so we SUCESS tokens on this token
	sbody := string(body)
	if strings.Contains(sbody, `"plugin_output`) {
		return body, nil
	}
	if strings.Contains(sbody, `"output`) {
		return body, nil
	}

	// ////////
	// Full message from cloud.tenable.io:
	//  {"statusCode":401,"error":"Unauthorized","message":"Invalid Credentials"}
	if strings.Contains(sbody, `"statusCode":401`) || strings.Contains(sbody, `{"error":"Invalid Credentials"}`) {
		err := errors.New("your secretKey and accessKey (credentials) are invalid")
		return nil, err
	}

	if strings.Contains(sbody, `{"error":"Asset or host not found"}`) {
		//
	} else if strings.Contains(sbody, `{"error":"You need to log in to perform this request"}`) || strings.Contains(sbody, "504 Gateway Time-out") || strings.Contains(sbody, `"statusCode":504`) || strings.Contains(sbody, `Please retry request`) || strings.Contains(sbody, `Please wait a moment`) {
		msg := fmt.Sprintf("FAILED: GET '%s' Body:'%s'", endPoint, body)
		err := errors.New(msg)
		return nil, err
	}

	return body, nil
}
func (p *Portal) POST(endPoint string, postData string, postType string) (body []byte, err error) {
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer([]byte(postData)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-ApiKeys", p.TenableXHeader())
	req.Header.Set("Content-Type", postType)

	resp, err := client.Do(req) // <-------HTTPS GET Request!
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
func (p *Portal) DELETE(endPoint string) (err error) {

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("DELETE", endPoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-ApiKeys", p.TenableXHeader())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), `"error"`) {
		err := errors.New(fmt.Sprintf("Cannot delete from Tenable.IO, feature not yet implemented:%s", string(body)))
		return err
	}

	return
}

func (p *Portal) PUT(endPoint string) (body []byte, err error) {

	var req *http.Request
	client := &http.Client{Transport: tr}
	req, err = http.NewRequest("PUT", endPoint, nil)
	if err != nil {
		return
	}

	req.Header.Add("X-ApiKeys", p.TenableXHeader())

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	return
}
