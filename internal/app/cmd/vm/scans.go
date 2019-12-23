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

	for i := range scans {
		details, err := a.ScanDetails(&scans[i], true, true)
		if err != nil {
			log.Fatalf("error: couldn't retrieve details: %v: %+v", scans[i], err)
		}
		cli.Println(fmt.Sprintf(" Name:\t\t\t%+v", details.Scan.Name))
		cli.Println(fmt.Sprintf(" ScanID:\t\t%+v", details.Scan.ScanID))
		cli.Println(fmt.Sprintf(" ScheduleUUID:\t\t%+v", details.Scan.ScheduleUUID))
		cli.Println(fmt.Sprintf(" ScanType:\t\t%+v", details.ScanType))
		cli.Println(fmt.Sprintf(" Last Status:\t\t%+v", details.Status))
		cli.Println(fmt.Sprintf(" Last StartTime:\t%+v", details.ScanStart))
		cli.Println(fmt.Sprintf(" Last EndTime:\t\t%+v", details.ScanEnd))
		cli.Println(fmt.Sprintf(" Timestamp:\t\t%+v", details.Timestamp))
		cli.Println(fmt.Sprintf(" RRules:\t\t%+v", details.Scan.RRules))
		cli.Println(fmt.Sprintf(" HistoryCount:\t\t%+v", details.HistoryCount))
		cli.Println(fmt.Sprintf(" HostCount:\t\t%+v", details.HostCount))
		cli.Println(fmt.Sprintf(" PluginTotalCount:\t%+v", details.PluginTotalCount))
		cli.Println(fmt.Sprintf(" PluginCriticalCount:\t%+v", details.PluginCriticalCount))
		cli.Println(fmt.Sprintf(" PluginHighCount:\t%+v", details.PluginHighCount))
		cli.Println(fmt.Sprintf(" PluginMediumCount:\t%+v", details.PluginMediumCount))
		cli.Println(fmt.Sprintf(" PluginLowCount:\t%+v", details.PluginLowCount))
		cli.Println(fmt.Sprintf(" ----------------------\n"))
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
