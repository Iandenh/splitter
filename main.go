package main

import (
	"flag"
	"fmt"
	"os"
	"splitter/config"
	"splitter/listener"
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

	l := listener.New(c.OriginHostName, c.RewriteHost, c.Port, c.Upstreams)

	l.Start()
}
