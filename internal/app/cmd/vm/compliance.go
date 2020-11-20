package vm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
)

// ComplianceList starts, status loops, gets and converts compliance scan results to CSV.
func (vm *VM) ComplianceList(cmd *cobra.Command, args []string) {
	a := client.NewAdapter(vm.Config, vm.Metrics)

	scans, err := a.Scans(true, true)
	if err != nil {
		log.Errorf("error: couldn't scans list: %v", err)
		return
	}
	scans = vm.FilterScans(a, &scans)

	histid := vm.Config.VM.HistoryID

	var format = "csv"
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

		_, err := a.ScansExportStart(&s, histid, format, chapters, true, true)
		if err != nil {
			log.Errorf("error: couldn't start export-scans: %v", err)
			panic(err)
		}

		export, err := a.ScansExportStatus(&s, histid, format, chapters, true, true) //Use cache=true,true don't clobe last READY with "ERROR"
		if err != nil {
			panic(err)
		}

		var defaultRetryIntervals = []int{1000}

		var maxattempts int
		for maxattempts = len(defaultRetryIntervals); maxattempts >= 0; maxattempts-- {
			if export.Status == "READY" {
				break
			}

			export, err = a.ScansExportStatus(&s, histid, format, chapters, false, true)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(defaultRetryIntervals[maxattempts-1]) * time.Millisecond)
		}
		if maxattempts <= 0 {
			panic(fmt.Sprintf("Status stuck on '%s' scan export '%s'.", export.Status, export.FileUUID))
		}

		filename, err := a.ScansExportLargeGet(&s, histid, format, chapters)
		if err != nil {
			panic(fmt.Sprintf("Failed to Export-Scans Large Get '%s' scan export '%s'.", export.Status, export.FileUUID))
		}

		var tgt = fmt.Sprintf("compliance.%s.history.%s.offset.%d.csv", s.ScanID, histid, offset)
		vm.copyToFile(filename, tgt)
	}

	return
}
