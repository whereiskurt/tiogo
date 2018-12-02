package vm

import (
	"errors"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/adapter"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
	"sync"
)

type Agent struct {
	Config *app.Config
	// Convenience functions for logging out.
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})
	// Worker *sync.WaitGroup
	Worker map[string]*sync.WaitGroup
}

func NewAgent(c *app.Config) (a *Agent) {
	a = new(Agent)
	a.Config = c
	a.Errorf = a.Config.Logger.Errorf
	a.Debugf = a.Config.Logger.Debugf
	a.Warnf = a.Config.Logger.Warnf
	a.Infof = a.Config.Logger.Infof

	a.Worker = make(map[string]*sync.WaitGroup)
	a.Worker["Agent"] = new(sync.WaitGroup)

	return
}

func (cmd *Agent) Main(cli *ui.CLI) (err error) {
	config := cmd.Config
	a := adapter.NewAdapter(config)

	var ss []dao.Scanner
	ss, err = a.Scanners()
	if err != nil {
		return
	}
	scanner, err := a.AgentScanner(ss)
	if err != nil {
		return
	}
	scanner.Agents, err = a.Agents(scanner)
	if err != nil {
		return
	}

	if config.VM.AgentGroupMode == true {
		err = cmd.AssignGroup(a, cli, scanner)

		// } else if config.VM.DetailView == true {
		// 	err = cmd.Detail(a, cli)

	} else if config.VM.ListView == true {
		err = cmd.List(a, cli, scanner)

	}

	return
}

func (cmd *Agent) AssignGroup(a *adapter.Adapter, cli *ui.CLI, scanner dao.Scanner) (err error) {
	cmd.Infof("agent.AgentGroupMode")

	// Must pass have AgentGroupName
	var groupName = a.Config.VM.AgentGroupName
	if groupName == "" {
		err = errors.New("error: no agent group name specified")
		return
	}

	// Lookup the AgentGroup basedon the groupName
	group, found, lkperr := a.MatchAgentGroup(scanner, groupName)
	if lkperr != nil {
		err = lkperr
		return
	}

	// Create the Agent Group if it doesn't exist
	if !found {
		group, err = a.CreateAgentGroup(scanner, groupName)
		if err != nil {
			return
		}
	}

	// Assign each Agent to the Group
agents:
	for _, ag := range scanner.Agents {
		for k := range ag.Groups {
			if k == groupName {
				cmd.Errorf("warn: agent already a member of group: %v :%v", ag, groupName)
				continue agents
			}
		}
		err = a.AssignAgentGroup(scanner, ag, group)
		if err != nil {
			cmd.Errorf("error: couldn't assign agent group for: %v : %v :%v :%v", scanner, ag, group, err)
		}
	}

	return
}

func (cmd *Agent) List(a *adapter.Adapter, cli *ui.CLI, scanner dao.Scanner) (err error) {
	cmd.Infof("agent.ListView")

	groupname := cmd.Config.VM.AgentGroupName
	emptygroup := cmd.Config.VM.NoAgentGroupName

	p := make(map[string]interface{})

	var agents []dao.ScannerAgent

	// IF we only want matching group names...
	if groupname != "" {
	AgentLoop:
		for _, agent := range scanner.Agents {
			for _, group := range agent.Groups {
				if groupname == group.Name {
					agents = append(agents, agent)
					continue AgentLoop
				}
			}
		}
	} else if emptygroup == true {
		for _, agent := range scanner.Agents {
			if len(agent.Groups) == 0 {
				agents = append(agents, agent)
			}
		}

	} else {
		agents = scanner.Agents
	}

	groups := a.AgentGroupNames(agents)

	p["Agents"] = agents
	p["AgentGroups"] = groups
	cli.Draw.CSV.Template("AgentsDefault", p)

	return
}

// func (cmd *Agent) Detail(a *adapter.Adapter, cli *ui.CLI) (err error) {
// 	cmd.Infof("agent.DetailView")
//
// 	return
// }
