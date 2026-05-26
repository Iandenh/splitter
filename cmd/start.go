package cmd

import (
	"github.com/Iandenh/splitter/config"
	"github.com/Iandenh/splitter/listener"
	"github.com/spf13/cobra"
)

type startCmd struct {
	cmd    *cobra.Command
	config string
}

func newStartCmd() *startCmd {
	root := &startCmd{}

	cmd := &cobra.Command{
		Use:           "start",
		Short:         "Start the server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			c := config.Load(root.config)

			l := listener.New(c.OriginHostName, c.RewriteHost, c.Port, c.Upstreams)

			go l.Start()

			select {
			case <-cmd.Context().Done():
				return
			}
		},
	}

	root.cmd = cmd

	cmd.Flags().StringVarP(&root.config, "config", "f", "splitter.yaml", "Path to config file")

	return root
}
