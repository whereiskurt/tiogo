package vm

import (
	"encoding/json"
	"fmt"

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

	//TODO: Make this a method :-P
	id := vm.Config.VM.ID
	uuid := vm.Config.VM.UUID
	name := vm.Config.VM.Name
	regex := vm.Config.VM.Regex
	if id != "" {
		scans = a.Filter.ScanByID(scans, id)
	} else if uuid != "" {
		scans = a.Filter.ScanByScheduleUUID(scans, uuid)
	} else if name != "" {
		scans = a.Filter.ScanByName(scans, name)
	} else if regex != "" {
		scans = a.Filter.ScanByRegex(scans, regex)
	}

	if a.Config.VM.OutputJSON {
		// Convert structs to JSON.
		data, err := json.Marshal(scans)
		if err != nil {
			log.Fatalf("error: couldn't marshal scan data to JSON: %v", err)
		}
		cli.Println(fmt.Sprintf("%s\n", data))

	} else if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		cli.Println(cli.Render("ScansListCSV", map[string]interface{}{"Scans": scans}))
	}

	return
}

// ScansDetail is invoked by Cobra with commandline args passed.
func (vm *VM) ScansDetail(cmd *cobra.Command, args []string) {
	log := vm.setupLog()
	log.Infof("tiogo scan detail command:")

	cli := ui.NewCLI(vm.Config)
	a := client.NewAdapter(vm.Config, vm.Metrics)

	scans, err := a.Scans(true, true)
	if err != nil {
		log.Errorf("error: couldn't scans list: %v", err)
		return
	}

	id := vm.Config.VM.ID
	uuid := vm.Config.VM.UUID
	name := vm.Config.VM.Name
	regex := vm.Config.VM.Regex
	if id != "" {
		scans = a.Filter.ScanByID(scans, id)
	} else if uuid != "" {
		scans = a.Filter.ScanByScheduleUUID(scans, uuid)
	} else if name != "" {
		scans = a.Filter.ScanByName(scans, name)
	} else if regex != "" {
		scans = a.Filter.ScanByRegex(scans, regex)
	}

	if len(scans) == 0 {
		log.Errorf("error: couldn't match a scans")
		return
	}

	cli.DrawGopher()

	for _, s := range scans {
		details, err := a.ScanDetails(s, true, true)
		if err != nil {
			log.Fatalf("error: couldn't retrieve details: %v", err)
		}
		cli.Println(fmt.Sprintf("details:\n%+v\n", details))
	}

	return
}

// ScansHosts is invoked by Cobra with commandline args passed.
func (vm *VM) ScansHosts(cmd *cobra.Command, args []string) {
	return
}

// ScansPlugins is invoked by Cobra with commandline args passed.
func (vm *VM) ScansPlugins(cmd *cobra.Command, args []string) {
	return
}

// ScansQuery is invoked by Cobra with commandline args passed.
func (vm *VM) ScansQuery(cmd *cobra.Command, args []string) {
	return
}
