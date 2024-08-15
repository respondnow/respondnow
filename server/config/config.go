package config

import (
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"
)

type Config struct {
	Address          string `json:"address"`
	ServerURL        string `json:"serverURL"`
	SkipSecureVerify bool   `json:"skipSecureVerify"`
}

func New(configFilePath string) (*Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if json.Valid(data) {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
