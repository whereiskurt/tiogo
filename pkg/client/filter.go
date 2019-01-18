package client

import (
	"github.com/whereiskurt/tiogo/pkg/config"
)

// Filter is used adapter to remove unneeded results (ie. only matching Gophers)
type Filter struct {
	Config *config.Config
}

// NewFilter loops through in[] and keeps/skips matching items based on attributes.
func NewFilter(config *config.Config) (filter *Filter) {
	filter = new(Filter)
	filter.Config = config
	return
}
