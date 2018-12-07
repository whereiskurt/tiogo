package ui

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/whereiskurt/tiogo/internal/pkg/dao"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Draw struct {
	CLI    *CLI
	CSV    *CSV
	JSON   *JSON
	Errorf func(fmt string, args ...interface{})
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
}

type JSON struct {
	Draw   *Draw
	Output *os.File
}
type CSV struct {
	Draw   *Draw
	Output *os.File
}

func NewCSV(d *Draw) (csv *CSV) {
	csv = new(CSV)
	csv.Draw = d
	csv.Output = d.CLI.Output
	return
}
func NewJSON(d *Draw) (j *JSON) {
	j = new(JSON)
	j.Draw = d
	j.Output = d.CLI.Output
	return
}

func NewDraw(c *CLI) (d *Draw) {
	d = new(Draw)
	d.CLI = c
	d.CSV = NewCSV(d)
	d.JSON = NewJSON(d)
	d.Errorf = c.Config.Logger.Errorf
	d.Infof = c.Config.Logger.Infof
	d.Debugf = c.Config.Logger.Debugf
	d.Warnf = c.Config.Logger.Warnf
	return
}

func (d *Draw) Banner() {
	d.Infof(`
  _   _                            _   ___  
 | |_(_) ___   __ _  ___    __   _/ | / _ \ 
 | __| |/ _ \ / _' |/ _ \   \ \ / / || | | |
 | |_| | (_) | (_| | (_) |   \ V /| || |_| |
  \__|_|\___(_)__, |\___/     \_/ |_(_)___/ 
              |___/                         
                           tio.go version 1.0`)
	// http://patorjk.com/software/taag/#p=author&f=Ivrit&t=tio-cli%20%20v0.5
	return
}

func (d *Draw) Gopher() {
	fmt.Printf(`
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
@https://gist.github.com/belbomemo`)
	return
}

func (d *Draw) Version() {

	fmt.Printf(`
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
@https://gist.github.com/belbomemo`)
	return
}

func (c *CSV) Scans(scans []dao.Scan) (err error) {
	if len(scans) == 0 {
		return
	}
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

func GroupMembership(groups []string, scanagent map[string]dao.ScannerAgentGroup) (members []string) {
	for g := range groups {
		if _, ok := scanagent[groups[g]]; ok {
			members = append(members, "1")
		} else {
			members = append(members, "0")
		}
	}
	return
}

func (c *CSV) Template(tname string, sh interface{}) (err error) {
	t := template.New("")

	output := c.Draw.CLI.Output

	t.Funcs(
		template.FuncMap{
			"StringsJoin":     strings.Join,
			"StringsSplit":    strings.Split,
			"Contains":        strings.Contains,
			"CSVString":       CSVString,
			"Base64":          Base64,
			"GroupMembership": GroupMembership,
		},
	).ParseGlob("config/template/csv/*.tmpl")

	if err != nil {
		c.Draw.Errorf("%v", err)
		return
	}

	err = t.ExecuteTemplate(output, tname, sh)
	if err != nil {
		c.Draw.Errorf("error in Execute template: %v", err)
	}
	return
}

func (j *JSON) ScanHistory(sh dao.ScanHistory) (err error) {
	output := j.Draw.CLI.Output

	bb, err := json.Marshal(sh)
	if err != nil {
		j.Draw.Errorf("%v", err)
		return
	}

	bb, err = PrettyPrintJSON(bb)
	if err != nil {
		return
	}
	output.WriteString(string(bb))

	return
}

func PrettyPrintJSON(bb []byte) (pp []byte, err error) {
	var pretty bytes.Buffer
	raw := bb
	cmd := exec.Command("jq", ".")
	cmd.Stdin = strings.NewReader(string(raw))
	cmd.Stdout = &pretty
	err = cmd.Run()
	if err == nil {
		pp = []byte(pretty.String())
	}
	return
}
