package vm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

func (vm *VM) ScannersHelp(cmd *cobra.Command, args []string) {
	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("scannersUsage", nil))
	return
}

func (vm *VM) ScannersList(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Infof("tiogo scanners list command:")
	cli := ui.NewCLI(vm.Config)
	cli.DrawGopher()
	return
}

func (vm *VM) setupLog() *log.Logger {
	vm.Config.Log.SetFormatter(&log.TextFormatter{})
	vm.Config.VM.EnableLogging()
	return vm.Config.Log
}
