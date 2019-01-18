package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

var (
	// ReleaseVersion is set by a --ldflags during a build/release
	ReleaseVersion = "v1.0.0-development"
	// GitHash is set by a --ldflags during a build/release
	GitHash = "0xhashhash"
)

// Version holds the config and CLI references.
type VM struct {
	Config *config.Config
}

// Version just outputs a gopher.
func (v *VM) Help(cmd *cobra.Command, args []string) {
	fmt.Printf(spew.Sprintf("%v", args))

	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	cli := ui.NewCLI(v.Config)
	cli.DrawGopher()
	return
}

func (v *VM) Scanners(cmd *cobra.Command, args []string) {
	fmt.Printf("scanners!")

	cli := ui.NewCLI(v.Config)
	cli.DrawGopher()
	return
}

// NewVersion holds a configuration and command line interface reference (for log out, etc.)
func NewVM(c *config.Config) (v VM) {
	v.Config = c
	return
}
