package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/taak-todo/api/internal/server"
)

type ConfigValidationError struct {
	InvalidFields []string
}

func (cve ConfigValidationError) Error() string {
	return fmt.Sprintf("missing or invalid config fields: %s", strings.Join(cve.InvalidFields, ","))
}

type Config struct {
	Server server.Config
}

func NewConfigFromPath(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := new(Config)
	decoder := toml.NewDecoder(file)

	return config, decoder.Decode(config)
}

func (c Config) Validate() error {
	var invalidFields []string

	if c.Server.Addr == "" {
		invalidFields = append(invalidFields, "Server.Addr")
	}

	if len(invalidFields) == 0 {
		return nil
	}

	return ConfigValidationError{
		InvalidFields: invalidFields,
	}
}
