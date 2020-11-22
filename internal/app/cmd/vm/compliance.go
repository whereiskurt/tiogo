package vm

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
)

// ComplianceList starts, status loops, gets and converts compliance scan results to CSV.
func (vm *VM) ComplianceList(cmd *cobra.Command, args []string) {
	a := client.NewAdapter(vm.Config, vm.Metrics)

	logger := vm.setupLog()
	scans, err := a.Scans(true, true)
	if err != nil {
		logger.Fatalf("error: couldn't scans list: %v", err)
	}
	scans = vm.FilterScans(a, &scans)

	histid := vm.Config.VM.HistoryID

	var format = "csv"
	var chapters = vm.Config.VM.Chapters

	var offset = 0
	if vm.Config.VM.Offset != "" && histid == "" {
		offset, err = strconv.Atoi(vm.Config.VM.Offset)
		if err != nil {
			logger.Fatalf("error: couldn't convert offset '%s': %v", vm.Config.VM.Offset, err)
		}
	}

	if len(scans) == 0 {
		logger.Errorf("error: history id didn't match a scan: %+v", vm.Config)
		return
	} else if len(scans) > 1 && histid != "" {
		logger.Errorf("error: histid doesn't limit to one scan")
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
				logger.Errorf("error: scan %v has less run histories than offset '%d' requires", s.ScanID, offset)
				continue
			}
		}

		var tgtFilename = fmt.Sprintf("compliance.%s.history.%s.csv", s.ScanID, histid)
		_, err := a.ScansExportStart(&s, histid, format, chapters, true, true)
		if err != nil {
			logger.Errorf("error: couldn't start export-scans: %v", err)
			panic(err)
		}

		var defaultRetryIntervals = []int{500, 1000, 2000, 2000, 2000, 5000, 5000, 5000, 10000, 20000, 20000, 30000}
		var maxattempts, sleptsec int
		export, err := a.ScansExportStatus(&s, histid, format, chapters, true, true) //Use cache=true,true don't clobe last READY with "ERROR"
		if err != nil {
			panic(err)
		}

		for maxattempts = len(defaultRetryIntervals); maxattempts >= 0; maxattempts-- {
			if export.Status == "READY" {
				break
			}

			export, err = a.ScansExportStatus(&s, histid, format, chapters, false, true)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(defaultRetryIntervals[maxattempts-1]) * time.Millisecond)
			sleptsec += defaultRetryIntervals[maxattempts-1] / 1000
		}
		if maxattempts <= 0 {
			panic(fmt.Sprintf("Status stuck on '%s' scan export '%s' after %d attempts and %dsecs.", export.Status, export.FileUUID, maxattempts, sleptsec))
		}

		filename, err := a.ScansExportLargeGet(&s, histid, format, chapters)
		if err != nil {
			panic(fmt.Sprintf("Failed to Export-Scans Large Get '%s' scan export '%s'.", export.Status, export.FileUUID))
		}

		//vm.copyToFile(filename, tgtFilename)
		err = vm.TimestampCSVRows(filename, tgtFilename)
		if err != nil {
			logger.Errorf("Deleting file: %s", tgtFilename)
			logger.Errorf("Couldn't write: %s to %s: %v", filename, tgtFilename, err)
			os.Remove(tgtFilename)
		}

		//Keep only X historicals
		vm.CleanupFiles(`.`, fmt.Sprintf(`compliance\.%s\.history\.\d+\.csv`, s.ScanID), 2)
	}

	return
}

// TimestampCSVRows Reads a CSV and adds a "_time" to the header and a Splunk friendly DTS to each row.
func (vm *VM) TimestampCSVRows(sourceName, targetName string) error {

	src, _ := os.Open(sourceName)
	defer src.Close()
	sread := csv.NewReader(src)

	// Get a header row and convert to a hash
	headers, err := sread.Read()
	if err != nil || len(headers) == 0 {
		return fmt.Errorf("Cannot read CSV: %s", sourceName)
	}

	// _time setup
	headers = append([]string{"_time"}, headers...)
	ttime := []string{fmt.Sprintf(`%s`, time.Now().UTC().Format("2006-01-02T15:04:05"))}

	tgt, err := os.Create(targetName)
	defer tgt.Close()
	if err != nil {
		return fmt.Errorf("Cannot create CSV to write out to: %s", targetName)
	}

	csvwriter := csv.NewWriter(tgt)
	defer csvwriter.Flush()

	err = csvwriter.Write(headers)
	if err != nil {
		return fmt.Errorf("Failed to write header row to CSV file: %+v", headers)
	}

	for cnt := 0; ; cnt++ {
		row, err := sread.Read()
		if err != nil || len(row) == 0 {
			break
		}
		// Add ttime to the first value
		row = append(ttime, row...)
		err = csvwriter.Write(row)
		if err != nil {
			return fmt.Errorf("Failed to write row to CSV file: %+v", row)
		}
	}
	csvwriter.Flush()

	tgt.Close()
	return nil
}
