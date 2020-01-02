package config

import (
	"fmt"
	home "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

// These defaults are needed to configure Viper/Cobra

// DefaultConfigType the file extension for the configuration files (using YAML instead of XML)
const DefaultConfigType = "yaml"

// DefaultHomeFilename base filename in the $HOMEDIR for the '.tiogo.v1.yaml'
const DefaultHomeFilename = ".tiogo.v1"

// DefaultConfigFolder the namefolder with the default.tiogo.v1.yaml file
const DefaultConfigFolder = "./config/"

// DefaultConfigFilename base filename for the 'default.tiogo.v1'
const DefaultConfigFilename = "default.tiogo.v1"

const DefaultTemplateFolder = "./config/template/"

// Sensible defaults even with out a configuration file present
const DefaultVerboseLevel = "3"
const DefaultServerListenPort = "10101"
const DefaultClientBaseURL = "http://localhost:" + DefaultServerListenPort
const DefaultServerBaseURL = "https://cloud.tenable.com"

// Used by the *_test to the set defaults

// DefaultClientCacheFolder stores default client cache file location
const DefaultClientCacheFolder = ".tiogo/cache/client/"

const DefaultClientCacheResponse = true
const DefaultClientCacheLookup = true

const DefaultLogFolder = "log/"
const DefaultServerMetricsFolder = "log/metrics/server/"
const DefaultMetricsListenPort = "22222"
const DefaultClientMetricsFolder = "log/metrics/client/"

// DefaultServerCacheFolder  stores default server cache file location
const DefaultServerCacheFolder = ".tiogo/cache/server/"
const DefaultServerCacheResponse = true
const DefaultServerCacheLookup = true

// SetToDefaults will use local values to set reasonable defaults
func (c *Config) SetToDefaults() {
	// Find the User's home folder
	folder, err := home.Dir()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to detect home directory: %v", err))
	} else {
		c.HomeFolder = folder
	}
	c.HomeFilename = DefaultHomeFilename

	c.VM.BaseURL = DefaultClientBaseURL

	c.LogFolder = filepath.Join(DefaultLogFolder)
	os.MkdirAll(c.LogFolder, 0777)

	c.VM.CacheFolder = filepath.Join([]string{c.HomeFolder, DefaultClientCacheFolder}...)
	c.VM.CacheResponse = DefaultClientCacheResponse
	c.VM.MetricsFolder = filepath.Join(DefaultClientMetricsFolder)

	c.Server.ServiceBaseURL = DefaultServerBaseURL
	c.Server.CacheFolder = filepath.Join([]string{c.HomeFolder, DefaultServerCacheFolder}...)
	c.Server.CacheResponse = DefaultServerCacheResponse
	c.Server.MetricsFolder = filepath.Join(DefaultServerMetricsFolder)
	c.Server.MetricsListenPort = DefaultMetricsListenPort
	c.Server.ListenPort = DefaultServerListenPort

	c.VerboseLevel = DefaultVerboseLevel
	c.ConfigFolder = DefaultConfigFolder
	c.ConfigFilename = DefaultConfigFilename
	c.TemplateFolder = DefaultTemplateFolder

	c.VM.ExportLimit = "5000" // Default asset and vulnerability export size (num_assets and chunk_size) ;-)

	c.DefaultServerStart = true

}

// SetLogFilename will set the ServerConfig log and duplicate to STDOUT
func (c *ServerConfig) SetLogFilename(filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// If DEBUG level is on we output log to STDOUT as well.
	mw := io.MultiWriter(f)
	if c.Log.IsLevelEnabled(log.TraceLevel) {
		mw = io.MultiWriter(os.Stdout, f)
	}
	c.Log.SetOutput(mw)
}

// SetLogFilename will set the VMConfig log and duplicate to STDOUT
func (c *VMConfig) SetLogFilename(filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// If DEBUG level is on we output log to STDOUT as well.
	mw := io.MultiWriter(f)
	if c.Log.IsLevelEnabled(log.TraceLevel) {
		mw = io.MultiWriter(os.Stdout, f)
	}
	c.Log.SetOutput(mw)
}

// LogFilename will set the ServerConfig logging filename
func (c *ServerConfig) LogFilename() string {
	dts := time.Now().Format("20060102")
	name := fmt.Sprintf("server.%s.log", dts)
	file := filepath.Join(".", c.Config.LogFolder, name)
	return file
}

// LogFilename will set the VMConfig logging filename
func (c *VMConfig) LogFilename() string {
	dts := time.Now().Format("20060102")
	name := fmt.Sprintf("client.%s.log", dts)
	file := filepath.Join(c.Config.LogFolder, name)
	return file
}
