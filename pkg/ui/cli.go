package ui

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/config"
	"io/ioutil"
	"strings"
	"text/template"
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

func AgentGroupsHeader(ag []client.AgentGroup) (s string) {
	var buf bytes.Buffer

	var ss []string
	for i := range ag {
		ss = append(ss, ag[i].Name)
	}

	w := csv.NewWriter(&buf)
	if err := w.Write(ss); err != nil {
		return
	}
	w.Flush()

	s = buf.String()
	s = strings.TrimSpace(s)
	return
}

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

func Base64(raw string) (encoded string) {
	encoded = string(base64.StdEncoding.EncodeToString([]byte(raw)))
	return
}

// GroupMembership returns an array of 1-or-0 for each agent group (ie. ["0","0","1","1","0"])
func GroupMembership(agent map[string]client.AgentGroup, groups []client.AgentGroup) (members []string) {
	for _, g := range groups {
		// Check if g.Name is list of groups
		if _, ok := agent[g.Name]; ok {
			members = append(members, "1")
		} else {
			members = append(members, "0")
		}
	}
	return members
}

// Render will output the UI templates as per the config bind the data.
func (cli *CLI) Render(name string, data interface{}) (usage string) {
	var raw bytes.Buffer
	var err error

	var log = cli.Config.VM.Log
	// TODO: Replace this with an 'index' concept - needs to be generated. vfsgen types/methods not visible.
	var templateFiles []string
	templateFiles = append(templateFiles, "template/client/table.tmpl")
	templateFiles = append(templateFiles, "template/client/csv.tmpl")
	templateFiles = append(templateFiles, "template/cmd/tiogo.tmpl")
	templateFiles = append(templateFiles, "template/cmd/vm/vm.tmpl")

	t := template.New("")
	for _, f := range templateFiles {
		file, err := config.TemplateFolder.Open(fmt.Sprintf("%s", f))
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Errorf("Couldn't load template file: %s: %s", fmt.Sprintf("%s", f), err)
			return "Error: couldn't produce usage."
		}

		t, err = t.Funcs(
			template.FuncMap{
				"Gopher":            Gopher,
				"AgentGroupsHeader": AgentGroupsHeader,
				"GroupMembership":   GroupMembership,
				"StringsJoin":       strings.Join,
				"StringsSplit":      strings.Split,
				"Contains":          strings.Contains,
				"CSVString":         CSVString,
				"Base64":            Base64,
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

func (cli *CLI) Println(line string) {
	fmt.Println(line)
	return
}

func (cli *CLI) Fatal(line string) {
	fmt.Println(line)
	cli.Config.VM.Log.Fatalf(line)
	return
}

func (cli *CLI) Fatalf(line string, params ...interface{}) {
	cli.Fatal(fmt.Sprintf(line, params...))
	return
}

func (cli *CLI) Error(line string) {
	cli.Config.VM.Log.Error(line)
	fmt.Println(line)
	return
}

func (cli *CLI) Errorf(line string, params ...interface{}) {
	cli.Error(fmt.Sprintf(line, params...))
	return
}
