package vm

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

func (vm *VM) ScannerList(cmd *cobra.Command, args []string) {
	vm.Config.Log.SetFormatter(&log.TextFormatter{})
	vm.Config.VM.EnableLogging()

	vm.Config.Log.Infof("tiogo scanners list command:")

	cli := ui.NewCLI(vm.Config)
	cli.DrawGopher()
	return
}
