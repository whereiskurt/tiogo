package internal

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whereiskurt/tiogo/internal/app/cmd"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	// ApplicationName is referenced for the usage help.
	ApplicationName = "tiogo"
	// CommandList entry[0] becomes default when a 'command' is omitted
	CommandList = []string{"vm", "server"}
)

// App is created from package main. App handles the configuration and cobra/viper.
type App struct {
	Config       *config.Config
	Metrics      *metrics.Metrics
	RootCmd      *cobra.Command
	DefaultUsage string
}

// NewApp constructs the command line and configuration
func NewApp(config *config.Config, mmetrics *metrics.Metrics) (a App) {
	a.Config = config
	a.Metrics = mmetrics
	a.RootCmd = new(cobra.Command)
	a.DefaultUsage = a.usageTemplate("Usage", nil)

	// Ensure before any command is run we Unmarshal and Validate the Config values.
	// NOTE: we need to set the PreRun BEFORE making other commands below.
	a.RootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		a.Config.UnmarshalViper()  // copy values from cobra
		cmd.ParseFlags(os.Args)    // parse commandline for parameters
		a.Config.ValidateOrFatal() // and validate.
	}

	makeString("VerboseLevel", &a.Config.VerboseLevel, []string{"level"}, a.RootCmd)
	makeBool("VerboseLevel1", &a.Config.VerboseLevel1, []string{"s", "silent"}, a.RootCmd)
	makeBool("VerboseLevel2", &a.Config.VerboseLevel2, []string{"q", "quiet"}, a.RootCmd)
	makeBool("VerboseLevel3", &a.Config.VerboseLevel3, []string{"v", "info"}, a.RootCmd)
	makeBool("VerboseLevel4", &a.Config.VerboseLevel4, []string{"debug"}, a.RootCmd)
	makeBool("VerboseLevel5", &a.Config.VerboseLevel5, []string{"trace"}, a.RootCmd)

	ver := cmd.NewVM(a.Config)
	vmCmd := makeCommand("vm", ver.Help, a.RootCmd)
	makeString("ID", &a.Config.VM.ID, []string{"i", "id"}, vmCmd)
	makeString("Name", &a.Config.VM.Name, []string{"n", "name"}, vmCmd)
	makeString("Regex", &a.Config.VM.Regex, []string{"regex"}, vmCmd)
	makeString("JQex", &a.Config.VM.JQex, []string{"jqex"}, vmCmd)
	makeBool("CSV", &a.Config.VM.OutputCSV, []string{"csv"}, vmCmd)
	makeBool("JSON", &a.Config.VM.OutputJSON, []string{"json"}, vmCmd)

	_ = makeCommand("help", ver.Help, vmCmd)

	sListCmd := makeCommand("scanners", ver.Scanners, vmCmd)
	_ = makeCommand("list", ver.Scanners, sListCmd)

	a.RootCmd.SetUsageTemplate(a.DefaultUsage)
	a.RootCmd.SetHelpTemplate(a.DefaultUsage)

	return
}

// InvokeCLI passes control over to the root cobra command.
func (a *App) InvokeCLI() {

	setDefaultRootCmd()

	// Call Cobra Execute which will PreRun and select the Command to execute.
	_ = a.RootCmd.Execute()

	return
}

// usageTemplate renders the usage/help/man pages for a cmd
func (a *App) usageTemplate(name string, data interface{}) (usage string) {
	var raw bytes.Buffer

	var templateFiles []string
	templateFiles = append(templateFiles, CommandList...)
	templateFiles = append(templateFiles, ApplicationName)

	t := template.New("")
	for _, f := range templateFiles {
		name := fmt.Sprintf("/template/cmd/%s.tmpl", f)
		file, err := config.TemplateFolder.Open(name)
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		t, err = t.Funcs(
			template.FuncMap{
				"Gopher": ui.Gopher,
			},
		).Parse(string(content))
	}

	var err error

	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		log.Printf("error execute template for usage: %v", err)
		return
	}

	usage = raw.String()
	return
}

func setDefaultRootCmd() {
	if len(os.Args) < 2 {
		return
	}
	// Check if the first arg is a root command
	arg := strings.ToLower(os.Args[1])

	// If the first argument isn't one we were expecting, shove CommandList[0] in.
	if !contains(CommandList, arg) {
		// If no root command passed inject the root[0] as default
		rest := os.Args[1:]
		os.Args = []string{os.Args[0], CommandList[0]} // Implant the Default ahead of the rest
		os.Args = append(os.Args, rest...)
	}

	return
}
func contains(a []string, x string) bool {
	for i := range a {
		if x == a[i] {
			return true
		}
	}
	return false
}
func makeCommand(s string, run func(*cobra.Command, []string), parent *cobra.Command) (child *cobra.Command) {
	alias := []string{fmt.Sprintf("%ss", s)} // Add a pluralized alias
	child = &cobra.Command{Use: s, Run: run, PreRun: parent.PreRun, Aliases: alias}
	parent.AddCommand(child)
	return
}
func makeBool(name string, ref *bool, aliases []string, cob *cobra.Command) {
	cob.PersistentFlags().BoolVar(ref, name, *ref, "")
	_ = viper.BindPFlag(name, cob.PersistentFlags().Lookup(name))
	if len(aliases) > 0 && len(aliases[0]) == 1 {
		cob.PersistentFlags().Lookup(name).Shorthand = aliases[0]
	}
	for _, alias := range aliases {
		cob.PersistentFlags().BoolVar(ref, alias, *ref, "")
		cob.PersistentFlags().Lookup(alias).Hidden = true
		viper.RegisterAlias(alias, name)
	}

	return
}
func makeString(name string, ref *string, aliases []string, cob *cobra.Command) {
	cob.PersistentFlags().StringVar(ref, name, *ref, "")
	_ = viper.BindPFlag(name, cob.PersistentFlags().Lookup(name))
	if len(aliases) > 0 && len(aliases[0]) == 1 {
		cob.PersistentFlags().Lookup(name).Shorthand = aliases[0]
	}
	for _, alias := range aliases {
		cob.PersistentFlags().StringVar(ref, alias, *ref, "")
		cob.PersistentFlags().Lookup(alias).Hidden = true
		viper.RegisterAlias(alias, name)
	}

	return
}
