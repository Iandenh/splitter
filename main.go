package main

import (
	"flag"
	"fmt"
	"os"
	"splitter/admin"
	"splitter/config"
	"splitter/listener"
	"splitter/upstream"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "", "Config file to load")
}
func main() {
	flag.Parse()

	if configFilePath == "" {
		fmt.Println("No Config loaded")
		os.Exit(0)
	}

	c := config.Load(configFilePath)

	if c.AdminEnabled {
		go admin.Start(c.AdminPort)
	}

	for _, upstreamUrl := range c.Upstreams {
		upstream.AddUpstream(upstream.Upstream{
			Url: upstreamUrl,
		})
	}

	l := listener.Listener{
		OriginHostName: c.OriginHostName,
		RewriteHost:    c.RewriteHost,
		Port:           c.Port,
	}
	l.Start()
}
