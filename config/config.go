package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	OriginHostName string   `yaml:"originHostName"`
	RewriteHost    bool     `yaml:"rewriteHost"`
	Port           int      `yaml:"port"`
	Upstreams      []string `yaml:"upstreams"`
}

func Load(filePath string) Config {
	c := Config{
		RewriteHost: false,
		Port:        1234,
	}

	if filePath == "" {
		return c
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatal(err)
	}

	return c
}
