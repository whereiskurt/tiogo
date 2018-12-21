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
	RootCmd cobra.Command
}

func NewApp() (a App) {
	// The RootCmd is the Top Level cobra command each command is attached to (e.g. 'vm','version', 'server')
	// Once a command is made the subcommands are attached (e.g. 'scan','agent' ..)
	a.RootCmd = cobra.Command{PreRun: a.ReflectViper}

	// Make the Config to hold the values parsed by Viper
	a.Config = app.NewConfig(&a.RootCmd)

	// Build Commands and their Subcommands
	// Make the "version" command
	ver := cmd.NewVersion(a.Config)
	// The 'version' command has no subcommands and attaches to the RootCmd
	_ = MakeCommand("version", ver.Version, &a.RootCmd)

	// Make the "vm" command attached RootCmd and and subcommands attached to 'vmcmd'
	vm := cmd.NewVM(a.Config)
	vmcmd := MakeCommand("vm", vm.Help, &a.RootCmd)
	// Attach the subcommands to 'vm'
	_ = MakeCommand("scan", vm.Scan, vmcmd)
	_ = MakeCommand("host", vm.Host, vmcmd)
	_ = MakeCommand("plugin", vm.Plugin, vmcmd)
	_ = MakeCommand("tag", vm.Tag, vmcmd)
	_ = MakeCommand("asset", vm.Asset, vmcmd)
	_ = MakeCommand("agent", vm.Agent, vmcmd)
	_ = MakeCommand("vuln", vm.Vuln, vmcmd)
	// Make the config for the VM command/subcommands Parsed by Viper
	a.Config.VM = app.NewVMConfig(a.Config, vmcmd)

	return
}

// Main executes the cobra.RootCmd.Execute() method on the root command .
// If os.Args are missing, we show help. The default root command is 'vm'.
func (a *App) Main() {
	cli := ui.NewCLI(a.Config)

	if len(os.Args) < 2 {
		// Application help requested
		cli.Draw.Gopher()
		return
	}

	EnsureRootCmd()

	err := a.RootCmd.Execute()
	if err != nil {
		cli.Config.Logger.Errorf("error: %v", err)
	}

	if a.Config.PerfProfile == true {
		pprof.StopCPUProfile()
	}

	return
}

// ReflectFromViper will copy Viper values from it's config, envs, cli, into our app.Config struct.
// Acts as a 'data transfer' pattern moving from Cobra/Viper -> app.Config
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

// EnsureRootCmd will assume 'vm' (roots[0]) is the command if none is passed.
// This allows shorter invocation without specifying 'vm' everytime.
func EnsureRootCmd() {
	// TODO: Potentially make the default word (ie. 'vm') configurable?
	roots := []string{"vm", "version", "webapp", "server"}

	// Check if the first arg is a root command
	lc := strings.ToLower(os.Args[1])
	if !Contains(roots, lc) {
		// If no root command passed inject the root[0] as default
		rest := os.Args[1:]
		os.Args = []string{os.Args[0], roots[0]} // Implant the Default ahead of the rest
		os.Args = append(os.Args, rest...)
	}
}
func Contains(a []string, x string) bool {
	for i := range a {
		if x == a[i] {
			return true
		}
	}
	return false
}

func MakeCommand(s string, run func(*cobra.Command, []string), parent *cobra.Command) (child *cobra.Command) {
	alias := []string{fmt.Sprintf("%ss", s)} // Add a pluralized alias
	child = &cobra.Command{Use: s, Run: run, PreRun: parent.PreRun, Aliases: alias}
	parent.AddCommand(child)
	return
}
