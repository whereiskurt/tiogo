package proxy

import (
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Stop visits the specific '/shutdown' URL beginning the clean server shutdown
func Stop(config *config.Config, metrics *metrics.Metrics) {
	clientLog := config.VM.EnableLogging()

	clientLog.Infof("Sending shutdown to server ...")

	url := fmt.Sprintf("%s/shutdown", config.VM.BaseURL)

	s := tenable.NewService(config.VM.BaseURL, config.VM.SecretKey, config.VM.AccessKey, clientLog)
	t := tenable.NewTransport(&s)

	body, status, err := t.Get(url, false, false)

	if err != nil {
		clientLog.Infof("Server at '%s' was not running or cannot be reached: error: '%v': %s: %d", url, err, body, status)
		return
	}

	clientLog.Info("Successfully shutdown command...")
	return
}
