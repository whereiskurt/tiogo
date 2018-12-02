package dao

import "sync"

type IdType string
type IdKey string
type IdValue string

// Anon changes
type Anon struct {
	Workers    *sync.WaitGroup
	ThreadSafe *sync.Mutex
	// RemappedId is used to create a map of GUIDs to new anonimized GUIDs
	RemappedId map[IdType]map[IdKey]IdValue
}

func NewAnon() (a *Anon) {
	a = new(Anon)
	a.RemappedId = make(map[IdType]map[IdKey]IdValue)
	return
}
