package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// ScansList is invoked by Cobra with commandline args passed.
func (vm *VM) ScansList(cmd *cobra.Command, args []string) {
	logger := vm.setupLog()
	cli := ui.NewCLI(vm.Config)

	logger.Infof("Starting scan list ...")
	a := client.NewAdapter(vm.Config, vm.Metrics)

	maxkeep, err := strconv.Atoi(vm.Config.VM.MaxKeep)
	if err != nil {
		logger.Fatalf("error: couldn't convert maxkeep '%s': %v", vm.Config.VM.MaxKeep, err)
	}

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

		dts := time.Now().Format("20060102T150405")
		filename := fmt.Sprintf("scanlist.%s.csv", dts)
		csv := cli.Render("ScansListCSV", map[string]interface{}{"Scans": scans})

		// NOTE: Using ioutil.WriteFile is OK for smaller files (less than 100MBs)
		err = ioutil.WriteFile(filename, []byte(csv), 0644)
		if err != nil {
			logger.Fatalf("can't write to file '%s': %+v", filename, err)
		}

		cleanTemplate := `scanlist\.\d+T\d+.csv`
		logger.Infof("keeping a maximum '%d' for template '%s'", maxkeep, cleanTemplate)
		vm.CleanupFiles(`.`, cleanTemplate, maxkeep)
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

// ScansGet starts, status loops, gets scan export - simplified export-scans
func (vm *VM) ScansGet(cmd *cobra.Command, args []string) {
	a := client.NewAdapter(vm.Config, vm.Metrics)

	logger := vm.setupLog()
	scans, err := a.Scans(true, true)
	if err != nil {
		logger.Fatalf("error: couldn't scans list: %v", err)
	}
	scans = vm.FilterScans(a, &scans)

	logger.Infof("Starting scan get ...")

	histid := vm.Config.VM.HistoryID

	var format = "csv"
	var chapters = vm.Config.VM.Chapters

	var offset int
	if vm.Config.VM.Offset != "" && histid == "" {
		offset, err = strconv.Atoi(vm.Config.VM.Offset)
		if err != nil {
			logger.Fatalf("error: couldn't convert offset '%s': %v", vm.Config.VM.Offset, err)
		}
	}

	var maxdepth int
	maxdepth, err = strconv.Atoi(vm.Config.VM.MaxDepth)
	if err != nil {
		logger.Fatalf("error: couldn't convert maxdepth '%s': %v", vm.Config.VM.MaxDepth, err)
	}

	var maxkeep int
	maxkeep, err = strconv.Atoi(vm.Config.VM.MaxKeep)
	if err != nil {
		logger.Fatalf("error: couldn't convert maxkeep '%s': %v", vm.Config.VM.MaxKeep, err)
	}

	if len(scans) == 0 {
		logger.Errorf("error: history id didn't match a scan: %+v", vm.Config)
		return
	} else if len(scans) > 1 && histid != "" {
		logger.Errorf("error: histid doesn't limit to one scan")
		return
	}
	// TODO: Functional decompose this massive block! :-)
SCANS:
	for _, s := range scans {
		det, err := a.ScanDetails(&s, true, true)
		if err != nil {
			logger.Errorf("error: can't retrieve scan details for %s at offset %d", s.ScanID, offset)
			continue SCANS
		} else if len(det.History) == 0 {
			logger.Warnf("warn: scan %v has not run yet", s.ScanID, offset)
			continue SCANS
		}

		// This scan has no history ie. no previous scans
		if len(det.History) <= offset {
			logger.Errorf("error: scan %v has less run histories than offset '%d' requires", s.ScanID, offset)
			continue SCANS
		}
	DEPTHS:
		for depth := 0; depth < maxdepth; depth++ {
			if len(det.History) <= offset+depth {
				logger.Infof("scan %s has only %d histories (%d doesn't exist)", s.ScanID, len(det.History), offset+depth)
				break DEPTHS
			}
			histid = det.History[offset+depth].HistoryID

			var tgtFilename = fmt.Sprintf("scan.%s.history.%s.csv", s.ScanID, histid)

			if _, err := os.Stat(tgtFilename); err == nil {
				logger.Infof("skipping scan:%s history:%s: file already downloaded: %s", s.ScanID, histid, tgtFilename)
				continue DEPTHS
			}

			_, err = a.ScansExportStart(&s, histid, format, chapters, true, true)
			if err != nil {
				logger.Errorf("error: cannot start export scans from Tenable.io, skipping scan: +v", err)
				continue SCANS
			}

			var sleepStatusCheckIntervals = []int{500, 1000, 2000, 2500, 3000, 3500, 5000, 5000, 5000, 10000, 20000, 20000, 30000, 30000, 30000, 30000}
			var maxattempts, sleptsec int
			export, err := a.ScansExportStatus(&s, histid, format, chapters, true, true) //Use cache=true,true don't clobe last READY with "ERROR"
			if err != nil {
				logger.Errorf("error: cannot get export status from Tenable.io, skipping scan: +v", err)
				continue SCANS
			}
		STATUSCHECK:
			for maxattempts = len(sleepStatusCheckIntervals); maxattempts >= 0; maxattempts-- {
				if export.Status == "READY" {
					break STATUSCHECK
				}
				export, err = a.ScansExportStatus(&s, histid, format, chapters, false, true)
				if err != nil {
					logger.Errorf("error: cannot get export status from Tenable.io, skipping scan: +v", err)
					continue SCANS
				}
				logger.Infof("scan %s uuid: %s download not ready (%s) sleeping %dms...", s.ScanID, export.FileUUID, export.Status, sleepStatusCheckIntervals[maxattempts-1])
				time.Sleep(time.Duration(sleepStatusCheckIntervals[maxattempts-1]) * time.Millisecond)
				sleptsec += sleepStatusCheckIntervals[maxattempts-1] / 1000
			}
			if maxattempts <= 0 {
				logger.Errorf(fmt.Sprintf("Status stuck on '%s' scan export '%s' after %d attempts and %dsecs, skipping scan.", export.Status, export.FileUUID, maxattempts, sleptsec))
				continue SCANS
			}

			logger.Infof("beginning download for %s uuid: %s", s.ScanID, export.FileUUID)
			filename, err := a.ScansExportLargeGet(&s, histid, format, chapters)
			if err != nil {
				logger.Errorf(fmt.Sprintf("Failed to Export-Scans Large Get '%s' scan export '%s' : %v, skipping scan.", export.Status, export.FileUUID, err))
				continue SCANS
			}

			//vm.copyToFile(filename, tgtFilename)
			logger.Infof("processing download with _time for %s uuid: %s", s.ScanID, export.FileUUID)
			err = vm.TimestampCSVRows(filename, tgtFilename)
			if err != nil {
				os.Remove(tgtFilename)
				logger.Fatalf("error: couldn't apply timestatmps write: %s to %s: %v, (out of disk space?)", filename, tgtFilename, err)
			}
		}

		//Keep only X historicals for this ScanId
		cleanTemplate := fmt.Sprintf(`scan\.%s\.history\.\d+\.csv`, s.ScanID)
		logger.Infof("keeping a maximum '%d' for template '%s'", maxkeep, cleanTemplate)
		vm.CleanupFiles(`.`, cleanTemplate, maxkeep)
	}

	return
}
