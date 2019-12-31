// +build !release

// This package is only used during development to load the resources from disk.
// Using go generate will create a 'template_generate.go' will set the TemplateFolder differently

package cmd

import (
	"net/http"
	"path"
	"runtime"
)

// CmdHelpEmbed virtual
var CmdHelpEmbed http.FileSystem

func init() {
	// This needs to be set to an absolute folder path, so we derive it. :-)
	_, filename, _, _ := runtime.Caller(0)

	CmdHelpEmbed = http.Dir(path.Join(path.Dir(filename), "./"))

}
