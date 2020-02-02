package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/taak-todo/api/internal/server"
)

type Application struct {
	Router http.Handler
}

func NewApplication() *Application {
	app := new(Application)
	routes := []server.Route{
		{Method: "GET", Path: "/health", Handler: app.V1HealthHandler},
		{Method: "GET", Path: "/v1/health", Handler: app.V1HealthHandler},
	}

	app.Router = server.NewRouter(routes)
	return app
}

func (app *Application) V1HealthHandler(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write([]byte("ok"))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"request_id": server.RequestID(req),
		}).Error("Failed to write response")
	}
}
