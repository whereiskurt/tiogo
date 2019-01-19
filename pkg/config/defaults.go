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
const DefaultHomeFilename = ".tiogo"
const DefaultConfigFilename = "default.tiogo"
const DefaultConfigType = "yaml"
const DefaultConfigFolder = "./config/"
const DefaultTemplateFolder = "./config/template/"

// Sensible defaults even with out a configuration file present
const DefaultVerboseLevel = "3"
const DefaultServerListenPort = "10101"
const DefaultMetricsListenPort = "22222"

const DefaultClientBaseURL = "http://localhost:" + DefaultServerListenPort
const DefaultServerBaseURL = "https://cloud.tenable.com"

// Used by the *_test to the set defaults
// DefaultClientCacheFolder stores default client cache file location
const DefaultClientCacheFolder = ".tiogo/cache/client/"
const DefaultClientCacheResponse = true

const DefaultLogFolder = "./log/"
const DefaultServerMetricsFolder = "./log/metrics/server/"
const DefaultClientMetricsFolder = "./log/metrics/client/"

// DefaultServerCacheFolder  stores default server cache file location
const DefaultServerCacheFolder = ".tiogo/cache/server/"
const DefaultServerCacheResponse = true

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

	c.LogFolder = DefaultLogFolder
	c.VM.CacheFolder = c.HomeFolder + "/" + DefaultClientCacheFolder
	c.VM.CacheResponse = DefaultClientCacheResponse
	c.Server.BaseURL = DefaultServerBaseURL

	c.Server.CacheFolder = c.HomeFolder + "/" + DefaultServerCacheFolder
	c.Server.CacheResponse = DefaultServerCacheResponse
	c.Server.MetricsFolder = DefaultServerMetricsFolder

	c.Server.MetricsListenPort = DefaultMetricsListenPort
	c.Server.ListenPort = DefaultServerListenPort
	c.VerboseLevel = DefaultVerboseLevel
	c.ConfigFolder = DefaultConfigFolder
	c.ConfigFilename = DefaultConfigFilename
	c.TemplateFolder = DefaultTemplateFolder
	c.VM.MetricsFolder = DefaultClientMetricsFolder

}
func (c *Config) SetLogFilename(filename string) {
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
