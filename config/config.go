package config

import (
	"log"
	"os"
	"sigs.k8s.io/yaml"
)

type Config struct {
	OriginHostName string   `json:"originHostName"`
	RewriteHost    bool     `json:"rewriteHost"`
	Port           int      `json:"port"`
	Upstreams      []string `json:"upstreams"`
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
