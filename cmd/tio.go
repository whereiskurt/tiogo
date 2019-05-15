package main

import (
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
)

func main() {
	c := config.NewConfig()
	m := metrics.NewMetrics()

	a := internal.NewApp(c, m)

	a.InvokeCLI()

	return
}
