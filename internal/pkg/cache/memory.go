package cache

import "github.com/karlseguin/ccache"

const ItemsToPrune = 50
const MaxSize = 500000

type Memory struct {
	Cache *ccache.Cache
}

func NewMemoryCache() (m *Memory) {
	m = new(Memory)
	m.Cache = ccache.New(ccache.Configure().MaxSize(MaxSize).ItemsToPrune(ItemsToPrune))

	return
}
