package vm

import (
	"encoding/json"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/adapter"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
	"sync"
)

type Plugin struct {
	Config *app.Config
	// Convenience functions for logging out.
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})
	Worker *sync.WaitGroup
}

func NewPlugin(c *app.Config) (p *Plugin) {
	p = new(Plugin)
	p.Config = c
	p.Errorf = c.Logger.Errorf
	p.Debugf = c.Logger.Debugf
	p.Warnf = c.Logger.Warnf
	p.Infof = c.Logger.Infof
	p.Worker = new(sync.WaitGroup)
	return
}

func (cmd *Plugin) Main(cli *ui.CLI) (err error) {
	config := cmd.Config
	a := adapter.NewAdapter(config)

	if config.VM.DetailView == true {
	} else if config.VM.ListView == true {
	}

	pp, perr := a.Plugins()
	if perr != nil {
		cmd.Errorf("failed in Plugins command: %+v", perr)
		err = perr
		return
	}

	var bb []byte
	if config.OutputJSONMode == true {
		bb, err = json.Marshal(pp)
		if err != nil {
			return
		}

		bb, err = ui.PrettyPrintJSON(bb)
		if err != nil {
			return
		}

		cli.Draw.CLI.Println(string(bb))

	}

	return
}
