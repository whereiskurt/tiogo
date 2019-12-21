package vm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// ExportAssetsStart is invoked by Cobra with commandline args passed.
func (vm *VM) ExportAssetsStart(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)

	uuid, err := a.ExportAssetsStart()
	if err != nil || uuid == "" {
		log.Errorf("error: couldn't start export-assets: %v", err)
		return
	}

	// Size of the export limit for records
	size := vm.Config.VM.ExportLimit
	log.Infof("successfully started export-assets: %s with using chunk_size of '%s' ", uuid, size)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportAssetsStart", map[string]string{"ExportUUID": uuid, "Limit": size}))

	folder := filepath.Join(a.Config.VM.CacheFolder, "service", "export", "assets", uuid)
	log.Infof("Creating folder: %s", folder)
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		log.Errorf("Could't make cache folder for future status lookup: %s", err)
	}

	return
}

// ExportAssetsStatus is invoked by Cobra with commandline args passed.
func (vm *VM) ExportAssetsStatus(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID

	if uuid == "" {
		log.Infof("export uuid was not specified - will use attempt to lookup from last 'start' call")

		var err error
		uuid, err = a.LastCachedExportAssetUUID()
		if err != nil {
			log.Errorf("error: cannot get export uuid: %v", err)
			return
		}
	}

	// We will use the cache hit to check for "FINISHED"
	// Check the cached response first
	status, err := a.ExportAssetsStatus(uuid, true, true)
	if err != nil {
		log.Errorf("error: couldn't status export-assets: %v", err)
		return
	}

	// If the status isn't FINISHED, ask for another from the server
	if status.Status != "FINISHED" {
		log.Infof("making export-assets call UUID='%s' status='%s' ", uuid, status.Status)
		status, err = a.ExportAssetsStatus(uuid, false, true)
		if err != nil {
			log.Errorf("error: couldn't status export-assets: %v", err)
			return
		}
	}

	log.Infof("successfully got status export-assets UUID='%s' status='%s' ", uuid, status)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportAssetsStatus", map[string]string{"ExportUUID": uuid, "Status": status.Status}))

	return
}

// ExportAssetsGet is invoked by Cobra with commandline args passed.
func (vm *VM) ExportAssetsGet(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID
	chunks := vm.Config.VM.Chunk

	if uuid == "" {
		a.Config.VM.Log.Infof("export uuid was not specified - will use attempt to lookup from last 'start' call")
		var err error
		uuid, err = a.LastCachedExportAssetUUID()
		if err != nil {
			return
		}
	}

	// Fetches all of the chunks - this can e long running and return
	// large amounts of data. Files are stored in the client cache and
	// and can be copied out using export-assets --jq=. | gzip > all.assets.json.gz
	err := a.ExportAssetsGet(uuid, chunks)
	if err != nil {
		log.Errorf("error: couldn't get export-assets: %v", err)
		return
	}

	log.Infof("successfully got export-assets UUID='%s' chunks: %s", uuid, chunks)

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportAssetsGet", map[string]string{"ExportUUID": uuid, "Chunks": chunks}))

	return
}

// ExportAssetsQuery is invoked by Cobra with commandline args passed.
func (vm *VM) ExportAssetsQuery(cmd *cobra.Command, args []string) {
	vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)
	uuid := vm.Config.VM.UUID
	chunks := vm.Config.VM.Chunk
	jqex := vm.Config.VM.JQex

	if jqex == "" {
		jqex = ".[]"
		a.Config.VM.Log.Infof("query --jqex was not specified - will use default '%s", jqex)
	}

	// before := vm.Config.VM.BeforeDate
	// after := vm.Config.VM.AfterDate
	// between := fmt.Sprintf(`select( .last_found >= "%s" and .last_found <= "%s" )`, after, before)
	//
	// var sevs []string
	// if vm.Config.VM.Critical == true {
	// 	sevs = append(sevs, `.severity == "critical"`)
	// }
	// if vm.Config.VM.High == true {
	// 	sevs = append(sevs, `.severity == "high"`)
	// }
	// if vm.Config.VM.Medium == true {
	// 	sevs = append(sevs, `.severity == "medium"`)
	// }
	// if vm.Config.VM.Info == true {
	// 	sevs = append(sevs, `.severity == "info"`)
	// }
	// if len(sevs) > 0 {
	// 	sev := strings.Join(sevs, " or ")
	// 	between = fmt.Sprintf("select( (%s) and (%s) )", between, sev)
	// }
	//
	// jqex = ".[]|" + between + "|" + jqex

	jqex = ".|" + jqex

	if uuid == "" {
		a.Config.VM.Log.Infof("export uuid was not specified - will use attempt to lookup from last 'start' call")
		var err error
		uuid, err = a.LastCachedExportAssetUUID()
		if err != nil {
			return
		}
	}

	a.Config.VM.Log.Debugf("Exporting JSON for '%s', chunks='%s' with JQex=`%s`", uuid, chunks, jqex)

	_ = a.ExportAssetsQuery(uuid, chunks, jqex)

	defer os.Remove(a.Config.VM.JQExec)

	return
}

// ExportAssetsHelp is invoked by Cobra with commandline args passed.
func (vm *VM) ExportAssetsHelp(cmd *cobra.Command, args []string) {
	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	if vm.Config.VM.Log.IsLevelEnabled(log.DebugLevel) {
		fmt.Println(spew.Print(vm.Config))
	}

	cli := ui.NewCLI(vm.Config)
	fmt.Println(cli.Render("ExportAssetsHelp", nil))

	return
}
