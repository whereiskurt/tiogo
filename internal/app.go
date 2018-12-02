package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/app/cmd"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
	"os"
	"runtime/pprof"
	"strings"
)

type App struct {
	Config  *app.Config
	RootCmd *cobra.Command
}

func NewApp() (a *App) {
	a = new(App)

	a.RootCmd = &cobra.Command{PreRun: a.ReflectViper}
	a.Config = app.NewConfig(a.RootCmd)

	vm := cmd.NewVM(a.Config)
	vmcmd := a.BuildVMCommand(vm)

	a.Config.VM = app.NewVMConfig(a.Config, vmcmd)

	return
}

// Main executes the cobra.RootCmd.Execute() method on the root command .
// If os.Args are missing, we show help. The default root command is 'vuln'.
func (a *App) Main() {
	cli := ui.NewCLI(a.Config)

	if len(os.Args) < 2 {
		// Application help requested
		cli.Draw.Gopher()
		return
	}

	a.EnsureRootCmd()

	err := a.RootCmd.Execute()
	if err != nil {
		cli.Config.Logger.Errorf("error: %v", err)
	}

	if a.Config.PerfProfile == true {
		pprof.StopCPUProfile()
	}

	return
}

// ReflectFromViper will copy Viper values from config, envs, cli, into our app.Config struct.
// Acts almost as a 'data transfer' pattern moving from Viper -> app.Config
// The cobra.Commmand.PreRun ensures execution before command.Execute is run.
func (a *App) ReflectViper(cmd *cobra.Command, args []string) {

	// For each element of Config and Config.VM lookup the Viper.get() for the field name and
	// set the struct value. This allows us to not rely on a global Viper to exist in the app -
	// which is useful for code aspects not initialized with from the CLI
	app.ReflectFromViper(a.Config)
	app.ReflectFromViper(a.Config.VM)

	a.Config.CacheFolder = strings.TrimSuffix(a.Config.CacheFolder, "/")
	a.Config.LogFolder = strings.TrimSuffix(a.Config.LogFolder, "/")

	a.Config.Finalize()

	return
}

func (a *App) EnsureRootCmd() {
	roots := []string{"vm", "webapp", "server"}
	// Check if the first arg is a root command
	m := strings.ToLower(os.Args[1])
	if !Contains(roots, m) {
		// If no root command passed inject default
		rest := os.Args[1:]
		os.Args = []string{os.Args[0], roots[0]} // Implant the Default ahead of the rest
		os.Args = append(os.Args, rest...)
	}
}
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// BuildVMCommand creates the root "vuln" combra command and attaches sub-commands like "scan","host",etc.
func (a *App) BuildVMCommand(v *cmd.VM) (root *cobra.Command) {

	root = a.MakeCommand("vm", v.Help, a.RootCmd)

	a.AttachCommand("scan", v.Scan, root)
	a.AttachCommand("host", v.Host, root)
	a.AttachCommand("plugin", v.Plugin, root)
	a.AttachCommand("tag", v.Tag, root)
	a.AttachCommand("asset", v.Asset, root)
	a.AttachCommand("agent", v.Agent, root)
	a.AttachCommand("vuln", v.Vuln, root)

	return
}

func (a *App) MakeCommand(s string, run func(*cobra.Command, []string), parent *cobra.Command) (child *cobra.Command) {
	alias := []string{fmt.Sprintf("%ss", s)} // Add a pluralized alias
	child = &cobra.Command{Use: s, Run: run, PreRun: parent.PreRun, Aliases: alias}
	parent.AddCommand(child)
	return
}
func (a *App) AttachCommand(s string, run func(*cobra.Command, []string), parent *cobra.Command) {
	a.MakeCommand(s, run, parent)
	return
}
