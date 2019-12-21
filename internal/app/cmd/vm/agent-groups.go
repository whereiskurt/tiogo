package vm

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

// AgentGroupsList is invoked by Cobra with commandline args passed.
func (vm *VM) AgentGroupsList(cmd *cobra.Command, args []string) {
	log := vm.Config.Server.EnableLogging()
	log.Debugf("AgentsGroupsList started")

	a := client.NewAdapter(vm.Config, vm.Metrics)
	cli := ui.NewCLI(vm.Config)

	log.Debugf("AgentsGroupsList client won't use cache and will fetch fresh results from server")
	groups, err := a.AgentGroups(false, true)
	if err != nil {
		log.Errorf("error: couldn't agents list: %v", err)
		return
	}

	regex := vm.Config.VM.Regex
	if regex != "" {
		groups = a.Filter.AgentGroupsByRegex(groups, regex)
	}

	// Outputs
	if a.Config.VM.OutputJSON {
	}
	if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		fmt.Println(cli.Render("AgentGroupsListCSV", map[string]interface{}{"AgentGroups": groups}))
	}

	return
}
