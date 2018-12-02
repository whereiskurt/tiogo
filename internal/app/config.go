package app

import (
	"fmt"
	home "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
	"time"
)

const DefaultConfigFile = "default.tiogo"
const DefaultConfigFolder = "./config/"
const DefaultHomeConfigFile = ".tiogo"

type Config struct {
	HomeFolder     string // Params
	HomeConfig     string // Name of config file used by RootCmd/Viper/etc
	ConfigFolder   string // Name of config folder
	AccessKey      string // API crypto Access Key
	SecretKey      string // API crypto Secret Key
	BaseURL        string // Tenable.IO or on premise
	ExecutionDTS   string
	VerboseLevel   string
	CacheKey       string // Disk cache crypto key
	CacheFolder    string
	CachePretty    bool
	LogFilename    string
	LogFolder      string
	StdoutFilename string // Base filename to output
	StdoutFolder   string // Folder to write output to
	DefaultTZ      string // Time Zone to work from
	WorkerCount    string // How many workers goroutines?
	PerfProfile    bool

	// Modes
	CryptoCacheMode  bool
	ClobberCacheMode bool
	NoColourMode     bool // No colour
	OutputCSVMode    bool // WriteCSV
	OutputJSONMode   bool // JSON

	// Colour           *ui.Colour
	OutputFileHandle *os.File
	LogFileHandle    *os.File
	Logger           *Logger

	VM *VMConfig
}
type VMConfig struct {
	AppConfig *Config
	// CRUD Modes
	CreateMode   bool
	ReadOnlyMode bool
	UpdateMode   bool
	DeleteMode   bool

	// Actions
	TagMode          bool
	UntagMode        bool
	AgentGroupMode   bool
	AgentUngroupMode bool
	ExportMode       bool

	// Views for output
	SearchView bool // Views
	DetailView bool
	ListView   bool

	// PARAMS for VM
	Depth            string // ScanHistoryDetail histories to include
	ScanID           string // ScanHistoryDetail
	IgnoreScanID     string
	HistoryID        string // ScanHistory
	IgnoreHistoryID  string
	HostID           string // Host
	IgnoreHostID     string
	AssetUUID        string // Asset
	IgnoreAssetUUID  string
	PluginID         string // Plugins
	IgnorePluginID   string
	AgentUUID        string // Agents
	AgentGroupName   string
	NoAgentGroupName bool
	Regex            string // REGEX
	TagCategory      string // Tag
	TagValue         string
	TemplateName     string // Name of template to render
	ExportUUID       string // Name of template to render
}

func NewConfig(cob *cobra.Command) (c *Config) {
	c = new(Config)

	c.Logger = new(Logger)

	cobra.OnInitialize(c.InitViperConfig)
	// These defaults are needed for Viper and can't be passed.
	c.HomeConfig = DefaultHomeConfigFile
	c.ConfigFolder = DefaultConfigFolder

	dts := time.Now().UTC().Format("20060102")

	c.LogFilename = fmt.Sprintf("log.%s.log", dts)

	c.OutputFileHandle = os.Stdout
	c.LogFileHandle = os.Stderr

	if c.HomeFolder == "" {
		hdir, hErr := home.Dir()
		if hErr != nil {
			log.Fatal(fmt.Sprintf("failed to detect home directory: %v", hErr))
		} else {
			c.HomeFolder = hdir
		}
	}

	ParamBool(cob, "CryptoCacheMode", &c.CryptoCacheMode, []string{"useCryptoCacheMode"}, "True/False - Enables encryption on the cache folder.")
	ParamBool(cob, "CachePretty", &c.CachePretty, []string{"cachePretty", "jq"}, "Tries to pretty print the cache entries using jq.")
	ParamBool(cob, "PerfProfile", &c.PerfProfile, []string{"perf", "perfProfile"}, "Capture Golang perf log.")

	ParamString(cob, "CacheFolder", &c.CacheFolder, []string{}, "Cache folder - sensitive data.")
	ParamString(cob, "LogFolder", &c.LogFolder, []string{}, "Log storage folder..")
	ParamString(cob, "LogFilename", &c.LogFilename, []string{}, "Name of log file for execution")
	ParamString(cob, "WorkerCount", &c.WorkerCount, []string{"worker", "workers", "workerCount"}, "How many Go routines to have as workers?")
	ParamString(cob, "DefaultTZ", &c.DefaultTZ, []string{"defaultTimezone"}, "Default TZ from any scanner we don't know.")
	ParamString(cob, "VerboseLevel", &c.VerboseLevel, []string{"verbose"}, "Larger values more vebose (5 max.)")

	ParamBool(cob, "ClobberCacheMode", &c.ClobberCacheMode, []string{"clobberCache"}, "When cache lookups fail, but file exists, overwrite with a fresh request.")
	ParamBool(cob, "OutputCSVMode", &c.OutputCSVMode, []string{"csv"}, "AssetWorker in csv format (default:false)")
	ParamBool(cob, "OutputJSONMode", &c.OutputJSONMode, []string{"json"}, "AssetWorker in JSON format (default:false)")

	return
}
func NewVMConfig(c *Config, cob *cobra.Command) (v *VMConfig) {
	v = new(VMConfig)
	c.VM = v // Link v config to c config
	v.AppConfig = c

	ParamBool(cob, "ListView", &v.ListView, []string{"listview", "list"}, "Sets the view to output each element in a list ")
	ParamBool(cob, "DetailView", &v.DetailView, []string{"detailview", "detail", "details"}, "Sets the view to a detailed output of each element (default:false)")
	ParamBool(cob, "AgentGroupMode", &v.AgentGroupMode, []string{"agentgroupmode", "assign"}, "Assign Agents to groups (default:false)")
	ParamBool(cob, "AgentUngroupMode", &v.AgentUngroupMode, []string{"agentungroupmode", "unassign"}, "Unassign Agents froms groups (default:false)")
	ParamBool(cob, "ExportMode", &v.ExportMode, []string{"export"}, "Use the exports API to pull assets,vulns,scans (default:false)")
	ParamBool(cob, "NoAgentGroupName", &v.NoAgentGroupName, []string{"agentnogroupname", "nogroup", "ng"}, "Set when you only want empty groups (default:false)")

	ParamString(cob, "ScanID", &v.ScanID, []string{"scanid", "scan", "sid", "s"}, "The scans ids to include in output (default:[empty] which is all)")
	ParamString(cob, "PluginID", &v.PluginID, []string{"pluginid", "plugin", "pid", "p"}, "The plugin ids to include in output (default:[empty] which is all)")
	ParamString(cob, "Regex", &v.Regex, []string{"regex", "rex", "r"}, "A REGEX to match against.")
	ParamString(cob, "AgentGroupName", &v.AgentGroupName, []string{"agentgroupname", "agentgroup", "group", "g"}, "An agent group name to work with.")
	ParamString(cob, "TemplateName", &v.TemplateName, []string{"templatename", "template", "t"}, "The name of the template to draw.")
	ParamString(cob, "ExportUUID", &v.ExportUUID, []string{"exportuuid", "UUID", "uuid"}, "The export UUID to fetch.")
	ParamString(cob, "Depth", &v.Depth, []string{"depth", "d"}, "How many past histories to pull.")

	return
}

func (c *Config) InitViperConfig() {
	var err error

	c.ExecutionDTS = time.Now().UTC().Format(time.RFC3339)

	viper.AddConfigPath(c.ConfigFolder)
	viper.SetConfigName(DefaultConfigFile)
	_ = viper.ReadInConfig()

	viper.AddConfigPath(c.HomeFolder)
	viper.SetConfigName(c.HomeConfig)
	err = viper.MergeInConfig()
	if err != nil {
		log.Fatalf("fatal: couldn't load config from home folder: '%s' from '%s': %s", c.HomeConfig, c.HomeFolder, err)
		return
	}

	viper.AutomaticEnv()

	return
}
func (c *Config) Finalize() {

	if !(c.StdoutFolder == "" || c.StdoutFilename == "") {
		file := fmt.Sprintf("%s/%s", c.StdoutFolder, c.StdoutFilename)
		filemode := os.O_CREATE | os.O_WRONLY | os.O_APPEND

		// If the file doesn't exist, create it, or append to the file
		fh, err := os.OpenFile(file, filemode, 0644)
		if err != nil {
			log.Fatal(err)
		}
		c.OutputFileHandle = fh

	}

	if c.LogFolder == "" {
		c.LogFolder = "."
	}

	file := fmt.Sprintf("%s/%s", c.LogFolder, c.LogFilename)
	filemode := os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// If the file doesn't exist, create it, or append to the file
	fh, err := os.OpenFile(file, filemode, 0644)
	if err != nil {
		log.Fatal(err)
	}
	c.LogFileHandle = fh
	c.Logger = NewLogger(c.VerboseLevel, fh)

	if c.PerfProfile {
		f, perr := os.Create(fmt.Sprintf("%s/tiogo.prof", c.LogFolder))
		if perr != nil {
			log.Fatal(perr)
		}
		perr = pprof.StartCPUProfile(f)
		if perr != nil {
			log.Fatal(perr)
		}
	}

}

func ReflectFromViper(config interface{}) {
	opts := reflect.ValueOf(config).Elem()

	for i := 0; i < opts.NumField(); i++ {
		name := opts.Type().Field(i).Name
		val := reflect.ValueOf(config).Elem().FieldByName(name)

		if !val.IsValid() {
			continue
		}

		switch opts.Field(i).Kind() {
		case reflect.String:
			val.SetString(viper.GetString(name))
			break
		case reflect.Bool:
			val.SetBool(viper.GetBool(name))
			break

		default:
			// TODO: Add 'is struct' reflection concepts to support Auth/etc.
		}
	}
	return
}

func ParamBool(cob *cobra.Command, name string, ref *bool, aliases []string, desc string) {
	cob.PersistentFlags().BoolVar(ref, name, *ref, desc)
	_ = viper.BindPFlag(name, cob.PersistentFlags().Lookup(name))
	for i, alias := range aliases {
		cob.PersistentFlags().BoolVar(ref, alias, *ref, desc)
		if i > 0 {
			cob.PersistentFlags().Lookup(alias).Hidden = true
		}

		viper.RegisterAlias(alias, name)
	}
}
func ParamString(cob *cobra.Command, name string, ref *string, aliases []string, desc string) {
	cob.PersistentFlags().StringVar(ref, name, *ref, desc)
	_ = viper.BindPFlag(name, cob.PersistentFlags().Lookup(name))
	for i, alias := range aliases {
		cob.PersistentFlags().StringVar(ref, alias, *ref, desc)
		if i > 0 {
			cob.PersistentFlags().Lookup(alias).Hidden = true
		}
		viper.RegisterAlias(alias, name)
	}

}
