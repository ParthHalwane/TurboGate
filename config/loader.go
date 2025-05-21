package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
