package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	OriginHostName string   `json:"originHostName"`
	RewriteHost    bool     `json:"rewriteHost"`
	Port           int      `json:"port"`
	Upstreams      []string `json:"upstreams"`
}

func Load(filePath string) (Config, error) {
	c := Config{
		RewriteHost: false,
		Port:        1234,
	}

	if filePath == "" {
		return c, nil
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return c, err
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		return c, err
	}

	return c, nil
}
