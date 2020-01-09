package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/prometheus/common/log"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// DefaultCacheSkipOnHit determines whether to try for a cache hit - when 'true' then return the cached results
// else don't do a cache lookup at all.
var DefaultCacheSkipOnHit = true

// DefaultWriteOnReturn determines whether to write a fresh cache entry - when 'true' then the result is written to cache
// else don't write the result.
var DefaultWriteOnReturn = true

// Contexts extract all of the params related to their route
type contextMapKey string

func (c contextMapKey) String() string {
	return "pkg.server.context" + string(c)
}

// ContextMapKey is the key to the request context
var ContextMapKey = contextMapKey("ctxMap")

// ContextMap extract from request and type asserts it (helper function.)
func ContextMap(r *http.Request) map[string]string {
	return (r.Context().Value(ContextMapKey)).(map[string]string)
}

// InitialCtx runs for every route, sets the response to JSON for all responses and unpacks AccessKey&SecretKey
func InitialCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctxMap := make(map[string]string)

		xKeys := strings.Split(r.Header.Get("X-ApiKeys"), ";")
		for x := range xKeys {
			keys := strings.Split(xKeys[x], "=")
			switch {
			case strings.ToLower(keys[0]) == "accesskey":
				ctxMap["AccessKey"] = keys[1]

			case strings.ToLower(keys[0]) == "secretkey":
				ctxMap["SecretKey"] = keys[1]
			}
		}

		ctxMap["SkipOnHit"] = r.Header.Get("X-Cache-SkipOnHit")
		ctxMap["WriteOnReturn"] = r.Header.Get("X-Cache-WriteOnReturn")

		if ctxMap["SkipOnHit"] == "" {
			ctxMap["SkipOnHit"] = fmt.Sprintf("%v", DefaultCacheSkipOnHit)
		}
		if ctxMap["WriteOnReturn"] == "" {
			ctxMap["WriteOnReturn"] = fmt.Sprintf("%v", DefaultWriteOnReturn)
		}

		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AccessKey is the Tenable.io access key required in the header
func AccessKey(r *http.Request) string {
	return ContextMap(r)["AccessKey"]
}

// SkipOnHit pulls the param from the request/contextmap
func SkipOnHit(r *http.Request) string {
	return ContextMap(r)["SkipOnHit"]
}

// WriteOnReturn is the Tenable.io secret key required in the header
func WriteOnReturn(r *http.Request) string {
	return ContextMap(r)["WriteOnReturn"]
}

// SecretKey is the Tenable.io secret key required in the header
func SecretKey(r *http.Request) string {
	return ContextMap(r)["SecretKey"]
}

// ExportUUID is used for Vulns and Asset exports
func ExportUUID(r *http.Request) string {
	return ContextMap(r)["ExportUUID"]
}

// ScannerID pulls the param from the request/contextmap
func ScannerID(r *http.Request) string {
	return ContextMap(r)["ScannerID"]
}

// ScanUUID pulls the param from the request/contextmap
func ScanUUID(r *http.Request) string {
	return ContextMap(r)["ScanUUID"]
}

// GroupID pulls the param from the request/contextmap
func GroupID(r *http.Request) string {
	return ContextMap(r)["GroupID"]
}

// AgentID pulls the param from the request/contextmap
func AgentID(r *http.Request) string {
	return ContextMap(r)["AgentID"]
}

// Offset pulls the param from the request/contextmap
func Offset(r *http.Request) string {
	return ContextMap(r)["Offset"]
}

// ScanID pulls the param from the request/contextmap
func ScanID(r *http.Request) string {
	return ContextMap(r)["ScanID"]
}

// HistoryID pulls the param from the request/contextmap
func HistoryID(r *http.Request) string {
	return ContextMap(r)["HistoryID"]
}

// FileUUID pulls the param from the request/contextmap
func FileUUID(r *http.Request) string {
	return ContextMap(r)["FileUUID"]
}

// Limit pulls the param from the request/contextmap
func Limit(r *http.Request) string {
	return ContextMap(r)["Limit"]
}

// ChunkID pulls the param from the request/contextmap
func ChunkID(r *http.Request) string {
	return ContextMap(r)["ChunkID"]
}

// Format pulls the param from the request/contextmap
func Format(r *http.Request) (format string, chapters string) {
	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("couldn't read scan export body: %v", err)
		return "", ""
	}

	var body tenable.ScansExportStartPost
	err = json.Unmarshal(bb, &body)
	if err != nil {
		return "", ""
	}

	return body.Format, body.Chapters
}

// ExportCtx pulls the param from the request/contextmap
func ExportCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["ExportUUID"] = chi.URLParam(r, "ExportUUID")
		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ExportChunkCtx pulls the param from the request/contextmap
func ExportChunkCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["ChunkID"] = chi.URLParam(r, "ChunkID")
		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ScannersCtx pulls the param from the request/contextmap
func ScannersCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["ScannerID"] = chi.URLParam(r, "ScannerID")
		ctxMap["Offset"] = r.URL.Query().Get("offset")
		ctxMap["Limit"] = r.URL.Query().Get("limit")

		if ctxMap["Offset"] == "" {
			ctxMap["Offset"] = "0"
		}
		if ctxMap["Limit"] == "" {
			ctxMap["Limit"] = "5000"
		}
		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AgentGroupCtx pulls the param from the request/contextmap
func AgentGroupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["GroupID"] = chi.URLParam(r, "GroupID")
		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AgentCtx extracts AgentID from the parameters
func AgentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["AgentID"] = chi.URLParam(r, "AgentID")
		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ScanCtx extracts ScanID from the parameters
func ScanCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["ScanUUID"] = chi.URLParam(r, "ScanUUID")
		ctxMap["ScanID"] = chi.URLParam(r, "ScanUUID")

		ctxMap["HistoryID"] = r.URL.Query().Get("history_id")
		if ctxMap["HistoryID"] == "" {
			ctxMap["HistoryID"] = r.URL.Query().Get("HistoryID")
		}

		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ScansExportCtx extracts FileUUID from the parameters
func ScansExportCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxMap := r.Context().Value(ContextMapKey).(map[string]string)
		ctxMap["FileUUID"] = chi.URLParam(r, "FileUUID")

		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
