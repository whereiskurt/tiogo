package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/app/cmd/vm"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
)

type VM struct {
	Config *app.Config
	CLI    *ui.CLI
}

func NewVM(c *app.Config) (v *VM) {
	v = new(VM)
	v.Config = c
	return
}

func (v *VM) Scan(cmd *cobra.Command, args []string) {
	scan := vm.NewScan(v.Config)
	cli := ui.NewCLI(v.Config)

	cli.Draw.Banner()
	scan.Infof("Executing 'internal/app/cmd/vuln/scan.go' ...")

	err := scan.Main(cli)
	if err != nil {
		scan.Errorf("Failed: %s", err)
		return
	}
	scan.Infof("Success! :-)")
}
func (v *VM) Help(cmd *cobra.Command, args []string) {
	cli := v.CLI
	cli.Draw.Banner()
	cli.Config.Logger.Infof("Executing 'internal/app/cmd/vuln/help.go' ...")
}
func (v *VM) Asset(cmd *cobra.Command, args []string) {
	asset := vm.NewAsset(v.Config)
	cli := ui.NewCLI(v.Config)

	cli.Draw.Banner()
	asset.Infof("Executing 'internal/app/cmd/vuln/asset.go' ...")

	err := asset.Main(cli)
	if err != nil {
		asset.Errorf("Failed: %p", err)
		return
	}
	asset.Infof("Success! :-)")
}
func (v *VM) Agent(cmd *cobra.Command, args []string) {
	agent := vm.NewAgent(v.Config)
	cli := ui.NewCLI(v.Config)

	cli.Draw.Banner()
	agent.Infof("Executing 'internal/app/cmd/vuln/agent.go' ...")

	err := agent.Main(cli)
	if err != nil {
		agent.Errorf("Failed: %v", err)
		return
	}
	agent.Infof("Success! :-)")
}
func (v *VM) Plugin(cmd *cobra.Command, args []string) {
	plugin := vm.NewPlugin(v.Config)
	cli := ui.NewCLI(v.Config)

	cli.Draw.Banner()
	plugin.Infof("Executing 'internal/app/cmd/vuln/plugin.go' ...")

	err := plugin.Main(cli)
	if err != nil {
		plugin.Errorf("Failed: %p", err)
		return
	}
	plugin.Infof("Success! :-)")
}
func (v *VM) Vuln(cmd *cobra.Command, args []string) {
	vuln := vm.NewVuln(v.Config)
	cli := ui.NewCLI(v.Config)

	cli.Draw.Banner()
	vuln.Infof("Executing 'internal/app/cmd/vuln/vuln.go' ...")

	err := vuln.Main(cli)
	if err != nil {
		vuln.Errorf("Failed: %p", err)
		return
	}
	vuln.Infof("Success! :-)")
}
func (v *VM) Tag(cmd *cobra.Command, args []string)  {}
func (v *VM) Host(cmd *cobra.Command, args []string) {}
