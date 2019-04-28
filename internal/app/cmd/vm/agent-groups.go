package vm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

func (vm *VM) AgentGroupsList(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()
	log.Debugf("AgentsGroupsList started")

	a := client.NewAdapter(vm.Config, vm.Metrics)
	cli := ui.NewCLI(vm.Config)

	groups, err := a.AgentGroups()
	if err != nil {
		log.Errorf("error: couldn't agents list: %v", err)
		return
	}

	// Outputs
	if a.Config.VM.OutputJSON {
	}
	if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		fmt.Println(cli.Render("AgentGroupsListCSV", map[string]interface{}{"AgentGroups": groups}))
	}

	return
}
