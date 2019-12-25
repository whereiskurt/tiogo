package vm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// ScannersList is invoked by Cobra with commandline args passed.
func (vm *VM) ScannersList(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Infof("tiogo scanners list command:")
	cli := ui.NewCLI(vm.Config)

	a := client.NewAdapter(vm.Config, vm.Metrics)

	scanners, err := a.Scanners(true, true)
	if err != nil {
		log.Errorf("error: couldn't scanners list: %v", err)
		return
	}

	if a.Config.VM.OutputJSON {
		// No JSON output yet...
	}

	if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		fmt.Println(cli.Render("ScannersListCSV", map[string]interface{}{"Scanners": scanners}))
	}

	return
}

func (vm *VM) setupLog() *log.Logger {
	vm.Config.VM.Log.SetFormatter(&log.TextFormatter{})
	vm.Config.VM.EnableLogging()
	return vm.Config.VM.Log
}
