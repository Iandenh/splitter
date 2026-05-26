package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

import _ "embed"

//go:embed static/config.yaml
var exampleConfig []byte

type initCmd struct {
	cmd    *cobra.Command
	config string
}

func newInitCmd() *initCmd {
	root := &initCmd{}

	cmd := &cobra.Command{
		Use:           "init",
		Aliases:       []string{"i"},
		Short:         "Generates a new splitter configuration file",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(root.config); err == nil {
				return errors.New(root.config + " already exists, delete it and run the command again")
			}
			conf, err := os.OpenFile(root.config, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0o644)
			if err != nil {
				return err
			}
			defer conf.Close()

			if _, err := conf.Write(exampleConfig); err != nil {
				return err
			}

			done := []string{
				boldStyle.Render("done!"),
				"please edit", codeStyle.Render(root.config),
			}
			done = append(done, "accordingly.")
			cmd.Println(strings.Join(done, " "))

			return nil
		},
	}

	root.cmd = cmd
	cmd.Flags().StringVarP(&root.config, "config", "c", "splitter.yaml", "Load configuration from file")

	root.cmd = cmd

	return root
}
