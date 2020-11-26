package vm

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
)

// ComplianceGet starts, status loops, gets scan export - simplified export-scans
func (vm *VM) ComplianceGet(cmd *cobra.Command, args []string) {
	a := client.NewAdapter(vm.Config, vm.Metrics)

	logger := vm.setupLog()

	logger.Infof("Starting compliance get scan")

	scans, err := a.Scans(true, true)
	if err != nil {
		logger.Fatalf("error: couldn't scans list: %v", err)
	}
	scans = vm.FilterScans(a, &scans)

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
		logger.Errorf("error: history id didn't match a scan: %v+", vm.Config)
		return
	} else if len(scans) > 1 && histid != "" {
		logger.Errorf("error: histid doesn't limit to one scan")
		return
	}
	// TODO: Functional decompose this massive block! :-)
SCANS:
	for _, s := range scans {
		isComplianceScan := false
		det, err := a.ScanDetails(&s, true, true)

		if err != nil {
			logger.Errorf("error: can't retrieve scan details for %s at offset %d", s.ScanID, offset)
			continue SCANS
		} else if len(det.History) == 0 {
			logger.Warnf("warn: scan %v has not run yet", s.ScanID)
			continue SCANS
		}

		logger.Debugf("scan '%s' is of policy '%s' ", det.Scan.Name, det.PolicyName)

		switch det.PolicyName {
		case "Policy Compliance Auditing":
			isComplianceScan = true
		case "Mobile Device Scan":
			isComplianceScan = true
		case "Offline Config Audit":
			isComplianceScan = true
		default:
			logger.Debugf("scan '%s' does not match compliance policy scan type - needs compliance plugins", det.Scan.Name)

		}

		// If there are CompliancePlugins that ran, this is a complieance scan.
		if len(det.CompliancePlugin) > 0 {
			isComplianceScan = true
			logger.Debugf("scan '%s' has compliance plugins, ", det.Scan.Name)
		} else {
			logger.Debugf("scan '%s' is of policy '%s does not have any compliance plugins' ", det.Scan.Name, det.PolicyName)
		}

		if !isComplianceScan {
			logger.Infof("scan '%s' is of policy '%s is not compliance scan, skipping' ", det.Scan.Name, det.PolicyName)
			continue SCANS
		}

		// This scan has no history ie. no previous scans
		if len(det.History) <= offset {
			logger.Errorf("error: scan '%s' (%v) has less run histories than offset '%d' requires", det.Scan.Name, s.ScanID, offset)
			continue SCANS
		}

	DEPTHS:
		for depth := 0; depth < maxdepth; depth++ {

			if len(det.History)-1 < offset+depth {
				logger.Infof("scan %s has only %d histories (wanted %d)", s.ScanID, len(det.History), offset+depth)
				break DEPTHS
			}
			histid = det.History[offset+depth].HistoryID

			var tgtFilename = fmt.Sprintf("compliance.%s.history.%s.csv", s.ScanID, histid)

			if _, err := os.Stat(tgtFilename); err == nil {
				logger.Infof("skipping scan:%s history:%s: file already downloaded: %s", s.ScanID, histid, tgtFilename)
				continue DEPTHS
			}

			_, err = a.ScansExportStart(&s, histid, format, chapters, true, true)
			if err != nil {
				logger.Errorf("error: cannot start export scans from Tenable.io, skipping scan: %v+", err)
				continue SCANS
			}

			var sleepStatusCheckIntervals = []int{500, 1000, 2000, 2500, 3000, 3500, 5000, 5000, 5000, 10000, 20000, 20000, 30000, 30000, 30000, 30000}
			var maxattempts, sleptsec int
			export, err := a.ScansExportStatus(&s, histid, format, chapters, true, true) //Use cache=true,true don't clobe last READY with "ERROR"
			if err != nil {
				logger.Errorf("error: cannot get export status from Tenable.io, skipping scan: %v+", err)
				continue SCANS
			}
		STATUSCHECK:
			for maxattempts = len(sleepStatusCheckIntervals); maxattempts >= 0; maxattempts-- {
				if export.Status == "READY" {
					break STATUSCHECK
				}
				export, err = a.ScansExportStatus(&s, histid, format, chapters, false, true)
				if err != nil {
					logger.Errorf("error: cannot get export status from Tenable.io, skipping scan: %v+", err)
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

			logger.Infof("processing download with _time for '%s' (%s) uuid: %s", s.Name, s.ScanID, export.FileUUID)

			err = vm.ProcessCSVRow(filename, tgtFilename, func(header map[string]int, row []string) (shouldKeep bool) {

				compPlugin := det.CompliancePlugin
				//NOTE: MDM Mobile scans don't actually have compliance plugins, so the LEN is zero...
				if len(compPlugin) == 0 {
					shouldKeep = true
				} else if _, ok := compPlugin[row[header["Name"]]]; ok {
					shouldKeep = true
				}

				// This block reduces the CIS benchmark output to the pass/fail line, plus the remote value and policy value.
				// This filter reduces LARGE amounts of repeated texts
				if row[header["Plugin Family"]] == "Policy Compliance" {
					var testNameStatus = regexp.MustCompile(`(?msi)^(.+?\[(?:PASSED|FAILED|ERROR|WARNING)\])`)
					var remoteValue = regexp.MustCompile(`(?msi)^\s*(?:Remote value:\s*)(.+?)\s*(?:(?:^Policy value\s*:)|(?:^Solution\s*:))`)
					var policyValue = regexp.MustCompile(`(?msi)^\s*(?:Policy value:\s*)(.+?)\s*(?:(?:^Solution\s*:)|(?:^Reference\(s\)\s*:)|$)`)

					name := testNameStatus.FindStringSubmatch(row[header["Description"]])
					rv := remoteValue.FindStringSubmatch(row[header["Description"]])
					pv := policyValue.FindStringSubmatch(row[header["Description"]])

					if len(name) > 1 && len(rv) > 1 {
						row[header["Description"]] = name[1] + "\n"
						row[header["Description"]] = strings.ReplaceAll(row[header["Description"]], `"`, "")
						if len(rv) > 1 {
							row[header["Description"]] += "Remote value:" + rv[1] + "\n"
						}
						if len(pv) > 1 {
							row[header["Description"]] += "Policy value:" + pv[1] + "\n"
						}
					}
				}
				return shouldKeep
			})

			if err != nil {
				os.Remove(tgtFilename)
				logger.Fatalf("error: couldn't apply timestatmps write: %s to %s: %v, (out of disk space?)", filename, tgtFilename, err)
			}

		}

		//Keep only X historicals for this ScanId
		cleanTemplate := fmt.Sprintf(`compliance\.%s\.history\.\d+\.csv`, s.ScanID)
		logger.Infof("keeping a maximum '%d' for template '%s'", maxkeep, cleanTemplate)
		vm.CleanupFiles(`.`, cleanTemplate, maxkeep)
	}

	return
}
