package client

import (
	"00-newapp-template/pkg/cache"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// CacheLabel is the type for where to store the response
type CachePathLabel string

func (c CachePathLabel) String() string {
	return "adapter/" + string(c)
}

// Adapter is used to call ACME services and convert them to Gopher/Things in Go structures we like.
type Adapter struct {
	Config    *config.Config
	Metrics   *metrics.Metrics
	Unmarshal Unmarshal
	Filter    *Filter
	Convert   Converter
	Worker    *sync.WaitGroup
	DiskCache *cache.Disk
}

// NewAdapter manages calls the remote services, converts the results and manages a memory/disk cache.
func NewAdapter(config *config.Config, metrics *metrics.Metrics) (a *Adapter) {
	a = new(Adapter)
	a.Config = config
	a.Metrics = metrics
	a.Worker = new(sync.WaitGroup)
	a.Unmarshal = NewUnmarshal(config, metrics)
	a.Filter = NewFilter(config)
	a.Convert = NewConvert()
	if a.Config.Client.CacheResponse {
		a.DiskCache = cache.NewDisk(a.Config.Client.CacheFolder, a.Config.Client.CacheKey, a.Config.Client.CacheKey != "")
	}

	return
}

func (a *Adapter) diskStore(label CachePathLabel, obj interface{}) {
	j, err := json.Marshal(obj)
	if err == nil {
		_ = a.DiskCache.Store(fmt.Sprintf("%s.json", label), PrettyJSON(j))
	}
}

// PrettyJSON will look for 'jq' to pretty the json input
func PrettyJSON(json []byte) []byte {
	jq, err := exec.LookPath("jq")
	if err == nil {
		var pretty bytes.Buffer
		cmd := exec.Command(jq, ".")
		cmd.Stdin = strings.NewReader(string(json))
		cmd.Stdout = &pretty
		err := cmd.Run()
		if err == nil {
			json = []byte(pretty.String())
		}
	}
	return json
}

// GopherThings populates each gopher with their things
func (a *Adapter) GopherThings() map[string]Gopher {
	var matchOnThings = false
	if a.Config.Client.ThingID != "" || a.Config.Client.ThingName != "" || a.Config.Client.ThingDescription != "" {
		matchOnThings = true
	}

	a.Metrics.ClientInc("GopherThings", metrics.Methods.Service.Get)
	gopherThings := make(map[string]Gopher)

	gg := a.Gophers()
	for _, g := range gg {
		things := a.Things(g.ID)

		// If there are no 'things' for this gopher and we are filtering for a thing
		// don't add this 'gopher' to the results
		if len(things) == 0 && matchOnThings {
			continue
		}
		gopherThings[g.ID] = Gopher{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Things:      things,
		}
	}

	a.diskStore(CachePathLabel("GopherThings"), &gopherThings)
	return gopherThings
}

// Gophers returns all gophers with 'things' == nil
func (a *Adapter) Gophers() map[string]Gopher {
	a.Metrics.ClientInc("Gophers", metrics.Methods.Service.Get)

	rawGophers := a.Unmarshal.gophers()
	filtered := a.Filter.gophers(rawGophers)
	gophers := a.Convert.gophers(filtered)

	a.diskStore(CachePathLabel("Gophers"), &gophers)

	return gophers
}

// Things will return all things for a gopherID
func (a *Adapter) Things(gopherID string) map[string]Thing {
	a.Metrics.ClientInc("Things", metrics.Methods.Service.Get)

	rawThings := a.Unmarshal.things(gopherID)
	filtered := a.Filter.things(rawThings)
	things := a.Convert.things(filtered)

	label := CachePathLabel(fmt.Sprintf("Things/Gopher.%s", gopherID))
	a.diskStore(label, &things)

	return things
}

// DeleteGopher will delete the matching gopherID
func (a *Adapter) DeleteGopher(gopherID string) map[string]Gopher {
	a.Metrics.ClientInc("Gopher", metrics.Methods.Service.Delete)
	rawGophers := a.Unmarshal.deleteGopher(gopherID)
	gophers := a.Convert.gophers(rawGophers)
	return gophers
}

// DeleteThing will delete the Thing matching gopherID and thingID - could use FindGopherByThing instead of taking thingID
func (a *Adapter) DeleteThing(gopherID string, thingID string) map[string]Thing {
	a.Metrics.ClientInc("Thing", metrics.Methods.Service.Delete)
	rawThings := a.Unmarshal.deleteThing(gopherID, thingID)
	things := a.Convert.things(rawThings)
	return things
}

// FindGopherByThing returns the Gopher ID to the associated Thing by ID.
func (a *Adapter) FindGopherByThing(thingID string) string {

	allGophers := a.Unmarshal.gophers()
	gophers := a.Convert.gophers(allGophers)

	for g := range gophers {
		rawThings := a.Unmarshal.things(gophers[g].ID)
		things := a.Convert.things(rawThings)
		for t := range things {
			if string(things[t].ID) == thingID {
				return gophers[g].ID
			}
		}
	}

	return ""
}

// UpdateGopher uses the details in newGopher to update the Gopher
func (a *Adapter) UpdateGopher(newGopher Gopher) (gopher Gopher) {
	a.Metrics.ClientInc("Gopher", metrics.Methods.Service.Update)
	a.Unmarshal.updateGopher(newGopher)
	return
}

// AddGopher
func (a *Adapter) AddGopher(newGopher Gopher) Gopher {
	a.Metrics.ClientInc("Gopher", metrics.Methods.Service.Add)
	a.Unmarshal.addGopher(newGopher)
	return newGopher
}

// AddThing
func (a *Adapter) AddThing(newThing Thing) Thing {
	a.Metrics.ClientInc("Thing", metrics.Methods.Service.Add)
	a.Unmarshal.addThing(newThing)
	return newThing
}

// UpdateThing uses the details in newThing to update the Thing
func (a *Adapter) UpdateThing(newThing Thing) (thing Thing) {
	a.Metrics.ClientInc("Thing", metrics.Methods.Service.Update)

	if newThing.Gopher.ID == "" {
		newThing.Gopher.ID = a.FindGopherByThing(newThing.ID)
		if newThing.Gopher.ID == "" {
			return thing
		}
	}

	a.Unmarshal.updateThing(newThing)

	return

}
