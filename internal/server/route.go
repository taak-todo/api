package server

import (
	"net/http"
)

// Route is a single route that has a handler and a list of middleware for the handler.
type Route struct {
	Method      string
	Path        string
	Middlewares []func(http.Handler) http.Handler
	Handler     http.HandlerFunc
}
