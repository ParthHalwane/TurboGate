package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

func LoadConfig(path string) ([]Route, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Routes, nil
}
