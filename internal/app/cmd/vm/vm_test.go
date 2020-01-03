package vm_test

import (
	"testing"

	"github.com/whereiskurt/tiogo/internal/app/cmd/vm"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
)

func TestApplicationHelp(t *testing.T) {
	m := metrics.NewMetrics()
	c := config.NewConfig()

	//
	t.Run("tio vm help", func(t *testing.T) {
		vm := vm.NewVM(c, m)
		vm.Help(nil, []string{})
		vm.Help(nil, []string{"agents"})
		vm.Help(nil, []string{"agent-groups"})
		vm.Help(nil, []string{"scans"})
		vm.Help(nil, []string{"scanners"})
		vm.Help(nil, []string{"export-scans"})
		vm.Help(nil, []string{"export-vulns"})
		vm.Help(nil, []string{"export-assets"})
		vm.Help(nil, []string{"cache"})
		vm.Help(nil, []string{"unexpectedvalue"})
	})

}
