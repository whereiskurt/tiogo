package vm

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

var (
	// ReleaseVersion is set by a --ldflags during a build/release
	ReleaseVersion = "v0.3.2020-development"
	// GitHash is set by a --ldflags during a build/release
	GitHash = "0x0123abcd"
)

// VM holds the config and CLI references.
type VM struct {
	Config  *config.Config
	Metrics *metrics.Metrics
}

// NewVM holds a configuration and command line interface reference (for log out, etc.)
func NewVM(c *config.Config, m *metrics.Metrics) (v VM) {
	v.Config = c
	v.Metrics = m
	v.Config.VM.ReleaseVersion = ReleaseVersion
	v.Config.VM.GitHash = GitHash
	return
}

// Help command renders a template showing the help based on parameters
func (vm *VM) Help(cmd *cobra.Command, args []string) {

	cli := ui.NewCLI(vm.Config)

	versionMap := map[string]string{"ReleaseVersion": vm.Config.VM.ReleaseVersion, "GitHash": vm.Config.VM.GitHash}

	// Always output the Gopher and version number
	fmt.Fprintf(os.Stderr, cli.Render("CommandHeader", versionMap))

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, cli.Render("vmUsage", versionMap))
		return
	}

	helpType := strings.ToLower(args[0])
	switch helpType {
	case "scanners", "scanner":
		fmt.Fprintf(os.Stderr, cli.Render("scannersUsage", versionMap))
	case "agent-groups", "agent-group":
		fmt.Fprintf(os.Stderr, cli.Render("agentGroupsUsage", versionMap))
	case "agents", "agent":
		fmt.Fprintf(os.Stderr, cli.Render("agentsUsage", versionMap))
	case "scans", "scan":
		fmt.Fprintf(os.Stderr, cli.Render("scansUsage", versionMap))
	case "export-vulns", "export-vuln":
		fmt.Fprintf(os.Stderr, cli.Render("exportVulnsUsage", versionMap))
	case "export-assets", "export-asset":
		fmt.Fprintf(os.Stderr, cli.Render("exportAssetsUsage", versionMap))
	case "export-scans", "export-scan":
		fmt.Fprintf(os.Stderr, cli.Render("exportScansUsage", versionMap))
	case "cache":
		fmt.Fprintf(os.Stderr, cli.Render("cacheUsage", versionMap))
	default:
		fmt.Fprintf(os.Stderr, cli.Render("vmUsage", versionMap))
	}

	return
}

// CleanupFiles is used to keep a maximum amount of matching files
func (vm *VM) CleanupFiles(dirpath string, regex string, keep int) {
	// 1. Compile regular expression to match each filename against
	r, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 2. Read the current working directory file list
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 3. For every filename in the dir that matches the regular expression, store in matches
	var matches []os.FileInfo
	for _, file := range files {
		if r.MatchString(file.Name()) {
			matches = append(matches, file)
		}
	}

	//If there are more matches than files copies we want to keep
	if len(matches) >= keep {
		// Sort newest[0] to oldest[len(matches)-1]
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].ModTime().After(matches[j].ModTime())
		})
		// Delete files name at index maxoldest and beyond
		for i := keep; i < len(matches); i++ {
			os.Remove(matches[i].Name())
		}
	}
}

// FilterKeepAll returns true so they caller keeps all rows.
func (vm *VM) FilterKeepAll(header, row []string) (shouldKeep bool) {
	return true
}

// ProcessCSVRow Reads a CSV and adds a "_time" to the header and a Splunk friendly DTS to each row.
func (vm *VM) ProcessCSVRow(sourceName, targetName string, keepFilter func(header, row []string) (shouldKeep bool)) (err error) {

	src, _ := os.Open(sourceName)
	defer src.Close()
	sread := csv.NewReader(src)

	// Get a header row and convert to a hash
	headers, err := sread.Read()
	if err != nil || len(headers) == 0 {
		return fmt.Errorf("Cannot read CSV: %s", sourceName)
	}

	//TODO: Build row HASH

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
ROW:
	for cnt := 0; ; cnt++ {
		row, err := sread.Read()
		if err != nil || len(row) == 0 {
			break
		}
		// Add ttime to the first value
		row = append(ttime, row...)

		//ROW FILTER HOOK
		if !keepFilter(headers, row) {
			continue ROW
		}
		err = csvwriter.Write(row)
		if err != nil {
			return fmt.Errorf("Failed to write row to CSV file: %+v", row)
		}
	}
	csvwriter.Flush()

	tgt.Close()
	return nil
}
