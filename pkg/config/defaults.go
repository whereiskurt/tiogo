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
const defaultConfigFilename = "default.tiogo"
const defaultConfigType = "yaml"
const defaultConfigFolder = "./config/"
const defaultTemplateFolder = "./config/template/"
const defaultHomeFilename = ".tiogo"

// Sensible defaults even with out a configuration file present
const defaultVerboseLevel = "3"
const defaultServerListenPort = "10101"
const defaultMetricsListenPort = "22222"

const defaultBaseURL = "http://localhost:" + defaultServerListenPort

// Used by the *_test to the set defaults
// DefaultClientCacheFolder stores default client cache file location
const DefaultClientCacheFolder = "./.cache/client/"
const defaultClientCacheResponse = true

const defaultLogFolder = "./log/"
const defaultServerMetricsFolder = "./log/metrics/server/"
const defaultClientMetricsFolder = "./log/metrics/client/"

// DefaultServerCacheFolder  stores default server cache file location
const DefaultServerCacheFolder = "./.cache/server/"
const defaultServerCacheResponse = true

func (c *Config) SetToDefaults() {
	c.VM.BaseURL = defaultBaseURL

	c.LogFolder = defaultLogFolder
	c.VM.CacheFolder = DefaultClientCacheFolder
	c.VM.CacheResponse = defaultClientCacheResponse
	c.Server.CacheFolder = DefaultServerCacheFolder
	c.Server.CacheResponse = defaultServerCacheResponse
	c.Server.MetricsFolder = defaultServerMetricsFolder

	c.Server.MetricsListenPort = defaultMetricsListenPort
	c.Server.ListenPort = defaultServerListenPort
	c.VerboseLevel = defaultVerboseLevel
	c.ConfigFolder = defaultConfigFolder
	c.ConfigFilename = defaultConfigFilename
	c.TemplateFolder = defaultTemplateFolder
	c.VM.MetricsFolder = defaultClientMetricsFolder

	// Find the User's home folder
	folder, err := home.Dir()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to detect home directory: %v", err))
	} else {
		c.HomeFolder = folder
	}
	c.HomeFilename = defaultHomeFilename
}
func (c *Config) SetLogFilename(filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// If DEBUG level is on we output log to STDOUT as well.
	mw := io.MultiWriter(f)
	if c.Log.IsLevelEnabled(log.DebugLevel) {
		mw = io.MultiWriter(os.Stdout, f)
	}
	c.Log.SetOutput(mw)

	c.Log.SetFormatter(&log.TextFormatter{})
}

func (c *ServerConfig) LogFilename() string {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("server.%d.%s.log", pid, dts)
	file := filepath.Join(".", c.Config.LogFolder, name)
	return file
}
func (c *VMConfig) LogFilename() string {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("client.%d.%s.log", pid, dts)
	file := filepath.Join(".", c.Config.LogFolder, name)
	return file
}
