package vm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

type actionType string

type action struct {
	ExportScanStart  actionType
	ExportScanStatus actionType
	ExportScanGet    actionType
	ExportScanQuery  actionType
	ExportScanTag    actionType
	ExportScanUntag  actionType
}

var actions = action{
	ExportScanStart:  actionType("ExportScanStart"),
	ExportScanStatus: actionType("ExportScanStatus"),
	ExportScanGet:    actionType("ExportScanGet"),
	ExportScanQuery:  actionType("ExportScanQuery"),
	ExportScanTag:    actionType("ExportScanTag"),
	ExportScanUntag:  actionType("ExportScanUntag"),
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

// ExportScansTag will add an Asset Tag to each asset in a Scan
func (vm *VM) ExportScansTag(cmd *cobra.Command, args []string) {
	vm.exportScansAction(cmd, args, actions.ExportScanTag)
	return
}

// ExportScansUntag will remove an Asset Tag to each asset in a Scan
func (vm *VM) ExportScansUntag(cmd *cobra.Command, args []string) {
	vm.exportScansAction(cmd, args, actions.ExportScanUntag)
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

			filename, err := a.ScansExportLargeGet(&s, histid, format, chapters)
			if err != nil {
				log.Errorf("error: couldn't start export-scans: %v", err)
				continue
			}

			var tgt = fmt.Sprintf("scan.id.%v.history.%v.offset.%v.%s.%s", s.ScanID, histid, offset, chapters, format)
			vm.copyToFile(filename, tgt)

			var template = "ExportScansDownload"
			cli.Println(cli.Render(template, map[string]string{"Format": format, "CacheFilename": filename, "Filename": tgt, "ScanID": s.ScanID, "ScanName": s.Name, "HistoryID": histid, "Offset": fmt.Sprintf("%d", offset)}))

			break
		case actions.ExportScanQuery:

			export, err := a.ScansExportGet(&s, histid, format, chapters, true, true)
			var src = export.SourceFile.CachedFileName

			if err != nil {
				log.Errorf("error: couldn't get export-scans: %v", err)
				continue
			}

			var outputJSON = false
			if vm.Config.VM.OutputJSON == true {
				outputJSON = true
			}

			if outputJSON {
				j, err := json.Marshal(export)
				if err != nil {
					log.Errorf("Error marshalling to JSON: %+v", err)
					return
				}
				cli.Println(string(j))
			} else {
				var template = "ExportScansQuery"
				cli.Println(cli.Render(template, map[string]string{"Format": format, "CacheFilename": src, "FileUUID": export.SourceFile.FileUUID, "ScanID": s.ScanID, "ScanName": s.Name, "HistoryID": histid, "Offset": fmt.Sprintf("%d", offset)}))
			}

			break

		case actions.ExportScanTag:
			log.Debugf("Tagging based on export-scan ...")

			tags := vm.Config.VM.Tags
			if tags == "" || (len(strings.Split(tags, ":"))-1 != len(strings.Split(tags, ","))) {
				log.Fatalf("error: --tags cannot be empty, must pass --tags=category:value[,c2:v2,..cX:vX]")
				break
			}

			// Get the tag uuids for these scanned assets
			var taguuid []string
			for _, tt := range strings.Split(tags, ",") {
				cv := strings.Split(tt, ":")
				category, value := cv[0], cv[1]

				tv, err := a.TagValueCreate(category, value, true, true)
				if err != nil {
					log.Errorf("error: cannot get tag uuid for tagging")
					break
				}
				taguuid = append(taguuid, tv.UUID)
			}

			filename, err := a.ScansExportLargeGet(&s, histid, format, chapters)
			if err != nil {
				log.Errorf("failed to get exported scan file.")
				break
			}
			// Part 1: open the file and scan it.
			f, err := os.Open(filename)
			if err != nil {
				log.Errorf("failed to read exported.")
				break
			}
			defer f.Close()
			const pattern = `host-uuid">`
			const lenpattern = len(pattern)
			const lenuuid = 36
			const BufferSize = 64000

			buffer := make([]byte, BufferSize)
			var buf string
			var assetuuid []string
			for {
				bytesread, err := f.Read(buffer)
				if err != nil {
					if err != io.EOF {
						fmt.Println(err)
					}
					break
				}
				buf += string(buffer[:bytesread])
				for {
					loc := strings.Index(buf, pattern)
					if loc <= 0 || loc+lenpattern+lenuuid > len(buf) {
						break
					}
					offset := loc + lenpattern
					assetuuid = append(assetuuid, buf[offset:offset+lenuuid])
					buf = buf[offset:]
				}
				if len(buf) > BufferSize {
					buf = buf[BufferSize:]
				}
			}

			log.Debugf("Applying tags fors assets: %+v, Tag UUIDs:%+v", assetuuid, taguuid)
			a.TagBulkApply(assetuuid, taguuid)

			break
		case actions.ExportScanUntag:
			break
		}

	}
	return
}

func (vm *VM) copyToFile(srcName string, tgtName string) error {
	var log = vm.Config.VM.Log

	srcFile, err := os.Open(srcName)
	defer srcFile.Close()
	if err == nil {
		destFile, err := os.Create(tgtName)
		defer destFile.Close()
		if err == nil {
			_, err = io.Copy(destFile, srcFile)
			err = destFile.Sync()
			if err != nil {
				log.Warnf("Couldn't copy '%s' to local file '%s'", srcName, tgtName)
			}
		} else {
			log.Warnf("Can't copy file to '%s' : %v", srcName, tgtName)
		}
	} else {
		log.Fatalf("Outputted cache file does not exist: %v: src=%s, tgt=%s", err, srcName, tgtName)
	}

	return err
}
