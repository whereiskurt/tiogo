package middleware

import (
	"bytes"
	"net/http"
	"os/exec"
	"strings"
)

// PrettyResponseCtx runs for every route, sets the response to JSON for all responses and unpacks AccessKey&SecretKey
func PrettyResponseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = NewPrettyPrint(w)
		next.ServeHTTP(w, r)
	})
}

// NewPrettyPrint checks for jq
func NewPrettyPrint(w http.ResponseWriter) (p *PrettyPrintJSON) {
	p = new(PrettyPrintJSON)
	p.w = w

	jq, err := exec.LookPath("jq")
	if err == nil {
		p.jq = jq
	}

	return
}

// PrettyPrintJSON holds a reference to the ResponseWrite and where 'jq' exec is
type PrettyPrintJSON struct {
	w  http.ResponseWriter
	jq string
}

// Write is called, and we rewrite if jq is installed in exec path
func (j *PrettyPrintJSON) Write(bb []byte) (int, error) {

	if j.jq == "" {
		return j.w.Write(bb)
	}

	bb = j.Prettify(bb)

	return j.w.Write(bb)
}

// Prettify takes raw JSON bytes and outputs pretty from jq
func (j *PrettyPrintJSON) Prettify(bb []byte) []byte {
	var pretty bytes.Buffer
	cmd := exec.Command(j.jq, ".")
	cmd.Stdin = strings.NewReader(string(bb))
	cmd.Stdout = &pretty
	err := cmd.Run()
	if err == nil {
		bb = []byte(pretty.String())
	}
	return bb
}

// Header overrides Header from ResponseWriter
func (j *PrettyPrintJSON) Header() http.Header {
	return j.w.Header()
}

// WriteHeader overrides ResponseWriter
func (j *PrettyPrintJSON) WriteHeader(statusCode int) {
	j.w.WriteHeader(statusCode)
}
