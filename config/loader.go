package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path   string `yaml:"path" json:"path"`
	Target string `yaml:"target" json:"target"`
}

type Config struct {
	Routes []Route `yaml:"routes" json:"routes"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func SaveConfig(cfgPath string, config Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}
