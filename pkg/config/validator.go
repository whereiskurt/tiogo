package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// ValidateOrFatal will validate the string values inside of the Config after copying from Unmarshal or self-setting.
func (c *Config) ValidateOrFatal() {

	c.validateVerbosity()

	c.validateChunks()
	c.validateDateBounds()

	return
}

func (c *Config) validateChunks() {
	if c.VM.Chunk == "" {
		c.VM.Chunk = "ALL"
		return
	}

	var chunks []string
	for _, v := range strings.Split(c.VM.Chunk, ",") {
		// Expand chunks expressed as --chunk=1-100,102-103
		if strings.Contains(v, "-") {
			r := strings.Split(v, "-")
			lower, err := strconv.Atoi(r[0])
			if err != nil {
				log.Fatalf("error: invalid lower bound for chunk range:%s", r)
			}
			upper, err := strconv.Atoi(r[0])
			if err != nil {
				log.Fatalf("error: invalid lower bound for chunk range:%s", r)
			}

			var vv []string
			for i := lower; i <= upper; i++ {
				vv = append(vv, fmt.Sprintf("%d", i))
			}
			v = strings.Join(vv, ",")
		}
		chunks = append(chunks, v)
	}

}

// DateLayout is the default date expected generally anywhere a date is requested
var DateLayout = "2006-01-02 15:04:05 -0700 MST"

func (c *Config) validateDateBounds() {

	var defaultDays = 365

	var days = 0
	if c.VM.Days != "" {
		d, err := strconv.Atoi(c.VM.Days)
		if err != nil {
			log.Fatalf("error: invalid days: %d", days)
		}
		days = d
	}

	if c.VM.BeforeDate == "" && c.VM.AfterDate == "" {
		if c.VM.Days == "" {
			days = defaultDays
		}
		now := time.Now()
		c.VM.BeforeDate = now.Format(DateLayout)
		c.VM.AfterDate = now.AddDate(0, 0, -1*days).Format(DateLayout)
		return
	}

	if c.VM.BeforeDate != "" && c.VM.AfterDate != "" {
		if c.VM.Days != "" {
			log.Fatalf("Setting '--days' value with --before and --after parameters is not supported.")
		}
	}

	//TODO: This check is insufficient and is only used add 00:00:00 to the begin, and the only adds 23:59:59 to the end...
	if len(c.VM.AfterDate) < 11 {
		c.VM.AfterDate = fmt.Sprintf("%s 00:00:00 %s", c.VM.AfterDate, c.VM.DefaultTimezone)
		log.Debugf("Update AfterDate: %s", c.VM.AfterDate)
	}
	if len(c.VM.BeforeDate) < 11 {
		c.VM.BeforeDate = fmt.Sprintf("%s 23:59:59 %s", c.VM.BeforeDate, c.VM.DefaultTimezone)
		log.Debugf("Update BeforeDate: %s", c.VM.BeforeDate)
	}

	// If we are missing a BEFORE date
	if c.VM.BeforeDate == "" {
		if c.VM.Days != "" {
			log.Fatalf("Must set --days with --before ")
		}

		after, err := time.Parse(DateLayout, c.VM.AfterDate)
		if err != nil {
			log.Fatalf("error: invalid after date: %s: %s", after, err)
		}
		c.VM.BeforeDate = after.AddDate(0, 0, 1*days).Format(DateLayout)
		return
	}

	// If we are missing an AFTER date
	if c.VM.AfterDate == "" {
		if c.VM.Days != "" {
			log.Fatalf("Must set --days with --after")
		}
		before, err := time.Parse(DateLayout, c.VM.BeforeDate)
		if err != nil {
			log.Fatalf("error: invalid before date: %s: %s", before, err)
		}
		c.VM.AfterDate = before.AddDate(0, 0, -1*days).Format(DateLayout)
		return
	}

	return
}

func (c *Config) validateVerbosity() {
	if c.hasVerboseLevel() {
		switch {
		case c.VerboseLevel5:
			c.VerboseLevel = "5"
		case c.VerboseLevel4:
			c.VerboseLevel = "4"
		case c.VerboseLevel3:
			c.VerboseLevel = "3"
		case c.VerboseLevel2:
			c.VerboseLevel = "2"
		case c.VerboseLevel1:
			c.VerboseLevel = "1"
		}
	}

	switch c.VerboseLevel {
	case "5":
		c.VerboseLevel5 = true
		c.Server.Log.SetLevel(log.TraceLevel)
		c.VM.Log.SetLevel(log.TraceLevel)
	case "4":
		c.VerboseLevel4 = true
		c.Server.Log.SetLevel(log.DebugLevel)
		c.VM.Log.SetLevel(log.DebugLevel)
	case "3":
		c.VerboseLevel3 = true
		c.VM.Log.SetLevel(log.InfoLevel)
		c.Server.Log.SetLevel(log.InfoLevel)
	case "2":
		c.VerboseLevel2 = true
		c.VM.Log.SetLevel(log.WarnLevel)
		c.Server.Log.SetLevel(log.WarnLevel)
	case "1":
		c.VerboseLevel1 = true
		c.VM.Log.SetLevel(log.ErrorLevel)
		c.Server.Log.SetLevel(log.ErrorLevel)
	}

	if !c.hasVerboseLevel() {
		log.Fatalf("invalid VerboseLevel: '%s'", c.VerboseLevel)
	}

}
func (c *Config) hasVerboseLevel() bool {
	return c.VerboseLevel1 || c.VerboseLevel2 || c.VerboseLevel3 || c.VerboseLevel4 || c.VerboseLevel5
}
