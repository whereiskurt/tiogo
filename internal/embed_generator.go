// Will generate a 'templates_generate.go' with all of the files under this folder
// This is necessary because a binary program can be run from anywhere on the filesystem and
// may not have a relative folder './config/template/'.  Using vfsgen we create a static go file
// with contents of the templates embedded.  This is done with build tags.
package main

//go:generate go run embed_generator.go

import (
	"net/http"

	"github.com/shurcooL/vfsgen"
	"github.com/sirupsen/logrus"
)

// Runs VFSv
func main() {

	err := vfsgen.Generate(http.Dir("../config/"), vfsgen.Options{
		Filename:     "../pkg/config/vfs_generated.go",
		PackageName:  "config",
		BuildTags:    "release",
		VariableName: "BinaryEmbedFolder",
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	err = vfsgen.Generate(http.Dir("app/cmd/"), vfsgen.Options{
		Filename:     "app/cmd/vfs_generated.go",
		PackageName:  "cmd",
		BuildTags:    "release",
		VariableName: "CmdHelpEmbed",
	})
	if err != nil {
		logrus.Fatalln(err)
	}

}
