package vm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"os"
	"strings"
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
