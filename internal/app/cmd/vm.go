package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

var (
	// ReleaseVersion is set by a --ldflags during a build/release
	ReleaseVersion = "v0.0.1-development"
	// GitHash is set by a --ldflags during a build/release
	GitHash = "0xhashhash"
)

// Version holds the config and CLI references.
type VM struct {
	Config *config.Config
}

// Version just outputs a gopher.
func (vm *VM) Help(cmd *cobra.Command, args []string) {
	fmt.Printf(spew.Sprintf("%v", args))

	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	cli := ui.NewCLI(vm.Config)
	cli.DrawGopher()
	return
}

func (vm *VM) Scanners(cmd *cobra.Command, args []string) {
	vm.Config.Log.SetFormatter(&log.TextFormatter{})
	vm.Config.VM.EnableLogging()

	vm.Config.Log.Infof("tiogo scanners list command:")

	cli := ui.NewCLI(vm.Config)
	cli.DrawGopher()
	return
}

// NewVersion holds a configuration and command line interface reference (for log out, etc.)
func NewVM(c *config.Config) (v VM) {
	v.Config = c
	return
}
