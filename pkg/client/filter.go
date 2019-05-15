package client

import (
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/config"
	"regexp"
)

// Filter is used adapter to remove unneeded results (ie. only matching Gophers)
type Filter struct {
	Config *config.Config
}

// NewFilter loops through in[] and keeps/skips matching items based on attributes.
func NewFilter(config *config.Config) (filter *Filter) {
	filter = new(Filter)
	filter.Config = config
	return
}

func (f *Filter) AgentGroupsByRegex(agents []AgentGroup, regex string) (filtered []AgentGroup) {
	var r = regexp.MustCompile(regex)

	for _, a := range agents {
		if r.MatchString(a.Name) == true {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
func (f *Filter) AgentsByRegex(agents []ScannerAgent, regex string) (filtered []ScannerAgent) {
	var r = regexp.MustCompile(regex)

	for _, a := range agents {
		s := fmt.Sprintf("%+v", a)
		if r.MatchString(s) == true {
			filtered = append(filtered, a)
			// f.Config.Log.Debugf("AgentString:%s", s)
		}
	}

	return filtered
}
func (f *Filter) AgentsByName(agents []ScannerAgent, name string) (filtered []ScannerAgent) {
	for _, a := range agents {
		if a.Name == name {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func (f *Filter) KeepOnlyGroupMembers(agents []ScannerAgent, group string) (filtered []ScannerAgent) {
	shouldKeepMatch := true
	filtered = GroupMembership(agents, group, shouldKeepMatch)
	return filtered
}
func (f *Filter) SkipGroupMembers(agents []ScannerAgent, group string) (filtered []ScannerAgent) {
	shouldKeepMatch := false
	filtered = GroupMembership(agents, group, shouldKeepMatch)
	return filtered
}

func GroupMembership(agents []ScannerAgent, group string, shouldKeepMatch bool) (filtered []ScannerAgent) {
	for _, a := range agents {
		if _, ok := a.Groups[group]; ok == shouldKeepMatch {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
