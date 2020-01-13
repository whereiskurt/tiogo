// Will generate 'vfs_generated.go' files with everything under the specific path
// This is necessary because a binary program can be run from anywhere on the filesystem and
// may not have a relative folder './config/template/'.  Using vfsgen we create a static go file
// with contents of the templates embedded.  This is done with build tags.
package main

//go:generate go run generate_vfs_release.go

import (
	"net/http"

	"github.com/shurcooL/vfsgen"
	"github.com/sirupsen/logrus"
)

// Runs VFSv
func main() {

	// Embeded config and jq folders
	err := vfsgen.Generate(http.Dir("../config/"), vfsgen.Options{
		Filename:     "../pkg/config/vfs_generated.go",
		PackageName:  "config",
		BuildTags:    "release",
		VariableName: "BinaryEmbedFolder",
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	// Embed the CLI tmpl files i
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
