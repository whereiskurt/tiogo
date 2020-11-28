package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdproxy "github.com/whereiskurt/tiogo/internal/app/cmd/proxy"
	"github.com/whereiskurt/tiogo/internal/app/cmd/vm"
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	pkgproxy "github.com/whereiskurt/tiogo/pkg/proxy"
	"github.com/whereiskurt/tiogo/pkg/ui"
)

var (
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

		port := a.Config.Server.ListenPort
		shouldServer := (a.Config.DefaultServerStart == true) && !isProxyServerCmd() && cmdproxy.IsPortAvailable(port)
		if shouldServer {
			// Enable 'client' log file, since we are invoke the client.
			serverLog := a.Config.Server.EnableLogging()

			serverLog.Infof(fmt.Sprintf("Starting a proxy server for the client: %s:%s", a.Config.Server.CacheFolder, a.Config.Server.CacheKey))

			proxy := pkgproxy.NewServer(a.Config, a.Metrics, serverLog)
			proxy.EnableCache(a.Config.Server.CacheFolder, a.Config.Server.CacheKey)
			proxy.EnableDefaultRouter()
			go proxy.ListenAndServe()
			defer cmdproxy.Stop(a.Config, a.Metrics)
		}
	}
	flagB("VerboseLevel1", &a.Config.VerboseLevel1, []string{"s", "silent"}, a.RootCmd)
	flagB("VerboseLevel2", &a.Config.VerboseLevel2, []string{"q", "quiet"}, a.RootCmd)
	flagB("VerboseLevel3", &a.Config.VerboseLevel3, []string{"v", "info"}, a.RootCmd)
	flagB("VerboseLevel4", &a.Config.VerboseLevel4, []string{"debug"}, a.RootCmd)
	flagB("VerboseLevel5", &a.Config.VerboseLevel5, []string{"trace"}, a.RootCmd)

	flagS("CryptoKey", &a.Config.CryptoKey, []string{"cryto", "key"}, a.RootCmd)

	// Define the proxy command ie. "proxy start", "proxy stop"
	proxy := cmdproxy.NewServer(a.Config, a.Metrics)
	proxyCmd := command("proxy", proxy.ServerHelp, a.RootCmd)
	subcommand("help", proxy.ServerHelp, proxyCmd)
	subcommand("start", proxy.Start, proxyCmd)
	subcommand("stop", proxy.Stop, proxyCmd)

	// Define the 'vm' command ie 'vm scans list' 'vm agents list'
	app := vm.NewVM(a.Config, a.Metrics)

	vmCmd := command("vm", app.Help, a.RootCmd)
	subcommand("help", app.Help, vmCmd)
	flagS("ID", &a.Config.VM.ID, []string{"i", "id"}, vmCmd)
	flagS("UUID", &a.Config.VM.UUID, []string{"uuid"}, vmCmd)
	flagS("Name", &a.Config.VM.Name, []string{"n", "name"}, vmCmd)
	flagS("Regex", &a.Config.VM.Regex, []string{"regex"}, vmCmd)
	flagS("JQex", &a.Config.VM.JQex, []string{"jqex"}, vmCmd)

	flagS("MaxDepth", &a.Config.VM.MaxDepth, []string{"depth"}, vmCmd)
	flagS("MaxKeep", &a.Config.VM.MaxKeep, []string{"keep"}, vmCmd)

	flagB("CSV", &a.Config.VM.OutputCSV, []string{"csv"}, vmCmd)
	flagB("JSON", &a.Config.VM.OutputJSON, []string{"json"}, vmCmd)
	flagB("Critical", &a.Config.VM.Critical, []string{"critical", "crit"}, vmCmd)
	flagB("High", &a.Config.VM.High, []string{"high"}, vmCmd)
	flagB("Medium", &a.Config.VM.Medium, []string{"medium", "med"}, vmCmd)
	flagB("Info", &a.Config.VM.Info, []string{"info", "low"}, vmCmd)

	exportVulnsCmd := command("export-vuln", app.ExportVulnsHelp, vmCmd)
	subcommand("start", app.ExportVulnsStart, exportVulnsCmd)
	subcommand("status", app.ExportVulnsStatus, exportVulnsCmd)
	subcommand("get", app.ExportVulnsGet, exportVulnsCmd)
	subcommand("query", app.ExportVulnsQuery, exportVulnsCmd)
	flagS("ExportLimit", &a.Config.VM.ExportLimit, []string{"limit", "size", "export-limit"}, exportVulnsCmd)
	flagS("Chunk", &a.Config.VM.Chunk, []string{"chunk", "chunks"}, exportVulnsCmd)
	flagS("BeforeDate", &a.Config.VM.BeforeDate, []string{"before"}, exportVulnsCmd)
	flagS("AfterDate", &a.Config.VM.AfterDate, []string{"after"}, exportVulnsCmd)
	flagS("Days", &a.Config.VM.Days, []string{"days"}, exportVulnsCmd)

	exportAssetsCmd := command("export-asset", app.ExportAssetsHelp, vmCmd)
	subcommand("start", app.ExportAssetsStart, exportAssetsCmd)
	subcommand("status", app.ExportAssetsStatus, exportAssetsCmd)
	subcommand("get", app.ExportAssetsGet, exportAssetsCmd)
	subcommand("query", app.ExportAssetsQuery, exportAssetsCmd)
	flagS("ExportLimit", &a.Config.VM.ExportLimit, []string{"limit", "size", "export-limit"}, exportAssetsCmd)
	flagS("Chunk", &a.Config.VM.Chunk, []string{"chunk", "chunks"}, exportAssetsCmd)
	flagS("AfterDate", &a.Config.VM.AfterDate, []string{"after"}, exportAssetsCmd)
	flagS("Days", &a.Config.VM.Days, []string{"days"}, exportAssetsCmd)

	scannersCmd := command("scanner", app.ScannersList, vmCmd)
	subcommand("list", app.ScannersList, scannersCmd)

	agentsCmd := command("agent", app.AgentsList, vmCmd)
	subcommand("list", app.AgentsList, agentsCmd)
	subcommand("group", app.AgentsGroup, agentsCmd)
	subcommand("ungroup", app.AgentsUngroup, agentsCmd)
	flagB("WithoutGroupName", &a.Config.VM.WithoutGroupName, []string{"without-group", "no-groups"}, agentsCmd)
	flagS("GroupName", &a.Config.VM.GroupName, []string{"group", "groupname", "group-name"}, agentsCmd)

	aGroupsCmd := command("agent-group", app.AgentGroupsList, vmCmd)
	subcommand("list", app.AgentGroupsList, aGroupsCmd)

	cacheCmd := command("cache", app.CacheInfo, vmCmd)
	subcommand("list", app.CacheInfo, cacheCmd)

	auditLogCmd := command("audit", app.AuditLogV1List, vmCmd)
	subcommand("list", app.AuditLogV1List, auditLogCmd)

	//TODO: Make all safe by adding '--all' parameter to remove historical/export outputs too
	cacheClearCmd := command("clear", app.CacheClear, cacheCmd)
	subcommand("all", app.CacheClearAll, cacheClearCmd)
	subcommand("agents", app.CacheClearAgents, cacheClearCmd)
	subcommand("scans", app.CacheClearScans, cacheClearCmd)
	subcommand("exports", app.CacheClearExports, cacheClearCmd)

	scansCmd := command("scan", app.ScansList, vmCmd)
	subcommand("list", app.ScansList, scansCmd)
	subcommand("detail", app.ScansDetail, scansCmd)
	subcommand("host", app.ScansHosts, scansCmd)
	subcommand("plugin", app.ScansPlugins, scansCmd)
	subcommand("get", app.ScansGet, scansCmd)
	subcommand("query", app.ScansQuery, scansCmd)
	flagS("HistoryUUID", &a.Config.VM.HistoryUUID, []string{"history", "history_uuid"}, scansCmd)
	flagS("Offset", &a.Config.VM.Offset, []string{"offset"}, scansCmd)

	compCmd := command("compliance", app.ComplianceGet, vmCmd)
	subcommand("get", app.ComplianceGet, compCmd)
	flagS("Offset", &a.Config.VM.Offset, []string{"offset"}, compCmd)

	exportScansCmd := command("export-scans", app.ExportScansHelp, vmCmd)
	subcommand("start", app.ExportScansStart, exportScansCmd)
	subcommand("status", app.ExportScansStatus, exportScansCmd)
	subcommand("get", app.ExportScansGet, exportScansCmd)
	subcommand("query", app.ExportScansQuery, exportScansCmd)
	subcommand("tag", app.ExportScansTag, exportScansCmd)
	subcommand("untag", app.ExportScansUntag, exportScansCmd)

	flagS("HistoryUUID", &a.Config.VM.HistoryUUID, []string{"history", "history_uuid"}, exportScansCmd)
	flagS("Offset", &a.Config.VM.Offset, []string{"offset"}, exportScansCmd)
	flagB("PDF", &a.Config.VM.OutputPDF, []string{"pdf"}, exportScansCmd)
	flagS("Chapters", &a.Config.VM.Chapters, []string{"chapter"}, exportScansCmd)
	flagS("Tags", &a.Config.VM.Tags, []string{"tag", "tags"}, exportScansCmd)

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

// Usage outputs the help related to the usage of tio.go
func (a *App) Usage() string {
	if len(os.Args) < 2 {
		cli := ui.NewCLI(a.Config)
		versionMap := map[string]string{"ReleaseVersion": vm.ReleaseVersion, "GitHash": vm.GitHash}
		fmt.Fprintf(os.Stderr, cli.Render("CommandHeader", versionMap))
		fmt.Fprintf(os.Stderr, cli.Render("tioUsage", versionMap))
	}

	return "\x00"
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

func command(s string, run func(*cobra.Command, []string), parent *cobra.Command) *cobra.Command {
	alias := []string{fmt.Sprintf("%ss", s)} // Add a pluralized alias
	child := &cobra.Command{Use: s, Run: run, PreRun: parent.PreRun, Aliases: alias}
	parent.AddCommand(child)
	return child
}

func subcommand(s string, run func(*cobra.Command, []string), parent *cobra.Command) {
	command(s, run, parent)
	return
}

func flagB(name string, ref *bool, aliases []string, cob *cobra.Command) {
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
func flagS(name string, ref *string, aliases []string, cob *cobra.Command) {
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
