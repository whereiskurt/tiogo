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
const defaultLogFolder = "log/"

// defaultConfigType the file extension for the configuration files (using YAML instead of XML)
const defaultConfigType = "yaml"

// defaultHomeFilename base filename in the $HOMEDIR for the '.tiogo.v1.yaml'
const defaultHomeFilename = ".tiogo.v1"

// defaultConfigFolder the namefolder with the default.tiogo.v1.yaml file
const defaultConfigFolder = "./config/"

// defaultConfigFilename base filename for the 'default.tiogo.v1'
const defaultConfigFilename = "default.tiogo.v1"

// Sensible defaults even with out a configuration file present
const defaultVerboseLevel = "3"
const defaultServerListenPort = "10101"
const defaultClientBaseURL = "http://localhost:" + defaultServerListenPort
const defaultServerBaseURL = "https://cloud.tenable.com"

const defaultMetricsListenPort = "22222"
const defaultServerMetricsFolder = "log/metrics/server/"
const defaultClientMetricsFolder = "log/metrics/client/"

const defaultClientCacheFolder = ".tiogo/cache/client/"
const defaultClientCacheResponse = true
const defaultClientCacheLookup = true

const defaultServerCacheFolder = ".tiogo/cache/server/"
const defaultServerCacheResponse = true
const defaultServerCacheLookup = true

// SetToDefaults will use local values to set reasonable defaults
func (c *Config) SetToDefaults() {
	// Find the User's home folder
	folder, err := home.Dir()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to detect home directory: %v", err))
	} else {
		c.HomeFolder = folder
	}
	c.HomeFilename = defaultHomeFilename

	c.VM.BaseURL = defaultClientBaseURL

	c.LogFolder = filepath.Join(defaultLogFolder)
	os.MkdirAll(c.LogFolder, 0777)

	c.VM.CacheFolder = filepath.Join([]string{c.HomeFolder, defaultClientCacheFolder}...)
	c.VM.CacheResponse = defaultClientCacheResponse
	c.VM.MetricsFolder = filepath.Join(defaultClientMetricsFolder)

	c.Server.ServiceBaseURL = defaultServerBaseURL
	c.Server.CacheFolder = filepath.Join([]string{c.HomeFolder, defaultServerCacheFolder}...)
	c.Server.CacheResponse = defaultServerCacheResponse
	c.Server.MetricsFolder = filepath.Join(defaultServerMetricsFolder)
	c.Server.MetricsListenPort = defaultMetricsListenPort
	c.Server.ListenPort = defaultServerListenPort

	c.VerboseLevel = defaultVerboseLevel
	c.ConfigFolder = defaultConfigFolder
	c.ConfigFilename = defaultConfigFilename

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
