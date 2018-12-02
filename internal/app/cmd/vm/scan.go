package vm

import (
	"github.com/spf13/viper"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/adapter"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
	"log"
	"strconv"
	"sync"
)

type Scan struct {
	Config *app.Config
	// Convenience functions for logging out.
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})
	// Worker *sync.WaitGroup
	Worker map[string]*sync.WaitGroup
}

func NewScan(c *app.Config) (s *Scan) {
	s = new(Scan)
	s.Config = c
	s.Errorf = s.Config.Logger.Errorf
	s.Debugf = s.Config.Logger.Debugf
	s.Warnf = s.Config.Logger.Warnf
	s.Infof = s.Config.Logger.Infof

	s.Worker = make(map[string]*sync.WaitGroup)
	s.Worker["Scan"] = new(sync.WaitGroup)
	s.Worker["ScanHistory"] = new(sync.WaitGroup)

	return
}

func (cmd *Scan) Main(cli *ui.CLI) (err error) {
	config := cmd.Config
	a := adapter.NewAdapter(config)

	if config.VM.DetailView == true {
		err = cmd.DetailView(a, cli)
	} else if config.VM.ListView == true {
		err = cmd.ListView(a, cli)
	}

	return
}
func (cmd *Scan) ListView(a *adapter.Adapter, cli *ui.CLI) (err error) {
	var scans []dao.Scan
	scans, err = a.Scans()
	if err != nil {
		return
	}

	cli.Draw.CSV.Template("Scans", scans)

	return
}

func (cmd *Scan) DetailView(a *adapter.Adapter, cli *ui.CLI) (err error) {

	cs := cmd.MakeChan(a)
	csh := cmd.MakeHistoryChan(a, cs)

	go cmd.Output(cli, csh)
	cmd.WaitToClose(csh)

	return
}

// MakeChan has a single worker iterate over matched
func (cmd *Scan) MakeChan(a *adapter.Adapter) (cs chan dao.Scan) {
	cs = make(chan dao.Scan)

	cmd.Worker["Scan"].Add(1)

	go func() {
		defer close(cs)
		defer cmd.Worker["Scan"].Done()
		//
		scans, err := a.Scans()
		if err != nil {
			return
		}
		// Push the scan into the scan channel
		for i := range scans {
			cs <- scans[i]
		}
	}()

	return
}

// MakeHistoryChan will consume a Scan, complete the [Host,Plugin,Asset] Details, pushing onto ccsh
func (cmd *Scan) MakeHistoryChan(a *adapter.Adapter, cs <-chan dao.Scan) (ccsh chan dao.ScanHistory) {
	ccsh = make(chan dao.ScanHistory)

	t := viper.Get("Depth")
	q := viper.Get("dpth")
	log.Printf("t:%+v, q:%+v", t, q)

	depth, _ := strconv.Atoi(a.Config.VM.Depth)
	workers, _ := strconv.Atoi(a.Config.WorkerCount)

	for i := 0; i < workers; i++ {
		cmd.Worker["Scan"].Add(1)
		cmd.Worker["ScanHistory"].Add(1)
		go func() {
			defer cmd.Worker["Scan"].Done()
			defer cmd.Worker["ScanHistory"].Done()
			// For every scan in the scan channel
			for scan := range cs {
				sh, scanerr := a.ScanHistory(scan, depth) // Get histories of appropriate depth (at least 1!)
				if scanerr != nil {
					cmd.Errorf("failed to parse : %+v", scanerr)
					return
				}

				for j, detail := range sh.History { // For the histories
					for k, host := range detail.Host { // For every host
						// HostHandleFunc will populate the host with details of plugins,vulns, etc.
						err := a.HostHandleFunc(detail, &host)
						if err != nil {
							cmd.Errorf("failed to handle host:%v", err)
						}

						detail.Host[k] = host
					}
					sh.History[j] = detail
				}
				ccsh <- sh
			}
		}()
	}

	return
}

func (cmd *Scan) WaitToClose(csh chan dao.ScanHistory) {
	cmd.Worker["ScanHistory"].Wait()
	close(csh)
	cmd.Worker["Scan"].Wait()
}

func (cmd *Scan) Output(cli *ui.CLI, ccsh chan dao.ScanHistory) {
	config := cli.Config

	cmd.Worker["Scan"].Add(1)
	go func(pccsh <-chan dao.ScanHistory, pcli *ui.CLI) {
		defer cmd.Worker["Scan"].Done()
		if config.OutputCSVMode == true {
			// Draw header
			pcli.Draw.CSV.Template("ScanHistoryHeader", nil)
		} else if config.OutputJSONMode == true {
			pcli.Draw.CLI.Println("[")
		}

		i := 0
		for sh := range pccsh {
			if config.OutputCSVMode == true {
				pcli.Draw.CSV.Template("ScanHistory", sh)
			} else if config.OutputJSONMode == true {
				if i > 0 {
					pcli.Draw.CLI.Println(",")
				}
				pcli.Draw.JSON.ScanHistory(sh)
			}
			i = i + 1
		}

		if config.OutputCSVMode == true {
			// Draw header
		} else if config.OutputJSONMode == true {
			pcli.Draw.CLI.Println("]")
		}

	}(ccsh, cli)
}
