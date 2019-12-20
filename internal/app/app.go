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

	proxy := cmdproxy.NewServer(a.Config, a.Metrics)
	proxyCmd := makeCommand("proxy", proxy.ServerHelp, a.RootCmd)
	_ = makeCommand("help", proxy.ServerHelp, proxyCmd)
	_ = makeCommand("start", proxy.Start, proxyCmd)
	_ = makeCommand("stop", proxy.Stop, proxyCmd)

	app := vm.NewVM(a.Config, a.Metrics)
	appCmd := makeCommand("vm", app.Help, a.RootCmd)
	makeString("ID", &a.Config.VM.ID, []string{"i", "id"}, appCmd)
	makeString("UUID", &a.Config.VM.UUID, []string{"uuid"}, appCmd)
	makeString("Name", &a.Config.VM.Name, []string{"n", "name"}, appCmd)
	makeString("Regex", &a.Config.VM.Regex, []string{"regex"}, appCmd)
	makeString("JQex", &a.Config.VM.JQex, []string{"jqex"}, appCmd)
	makeBool("CSV", &a.Config.VM.OutputCSV, []string{"csv"}, appCmd)
	makeBool("JSON", &a.Config.VM.OutputJSON, []string{"json"}, appCmd)

	makeBool("Critical", &a.Config.VM.Critical, []string{"critical", "crit"}, appCmd)
	makeBool("High", &a.Config.VM.High, []string{"high"}, appCmd)
	makeBool("Medium", &a.Config.VM.Medium, []string{"medium", "med"}, appCmd)
	makeBool("Info", &a.Config.VM.Info, []string{"info", "low"}, appCmd)

	_ = makeCommand("help", app.Help, appCmd)

	exportVulns := makeCommand("export-vulns", app.ExportVulnsHelp, appCmd)
	makeString("ExportLimit", &a.Config.VM.ExportLimit, []string{"limit", "size", "export-limit"}, exportVulns)
	_ = makeCommand("start", app.ExportVulnsStart, exportVulns)
	_ = makeCommand("status", app.ExportVulnsStatus, exportVulns)
	_ = makeCommand("get", app.ExportVulnsGet, exportVulns)
	_ = makeCommand("query", app.ExportVulnsQuery, exportVulns)
	makeString("Chunk", &a.Config.VM.Chunk, []string{"chunk", "chunks"}, exportVulns)
	makeString("BeforeDate", &a.Config.VM.BeforeDate, []string{"before"}, exportVulns)
	makeString("AfterDate", &a.Config.VM.AfterDate, []string{"after"}, exportVulns)
	makeString("Days", &a.Config.VM.Days, []string{"days"}, exportVulns)

	exportAssetsCmd := makeCommand("export-assets", app.ExportAssetsHelp, appCmd)
	makeString("ExportLimit", &a.Config.VM.ExportLimit, []string{"limit", "size", "export-limit"}, exportAssetsCmd)
	_ = makeCommand("start", app.ExportAssetsStart, exportAssetsCmd)
	_ = makeCommand("status", app.ExportAssetsStatus, exportAssetsCmd)
	_ = makeCommand("get", app.ExportAssetsGet, exportAssetsCmd)
	_ = makeCommand("query", app.ExportAssetsQuery, exportAssetsCmd)
	makeString("Chunk", &a.Config.VM.Chunk, []string{"chunk", "chunks"}, exportAssetsCmd)

	sListCmd := makeCommand("scanners", app.ScannersList, appCmd)
	_ = makeCommand("list", app.ScannersList, sListCmd)

	aListCmd := makeCommand("agents", app.AgentsList, appCmd)
	_ = makeCommand("list", app.AgentsList, aListCmd)
	_ = makeCommand("group", app.AgentsGroup, aListCmd)
	_ = makeCommand("ungroup", app.AgentsUngroup, aListCmd)

	makeBool("WithoutGroupName", &a.Config.VM.WithoutGroupName, []string{"without-group", "no-groups"}, aListCmd)
	makeString("GroupName", &a.Config.VM.GroupName, []string{"group", "groupname", "group-name"}, aListCmd)

	aGroupsCmd := makeCommand("agent-groups", app.AgentGroupsList, appCmd)
	_ = makeCommand("list", app.AgentGroupsList, aGroupsCmd)

	//
	cacheCmd := makeCommand("cache", app.CacheInfo, appCmd)
	_ = makeCommand("list", app.CacheInfo, cacheCmd)
	cacheClearCmd := makeCommand("clear", app.CacheClear, cacheCmd)
	_ = makeCommand("all", app.CacheClearAll, cacheClearCmd)
	_ = makeCommand("agents", app.CacheClearAgents, cacheClearCmd)
	_ = makeCommand("scans", app.CacheClearScans, cacheClearCmd)

	a.RootCmd.SetUsageTemplate(a.DefaultUsage)
	a.RootCmd.SetHelpTemplate(a.DefaultUsage)

	return
}

// InvokeCLI passes control over to the root cobra command.
func (a *App) InvokeCLI() {
	// Enable 'client' log file, since we are invoke the client.
	clientLog := a.Config.VM.EnableLogging()
	serverLog := a.Config.Server.EnableLogging()

	//a.Config.IsServerPortAvailable()
	port := a.Config.Server.ListenPort
	shouldServer := (a.Config.DefaultServerStart == true) && !isProxyServerCmd() && cmdproxy.IsPortAvailable(port)

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
		defer cmdproxy.Stop(a.Config, a.Metrics)
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
