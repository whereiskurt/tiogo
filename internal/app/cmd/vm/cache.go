package vm

import "github.com/spf13/cobra"

import "github.com/whereiskurt/tiogo/pkg/tenable"

import "path/filepath"

import "os"

import "github.com/whereiskurt/tiogo/pkg/client"

// CacheInfo handles the .tiogo/[client|server] folders
func (vm *VM) CacheInfo(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("CacheInfo Requested started...")
	log.Debugf("  Client Folder: %s", vm.Config.VM.CacheFolder)
	log.Debugf("  Server Folder: %s", vm.Config.Server.CacheFolder)
	log.Debugf("")

	return
}

// CacheClear will empty the appropriate cache folders
func (vm *VM) CacheClear(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("CacheClear Requested started...")
	return
}

// CacheClearAll will empty ALL cache folders
func (vm *VM) CacheClearAll(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("CacheClear Requested started...")
	return
}

// CacheClearAgents will empty just folders for Agents
func (vm *VM) CacheClearAgents(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Debugf("Clearing Agents requested ...")

	scannerEndPoint := tenable.ServiceMap[tenable.EndPoints.ScannersList].CacheFilename

	cpath := filepath.Join(vm.Config.VM.CacheFolder, client.DefaultServiceFolder, scannerEndPoint)
	cdir := filepath.Dir(cpath)
	log.Debugf("  Removing Client Dir: %s", cdir)
	os.RemoveAll(cdir)

	spath := filepath.Join(vm.Config.Server.CacheFolder, scannerEndPoint)
	sdir := filepath.Dir(spath)
	log.Debugf("  Removing Server Dir: %s", sdir)
	os.RemoveAll(sdir)

	return
}

// CacheClearScans will empty just folders for Scans
func (vm *VM) CacheClearScans(cmd *cobra.Command, args []string) {
	log := vm.setupLog()

	log.Debugf("Clearing Scanners cache...")

	return
}
