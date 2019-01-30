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
	fmt.Println(cli.Render("ExportVulnsStart", map[string]string{"ExportUUID": uuid}))

	return
}

func (vm *VM) ExportVulnsStatus(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID

	if uuid == "" {
		log.Errorf("error: uuid not specified. Use '--uuid=XYZ' to specify the uuid of the vulnerability export")
		return
	}

	status, err := a.VulnsExportStatus(uuid)
	if err != nil {
		log.Errorf("error: couldn't status export-vulns: %v", err)
		return
	}

	log.Infof("successfully got status export-vulns UUID='%s' status='%s' ", uuid, status)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportVulnsStatus", map[string]string{"ExportUUID": uuid, "Status": status.Status}))

	return
}
func (vm *VM) ExportVulnsGet(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID
	chunks := vm.Config.VM.Chunk

	if uuid == "" {
		log.Errorf("error: uuid not specified. Use '--uuid=XYZ' to specify the uuid of the vulnerability export")
		return
	}

	if chunks == "" {
		chunks = "ALL"
		log.Infof("info:  Using --chunk='%s' -- no chunk range/value specified.", chunks)
	}

	// Fetches all of the chunks - this can e long running and return
	// large amounts of data. Files are stored in the client cache and
	// and can be copied out using export-vulns --jq=. | gzip > all.vulns.json.gz
	err := a.VulnsExportGet(uuid, chunks)
	if err != nil {
		log.Errorf("error: couldn't get export-vulns: %v", err)
		return
	}

	log.Infof("successfully got export-vulns UUID='%s' chunks: %s", uuid, chunks)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportVulnsGet", map[string]string{"ExportUUID": uuid, "Chunks": chunks}))

	return
}
func (vm *VM) ExportVulnsQuery(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID
	chunks := vm.Config.VM.Chunk

	jqex := vm.Config.VM.JQex
	if jqex == "" {
		jqex = "."
	}

	if uuid == "" {
		log.Errorf("error: uuid not specified. Use '--uuid=XYZ' to specify the uuid of the vulnerability export")
		return
	}

	if chunks == "" {
		chunks = "ALL"
		log.Infof("info:  Using --chunk='%s' -- no chunk range/value specified.", chunks)
	}

	a.VulnsExportQuery(uuid, chunks, jqex)

	return
}

func (vm *VM) ExportVulnsHelp(cmd *cobra.Command, args []string) {
	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	if vm.Config.Log.IsLevelEnabled(log.DebugLevel) {
		fmt.Println(spew.Print(vm.Config))
	}

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportVulnsHelp", nil))

	return
}
