package main

import (
	_ "embed"
	"os"

	"github.com/Iandenh/splitter/cmd"
	goversion "github.com/caarlos0/go-version"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(
		buildVersion(),
		os.Exit,
		os.Args[1:],
	)
}

//go:embed art.txt
var asciiArt string

func buildVersion() goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails("splitter", "Proxy incoming requests to multiple upstreams.", "https://github.com/iandenh/splitter"),
		goversion.WithASCIIName(asciiArt),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
		},
	)
}
