package server

import (
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Stop visits the specific '/shutdown' URL beginning the clean server shutdown
func Stop(config *config.Config, metrics *metrics.Metrics) {
	a := client.NewAdapter(config, metrics)
	a.Config.VM.EnableLogging()
	config.Log.Debugf("Sending shutdown command to: %s/shutdown", config.VM.BaseURL)

	url := fmt.Sprintf("%s/shutdown", config.VM.BaseURL)

	s := tenable.NewService(config.VM.BaseURL, config.VM.SecretKey, config.VM.AccessKey)
	t := tenable.NewTransport(&s)

	body, status, err := t.Get(url)
	if err != nil {
		config.Log.Infof("Server at '%s' was not running or cannot be reached: error: '%v'", url, err)
		return
	}

	fmt.Println(fmt.Sprintf("Success [%d]!\n%s", status, body))
	return
}