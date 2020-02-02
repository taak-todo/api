package server

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/apex/log"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type requestContextKey int

const (
	requestIDKey requestContextKey = iota
)

// RequestIDMiddleware creates random request ids and adds them to the request context.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		uuid, err := uuid.NewRandom()
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to create request id")
		}

		ctx := context.WithValue(req.Context(), requestIDKey, uuid)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

// RecovererMiddleware recovers panics in an HTTP handler and logs the panic via Chi.
func RecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			val := recover()
			if val != nil {
				rw.WriteHeader(http.StatusInternalServerError)

				logEntry := middleware.GetLogEntry(req)
				if logEntry != nil {
					logEntry.Panic(val, debug.Stack())
				}
			}
		}()

		next.ServeHTTP(rw, req)
	})
}
