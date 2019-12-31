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

// AgentGroupsByRegex filters matching agent groups
func (f *Filter) AgentGroupsByRegex(agents []AgentGroup, regex string) (filtered []AgentGroup) {
	var r = regexp.MustCompile(regex)

	for _, a := range agents {
		if r.MatchString(a.Name) == true {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// AgentsByRegex filters matching Agents matching regex
func (f *Filter) AgentsByRegex(agents []ScannerAgent, regex string) (filtered []ScannerAgent) {
	var r = regexp.MustCompile(regex)

	for _, a := range agents {
		s := fmt.Sprintf("%+v", a)
		if r.MatchString(s) == true {
			filtered = append(filtered, a)
		}
	}

	return filtered
}

// AgentsByName filters Agents matching name
func (f *Filter) AgentsByName(agents []ScannerAgent, name string) (filtered []ScannerAgent) {
	for _, a := range agents {
		if a.Name == name {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// KeepOnlyGroupMembers filters only agents who have the agent group
func (f *Filter) KeepOnlyGroupMembers(agents []ScannerAgent, group string) (filtered []ScannerAgent) {
	shouldKeepMatch := true
	filtered = GroupMembership(agents, group, shouldKeepMatch)
	return filtered
}

// SkipGroupMembers keeps only agents not in group
func (f *Filter) SkipGroupMembers(agents []ScannerAgent, group string) (filtered []ScannerAgent) {
	shouldKeepMatch := false
	filtered = GroupMembership(agents, group, shouldKeepMatch)
	return filtered
}

// GroupMembership is a generic filter - keeps based on shouldKeep and if group match
func GroupMembership(agents []ScannerAgent, group string, shouldKeepMatch bool) (filtered []ScannerAgent) {
	for _, a := range agents {
		if _, ok := a.Groups[group]; ok == shouldKeepMatch {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// ScanByID loops over each scans and keeps only matching
func (f *Filter) ScanByID(scans []Scan, id string) (filtered []Scan) {
	for _, a := range scans {
		if a.ScanID == id {
			filtered = append(filtered, a)
			break
		}
	}
	return filtered
}

// ScanByScheduleUUID loops over each scans and keeps only matching
func (f *Filter) ScanByScheduleUUID(scans []Scan, uuid string) (filtered []Scan) {
	for _, a := range scans {
		if a.ScheduleUUID == uuid {
			filtered = append(filtered, a)
			break
		}
	}
	return filtered
}

// ScanByName filters only the first matching by name
func (f *Filter) ScanByName(scans []Scan, name string) (filtered []Scan) {
	for _, s := range scans {
		if s.Name == name {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// ScanByRegex filters only the first matching by regex
func (f *Filter) ScanByRegex(scans []Scan, regex string) (filtered []Scan) {
	var r = regexp.MustCompile(regex)

	for _, a := range scans {
		s := fmt.Sprintf("%+v", a)
		if r.MatchString(s) == true {
			filtered = append(filtered, a)
		}
	}

	return filtered
}
