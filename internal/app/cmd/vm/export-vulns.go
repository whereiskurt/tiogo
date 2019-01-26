package vm

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

func (vm *VM) ExportVulnsStart(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)

	uuid, err := a.VulnsExportStart()
	if err != nil {
		log.Errorf("error: couldn't start export-vulns: %v", err)
		return
	}

	log.Infof("successfully started export-vulns: %s", uuid)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("exportVulnsStarted", map[string]string{"ExportUUID": uuid}))

	return
}

func (vm *VM) ExportVulnsStatus(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()
	log.Debug("ExportVulnsStatus")

	return
}

func (vm *VM) ExportVulnsHelp(cmd *cobra.Command, args []string) {
	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	if vm.Config.Log.IsLevelEnabled(log.DebugLevel) {
		fmt.Println(spew.Print(vm.Config))
	}

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("exportVulnsUsage", nil))

	return
}
