// +build !release

// This package is only used during development to load the resources from disk.
// Using go generate will create a 'template_generate.go' will set the TemplateFolder differently

package config

import (
	"net/http"
	"path"
	"runtime"
)

// TemplateFolder implements the http filesystem, but is overridden when we
// build with tags (go build -tags release) this file won't be built, but
// templates_generate.go will be.
var TemplateFolder http.FileSystem

func init() {
	// This needs to be set to an absolute folder path, so we derive it. :-)
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../config/")
	TemplateFolder = http.Dir(dir)
}
