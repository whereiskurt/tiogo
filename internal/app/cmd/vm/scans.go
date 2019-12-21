package vm

import (
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// ScansList is invoked by Cobra with commandline args passed.
func (vm *VM) ScansList(cmd *cobra.Command, args []string) {

	log := vm.setupLog()
	log.Infof("tiogo scanners list command:")

	cli := ui.NewCLI(vm.Config)

	a := client.NewAdapter(vm.Config, vm.Metrics)

	scans, err := a.Scans(true, true)
	if err != nil {
		log.Errorf("error: couldn't scans list: %v", err)
		return
	}

	if a.Config.VM.OutputJSON {
		// No JSON output yet...
	}

	if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		cli.Println(cli.Render("ScansListCSV", map[string]interface{}{"Scans": scans}))
	}

	return
}

// ScansDetail is invoked by Cobra with commandline args passed.
func (vm *VM) ScansDetail(cmd *cobra.Command, args []string) {
	return
}

// ScansHistory is invoked by Cobra with commandline args passed.
func (vm *VM) ScansHistory(cmd *cobra.Command, args []string) {
	return
}
