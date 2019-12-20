package vm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"strings"
)

var (
	// ReleaseVersion is set by a --ldflags during a build/release
	ReleaseVersion = "v0.0.1-development"
	// GitHash is set by a --ldflags during a build/release
	GitHash = "0xhashhash"
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
	return
}

// The help command renders a template showing the help based on parameters
func (vm *VM) Help(cmd *cobra.Command, args []string) {

	cli := ui.NewCLI(vm.Config)
	versionMap := map[string]string{"ReleaseVersion": ReleaseVersion, "GitHash": GitHash}

	if len(args) == 0 {
		fmt.Println(cli.Render("vmUsage", versionMap))
		return
	}

	helpType := strings.ToLower(args[0])
	switch helpType {
	case "scanners":
		fmt.Println(cli.Render("scannersUsage", versionMap))

	case "agent-groups":
		fmt.Print(cli.Render("agentGroupsUsage", versionMap))
	case "agents":
		fmt.Print(cli.Render("agentsUsage", versionMap))
	case "export-vulns":
		fmt.Print(cli.Render("exportVulnsUsage", versionMap))

	default:
		fmt.Println(cli.Render("vmUsage", versionMap))
	}

	return
}
