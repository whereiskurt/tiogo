package vm

import (
	"github.com/whereiskurt/tiogo/internal/app"
	"github.com/whereiskurt/tiogo/internal/pkg/ui"
)

type Tag struct {
	Config *app.Config
	// Convenience functions for logging out.
	Infof  func(fmt string, args ...interface{})
	Debugf func(fmt string, args ...interface{})
	Warnf  func(fmt string, args ...interface{})
	Errorf func(fmt string, args ...interface{})
}

func NewTag(c *app.Config) (t *Tag) {
	t = new(Tag)
	t.Config = c
	t.Errorf = t.Config.Logger.Errorf
	t.Debugf = t.Config.Logger.Debugf
	t.Warnf = t.Config.Logger.Warnf
	t.Infof = t.Config.Logger.Infof
	return
}

func (cmd *Tag) Main(cli *ui.CLI) (err error) {
	// config := t.Config
	// a := adapter.NewAdapter(config)
	// a.Tags()

	return
}
