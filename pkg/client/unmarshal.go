package client

import (
	"github.com/whereiskurt/tiogo/pkg/config"
	"github.com/whereiskurt/tiogo/pkg/metrics"
	"github.com/whereiskurt/tiogo/pkg/tenable"
	"path/filepath"
)

// Unmarshal holds the config - needed for Service.... TODO: Remove config and take Service
type Unmarshal struct {
	Config  *config.Config
	Metrics *metrics.Metrics
}

// NewUnmarshal calls the ACME EndPoints and returns ACME JSONs to the adapter
func NewUnmarshal(config *config.Config, metrics *metrics.Metrics) (u Unmarshal) {
	u.Config = config
	u.Metrics = metrics
	return
}

func (u *Unmarshal) service() (s tenable.Service) {
	s = tenable.NewService(u.Config.VM.BaseURL, u.Config.VM.SecretKey, u.Config.VM.AccessKey)
	s.EnableMetrics(u.Metrics)

	if u.Config.VM.CacheResponse {
		serviceCacheFolder := filepath.Join(".", u.Config.VM.CacheFolder, "service/")
		s.EnableCache(serviceCacheFolder, u.Config.VM.CacheKey)
	}
	s.SetLogger(u.Config.Log)

	return
}

// func (u *Unmarshal) updateThing(thing Thing) (tt tenable.Thing) {
// 	s := u.service()
//
// 	var t = tenable.Thing{
// 		GopherID:    json.Number(thing.Gopher.ID),
// 		ID:          json.Number(thing.ID),
// 		Description: thing.Description,
// 		Name:        thing.Name,
// 	}
// 	tt = s.UpdateThing(t)
// 	return
// }
