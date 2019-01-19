package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

// Version holds the config and CLI references.
type VM struct {
	Config  *config.Config
	Metrics *metrics.Metrics
}

// Version just outputs a gopher.
func (vm *VM) Help(cmd *cobra.Command, args []string) {

	fmt.Printf("tiogo version %s (%s)", ReleaseVersion, GitHash)
	cli := ui.NewCLI(vm.Config)
	if len(args) == 0 {
		cli.DrawGopher()
		fmt.Println(cli.Render("vmUsage", nil))
		return
	}

	helpType := strings.ToLower(args[0])
	switch helpType {
	case "scanners":
		fmt.Println(cli.Render("scannersUsage", nil))
	case "scans":
		fmt.Println(cli.Render("scansUsage", nil))
	case "agent-groups":
		fmt.Println(cli.Render("agentGroupsUsage", nil))
	case "agents":
		fmt.Println(cli.Render("agentsUsage", nil))
	case "users":
		fmt.Println(cli.Render("usersUsage", nil))
	case "user-groups":
		fmt.Println(cli.Render("userGroupsUsage", nil))
	case "target-groups":
		fmt.Println(cli.Render("targetGroupsUsage", nil))
	case "export-vulns":
		fmt.Println(cli.Render("exportVulnsUsage", nil))
	case "export-assets":
		fmt.Println(cli.Render("exportAssetsUsage", nil))
	default:

	}

	return
}

func (vm *VM) Scanners(cmd *cobra.Command, args []string) {
	vm.Config.Log.SetFormatter(&log.TextFormatter{})
	vm.Config.VM.EnableLogging()

	vm.Config.Log.Infof("tiogo scanners list command:")

	cli := ui.NewCLI(vm.Config)
	cli.DrawGopher()
	return
}

// NewVersion holds a configuration and command line interface reference (for log out, etc.)
func NewVM(c *config.Config, m *metrics.Metrics) (v VM) {
	v.Config = c
	v.Metrics = m
	return
}
