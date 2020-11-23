package main

import (
	internal "github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
)

func main() {
	c := config.NewConfig()
	m := metrics.NewMetrics()

	a := internal.NewApp(c, m)

	a.InvokeCLI()

	// This will retun zero to the OS. Unless, a log.Fatalf() was called along the way
	return
}
