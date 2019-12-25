package vm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// CacheInfo handles the .tiogo/[client|server] folders
func (vm *VM) CacheInfo(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf(fmt.Sprintf("CacheInfo Requested started..."))

	cli := ui.NewCLI(vm.Config)

	cli.Println(fmt.Sprintf("An interface into the Tenable.io API using Go!	"))
	cli.Println(fmt.Sprintf("Version %s %s ", vm.Config.VM.ReleaseVersion, vm.Config.VM.GitHash))
	cli.Println("")
	cli.Println(fmt.Sprintf("Local Cache Information"))
	cli.DrawGopher()
	cli.Println(fmt.Sprintf("  Client Cache Folder: %s", vm.Config.VM.CacheFolder))
	cli.Println(fmt.Sprintf("  Server Cache Folder: %s", vm.Config.Server.CacheFolder))
	cli.Println("")
	return
}

// CacheClear will empty the appropriate cache folders
func (vm *VM) CacheClear(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("CacheClear Requested with no parameters - showing help ...")

	// Render help when nothing is passed in
	vm.Help(cmd, []string{"cache"})

	return
}

// CacheClearAll will empty ALL cache folders
func (vm *VM) CacheClearAll(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("Clear ALL cache entries...")

	cli := ui.NewCLI(vm.Config)
	cli.DrawGopher()
	cli.Println(fmt.Sprintf("Working to delete ALL cache entries..."))

	//TODO: Add a '--all' flag or prompt for "Are you sure? (Y/N)" prompt before deletes
	vm.CacheClearAgents(cmd, args)
	vm.CacheClearScans(cmd, args)
	vm.CacheClearExports(cmd, args)

	cli.Println(fmt.Sprintf("Done."))

	return
}

// CacheClearExports will empty  folders for all past exports (vulns/assets/etc.)
func (vm *VM) CacheClearExports(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("Clearing cached Exports requested ...")

	vm.clearExports(tenable.EndPoints.VulnsExportStart)
	vm.clearExports(tenable.EndPoints.AssetsExportStart)
}

func (vm *VM) clearExports(e tenable.EndPointType) {
	log := vm.setupLog()
	endPoint := tenable.ServiceMap[e].CacheFilename

	// Client Cache folder
	vpath := filepath.Join(vm.Config.VM.CacheFolder, client.DefaultServiceFolder, endPoint)
	vdir := filepath.Dir(vpath)
	log.Debugf("Delete cache folder: os.RemoveAll(%s)", vdir)
	os.RemoveAll(vdir)

	// Server Cache folder
	vpath = filepath.Join(vm.Config.Server.CacheFolder, endPoint)
	vdir = filepath.Dir(vpath)
	log.Debugf("Delete cache folder: os.RemoveAll(%s)", vdir)
	os.RemoveAll(vdir)
}

// CacheClearAgents will empty just folders for Agents
func (vm *VM) CacheClearAgents(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("Clearing Agents requested ...")

	scannerEndPoint := tenable.ServiceMap[tenable.EndPoints.ScannersList].CacheFilename

	cpath := filepath.Join(vm.Config.VM.CacheFolder, client.DefaultServiceFolder, scannerEndPoint)
	cdir := filepath.Dir(cpath)
	log.Debugf("Delete cache folder: os.RemoveAll(%s)", cdir)
	os.RemoveAll(cdir)

	cpath = filepath.Join(vm.Config.Server.CacheFolder, scannerEndPoint)
	cdir = filepath.Dir(cpath)
	log.Debugf("Delete cache folder: os.RemoveAll(%s)", cdir)
	os.RemoveAll(cdir)

	return
}

// CacheClearScans will empty just folders for Scans
func (vm *VM) CacheClearScans(cmd *cobra.Command, args []string) {
	log := vm.setupLog()

	log.Debugf("Clearing Scanners cache...")

	scannerEndPoint := tenable.ServiceMap[tenable.EndPoints.ScansList].CacheFilename

	cpath := filepath.Join(vm.Config.VM.CacheFolder, client.DefaultServiceFolder, scannerEndPoint)
	cdir := filepath.Dir(cpath)
	log.Debugf("Delete cache folder: os.RemoveAll(%s)", cdir)
	os.RemoveAll(cdir)

	cpath = filepath.Join(vm.Config.Server.CacheFolder, scannerEndPoint)
	cdir = filepath.Dir(cpath)
	log.Debugf("Delete cache folder: os.RemoveAll(%s)", cdir)
	os.RemoveAll(cdir)

	return
}
