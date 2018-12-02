package app

import "sync"

type StatType string

type Stats struct {
	ThreadSafe *sync.Mutex
	statMap    map[StatType]interface{}
	counterMap map[StatType]int
}

func NewStats() (s *Stats) {
	s = new(Stats)
	s.ThreadSafe = new(sync.Mutex)
	s.statMap = make(map[StatType]interface{})
	s.counterMap = make(map[StatType]int)
	return
}

func (s Stats) Tick(key StatType) {
}
func (s Stats) Tock(key StatType) {
}

func (s Stats) Count(key StatType) {
	s.ThreadSafe.Lock()
	s.counterMap[key]++ // Not thread safe :-)
	s.ThreadSafe.Unlock()
	return
}

func (s Stats) CountMap() (c map[StatType]int) {
	s.ThreadSafe.Lock()
	c = s.counterMap // Make a copy.
	s.ThreadSafe.Unlock()
	return
}
