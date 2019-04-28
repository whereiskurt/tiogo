package ui

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/client"
	"github.com/whereiskurt/tiogo/pkg/config"
	"io/ioutil"
	"log"
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
		ss =append(ss, ag[i].Name)
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

// GroupMembership returns a 1-or-0 entry for each agent's group membership (ie. [0,0,1,1,0])
// - 'groups' is a list of all scanner groups ("Group 1, Group 2, Group 3, Group 4, Group 5")
// - 'agentGroups' is map of group names containing the agents groups agent is a part of (Group 3=>{ID:1,Name:Group 3})
// - 'members' is all the groups and the membership for the agent (ie. [0,0,1,1,0])
func GroupMembership(groups []client.AgentGroup, agentGroups map[string]client.AgentGroup) (members []string) {
	for _, name := range groups {
		if _, ok := agentGroups[name.Name]; ok {
			members = append(members, "1")
		} else {
			members = append(members, "0")
		}
	}
	return
}


// Render will output the UI templates as per the config bind the data.
func (cli *CLI) Render(name string, data interface{}) (usage string) {
	var raw bytes.Buffer
	var err error

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
			log.Fatal(err)
		}

		t, err = t.Funcs(
			template.FuncMap{
				"Gopher": Gopher,
				"AgentGroupsHeader": AgentGroupsHeader,
				"GroupMembership": GroupMembership,
				"StringsJoin":     strings.Join,
				"StringsSplit":    strings.Split,
				"Contains":        strings.Contains,
				"CSVString":       CSVString,
				"Base64":          Base64,
			},
		).Parse(string(content))
	}

	if err != nil {
		log.Fatalf("couldn't load template: %v", err)
	}

	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		log.Fatalf("error in Execute template: %v", err)
	}

	usage = raw.String()
	return
}
