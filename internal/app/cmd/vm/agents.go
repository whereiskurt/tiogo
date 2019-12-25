package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

func (vm *VM) AgentsList(cmd *cobra.Command, args []string) {
	log := vm.Config.VM.EnableLogging()
	log.Debugf("AgentsList started")

	a := client.NewAdapter(vm.Config, vm.Metrics)
	cli := ui.NewCLI(vm.Config)

	regex := vm.Config.VM.Regex
	name := vm.Config.VM.Name
	groupName := vm.Config.VM.GroupName

	agents, agentGroups, err := vm.list(cli, a)
	if err != nil {
		cli.Fatalf("%s", err)
		return
	}

	// Reduce our agents to only ones matching a REGEX or name.
	filter := a.Filter
	if regex != "" {
		agents = filter.AgentsByRegex(agents, regex)
	} else if name != "" {
		agents = filter.AgentsByName(agents, name)
	}

	// Reduce agents to just ones matching group membership
	if groupName != "" {
		agents = filter.KeepOnlyGroupMembers(agents, groupName)
	}

	// TODO: Add a switch to override this. :-)
	// Rewrite the agentGroups for groups just found in agents[]
	if regex != "" {
		// Collect unique AgentGroups in map, indexed by name
		var ag = make(map[string]client.AgentGroup)
		for _, a := range agents {
			for _, g := range a.Groups {
				ag[g.Name] = g
			}
		}
		// Convert map to a list agentGroup list
		agentGroups = make([]client.AgentGroup, 0)
		for k := range ag {
			agentGroups = append(agentGroups, ag[k])
		}
	}

	// Outputs
	if a.Config.VM.OutputJSON {
		j, _ := json.Marshal(agents)
		fmt.Println(fmt.Sprintf("%s", j))
	}

	if a.Config.VM.OutputCSV || !a.Config.VM.OutputJSON {
		fmt.Println(cli.Render("AgentsListCSV", map[string]interface{}{"Agents": agents, "AgentGroups": agentGroups}))
	}

	return
}
func (vm *VM) AgentsGroup(cmd *cobra.Command, args []string) {
	vm.action(filterRemoveWithoutGroup, group)
}
func (vm *VM) AgentsUngroup(cmd *cobra.Command, args []string) {
	vm.action(filterKeepWithGroup, ungroup)
}

func (vm *VM) action(filterFunc func(*client.Adapter, ui.CLI, []client.ScannerAgent, string) []client.ScannerAgent, groupFunc func(*client.Adapter, ui.CLI, client.ScannerAgent, *client.AgentGroup)) {
	a := client.NewAdapter(vm.Config, vm.Metrics)
	cli := ui.NewCLI(vm.Config)

	groupName := vm.Config.VM.GroupName
	if groupName == "" {
		err := errors.New(fmt.Sprint("error: must provide group name to group agents: missing --group"))
		cli.Fatalf("%s", err)
	}

	regex := vm.Config.VM.Regex
	name := vm.Config.VM.Name
	if regex == "" && name == "" {
		err := errors.New(fmt.Sprint("error: must set either --name or --regex not both"))
		cli.Fatalf("%s", err)
	}

	// 2) Get Agents and Groups:
	agents, agentGroups, err := vm.list(cli, a)
	if err != nil {
		cli.Fatalf("error: %s", err)
	}

	group := lookupGroup(cli, agentGroups, groupName)

	agents = filterFunc(a, cli, agents, groupName)
	if regex != "" {
		agents = a.Filter.AgentsByRegex(agents, regex)
	} else {
		agents = a.Filter.AgentsByName(agents, name)
	}

	for _, agent := range agents {
		groupFunc(a, cli, agent, group)
	}

	// Update the cache :-)
	agents, err = a.Agents(false, true)

	return
}

func (vm *VM) list(cli ui.CLI, a *client.Adapter) ([]client.ScannerAgent, []client.AgentGroup, error) {
	regex := vm.Config.VM.Regex
	name := vm.Config.VM.Name
	if name != "" && regex != "" {
		err := errors.New(fmt.Sprint("error: cannot have both name parameters --name and --regex"))
		cli.Fatalf("%s", err)
	}

	agents, err := a.Agents(true, true)
	if err != nil {
		err := errors.New(fmt.Sprintf("error: couldn't agents list: %v", err))
		return nil, nil, err
	}

	agentGroups, err := a.AgentGroups()
	if err != nil {
		err := errors.New(fmt.Sprintf("error: couldn't agent groups list: %v", err))
		return nil, nil, err
	}

	// Invoke with --trace outputs these lines.
	log.Debugf("Total agents:%d, Total Agent Groups: %d", len(agents), len(agentGroups))

	return agents, agentGroups, nil
}

func lookupGroup(cli ui.CLI, agentGroups []client.AgentGroup, groupName string) *client.AgentGroup {
	// 3) Check the Group Name passed is an actual agent group
	var group *client.AgentGroup = nil
	for g := range agentGroups {
		if agentGroups[g].Name == groupName {
			group = &agentGroups[g]
			return group
		}
	}
	cli.Fatalf(`error: no group name matching: "%s"`, groupName)
	return nil // Fatal never gets here
}

func group(a *client.Adapter, cli ui.CLI, agent client.ScannerAgent, group *client.AgentGroup) {
	cli.Println(fmt.Sprintf("Adding '%s'(ID:%s) to group '%s'(ID: %s) ...", agent.Name, agent.ID, group.Name, group.ID))
	err := a.AgentAssignGroup(agent.ID, group.ID, agent.Scanner.ID)
	if err != nil {
		err := errors.New(fmt.Sprintf("  error: failed to add agent to group: %s", err))
		cli.Errorf("%s", err)
	}
}
func ungroup(a *client.Adapter, cli ui.CLI, agent client.ScannerAgent, group *client.AgentGroup) {
	cli.Println(fmt.Sprintf("Removing '%s'(ID:%s) from group '%s'(ID: %s) ...", agent.Name, agent.ID, group.Name, group.ID))
	err := a.AgentUnassignGroup(agent.ID, group.ID, agent.Scanner.ID)
	if err != nil {
		err := errors.New(fmt.Sprintf("  error: failed to remove agent to group: %s", err))
		cli.Errorf("%s", err)
	}
}
func filterKeepWithGroup(a *client.Adapter, cli ui.CLI, agents []client.ScannerAgent, groupName string) []client.ScannerAgent {
	// 3) Filter Agents
	// Filter agent that are to non-group members - don't reassign if assigned
	agents = a.Filter.KeepOnlyGroupMembers(agents, groupName)

	if len(agents) == 0 {
		err := errors.New(fmt.Sprint("error: no agents candidate to be removed from group."))
		cli.Fatalf("%s", err)
	}
	return agents
}
func filterRemoveWithoutGroup(a *client.Adapter, cli ui.CLI, agents []client.ScannerAgent, groupName string) []client.ScannerAgent {
	// 3) Filter Agents
	// Filter agent that are to non-group members - don't reassign if assigned
	agents = a.Filter.SkipGroupMembers(agents, groupName)

	if len(agents) == 0 {
		err := errors.New(fmt.Sprint("error: no agents candidate to be put in group."))
		cli.Fatalf("%s", err)
	}
	return agents
}
