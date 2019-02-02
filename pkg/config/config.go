package config

import (
	"bufio"
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config holds all parameters for the application and is structured based on the command hierarchy
type Config struct {
	// Global parameters available to all commands
	Context        context.Context
	HomeFolder     string
	HomeFilename   string
	ConfigFolder   string
	ConfigFilename string
	TemplateFolder string
	LogFolder      string
	Log            *log.Logger
	VerboseLevel   string
	VerboseLevel1  bool
	VerboseLevel2  bool
	VerboseLevel3  bool
	VerboseLevel4  bool
	VerboseLevel5  bool

	VM     VMClient
	Server ServerConfig
}

// ClientConfig are all of the params for the Client Command
type VMClient struct {
	Config        *Config
	BaseURL       string
	AccessKey     string
	SecretKey     string
	CacheKey      string
	CacheFolder   string
	CacheResponse bool
	MetricsFolder string
	OutputCSV     bool
	OutputJSON    bool

	ID       string
	Name     string
	Regex    string
	JQex     string
	UUID     string
	Critical bool
	High     bool
	Medium   bool
	Info     bool

	Chunk      string
	BeforeDate string // Date bounding YYYY-MM-DD
	AfterDate  string // Date bounded for YYYY-MM-DD
	Days       string // Number of days to include either before or after

	Category         string // Asset Category
	Tag              string // Asset Tagging
	Members          string
	GroupName        string
	WithoutGroupName bool
	AsTargetGroups   bool
	Size             string
	NeverRun         bool
	Onetime          bool
	Date             string
	Time             string
	EmailAddress     string
}

// String allows us to mask the Keys we don't want to reveal
func (c *Config) String() string {
	var safeConfig = new(Config)

	spew.Config.MaxDepth = 2

	// With DisableMethods, String() will be recursively called on the *Config sub-elements and blow the stack. :-)
	spew.Config.DisableMethods = true

	// Copy config that was passed
	*safeConfig = *c

	// Overwrite sensitive values with the masked value
	mask := "[**MASKED**]"
	safeConfig.VM.AccessKey = mask
	safeConfig.VM.SecretKey = mask
	safeConfig.VM.CacheKey = mask
	safeConfig.Server.AccessKey = mask
	safeConfig.Server.SecretKey = mask
	safeConfig.Server.CacheKey = mask

	s := spew.Sdump(safeConfig)

	return s
}

// ServerConfig are all of the params for the Client Command
type ServerConfig struct {
	Config            *Config
	ListenPort        string
	ServiceBaseURL    string
	AccessKey         string
	SecretKey         string
	CacheKey          string
	CacheFolder       string
	CacheResponse     bool
	MetricsListenPort string
	MetricsFolder     string
}

// NewConfig returns config that has default values set and is hooked to cobra/viper (if invoked.)
func NewConfig() (c *Config) {
	c = new(Config)
	c.SetToDefaults()

	c.Context = context.Background()
	c.Log = log.New()

	cobra.OnInitialize(func() {
		for i := range os.Args {
			if strings.ToLower(os.Args[i]) == "help" {
				return
			}
		}
		// Only read configuration when not invoked with 'help'
		c.readWithViper()
	})

	// Provide access to c variables - ie. log!
	c.VM.Config = c
	c.Server.Config = c

	return
}

func (c *VMClient) DumpMetrics() {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("client.%d.%s.prom", pid, dts)
	file := filepath.Join(".", c.MetricsFolder, name)
	metrics.DumpMetricsToFile(file)
}
func (c *ServerConfig) DumpMetrics() {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("server.%d.%s.prom", pid, dts)
	file := filepath.Join(".", c.MetricsFolder, name)
	metrics.DumpMetricsToFile(file)
}

func (c *VMClient) EnableLogging() *log.Logger {
	filename := c.LogFilename()
	c.Config.SetLogFilename(filename)
	return c.Config.Log
}
func (c *ServerConfig) EnableLogging() {
	filename := c.LogFilename()
	dirName := filepath.Dir(c.LogFilename())

	err := os.MkdirAll(dirName, 0777)
	if err != nil {
		log.Fatalf("error: making folder for log folder: '%s'", err)
	}

	c.Config.SetLogFilename(filename)
}

// UnmarshalViper copies all of the cobra/viper config data into our Config struct
// This is the delineation between cobra/viper and using our Config struct.
func (c *Config) UnmarshalViper() {
	// Copy everything from the Viper into our Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return
}

func (c *Config) readWithViper() {
	var err error

	defaultFilename := DefaultConfigFilename + "." + DefaultConfigType
	f, err := TemplateFolder.Open(defaultFilename)
	defer f.Close()
	err = viper.ReadConfig(f)
	if err != nil {
		c.Log.Fatalf("fatal: couldn't read default config: %s", err)
	}

	filename := filepath.Join(c.HomeFolder, c.HomeFilename)
	filename = filename + "." + DefaultConfigType

	viper.AddConfigPath(c.HomeFolder)
	viper.SetConfigName(c.HomeFilename)
	err = viper.MergeInConfig()
	if err != nil {
		// First run, try and get user inputted configuration
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			c.userInputConfiguration(filename)
		}
		err = viper.MergeInConfig()
		if err != nil {
			c.Log.Warnf("warning: couldn't viper MergeInConfig after default config: %s", err)
			return
		}
	}

	viper.AutomaticEnv()
	return
}

func (c *Config) userInputConfiguration(filename string) bool {

	// home := c.HomeFolder

	t := time.Now()
	ts := fmt.Sprintf("%v", t)
	z := strings.Split(ts, " ")
	tzDefault := fmt.Sprintf("%s %s", z[2], z[3])

	fmt.Println()
	fmt.Println(fmt.Sprintf("WARN: "+"No configuration file '%s' found", filename))
	fmt.Println()
	fmt.Print(fmt.Sprintf("Need Tenable.IO access keys and secret keys for API usage."))
	fmt.Println()
	fmt.Println(fmt.Sprintf("You must provide the X-ApiKeys '" + "accessKey" + "' and '" + "secretKey" + "' to the Tenable.IO API."))
	fmt.Println(fmt.Sprintf("For complete details see: https://cloud.tenable.com/api#/authorization"))
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(fmt.Sprintf("Enter required Tenable.io" + "'AccessKey'" + ": "))
	c.VM.AccessKey, _ = reader.ReadString('\n')
	c.VM.AccessKey = strings.TrimSpace(c.VM.AccessKey)
	if len(c.VM.AccessKey) != 64 {
		c.Log.Warnf(fmt.Sprintf("Invalid accessKey '%s' length %d not 64.\n\n", c.VM.AccessKey, len(c.VM.AccessKey)))
		return false
	}

	fmt.Print(fmt.Sprintf("Enter required Tenable.io" + "'SecretKey'" + ": "))
	c.VM.SecretKey, _ = reader.ReadString('\n')
	c.VM.SecretKey = strings.TrimSpace(c.VM.SecretKey)
	if len(c.VM.SecretKey) != 64 {
		c.Log.Warnf(fmt.Sprintf("Invalid secretKey '%s' length %d not 64.\n\n", c.VM.SecretKey, len(c.VM.SecretKey)))
		return false
	}
	c.VM.CacheKey = fmt.Sprintf("%s%s", c.VM.AccessKey[:16], c.VM.SecretKey[:16])

	fmt.Println()
	fmt.Print(fmt.Sprintf("Save configuration file? [yes or default:no]: "))
	shouldSave, _ := reader.ReadString('\n')
	fmt.Println()

	if len(shouldSave) > 0 && strings.ToUpper(shouldSave)[0] == 'Y' {
		// Sensible defaults - users homedir are usually writeable for a cache

		fmt.Println(fmt.Sprintf("Creating default configuration file '%s' ...", filename))

		file, err := os.Create(filename)
		if err != nil {
			c.Log.Warnf(fmt.Sprintf("Cannot create default configuration file '%s':%s", filename, err))
			return false
		}
		defer file.Close()

		fmt.Fprintf(file, "#################################################\n")
		fmt.Fprintf(file, "## Successfully created by tiogo commandline tool\n")
		fmt.Fprintf(file, "#################################################\n")
		fmt.Fprintf(file, "\n")
		fmt.Fprintf(file, "VM:\n")
		fmt.Fprintf(file, "  ServiceBaseURL: %s\n", c.VM.BaseURL)
		fmt.Fprintf(file, "  AccessKey: %s\n", c.VM.AccessKey)
		fmt.Fprintf(file, "  SecretKey: %s\n", c.VM.SecretKey)
		fmt.Fprintf(file, "  CacheKey: %s%s\n", c.VM.AccessKey[:16], c.VM.SecretKey[:16])
		fmt.Fprintf(file, "  CacheFolder: %s\n", c.VM.CacheFolder)
		fmt.Fprintf(file, "  DefaultTimezone: %s\n", tzDefault)
		fmt.Fprintf(file, "\n")
		fmt.Fprintf(file, "Server:\n")
		fmt.Fprintf(file, "  ServiceBaseURL: %s\n", c.Server.ServiceBaseURL)
		fmt.Fprintf(file, "  CacheFolder: %s\n", c.Server.CacheFolder)
		fmt.Fprintf(file, "  ListenPort: %s\n", c.Server.ListenPort)
		fmt.Fprintf(file, "\n")
		fmt.Println(fmt.Sprintf("\n\nDone! \nSuccessfully created '%s'", filename))
		fmt.Println()

		return true
	}

	return false
}
