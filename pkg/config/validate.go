package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// ValidateOrFatal will validate the string values inside of the Config after copying from Unmarshal or self-setting.
func (c *Config) ValidateOrFatal() {
	c.validateVerbosity()

	c.validateChunks()
	c.validateDateBounds()

	// var err error
	// err = os.MkdirAll(path.Dir(c.VM.MetricsFolder), 0777)
	// if err != nil {
	// 	log.Fatalf("error: making folder for client metrics: '%s'", err)
	// }
	// err = os.MkdirAll(path.Dir(c.Server.MetricsFolder), 0777)
	// if err != nil {
	// 	log.Fatalf("error: making folder for server metrics: '%s'", err)
	// }
	// err = os.MkdirAll(path.Dir(c.LogFolder), 0777)
	// if err != nil {
	// 	log.Fatalf("error: making folder for log folder: '%s'", err)
	// }

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
				log.Fatal("error: invalid lower bound for chunk range:%s", r)
			}
			upper, err := strconv.Atoi(r[0])
			if err != nil {
				log.Fatal("error: invalid lower bound for chunk range:%s", r)
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

var DateLayout = "2006-01-_2 15:04:05 -0700 MST"

func (c *Config) validateDateBounds() {
	now := time.Now()

	var days = 0
	if c.VM.Days != "" {
		d, err := strconv.Atoi(c.VM.Days)
		if err != nil {
			log.Fatalf("error: invalid days: %s", days)
		}
		days = d
	}

	if c.VM.BeforeDate == "" && c.VM.AfterDate == "" {
		if c.VM.Days == "" {
			days = 365
		}
		c.VM.BeforeDate = now.Format(DateLayout)
		c.VM.AfterDate = now.AddDate(0, 0, -1*days).Format(DateLayout)
		return
	}

	if c.VM.BeforeDate != "" && c.VM.AfterDate != "" {
		if c.VM.Days != "" {
			log.Fatalf("Setting '--days' value with --before and --after parameters is not supported.")
		}
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
