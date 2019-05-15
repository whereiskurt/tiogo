package internal

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdproxy "github.com/whereiskurt/tiogo/internal/app/cmd/proxy"
	"github.com/whereiskurt/tiogo/internal/app/cmd/vm"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	pkgproxy "github.com/whereiskurt/tiogo/pkg/proxy"
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
	CommandList = []string{"vm", "proxy"}
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
	a.DefaultUsage = a.Usage()

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

	vmServer := cmdproxy.NewServer(a.Config, a.Metrics)
	vmServerCmd := makeCommand("proxy", vmServer.ServerHelp, a.RootCmd)
	_ = makeCommand("help", vmServer.ServerHelp, vmServerCmd)
	_ = makeCommand("start", vmServer.Start, vmServerCmd)
	_ = makeCommand("stop", vmServer.Stop, vmServerCmd)

	vmApp := vm.NewVM(a.Config, a.Metrics)
	vmCmd := makeCommand("vm", vmApp.Help, a.RootCmd)
	makeString("ID", &a.Config.VM.ID, []string{"i", "id"}, vmCmd)
	makeString("UUID", &a.Config.VM.UUID, []string{"uuid"}, vmCmd)
	makeString("Name", &a.Config.VM.Name, []string{"n", "name"}, vmCmd)
	makeString("Regex", &a.Config.VM.Regex, []string{"regex"}, vmCmd)
	makeString("JQex", &a.Config.VM.JQex, []string{"jqex"}, vmCmd)
	makeBool("CSV", &a.Config.VM.OutputCSV, []string{"csv"}, vmCmd)
	makeBool("JSON", &a.Config.VM.OutputJSON, []string{"json"}, vmCmd)

	makeBool("Critical", &a.Config.VM.Critical, []string{"critical", "crit"}, vmCmd)
	makeBool("High", &a.Config.VM.High, []string{"high"}, vmCmd)
	makeBool("Medium", &a.Config.VM.Medium, []string{"medium", "med"}, vmCmd)
	makeBool("Info", &a.Config.VM.Info, []string{"info", "low"}, vmCmd)

	_ = makeCommand("help", vmApp.Help, vmCmd)

	exportVulnsCmd := makeCommand("export-vulns", vmApp.ExportVulnsHelp, vmCmd)
	makeString("ExportLimit", &a.Config.VM.ExportLimit, []string{"limit", "size", "export-limit"}, exportVulnsCmd)
	_ = makeCommand("start", vmApp.ExportVulnsStart, exportVulnsCmd)
	_ = makeCommand("status", vmApp.ExportVulnsStatus, exportVulnsCmd)
	_ = makeCommand("get", vmApp.ExportVulnsGet, exportVulnsCmd)
	_ = makeCommand("query", vmApp.ExportVulnsQuery, exportVulnsCmd)
	makeString("Chunk", &a.Config.VM.Chunk, []string{"chunk", "chunks"}, exportVulnsCmd)
	makeString("BeforeDate", &a.Config.VM.BeforeDate, []string{"before"}, exportVulnsCmd)
	makeString("AfterDate", &a.Config.VM.AfterDate, []string{"after"}, exportVulnsCmd)
	makeString("Days", &a.Config.VM.Days, []string{"days"}, exportVulnsCmd)

	exportAssetsCmd := makeCommand("export-assets", vmApp.ExportAssetsHelp, vmCmd)
	makeString("ExportLimit", &a.Config.VM.ExportLimit, []string{"limit", "size", "export-limit"}, exportAssetsCmd)
	_ = makeCommand("start", vmApp.ExportAssetsStart, exportAssetsCmd)
	_ = makeCommand("status", vmApp.ExportAssetsStatus, exportAssetsCmd)
	_ = makeCommand("get", vmApp.ExportAssetsGet, exportAssetsCmd)
	_ = makeCommand("query", vmApp.ExportAssetsQuery, exportAssetsCmd)
	makeString("Chunk", &a.Config.VM.Chunk, []string{"chunk", "chunks"}, exportAssetsCmd)

	sListCmd := makeCommand("scanners", vmApp.ScannersList, vmCmd)
	_ = makeCommand("list", vmApp.ScannersList, sListCmd)

	aListCmd := makeCommand("agents", vmApp.AgentsList, vmCmd)
	_ = makeCommand("list", vmApp.AgentsList, aListCmd)
	_ = makeCommand("group", vmApp.AgentsGroup, aListCmd)
	_ = makeCommand("ungroup", vmApp.AgentsUngroup, aListCmd)

	makeBool("WithoutGroupName", &a.Config.VM.WithoutGroupName, []string{"without-group", "no-groups"}, aListCmd)
	makeString("GroupName", &a.Config.VM.GroupName, []string{"group", "groupname", "group-name"}, aListCmd)

	aGroupsCmd := makeCommand("agent-groups", vmApp.AgentGroupsList, vmCmd)
	_ = makeCommand("list", vmApp.AgentGroupsList, aGroupsCmd)

	a.RootCmd.SetUsageTemplate(a.DefaultUsage)
	a.RootCmd.SetHelpTemplate(a.DefaultUsage)

	return
}

// InvokeCLI passes control over to the root cobra command.
func (a *App) InvokeCLI() {
	// Enable 'client' log file, since we are invoke the client.
	clientLog := a.Config.VM.EnableLogging()
	serverLog := a.Config.Server.EnableLogging()

	port := a.Config.Server.ListenPort
	shouldServer := (a.Config.DefaultServerStart == true) && !isProxyServerCmd() && cmdproxy.PortAvailable(port)

	setDefaultRootCmd()

	if shouldServer {
		clientLog.Infof("Starting a proxy server for the client...")
		proxy := pkgproxy.NewServer(a.Config, a.Metrics, serverLog)
		proxy.EnableDefaultRouter()
		go proxy.ListenAndServe()
	}

	// Call Cobra Execute which will PreRun and select the Command to execute.
	_ = a.RootCmd.Execute()

	if shouldServer {
		cmdproxy.Stop(a.Config, a.Metrics)
	}

	return
}

func (a *App) Usage() string {
	versionMap := map[string]string{"ReleaseVersion": vm.ReleaseVersion, "GitHash": vm.GitHash}
	return a.commandUsageTmpl("tioUsage", versionMap)
}

// usageTemplate renders the usage/help/man pages for a cmd
func (a *App) commandUsageTmpl(name string, data interface{}) string {
	var err error
	var templateFiles []string

	templateFiles = append(templateFiles, ApplicationName)

	t := template.New("")

	file, err := config.TemplateFolder.Open("/template/cmd/tiogo.tmpl")
	if err != nil {
		log.Fatal(err)
		return ""
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	t, err = t.Funcs(
		template.FuncMap{
			"Gopher": ui.Gopher,
		},
	).Parse(string(content))

	var raw bytes.Buffer
	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		a.Config.VM.Log.Fatalf("error executing help usage template for tiogo: %v", err)
		return ""
	}

	return raw.String()
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

func isProxyServerCmd() bool {
	return contains(os.Args, "proxy")
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
