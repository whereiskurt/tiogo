package internal_test

import (
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"os"
	"testing"
)

func TestApplicationHelp(t *testing.T) {

	var m = metrics.NewMetrics()
	c := config.NewConfig()

	t.Run("tio help", func(t *testing.T) {
		os.Args = []string{"tio", "help"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})
	t.Run("tio vm help", func(t *testing.T) {
		os.Args = []string{"tio", "vm", "help"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})

	t.Run("tio vm agents help", func(t *testing.T) {
		os.Args = []string{"tio", "vm", "agents", "help"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})
	t.Run("tio vm help agents", func(t *testing.T) {
		os.Args = []string{"tio", "vm", "help", "agents"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})

	t.Run("tio vm scanners help", func(t *testing.T) {
		os.Args = []string{"tio", "vm", "scanners", "help"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})
	t.Run("tio vm help scanners", func(t *testing.T) {
		os.Args = []string{"tio", "vm", "help", "scanners"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})

}
