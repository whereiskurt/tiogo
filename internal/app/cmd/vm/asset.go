package vm

import (
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/adapter"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
	"sync"
)

type Asset struct {
	Config *app.Config
	// Convenience functions for logging out.
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})

	Worker map[string]*sync.WaitGroup
}

func NewAsset(c *app.Config) (a *Asset) {
	a = new(Asset)
	a.Config = c
	a.Errorf = a.Config.Logger.Errorf
	a.Debugf = a.Config.Logger.Debugf
	a.Warnf = a.Config.Logger.Warnf
	a.Infof = a.Config.Logger.Infof

	a.Worker = make(map[string]*sync.WaitGroup)
	a.Worker["Asset"] = new(sync.WaitGroup)

	return
}

func (cmd *Asset) Main(cli *ui.CLI) (err error) {

	config := cmd.Config
	adapt := adapter.NewAdapter(config)

	// Use Scans to drive the Assets included in the output
	s := NewScan(config)

	if config.VM.ExportMode {
		exportUUID := config.VM.ExportUUID
		if exportUUID == "" {
			adapt.AssetExport()

		} else {
			adapt.AssetExportDownload(exportUUID)
		}

	} else if config.VM.ListView || config.VM.DetailView {
		cs := s.MakeChan(adapt)
		csh := s.MakeHistoryChan(adapt, cs)
		go cmd.AssetWorker(csh, cli)
		s.WaitToClose(csh)
		cmd.Worker["Asset"].Wait()

	}

	return
}

// AssetWorker will draw the template for every host that has asset details from every scan history.
func (cmd *Asset) AssetWorker(ccsh chan dao.ScanHistory, cli *ui.CLI) {
	config := cmd.Config

	tname := "DefatulAsset" // Default template to draw.
	if config.VM.TemplateName != "" {
		tname = config.VM.TemplateName
	}

	cmd.Worker["Asset"].Add(1)

	go func(pccsh <-chan dao.ScanHistory, pcli *ui.CLI) {
		defer cmd.Worker["Asset"].Done()
		for sh := range pccsh { // For each ScanHistory record
			for _, det := range sh.History { // For each of the History (aka historical scan)
				for _, h := range det.Host { // For each Host of the historical scan
					if !h.HasAsset() { // If we don't have asset details? (some don't... :-/ )
						continue
					}
					pcli.Draw.CSV.Template(tname, h.Asset)
				}
			}
		}
	}(ccsh, cli)
}
