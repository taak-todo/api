package main

import (
	"os"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/spf13/pflag"
	"github.com/taak-todo/api/internal/server"
)

func main() {
	log.SetHandler(jsonhandler.New(os.Stdout))

	var configPath string

	pflag.StringVarP(&configPath, "config", "c", "", "path to config file(required)")
	pflag.Parse()

	if configPath == "" {
		pflag.Usage()
		os.Exit(2)
	}

	err := run(configPath)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed running the taak api")
		os.Exit(1)
	}
}

func run(configPath string) error {
	config, err := NewConfigFromPath(configPath)
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return err
	}

	application := NewApplication()
	server := server.New(application.Router, config.Server)

	log.WithField("addr", config.Server.Addr).Info("Starting server")
	err = server.Serve()
	if err != nil {
		return err
	}

	log.Info("Server stopped")
	return nil
}
