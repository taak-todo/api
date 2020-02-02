package server

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/go-chi/chi/middleware"
)

// RequestLogger implements chi.LogFormatter to log requests via Apex log.
type RequestLogger struct{}

func (rl RequestLogger) NewLogEntry(req *http.Request) middleware.LogEntry {
	return RequestLogEntry{Request: req}
}

// RequestLogEntry handles logging requests and panics from Chi to Apex log.
type RequestLogEntry struct {
	Request *http.Request
}

func (rle RequestLogEntry) Write(status, bytes int, elapsed time.Duration) {
	log.WithFields(log.Fields{
		"request_id":     RequestID(rle.Request),
		"status_code":    status,
		"content_length": bytes,
		"response_time":  elapsed,
		"method":         rle.Request.Method,
		"path":           rle.Request.URL.Path,
		"query_params":   rle.Request.URL.RawQuery,
	}).Info("Request")
}

func (rle RequestLogEntry) Panic(val interface{}, stack []byte) {
	err, ok := val.(error)
	if ok {
		val = err.Error()
	}

	log.WithFields(log.Fields{
		"error":        val,
		"request_id":   RequestID(rle.Request),
		"method":       rle.Request.Method,
		"path":         rle.Request.URL.Path,
		"query_params": rle.Request.URL.RawQuery,
	}).Error("Request panic")
}
