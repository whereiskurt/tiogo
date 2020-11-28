package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// AuditLogV1List will output all of the Audit Log details available
func (vm *VM) AuditLogV1List(cmd *cobra.Command, args []string) {
	logger := vm.setupLog()
	cli := ui.NewCLI(vm.Config)
	a := client.NewAdapter(vm.Config, vm.Metrics)

	dts := time.Now().Format("20060102T150405")
	maxkeep, err := strconv.Atoi(vm.Config.VM.MaxKeep)
	if err != nil {
		logger.Fatalf("error: couldn't convert maxkeep '%s': %v", vm.Config.VM.MaxKeep, err)
	}

	logger.Infof("Starting audit log v1 list ...")
	events, err := a.AuditLogV1(true, true)
	if err != nil {
		logger.Fatalf("error: failed to fetch audit logs: %v", err)
	}

	var content, saveToFilename string
	if a.Config.VM.OutputJSON {
		saveToFilename = fmt.Sprintf("auditlogv1.%s.json", dts)
		j, err := json.MarshalIndent(events, "", "\t")
		if err != nil {
			logger.Fatalf("error: couldn't marshal scan data to JSON: %v", err)
		}
		content = string(j)
	} else {
		saveToFilename = fmt.Sprintf("auditlogv1.%s.csv", dts)
		header := cli.Render("AuditLogV1HeaderCSV", map[string]interface{}{})
		body := cli.Render("AuditLogV1CSV", map[string]interface{}{"Events": events})
		content = fmt.Sprintf("%s\n%s", header, body)
	}

	err = ioutil.WriteFile(saveToFilename, []byte(content), 0644)
	if err != nil {
		logger.Fatalf("can't write to file '%s': %+v", saveToFilename, err)
	}
	logger.Infof("Wrote audit log v1 to '%s' ...", saveToFilename)

	//Keep only X historicals for auditlogs
	cleanTemplate := fmt.Sprintf(`auditlogv1.\d+T\d+\.csv`)
	vm.CleanupFiles(`.`, cleanTemplate, maxkeep)
	logger.Infof("keeping a maximum '%d' for template '%s'", maxkeep, cleanTemplate)

	cleanTemplate = fmt.Sprintf(`auditlogv1.\d+T\d+\.json`)
	vm.CleanupFiles(`.`, cleanTemplate, maxkeep)
	logger.Infof("keeping a maximum '%d' for template '%s'", maxkeep, cleanTemplate)

	return
}
