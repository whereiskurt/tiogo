package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/server"
)

func Start(config *config.Config, metrics *metrics.Metrics) {
	config.Server.EnableLogging()

	config.Log.Debugf("server.Start called with -> config(%+v) and metrics->(%+v)", config, metrics)
	l := config.Log.WithFields(log.Fields{
		"cache": config.Server.CacheFolder,
		"port":  config.Server.ListenPort,
	})

	s := server.NewServer(config, metrics)

	s.EnableDefaultRouter()

	l.Info("starting server")
	_ = s.ListenAndServe()
	l.Info("server stopped.")

	l.Info("dumping metrics for server")
	config.Server.DumpMetrics()

	return
}
