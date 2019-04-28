package ui

import (
	"bytes"
	"fmt"
	"github.com/whereiskurt/tiogo/pkg/config"
	"io/ioutil"
	"log"
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
