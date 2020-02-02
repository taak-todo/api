package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter creates a http.Handler to route endpoints for a JSON API.
func NewRouter(routes []Route) http.Handler {
	router := chi.NewRouter()

	router.Use(
		RequestIDMiddleware,
		middleware.RequestLogger(RequestLogger{}),
		RecovererMiddleware,
		middleware.AllowContentType("application/json"),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Timeout(4*time.Second),
	)

	router.MethodNotAllowed(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})

	router.NotFound(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.With(route.Middlewares...).Get(route.Path, route.Handler)
		default:
			router.With(route.Middlewares...).Get(route.Path, route.Handler)
		}
	}

	return router
}
