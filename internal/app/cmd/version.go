package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
)

type Version struct {
	Config *app.Config
	CLI    *ui.CLI
}

func NewVersion(c *app.Config) (v *Version) {
	v = new(Version)
	v.Config = c
	return
}

func (v *Version) Version(cmd *cobra.Command, args []string) {

	cli := ui.NewCLI(v.Config)

	cli.Draw.Version()

	return
}
