package vm

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"strconv"
)

type actionType string

type action struct {
	ExportScanStart  actionType
	ExportScanStatus actionType
	ExportScanGet    actionType
	ExportScanQuery  actionType
}

var actions = action{
	ExportScanStart:  actionType("ExportScanStart"),
	ExportScanStatus: actionType("ExportScanStatus"),
	ExportScanGet:    actionType("ExportScanGet"),
	ExportScanQuery:  actionType("ExportScanQuery"),
}

// ExportScansStart request Tenable.io create a new export for the scan
func (vm *VM) ExportScansStart(cmd *cobra.Command, args []string) {
	vm.exportScansAction(cmd, args, actions.ExportScanStart)
	return
}

// ExportScansStatus check export status of Tenable.io scan export
func (vm *VM) ExportScansStatus(cmd *cobra.Command, args []string) {
	vm.exportScansAction(cmd, args, actions.ExportScanStatus)
	return
}

// ExportScansGet download the exported status of Tenable.io scan export
func (vm *VM) ExportScansGet(cmd *cobra.Command, args []string) {
	vm.exportScansAction(cmd, args, actions.ExportScanGet)
	return
}

// ExportScansQuery download the exported status of Tenable.io scan export
func (vm *VM) ExportScansQuery(cmd *cobra.Command, args []string) {
	vm.exportScansAction(cmd, args, actions.ExportScanQuery)
	return
}

// ExportScansHelp outputs the cli help export-scans command
func (vm *VM) ExportScansHelp(cmd *cobra.Command, args []string) {
	cli := ui.NewCLI(vm.Config)

	cli.DrawVersion(ReleaseVersion, GitHash)
	cli.DrawGopher()
	cli.Println(cli.Render("exportScansUsage", nil))

	return
}

func (vm *VM) exportScansAction(cmd *cobra.Command, args []string, action actionType) {
	log := vm.Config.VM.EnableLogging()

	a := client.NewAdapter(vm.Config, vm.Metrics)

	scans, err := a.Scans(true, true)
	if err != nil {
		log.Errorf("error: couldn't scans list: %v", err)
		return
	}
	scans = vm.FilterScans(a, &scans)

	histid := vm.Config.VM.HistoryID

	var format = "nessus"
	if vm.Config.VM.OutputPDF == true {
		format = "pdf"
	} else if vm.Config.VM.OutputCSV == true {
		format = "csv"
	}

	var chapters = vm.Config.VM.Chapters

	var offset = 0
	if vm.Config.VM.Offset != "" && histid == "" {
		offset, err = strconv.Atoi(vm.Config.VM.Offset)
		if err != nil {
			log.Errorf("error: couldn't convert offset '%s': %v", vm.Config.VM.Offset, err)
			return
		}
	}

	if len(scans) == 0 {
		log.Errorf("error: history id didn't match a scan")
		return
	} else if len(scans) > 1 && histid != "" {
		log.Errorf("error: histid doesn't limit to one scan")
		return
	}

	cli := ui.NewCLI(vm.Config)
	cli.DrawVersion(ReleaseVersion, GitHash)
	cli.DrawGopher()

	for _, s := range scans {
		// If we need to get the latest histid for this scan s
		if histid == "" {
			det, err := a.ScanDetails(&s, true, true)
			if err != nil {
				continue
			}

			// This scan has no history ie. no previous scans
			if len(det.History) > offset {
				histid = det.History[offset].HistoryID
			} else {
				log.Errorf("error: scan %v has less run histories than offset '%d' requires", s.ScanID, offset)
				continue
			}

		}

		switch action {
		case actions.ExportScanStart:
			export, err := a.ScansExportStart(&s, histid, format, chapters, true, true)
			if err != nil {
				log.Errorf("error: couldn't start export-scans: %v", err)
				continue
			}

			cli.Println(cli.Render("ExportScansStart", map[string]string{"Chapters": chapters, "Format": format, "FileUUID": export.FileUUID, "ScanID": s.ScanID, "HistoryID": histid, "Offset": fmt.Sprintf("%d", offset)}))
			break
		case actions.ExportScanStatus:
			export, err := a.ScansExportStatus(&s, histid, format, chapters, true, true) //Use cache=true,true don't clobe last READY with "ERROR"
			if err != nil {
				log.Errorf("error: couldn't get status for export-scans: %v", err)
				continue
			}
			if export.Status != "READY" { // We'll try again if the status isn't ready.
				export, err = a.ScansExportStatus(&s, histid, format, chapters, false, true)
				if err != nil {
					log.Errorf("error: couldn't get status for export-scans: %v", err)
					continue
				}
			}

			cli.Println(cli.Render("ExportScansStatus", map[string]string{"Format": format, "FileUUID": export.FileUUID, "Status": export.Status, "ScanID": s.ScanID, "ScanName": s.Name, "HistoryID": histid, "Offset": fmt.Sprintf("%d", offset)}))
			break
		case actions.ExportScanGet:
			export, err := a.ScansExportGet(&s, histid, format, chapters, true, true)
			if err != nil {
				log.Errorf("error: couldn't start export-scans: %v", err)
				continue
			}
			var template = "ExportScansGet"
			if format == "pdf" {
				template = "ExportScansGetPDF"
			}
			cli.Println(cli.Render(template, map[string]string{"Format": format, "Filename": export.SourceFile.CachedFileName, "FileUUID": export.SourceFile.FileUUID, "ScanID": s.ScanID, "ScanName": s.Name, "HistoryID": histid, "Offset": fmt.Sprintf("%d", offset)}))
			break
		case actions.ExportScanQuery:
			export, err := a.ScansExportGet(&s, histid, format, chapters, true, true)
			if err != nil {
				log.Errorf("error: couldn't get export-scans: %v", err)
				continue
			}

			json, err := json.Marshal(export)
			if err != nil {
				log.Errorf("Error marshalling to JSON", err)
				return
			}

			cli.Println(string(json))
			break
		}

	}
	return
}
