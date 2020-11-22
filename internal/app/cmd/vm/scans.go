package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// ScansList is invoked by Cobra with commandline args passed.
func (vm *VM) ScansList(cmd *cobra.Command, args []string) {
	logger := vm.setupLog()
	cli := ui.NewCLI(vm.Config)

	a := client.NewAdapter(vm.Config, vm.Metrics)

	scans, err := a.Scans(true, true)
	if err != nil {
		logger.Fatalf("error: couldn't scans list: %v", err)
	}
	scans = vm.FilterScans(a, &scans)

	if a.Config.VM.OutputJSON {
		// Convert structs to JSON.
		data, err := json.Marshal(scans)
		if err != nil {
			logger.Fatalf("error: couldn't marshal scan data to JSON: %v", err)
		}
		cli.Println(fmt.Sprintf("%s\n", data))

	} else if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		vm.CleanupFiles(`.`, `scanlist\.\d+T\d+.csv`, 2)

		dts := time.Now().Format("20060102T150405")
		filename := fmt.Sprintf("scanlist.%s.csv", dts)
		csv := cli.Render("ScansListCSV", map[string]interface{}{"Scans": scans})

		// NOTE: Using ioutil.WriteFile is OK for smaller files
		err = ioutil.WriteFile(filename, []byte(csv), 0644)
		if err != nil {
			logger.Fatalf("can't write file: %+v", err)
		}
	}

	return
}

// FilterScans uses cli arguments to reduce scans to filtered
func (vm *VM) FilterScans(a *client.Adapter, scans *[]client.Scan) (filtered []client.Scan) {
	//TODO: Make this a method :-P
	id := vm.Config.VM.ID
	uuid := vm.Config.VM.UUID
	name := vm.Config.VM.Name
	regex := vm.Config.VM.Regex
	histid := vm.Config.VM.HistoryID

	if id != "" {
		filtered = a.Filter.ScanByID(*scans, id)
	} else if uuid != "" {
		filtered = a.Filter.ScanByScheduleUUID(*scans, uuid)
	} else if name != "" {
		filtered = a.Filter.ScanByName(*scans, name)
	} else if regex != "" {
		filtered = a.Filter.ScanByRegex(*scans, regex)
	} else {
		filtered = append(filtered, *scans...)
	}

	if histid != "" {
		var reduced []client.Scan
		for _, s := range filtered {
			det, err := a.ScanDetails(&s, true, true)
			if err != nil {
				continue
			}
			for _, h := range det.History {
				if h.HistoryID == histid {
					reduced = append(reduced, s)
					break
				}
			}

		}
		filtered = reduced
	}

	return
}

// ScansDetail is invoked by Cobra with commandline args passed.
func (vm *VM) ScansDetail(cmd *cobra.Command, args []string) {
	logger := vm.setupLog()
	logger.Infof("tiogo scan detail command:")

	cli := ui.NewCLI(vm.Config)

	a := client.NewAdapter(vm.Config, vm.Metrics)

	scans, err := a.Scans(true, true)
	if err != nil {
		logger.Errorf("error: couldn't scans list: %v", err)
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
		logger.Errorf("error: couldn't match a scans")
		return
	}

	for i := range scans {
		details, err := a.ScanDetails(&scans[i], true, true)
		if err != nil {
			logger.Fatalf("error: couldn't retrieve details: %v: %+v", scans[i], err)
		}
		cli.Println(fmt.Sprintf(" Name:\t\t\t%+v", details.Scan.Name))
		cli.Println(fmt.Sprintf(" Enabled:\t\t%+v", details.Scan.Enabled))
		cli.Println(fmt.Sprintf(" ScanID:\t\t%+v", details.Scan.ScanID))
		cli.Println(fmt.Sprintf(" ScheduleUUID:\t\t%+v", details.Scan.ScheduleUUID))
		cli.Println(fmt.Sprintf(" ScannerName:\t\t%+v", details.ScannerName))
		cli.Println(fmt.Sprintf(" ScanType:\t\t%+v", details.ScanType))
		cli.Println(fmt.Sprintf(" Target:\t\t%+v", details.Targets))
		cli.Println(fmt.Sprintf(" AgentGroups:\t\t%+v", details.AgentGroup))
		cli.Println(fmt.Sprintf(" HistoryCount:\t\t%+v", details.HistoryCount))
		cli.Println(fmt.Sprintf(" Last Status:\t\t%+v", details.Status))
		cli.Println(fmt.Sprintf(" Last StartTime:\t%+v", details.ScanStart))
		cli.Println(fmt.Sprintf(" Last EndTime:\t\t%+v", details.ScanEnd))
		cli.Println(fmt.Sprintf(" Timestamp:\t\t%+v", details.Timestamp))
		cli.Println(fmt.Sprintf(" RRules:\t\t%+v", details.Scan.RRules))
		cli.Println(fmt.Sprintf(" Hosts Scanned:\t\t%+v", details.HostCount))
		cli.Println(fmt.Sprintf(" CriticalCount:\t\t%+v", details.PluginCriticalCount))
		cli.Println(fmt.Sprintf(" HighCount:\t\t%+v", details.PluginHighCount))
		cli.Println(fmt.Sprintf(" MediumCount:\t\t%+v", details.PluginMediumCount))
		cli.Println(fmt.Sprintf(" LowCount:\t\t%+v", details.PluginLowCount))
		cli.Println(fmt.Sprintf(" InfoCount:\t\t%+v", details.PluginInfoCount))
		cli.Println(fmt.Sprintf(" ====================================\n\t\t\t%+v total", details.PluginTotalCount))
		cli.Println(fmt.Sprintf(" \n--------------------------------------------\n"))
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
