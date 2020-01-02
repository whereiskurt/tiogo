package ui

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"

	"github.com/whereiskurt/tiogo/internal/app/cmd"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/config"
)

// CLI makes the text output to the terminal.
type CLI struct {
	Config *config.Config
}

// NewCLI takes a configuration used for describing how to output.
func NewCLI(c *config.Config) (cli CLI) {
	cli.Config = c
	return
}

// DrawGopher outputs a text gopher to stdout
func (cli *CLI) DrawGopher() {
	fmt.Println(Gopher())
	return
}

// Gopher returns a string printable gopher! Thanks to belbomemo!
func Gopher() string {
	gopher := `
	         ,_---~~~~~----._         
	  _,,_,*^____      _____''*g*\"*, 
	 / __/ /'     ^.  /      \ ^@q   f 
	[  @f | @))    |  | @))   l  0 _/  
	 \'/   \~____ / __ \_____/    \   
	  |           _l__l_           I   
	  }          [______]           I  
	  ]            | | |            |  
	  ]             ~ ~             |  
	  |                            |   
	
	[[@https://gist.github.com/belbomemo]]
	`
	return gopher
}

// AgentGroupHeader will output a CSV header for AgentGroup[] passed
func AgentGroupHeader(ag []client.AgentGroup) (s string) {
	var buf bytes.Buffer

	var ss []string
	for i := range ag {
		ss = append(ss, ag[i].Name)
	}
	sort.Strings(ss)

	w := csv.NewWriter(&buf)
	if err := w.Write(ss); err != nil {
		return
	}
	w.Flush()

	s = buf.String()
	s = strings.TrimSpace(s)
	return
}

// AgentGroupNameJoin takes an AgentGroup map and joins with sep
func AgentGroupNameJoin(agent map[string]client.AgentGroup, sep string) string {
	if agent == nil || len(agent) == 0 {
		return ""
	}

	var ss []string
	for _, v := range agent {
		ss = append(ss, v.Name)
	}
	sort.Strings(ss)
	return strings.Join(ss, sep)
}

// AgentGroupMembership returns an array of 1-or-0 for each agent group (ie. ["0","0","1","1","0"])
func AgentGroupMembership(agent map[string]client.AgentGroup, groups []client.AgentGroup) (members []string) {

	var ss []string
	for _, g := range groups {
		ss = append(ss, g.Name)
	}
	sort.Strings(ss)

	for _, g := range ss {
		// Check if g.Name is list of groups
		if _, ok := agent[g]; ok {
			members = append(members, "1")
		} else {
			members = append(members, "0")
		}
	}
	return members
}

// CSVString takes ss[] strings and outputs a CSV string
func CSVString(ss []string) (s string) {
	var buf bytes.Buffer

	w := csv.NewWriter(&buf)
	if err := w.Write(ss); err != nil {
		return
	}
	w.Flush()

	s = buf.String()
	s = strings.TrimSpace(s)
	return
}

// Base64 takes a raw string and Base64 encodes it
func Base64(raw string) (encoded string) {
	encoded = string(base64.StdEncoding.EncodeToString([]byte(raw)))
	return
}

// Render will output the UI templates as per the config bind the data.
func (cli *CLI) Render(name string, data interface{}) (usage string) {
	var raw bytes.Buffer
	var err error
	var log = cli.Config.VM.Log
	var templateFiles []string

	templateFiles = append(templateFiles, "tio.tmpl")

	templateFiles = append(templateFiles, "vm/agent-groups.tmpl")
	templateFiles = append(templateFiles, "vm/agents.tmpl")
	templateFiles = append(templateFiles, "vm/cache.tmpl")
	templateFiles = append(templateFiles, "vm/export-assets.tmpl")
	templateFiles = append(templateFiles, "vm/export-scans.tmpl")
	templateFiles = append(templateFiles, "vm/export-vulns.tmpl")
	templateFiles = append(templateFiles, "vm/scanners.tmpl")
	templateFiles = append(templateFiles, "vm/scans.tmpl")
	templateFiles = append(templateFiles, "vm/vm.tmpl")

	templateFiles = append(templateFiles, "proxy/server.tmpl")
	templateFiles = append(templateFiles, "proxy/start.tmpl")
	templateFiles = append(templateFiles, "proxy/stop.tmpl")

	t := template.New("")
	for _, f := range templateFiles {
		file, err := cmd.CmdHelpEmbed.Open(fmt.Sprintf("%s", f))
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Errorf("Couldn't load template file: %s: %s", fmt.Sprintf("%s", f), err)
			return "Error: couldn't produce usage."
		}

		t, err = t.Funcs(
			template.FuncMap{
				"Gopher":               Gopher,
				"AgentGroupHeader":     AgentGroupHeader,
				"AgentGroupMembership": AgentGroupMembership,
				"AgentGroupNameJoin":   AgentGroupNameJoin,
				"StringsSplit":         strings.Split,
				"ToUpper":              strings.ToUpper,
				"ToLower":              strings.ToLower,
				"Contains":             strings.Contains,
				"CSVString":            CSVString,
				"Base64":               Base64,
			},
		).Parse(string(content))
	}

	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		log.Fatalf("error in Execute template: %v", err)
	}

	usage = raw.String()
	return
}

// Println wraps printing to STDOUT
func (cli *CLI) Println(line string) {
	fmt.Println(line)
	return
}

// Fatal wraps log.Fatalf
func (cli *CLI) Fatal(line string) {
	fmt.Println(line)
	cli.Config.VM.Log.Fatalf(line)
	return
}

// Fatalf wraps log.Fatal
func (cli *CLI) Fatalf(line string, params ...interface{}) {
	cli.Fatal(fmt.Sprintf(line, params...))
	return
}

// Error wraps log.Error
func (cli *CLI) Error(line string) {
	cli.Config.VM.Log.Error(line)
	fmt.Println(line)
	return
}

// Errorf wraps log.Errorf
func (cli *CLI) Errorf(line string, params ...interface{}) {
	cli.Error(fmt.Sprintf(line, params...))
	return
}
