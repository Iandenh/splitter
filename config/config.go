package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	AdminEnabled   bool     `yaml:"adminEnabled"`
	AdminPort      int      `yaml:"adminPort"`
	OriginHostName string   `yaml:"originHostName"`
	RewriteHost    bool     `yaml:"rewriteHost"`
	Port           int      `yaml:"port"`
	Upstreams      []string `yaml:"upstreams"`
}

func Load(filePath string) Config {
	c := Config{
		AdminEnabled: true,
		RewriteHost:  false,
		AdminPort:    8888,
		Port:         1234,
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
