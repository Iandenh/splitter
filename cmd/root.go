package cmd

import (
	"context"
	"os"

	"charm.land/lipgloss/v2"
	goversion "github.com/caarlos0/go-version"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var (
	boldStyle = lipgloss.NewStyle().Bold(true)
	codeStyle = lipgloss.NewStyle().Italic(true)
)

type rootCmd struct {
	cmd     *cobra.Command
	verbose bool
	exit    func(int)
}

func Execute(version goversion.Info, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := fang.Execute(
		context.Background(),
		cmd.cmd,
		fang.WithVersion(cmd.cmd.Version),
		fang.WithColorSchemeFunc(fang.AnsiColorScheme),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		cmd.exit(1)
	}
}

func newRootCmd(version goversion.Info, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}

	cmd := &cobra.Command{
		Use:               "splitter",
		Version:           version.String(),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Example: `
# Initialize:
splitter init

# Run
splitter start

# Run with config
splitter start --config my-config.yaml
		`,
	}

	cmd.SetVersionTemplate("{{.Version}}")

	cmd.AddCommand(
		newInitCmd().cmd,
		newStartCmd().cmd,
	)

	root.cmd = cmd

	return root
}
