package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

// ValidateOrFatal will validate the string values inside of the Config after copying from Unmarshal or self-setting.
func (c *Config) ValidateOrFatal() {
	c.validateVerbosity()
	c.validateOutputMode()

	var err error
	err = os.MkdirAll(path.Dir(c.VM.MetricsFolder), 0777)
	if err != nil {
		log.Fatalf("error: making folder for client metrics: '%s'", err)
	}
	err = os.MkdirAll(path.Dir(c.Server.MetricsFolder), 0777)
	if err != nil {
		log.Fatalf("error: making folder for server metrics: '%s'", err)
	}
	err = os.MkdirAll(path.Dir(c.LogFolder), 0777)
	if err != nil {
		log.Fatalf("error: making folder for log folder: '%s'", err)
	}

	return
}

func (c *Config) validateOutputMode() {
	switch strings.ToLower(c.VM.OutputMode) {
	case "csv":
	case "json":
	case "xml":
	case "table":

	default:
		log.Fatalf("invalid OutputMode: '%s'", c.VM.OutputMode)
	}
}
func (c *Config) validateVerbosity() {
	if c.hasVerboseLevel() {
		switch {
		case c.VerboseLevel1:
			c.VerboseLevel = "1"
		case c.VerboseLevel2:
			c.VerboseLevel = "2"
		case c.VerboseLevel3:
			c.VerboseLevel = "3"
		case c.VerboseLevel4:
			c.VerboseLevel = "4"
		case c.VerboseLevel5:
			c.VerboseLevel = "5"
		}
	}

	switch c.VerboseLevel {
	case "5":
		c.VerboseLevel5 = true
		c.Log.SetLevel(log.TraceLevel)
	case "4":
		c.VerboseLevel4 = true
		c.Log.SetLevel(log.DebugLevel)
	case "3":
		c.VerboseLevel3 = true
		c.Log.SetLevel(log.InfoLevel)
	case "2":
		c.VerboseLevel1 = true
		c.Log.SetLevel(log.WarnLevel)
	case "1":
		c.VerboseLevel1 = true
		c.Log.SetLevel(log.ErrorLevel)
	}

	if !c.hasVerboseLevel() {
		log.Fatalf("invalid VerboseLevel: '%s'", c.VerboseLevel)
	}

}
func (c *Config) hasVerboseLevel() bool {
	return c.VerboseLevel1 || c.VerboseLevel2 || c.VerboseLevel3 || c.VerboseLevel4 || c.VerboseLevel5
}
