package config

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"os"
	"path/filepath"
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
	ID             string
	Name           string
	Regex          string
	JQex           string
	OutputCSV      string
	OutputJSON     string

	VM     VMConfig
	Server ServerConfig
}

// ClientConfig are all of the params for the Client Command
type VMConfig struct {
	Config        *Config
	BaseURL       string
	AccessKey     string
	SecretKey     string
	CacheKey      string
	CacheFolder   string
	CacheResponse bool
	MetricsFolder string
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
	AccessKey         string
	SecretKey         string
	CacheKey          string
	CacheFolder       string
	CacheResponse     bool
	MetricsListenPort string
	MetricsFolder     string
}

// NewConfig returns config that has default values set and is hooked to cobra/viper (if invoked.)
func NewConfig() (config *Config) {
	config = new(Config)
	config.SetToDefaults()

	config.Context = context.Background()
	config.Log = log.New()

	cobra.OnInitialize(func() {
		config.readWithViper()
	})

	// Provide access to config variables - ie. log!
	config.VM.Config = config
	config.Server.Config = config

	return
}

func (c *VMConfig) DumpMetrics() {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("client.%d.%s.prom", pid, dts)
	file := filepath.Join(".", c.MetricsFolder, name)
	metrics.DumpMetrics(file)
}
func (c *ServerConfig) DumpMetrics() {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("server.%d.%s.prom", pid, dts)
	file := filepath.Join(".", c.MetricsFolder, name)
	metrics.DumpMetrics(file)
}

func (c *VMConfig) EnableLogging() {
	filename := c.LogFilename()
	c.Config.SetLogFilename(filename)
}
func (c *ServerConfig) EnableLogging() {
	filename := c.LogFilename()
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

	f, err := TemplateFolder.Open(defaultConfigFilename + "." + defaultConfigType)
	err = viper.ReadConfig(f)
	if err != nil {
		c.Log.Fatalf("fatal: couldn't read in config: %s", err)
	}

	viper.AddConfigPath(c.HomeFolder)
	viper.SetConfigName(c.HomeFilename)
	err = viper.MergeInConfig()
	if err != nil {
		c.Log.Infof("warning: couldn't load configs from: %s or : %s: %s", c.HomeFolder, c.HomeFilename, err)
	}

	viper.AutomaticEnv()

	return
}
