package vm

import (
	"fmt"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/adapter"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
	"sync"
)

type Vuln struct {
	Config *app.Config
	// Convenience functions for logging out.
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})

	Worker map[string]*sync.WaitGroup
}

func NewVuln(c *app.Config) (a *Vuln) {
	a = new(Vuln)
	a.Config = c
	a.Errorf = a.Config.Logger.Errorf
	a.Debugf = a.Config.Logger.Debugf
	a.Warnf = a.Config.Logger.Warnf
	a.Infof = a.Config.Logger.Infof

	a.Worker = make(map[string]*sync.WaitGroup)
	a.Worker["VM"] = new(sync.WaitGroup)

	return
}

func (cmd *Vuln) Main(cli *ui.CLI) (err error) {
	config := cmd.Config
	adapt := adapter.NewAdapter(config)

	if config.VM.ExportMode {
		var assets []dao.VulnExportChunk
		exportUUID := config.VM.ExportUUID
		if exportUUID == "" {
			assets, err = adapt.VulnExport()
		} else {
			assets, err = adapt.VulnExportDownload(exportUUID)
		}

		cli.Println(fmt.Sprintf("Assets loaded from export: %d", len(assets)))

		p := make(map[string]interface{})
		p["Assets"] = assets

		cli.Draw.CSV.Template("VulnExportDefault", p)

	}

	return
}
