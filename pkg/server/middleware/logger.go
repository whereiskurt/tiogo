package middleware

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// https://github.com/go-chi/chi/blob/master/middleware/logger.go
//
// StructuredLogger is a simple, but powerful implementation of a custom structured
// logger backed on log. I encourage users to copy it, adapt it and make it their
// own. Also take a look at https://github.com/pressly/lg for a dedicated pkg based
// on this work, designed for context-based http routers.

// NewStructuredLogger implements the go-chi middle interface method necessary to be called
func NewStructuredLogger(logger *log.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{logger})
}

// structuredLogger keeps reference to the logrus :-)
type structuredLogger struct {
	Logger *log.Logger
}

// NewLogEntry is middleware that's called part of go-chi's pipeline
func (l *structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: log.NewEntry(l.Logger)}

	logFields := log.Fields{}
	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	logFields["http_scheme"] = "http"
	if r.TLS != nil {
		logFields["http_scheme"] = "https"
	}
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()
	logFields["uri"] = fmt.Sprintf("%s://%s%s", logFields["http_scheme"], r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Debug("request started")

	return entry
}

// StructuredLoggerEntry is a required inherited interface
type StructuredLoggerEntry struct {
	Logger log.FieldLogger
}

// Write is a required inherited method
func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Debug("request complete")
}

// Panic is a required inherited method
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
