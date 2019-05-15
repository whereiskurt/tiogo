package proxy

import (
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/proxy"
)

func Start(config *config.Config, metrics *metrics.Metrics) {

	serverLog := config.Server.EnableLogging()
	clientLog := config.VM.EnableLogging()

	clientLog.Infof("Starting a server ...")

	s := proxy.NewServer(config, metrics, serverLog)

	s.EnableDefaultRouter()

	s.ListenAndServe()

	config.Server.DumpMetrics()

	return
}
