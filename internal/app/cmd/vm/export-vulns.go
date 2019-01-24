package vm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

func (vm *VM) ExportVulnsStart(cmd *cobra.Command, args []string) {

}
func (vm *VM) ExportVulnsStatus(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	log.Debugf("Export Vulns Status: %+v", "WOOOOOOOOT!")

}

func (vm *VM) ExportVulnsHelp(cmd *cobra.Command, args []string) {

	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	cli := ui.NewCLI(vm.Config)
	if len(args) == 0 {
		cli.DrawGopher()
		fmt.Println(cli.Render("exportVulnsUsage", nil))
		return
	}
}
