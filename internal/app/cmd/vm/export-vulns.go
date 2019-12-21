package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

//ExportVulnsStart begin a download with a date since
func (vm *VM) ExportVulnsStart(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()
	since := vm.Config.VM.AfterDate

	a := client.NewAdapter(vm.Config, vm.Metrics)

	uuid, err := a.ExportVulnsStart()
	if err != nil {
		log.Errorf("error: couldn't start export-vulns: %v", err)
		return
	}

	log.Infof("successfully started export-vulns: %s with since date '%s' ", uuid, since)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportVulnsStart", map[string]string{"ExportUUID": uuid, "Since": since}))

	folder := filepath.Join(a.Config.VM.CacheFolder, "service", "export", "vulns", uuid)
	log.Infof("Creating folder: %s", folder)
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		log.Errorf("Could't make cache folder for future status lookup: %s", err)
	}

	return
}

//ExportVulnsStatus get the status for download uuid
func (vm *VM) ExportVulnsStatus(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID

	if uuid == "" {
		a.Config.VM.Log.Infof("export uuid was not specified - will use attempt to lookup from last 'start' call")

		var err error
		uuid, err = a.LastCachedExportVulnUUID()
		if err != nil {
			vm.Config.VM.Log.Errorf("error: cannot get export uuid: %v", err)
			return
		}
	}

	// We will use the cache hit to check for "FINISHED"
	// Check the cached response first
	status, err := a.ExportVulnsStatus(uuid, true, true)
	if err != nil {
		log.Errorf("error: couldn't get status export-vulns: %v", err)
		return
	}
	// If the status isn't FINISHED, ask for another from the server
	if status.Status != "FINISHED" {
		status, err = a.ExportVulnsStatus(uuid, false, true)
		if err != nil {
			log.Errorf("error: couldn't status export-vulns: %v", err)
			return
		}
	}

	log.Infof("successfully got status export-vulns UUID='%s' status='%s' ", uuid, status)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportVulnsStatus", map[string]string{"ExportUUID": uuid, "Status": status.Status}))

	return
}

//ExportVulnsGet for chunks files for download uuid
func (vm *VM) ExportVulnsGet(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID
	chunks := vm.Config.VM.Chunk

	if uuid == "" {
		a.Config.VM.Log.Infof("export uuid was not specified - will use attempt to lookup from last 'start' call")
		var err error
		uuid, err = a.LastCachedExportVulnUUID()
		if err != nil {
			return
		}
	}

	// Fetches all of the chunks - this can e long running and return
	// large amounts of data. Files are stored in the client cache and
	// and can be copied out using export-vulns --jq=. | gzip > all.vulns.json.gz
	err := a.ExportVulnsGet(uuid, chunks)
	if err != nil {
		log.Errorf("error: couldn't get export-vulns: %v", err)
		return
	}

	log.Infof("successfully got export-vulns UUID='%s' chunks: %s", uuid, chunks)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportVulnsGet", map[string]string{"ExportUUID": uuid, "Chunks": chunks}))

	return
}

//ExportVulnsQuery looks for matching jqex in chunks files for uuid
func (vm *VM) ExportVulnsQuery(cmd *cobra.Command, args []string) {
	vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID
	chunks := vm.Config.VM.Chunk
	jqex := vm.Config.VM.JQex

	if jqex == "" {
		jqex = ".[]"
		a.Config.VM.Log.Infof("query --jqex was not specified - will use default '%s", jqex)
	}

	before := vm.Config.VM.BeforeDate
	after := vm.Config.VM.AfterDate

	between := fmt.Sprintf(`select( .last_found >= "%s" and .last_found <= "%s" )`, after, before)

	var sevs []string
	if vm.Config.VM.Critical == true {
		sevs = append(sevs, `.severity == "critical"`)
	}
	if vm.Config.VM.High == true {
		sevs = append(sevs, `.severity == "high"`)
	}
	if vm.Config.VM.Medium == true {
		sevs = append(sevs, `.severity == "medium"`)
	}
	if vm.Config.VM.Info == true {
		sevs = append(sevs, `.severity == "info"`)
	}

	// If we have severity limits, add them to the between clause
	if len(sevs) > 0 {
		sev := strings.Join(sevs, " or ")
		between = fmt.Sprintf("select( (%s) and (%s) )", between, sev)
	}

	jqex = ".[]|" + between + "|" + jqex

	if uuid == "" {
		a.Config.VM.Log.Infof("export uuid was not specified - will use attempt to lookup from last 'start' call")
		var err error
		uuid, err = a.LastCachedExportVulnUUID()
		if err != nil {
			return
		}
	}

	a.Config.VM.Log.Debugf("Exporting JSON for '%s', chunks='%s' with JQex=`%s`", uuid, chunks, jqex)

	_ = a.ExportVulnsQuery(uuid, chunks, jqex)

	defer os.Remove(a.Config.VM.JQExec)

	return
}

// ExportVulnsHelp outputs the help template
func (vm *VM) ExportVulnsHelp(cmd *cobra.Command, args []string) {
	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	if vm.Config.VM.Log.IsLevelEnabled(log.DebugLevel) {
		fmt.Println(spew.Print(vm.Config))
	}

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("exportVulnsUsage", nil))

	return
}
