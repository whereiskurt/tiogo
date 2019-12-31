// Will generate a 'templates_generate.go' with all of the files under this folder
// This is necessary because a binary program can be run from anywhere on the filesystem and
// may not have a relative folder './config/template/'.  Using vfsgen we create a static go file
// with contents of the templates embedded.  This is done with build tags.
package main

//go:generate go run vfsgen_templates.go

import (
	"net/http"

	"github.com/shurcooL/vfsgen"
	"github.com/sirupsen/logrus"
)

func main() {
	outputFilename := "../pkg/config/vfsgenerate.go"
	err := vfsgen.Generate(http.Dir("./"), vfsgen.Options{
		Filename:     outputFilename,
		PackageName:  "config",
		BuildTags:    "release",
		VariableName: "BinaryEmbedFolder",
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	outputFilename = "../internal/app/cmd/vfsgenerate.go"
	err = vfsgen.Generate(http.Dir("../internal/app/cmd/"), vfsgen.Options{
		Filename:     outputFilename,
		PackageName:  "cmd",
		BuildTags:    "release",
		VariableName: "CmdHelpEmbed",
	})
	if err != nil {
		logrus.Fatalln(err)
	}

}
