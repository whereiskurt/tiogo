// Will generate a 'templates_generate.go' with all of the files under this folder
// This is necessary because a binary program can be run from anywhere on the filesystem and
// may not have a relative folder './config/template/'.  Using vfsgen we create a static go file
// with contents of the templates embedded.  This is done with build tags.
package main

//go:generate go run vfsgen_templates.go

import (
	"github.com/shurcooL/vfsgen"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

func main() {
	outputFilename := "pkg/config/templates_generate.go"

	// NOTE: If we run from the IDE with a right-click our cwd() is inside of config/templates
	cwd, _ := os.Getwd()
	templateFolder := "config/"

	// Check if we're running inside the config/template folder, and adjust relative paths.
	if strings.Contains(cwd, "template/") {
		templateFolder = "./"
		outputFilename = "../" + outputFilename
	}

	err := vfsgen.Generate(http.Dir(templateFolder), vfsgen.Options{
		Filename:     outputFilename,
		PackageName:  "config",
		BuildTags:    "release",
		VariableName: "TemplateFolder",
	})

	if err != nil {
		logrus.Fatalln(err)
	}

}
